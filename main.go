// * Download .zip from gob.mx
// * Save .zip to server
// * Extract .zip with its content (.csv)
// * Parse .csv rows and save to MySQL DB
// * Create API with router with search functionality
// * Create DigitalOcean Droplet (MySQL Server)
// * Deploy .go microservice to Heroku

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/mileusna/crontab"
)

var db *sqlx.DB

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	db = InitDB()
	defer db.Close()

	router := InitRouter()

	ctab := crontab.New() // create cron table

	ctab.MustAddJob("0 * * * *", func() {
		fmt.Println("Running cron job fetch data")
		go FetchData()
	})

	log.Println("App running at port 5000")
	http.ListenAndServe(":"+port, router)

}

func getStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var stats Stats
	err := db.Get(&stats, `SELECT COUNT(*) as total,
	sum(case when RESULTADO = 1 then 1 end) as confirmed, 
		sum(case when RESULTADO = 1 AND FECHA_DEF NOT LIKE "9999-%%-%%" then 1 end) as dead 
	FROM cases`)

	if err != nil {
		w.Write([]byte(`{"error": "ERROR Querying DB"}`))
		return
	}

	statsJSON, err := json.Marshal(stats)

	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprintf(`{"count": %s}`, statsJSON)))
}

func getData(w http.ResponseWriter, r *http.Request) {
	const defaultCount = 10

	requestedPage, err := strconv.ParseUint(r.URL.Query().Get("page"), 10, 64)
	requestedCount, err := strconv.ParseUint(r.URL.Query().Get("count"), 10, 64)
	if requestedPage <= 0 {
		requestedPage = 1
	}

	if requestedCount <= 0 || requestedCount > 100 {
		requestedCount = defaultCount
	}

	var covidCases []CovidCase
	var totalCases uint64
	w.Header().Set("Content-Type", "application/json")
	offset := (requestedPage - 1) * requestedCount
	sqlQuery := sq.Select("*").From("cases").OrderBy("FECHA_INGRESO DESC").Limit(requestedCount).Offset(offset)

	if r.URL.Query().Get("resultado") != "" {
		sqlQuery = sqlQuery.Where("RESULTADO", r.URL.Query().Get("resultado"))
	}
	sql, _, err := sqlQuery.ToSql()

	rows, err := db.Queryx(sql)

	if err != nil {
		w.Write([]byte(`{"error": "ERROR Querying DB"}`))
	}
	err = db.Get(&totalCases, "SELECT count(*) FROM cases")

	for rows.Next() {
		var covidCase CovidCase
		err = rows.StructScan(&covidCase)

		covidCases = append(covidCases, covidCase)
	}

	rowsJSON, err := json.Marshal(covidCases)

	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprintf(`
	{
		"cases": %s,
		"pagination": {
			"count": %d,
			"total": %d,
			"totalPages": %d,
		}
	}
	`, rowsJSON, requestedCount, totalCases, totalCases/requestedCount)))
}

func forceFetch(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Running fetch data task forcefully")
	go FetchData()
	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprintf("Running manual fetch data task. May take a few moments to update")))
}

func getMain(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("CovidMx API Running"))
}
