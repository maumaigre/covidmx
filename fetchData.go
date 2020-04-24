package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/artdarek/go-unzip"
	"github.com/sergi/go-diff/diffmatchpatch"
)

// FetchData downloads, unzips and renames file to data.csv
func FetchData() {

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

		fmt.Println("TEST")
		content, err := ioutil.ReadFile("./data_new/data.csv")
		if err != nil {
			log.Fatal(err)
		}

		file1 := string(content)

		content, err = ioutil.ReadFile("./data/data.csv")
		if err != nil {
			log.Fatal(err)
		}

		file2 := string(content)

		fmt.Println(len(file1), len(file2))

		dmp := diffmatchpatch.New()

		diffs := dmp.DiffMain(file1, file2, false)

		arr := DiffCSV(diffs)

		fmt.Println("print", arr)
		_ = ioutil.WriteFile("./data_new/diff.csv", []byte(arr), 0644)

		os.Remove("data/data.csv")

		os.Rename("data_new/data.csv", "data/data.csv")

		writeCSVToDB("./data_new/diff.csv", true)

		os.RemoveAll("./data_new/")

	} else {
		os.MkdirAll("./data", 0755)
		os.Rename("./data_new/data.csv", "./data/data.csv")
		os.RemoveAll("./data_new/")
		writeCSVToDB("./data/data.csv", false)
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

func writeCSVToDB(inputCsvFile string, diff bool) {
	recordFile, err := os.Open(inputCsvFile)
	if err != nil {
		fmt.Println("An error encountered ::", err)
	}
	// 2. Initialize the reader
	reader := csv.NewReader(recordFile)
	// 3. Read all the records
	records, _ := reader.ReadAll()

	stmt := `INSERT IGNORE INTO cases(FECHA_ACTUALIZACION,
		ID_REGISTRO,
		ORIGEN,
		SECTOR,
		ENTIDAD_UM,
		SEXO,
		ENTIDAD_NAC,
		ENTIDAD_RES,
		MUNICIPIO_RES,
		TIPO_PACIENTE,
		FECHA_INGRESO,
		FECHA_SINTOMAS,
		FECHA_DEF,
		INTUBADO,
		NEUMONIA,
		EDAD,
		NACIONALIDAD,
		EMBARAZO,
		HABLA_LENGUA_INDIG,
		DIABETES,
		EPOC,
		ASMA,
		INMUSUPR,
		HIPERTENSION,
		OTRA_COM,
		CARDIOVASCULAR,
		OBESIDAD,
		RENAL_CRONICA,
		TABAQUISMO,
		OTRO_CASO,
		RESULTADO,
		MIGRANTE,
		PAIS_NACIONALIDAD,
		PAIS_ORIGEN,
		UCI
	) VALUES `
	if err != nil {
		log.Fatal(err)
	}

	var slicedRecords [][]string
	if diff {
		slicedRecords = records
	} else {
		slicedRecords = records[1:]
	}

	for _, row := range slicedRecords {
		statement := fmt.Sprintf(`%s ("%s");`, stmt, strings.Join(row, `", "`))
		fmt.Println(statement)
		_, err := db.Exec(statement)
		if err != nil {
			log.Fatal(err)
		}
	}
}
