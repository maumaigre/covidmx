package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/artdarek/go-unzip"
)

// FetchData downloads, unzips and renames file to data.csv
func FetchData() {

	err := downloadFile("./data.zip", "http://187.191.75.115/gobmx/salud/datos_abiertos/datos_abiertos_covid19.zip")
	if err != nil {
		fmt.Println("ERROR Downloading file", err)
	}
	err = unzipFile("./data.zip", "./data_new/")

	if err != nil {
		fmt.Println("ERROR Unzipping file", err)
	}

	files, err := ioutil.ReadDir("./data_new")

	oldPath := fmt.Sprintf("./data_new/%s", files[0].Name())
	os.Rename(oldPath, "./data_new/data.csv")

	oldFilesDir, err := ioutil.ReadDir("./data")

	fmt.Println(oldFilesDir)

	if len(oldFilesDir) > 0 {
		equal := CompareMD5("./data/data.csv", "./data_new/data.csv")
		if !equal {
			fmt.Println("New CSV File detected, inserting to DB")
			writeCSVToDB("./data_new/data.csv")
		} else {
			fmt.Println("No changes detected to new CSV file. Skipping.")
		}
	} else {
		err = os.Mkdir("./data", 0777)
		fmt.Println("No previous CSV file detected, generating query to DB")
		writeCSVToDB("./data_new/data.csv")
	}

	fmt.Println("Finished fetchData process.")

	os.Rename("./data_new/data.csv", "./data/data.csv")

	os.Remove("./data.sql")
	os.Remove("./data.zip")
	os.RemoveAll("./data_new")

}

func downloadFile(filepath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func unzipFile(inputFile string, outputDirectory string) error {
	uz := unzip.New(inputFile, outputDirectory)
	err := uz.Extract()
	if err != nil {
		fmt.Println(err)
	}
	return err
}

func writeCSVToDB(inputCsvFile string) {
	db.Exec(`DELETE FROM cases`)

	recordFile, err := os.Open(inputCsvFile)
	if err != nil {
		fmt.Println("An error encountered ::", err)
	}
	reader := csv.NewReader(recordFile)
	records, _ := reader.ReadAll()

	newFile, _ := os.Create("data.sql")

	for _, row := range records[1:] {
		statement := fmt.Sprintf(`%s ("%s");`, InsertStatement, strings.Join(row, `", "`))
		newFile.WriteString(statement)
	}
	path, err := os.Getwd()
	cmd := exec.Command("/bin/sh", "./sql_import_data.sh", filepath.Join(path, newFile.Name()))
	err = cmd.Run()

	if err != nil {
		log.Fatalf("Error executing query.")
	}

	updateDailyNewStat()
}

func updateDailyNewStat() {
	var previousDailyStat DailyNewStat
	var currentStat Stats
	previousDailyStatSQLQuery := sq.Select("*").From("daily_new_stats").OrderBy("fecha_ingreso desc").Limit(1)

	previousDailyStatSQL, _, err := previousDailyStatSQLQuery.ToSql()

	if err != nil {
		fmt.Println("Error querying daily new stat")
	}

	sqlQuery := sq.Select(`COUNT(*) as total,
	sum(case when RESULTADO = 1 then 1 end) as confirmed, 
		sum(case when RESULTADO = 1 AND FECHA_DEF NOT LIKE "9999-%%-%%" then 1 end) as dead`).From("cases")

	sql, _, err := sqlQuery.ToSql()

	db.Get(&currentStat, sql)

	db.Get(&previousDailyStat, previousDailyStatSQL)

	dt := time.Now()
	dateFormatted := dt.Format("2006-01-02")

	newConfirmed := currentStat.Confirmed - previousDailyStat.Confirmed
	newDead := currentStat.Dead - previousDailyStat.NewDead
	newTotal := currentStat.Tested - previousDailyStat.Total

	insertNewStatSQLQuery := sq.Insert("daily_new_stats").Columns("id", "fecha_ingreso", "nuevos_confirmados",
		"nuevos_fallecidos", "nuevos_pruebas", "total_pruebas", "total_confirmados", "total_fallecidos").Values(
		nil, dateFormatted, newConfirmed, newDead, newTotal, currentStat.Tested, currentStat.Confirmed, currentStat.Dead)

	sql, _, err = insertNewStatSQLQuery.ToSql()

	_, err = db.Exec(sql)

	if err != nil {
		fmt.Println("Error adding new daily stat")
	}
}
