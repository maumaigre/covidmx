package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"

	"github.com/artdarek/go-unzip"
)

// FetchData downloads, unzips and renames file to data.csv
func FetchData() {
	var updatedValues [][]string
	var newValues [][]string

	err := downloadFile("data.zip", "http://187.191.75.115/gobmx/salud/datos_abiertos/datos_abiertos_covid19.zip")
	if err != nil {
		fmt.Println("ERROR Downloading file", err)
	}
	err = unzipFile("data.zip", "data_new/")

	if err != nil {
		fmt.Println("ERROR Unzipping file", err)
	}

	files, err := ioutil.ReadDir("data_new")

	oldPath := fmt.Sprintf("data_new/%s", files[0].Name())
	os.Rename(oldPath, "data_new/data.csv")

	files, err = ioutil.ReadDir("data")

	if len(files) >= 1 {

		recordFile, _ := os.Open("./data/data.csv")
		recordFile2, _ := os.Open("./data_new/data.csv")

		reader := csv.NewReader(recordFile)
		records, _ := reader.ReadAll()

		reader2 := csv.NewReader(recordFile2)

		records2, _ := reader2.ReadAll()

		slicedRecords1 := records[1:]
		slicedRecords2 := records2[1:]

		for i, row := range slicedRecords2 {
			if i <= len(slicedRecords1)-1 {
				if !reflect.DeepEqual(slicedRecords2[i][1:], slicedRecords1[i][1:]) {
					updatedValues = append(updatedValues, row)
				}
			} else {
				newValues = append(newValues, row)
			}
		}

		os.Remove("data/data.csv")

		os.Rename("data_new/data.csv", "data/data.csv")
		os.RemoveAll("./data_new/")

		if len(newValues) > 0 || len(updatedValues) > 0 {
			patchValuesToDB(newValues, updatedValues)
		} else {
			fmt.Println("No changes found in new CSV file. Skipping patch")
		}
	} else {
		fmt.Println("No old CSV File found, generating from scratch")
		os.MkdirAll("./data", 0755)
		os.Rename("./data_new/data.csv", "./data/data.csv")
		os.RemoveAll("./data_new/")
		writeCSVToDB("./data/data.csv")
	}
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
	// 2. Initialize the reader
	reader := csv.NewReader(recordFile)
	// 3. Read all the records
	records, _ := reader.ReadAll()

	for _, row := range records[1:] {
		statement := fmt.Sprintf(`%s ("%s");`, InsertStatement, strings.Join(row, `", "`))
		_, err := db.Exec(statement)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func patchValuesToDB(updatedValues [][]string, newValues [][]string) {
	fmt.Println(fmt.Sprintf("Patching %d new values and %d updated values", len(newValues), len(updatedValues)))
	for _, row := range newValues {
		statement := fmt.Sprintf(`%s ("%s");`, InsertStatement, strings.Join(row, `", "`))
		// fmt.Println(statement)
		_, err := db.Exec(statement)
		if err != nil {
			log.Fatal(err)
		}
	}

	for _, row := range updatedValues {
		statement := fmt.Sprintf(`UPDATE IGNORE cases SET RESULTADO = '%s', INTUBADO = '%s', NEUMONIA = '%s', FECHA_DEF = '%s', UCI = '%s'  WHERE ID_REGISTRO= '%s'`, row[30], row[13], row[14], row[12], row[34], row[1])
		// fmt.Println(statement)
		_, err := db.Exec(statement)
		if err != nil {
			log.Fatal(err)
		}
	}

}
