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

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
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

	// If cron job is due (present or past due_Date), retrieve data
	// go FetchData()

	log.Println("App running at port 5000")
	http.ListenAndServe(":"+port, router)

}

func getStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var count struct {
		Total     int `db:"total"`
		Confirmed int `db:"confirmed"`
		Dead      int `db:"dead"`
	}
	err := db.Get(&count, `SELECT COUNT(*) as total,
	sum(case when RESULTADO = 1 then 1 end) as confirmed, 
		sum(case when RESULTADO = 1 AND FECHA_DEF NOT LIKE "9999-%%-%%" then 1 end) as dead 
	FROM cases`)

	if err != nil {
		w.Write([]byte(`{"error": "ERROR Querying DB"}`))
		return
	}

	countJSON, err := json.Marshal(count)

	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprintf(`{"count": %s}`, countJSON)))
}

func getData(w http.ResponseWriter, r *http.Request) {
	var covidCases []CovidCase
	var totalCases int
	w.Header().Set("Content-Type", "application/json")
	rows, err := db.Queryx("SELECT * FROM cases WHERE resultado = 1 LIMIT 5")

	if err != nil {
		w.Write([]byte(`{"error": "ERROR Querying DB"}`))
	}
	err = db.Get(&totalCases, "SELECT count(*) FROM cases")

	for rows.Next() {
		var covidCase CovidCase
		err = rows.StructScan(&covidCase)

		covidCases = append(covidCases, covidCase)
		fmt.Println("a", covidCase)
	}

	rowsJSON, err := json.Marshal(covidCases)

	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprintf(`
	{
		"cases": %s,
		"pagination": {
			"total": %d,
		}
	}
	`, rowsJSON, totalCases)))
}
