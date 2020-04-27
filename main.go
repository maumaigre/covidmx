package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
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

	sqlQuery := sq.Select(`COUNT(*) as total,
	sum(case when RESULTADO = 1 then 1 end) as confirmed, 
		sum(case when RESULTADO = 1 AND FECHA_DEF NOT LIKE "9999-%%-%%" then 1 end) as dead`).From("cases")

	if r.URL.Query().Get("entidad_res") != "" {
		sqlQuery = sqlQuery.Where(fmt.Sprintf("ENTIDAD_RES = %s", r.URL.Query().Get("entidad_res")))
	}

	sql, _, err := sqlQuery.ToSql()

	db.Get(&stats, sql)

	if err != nil {
		w.Write([]byte(`{"error": "ERROR Querying DB"}`))
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
	var totalCases uint32
	w.Header().Set("Content-Type", "application/json")
	offset := (requestedPage - 1) * requestedCount
	sqlQuery := sq.Select("*").From("cases").OrderBy("FECHA_INGRESO DESC")
	countSQLQuery := sq.Select("COUNT(*)").From("cases")

	if r.URL.Query().Get("resultado") != "" {
		sqlQuery = sqlQuery.Where(fmt.Sprintf("RESULTADO = %s", r.URL.Query().Get("resultado")))
		countSQLQuery = countSQLQuery.Where(fmt.Sprintf("RESULTADO = %s", r.URL.Query().Get("resultado")))
	}

	if r.URL.Query().Get("entidad_res") != "" {
		sqlQuery = sqlQuery.Where(fmt.Sprintf("ENTIDAD_RES = %s", r.URL.Query().Get("entidad_res")))
		countSQLQuery = countSQLQuery.Where(fmt.Sprintf("ENTIDAD_RES = %s", r.URL.Query().Get("entidad_res")))
	}

	paginatedSQLQuery := sqlQuery.Limit(requestedCount).Offset(offset)

	sql, _, err := paginatedSQLQuery.ToSql()

	countSQL, _, err := countSQLQuery.ToSql()

	rows, err := db.Queryx(sql)

	db.Get(&totalCases, countSQL)

	if err != nil {
		w.Write([]byte(`{"error": "ERROR Querying DB"}`))
	}

	for rows.Next() {
		var covidCase CovidCase
		rows.StructScan(&covidCase)

		covidCases = append(covidCases, covidCase)
	}

	rowsJSON, err := json.Marshal(covidCases)

	w.WriteHeader(200)
	var totalPages int
	totalPages = int(math.Ceil(float64(totalCases) / float64(requestedCount)))
	w.Write([]byte(fmt.Sprintf(`
	{
		"cases": %s,
		"pagination": {
			"count": %d,
			"total": %d,
			"totalPages": %d,
		}
	}
	`, rowsJSON, requestedCount, totalCases, totalPages)))
}

func getDailyNewStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	sqlQuery := sq.Select("*").From("daily_new_stats").OrderBy("FECHA_INGRESO asc")

	sql, _, err := sqlQuery.ToSql()

	if err != nil {
		fmt.Println("Error querying daily new stats")
	}
	rows, err := db.Queryx(sql)

	var dailyNewStats []DailyNewStat
	for rows.Next() {
		var dailyNewStat DailyNewStat
		rows.StructScan(&dailyNewStat)
		dailyNewStats = append(dailyNewStats, dailyNewStat)
	}

	responseJSON, err := json.Marshal(dailyNewStats)
	w.WriteHeader(200)
	w.Write([]byte(responseJSON))

}

func getStateStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	sqlQuery := sq.Select("ENTIDAD_RES, COUNT(*) as CONFIRMADOS, sum(case when FECHA_DEF NOT LIKE '9999-%%-%%' then 1 end) as FALLECIDOS").From("cases").Where("RESULTADO = 1").GroupBy("ENTIDAD_RES")

	requestedOrder := r.URL.Query().Get("order")

	if requestedOrder == "asc" || requestedOrder == "desc" {
		sqlQuery = sqlQuery.OrderBy(fmt.Sprintf("count(*) %s", requestedOrder))
	}

	sql, _, err := sqlQuery.ToSql()

	if err != nil {
		fmt.Println("Error getting query")
	}

	rows, err := db.Queryx(sql)

	var cases []StateStat
	for rows.Next() {
		var entidadCaso StateStat
		rows.StructScan(&entidadCaso)
		cases = append(cases, entidadCaso)
	}

	responseJSON, err := json.Marshal(cases)
	w.WriteHeader(200)
	w.Write([]byte(responseJSON))
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
