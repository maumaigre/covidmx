// * Descargar .zip de gob.mx
// * Guardar .zip en servidor
// * Extraer .zip con sus contenidos (.csv)
// * Parsear .csv y guardar a DB MySQL -----
// * Crear REST API con rutas para busqueda y filtrado de datos

package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/artdarek/go-unzip"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var covidCase struct {
	fechaActualizacion string `mysql:"FECHA_ACTUALIZACION"`
	idRegistro         string `mysql:"ID_REGISTRO"`
	origen             int    `mysql:"ORIGEN"`
	sector             int    `mysql:"SECTOR"`
	entidadUm          int    `mysql:"ENTIDAD_UM"`
	sexo               int    `mysql:"SEXO"`
	entidadNac         int    `mysql:"ENTIDAD_NAC"`
	entidadRes         int    `mysql:"ENTIDAD_RES"`
	municipioRes       int    `mysql:"MUNICIPIO_RES"`
	tipoPaciente       int    `mysql:"TIPO_PACIENTE"`
	fechaIngreso       string `mysql:"FECHA_INGRESO"`
	fechaSintomas      string `mysql:"FECHA_SINTOMAS"`
	fechaDef           string `mysql:"FECHA_DEF"`
	intubado           int    `mysql:"INTUBADO"`
	neumonia           int    `mysql:"NEUMONIA"`
	edad               int    `mysql:"EDAD"`
	nacionalidad       int    `mysql:"NACIONALIDAD"`
	embarazo           int    `mysql:"EMBARAZO"`
	hablaLenguaIndig   int    `mysql:"HABLA_LENGUA_INDIG"`
	diabetes           int    `mysql:"DIABETES"`
	epoc               int    `mysql:"EPOC"`
	asma               int    `mysql:"ASMA"`
	inmusupr           int    `mysql:"INMUSUPR"`
	hipertension       int    `mysql:"HIPERTENSION"`
	otraCom            int    `mysql:"OTRA_COM"`
	cardiovascular     int    `mysql:"CARDIOVASCULAR"`
	obesidad           int    `mysql:"OBESIDAD"`
	renalCronica       int    `mysql:"RENAL_CRONICA"`
	tabaquismo         int    `mysql:"TABAQUISMO"`
	otroCaso           int    `mysql:"OTRO_CASO"`
	resultado          int    `mysql:"RESULTADO"`
	migrante           int    `mysql:"MIGRANTE"`
	paisNacionalidad   int    `mysql:"PAIS_NACIONALIDAD"`
	paisOrigen         int    `mysql:"PAIS_ORIGEN"`
	uci                int    `mysql:"UCI"`
}

func main() {
	db, err := sql.Open("mysql", "root:[]@/covid")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// If cron job is due (present or past due_Date), retrieve data
	// retrieveData(db)

	router := mux.NewRouter()
	router.HandleFunc("/", getData).Methods("GET")

	log.Println("App running at port 5000")
	http.ListenAndServe(":5000", router)
}

// DownloadFile will download a url to a local file
func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

// UnzipFile unzips the input file to the output directory
func UnzipFile(inputFile string, outputDirectory string) error {
	uz := unzip.New(inputFile, outputDirectory)
	err := uz.Extract()
	if err != nil {
		fmt.Println(err)
	}
	return err
}

func writeCSVToDB(inputCsvFile string, db *sql.DB) {
	recordFile, err := os.Open(inputCsvFile)
	if err != nil {
		fmt.Println("An error encountered ::", err)
	}
	// 2. Initialize the reader
	reader := csv.NewReader(recordFile)
	// 3. Read all the records
	records, _ := reader.ReadAll()

	stmt := `INSERT INTO cases(FECHA_ACTUALIZACION,
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

	for _, row := range records[1:] {
		statement := fmt.Sprintf(`%s ("%s");`, stmt, strings.Join(row, `", "`))
		fmt.Println(statement)
		_, err := db.Exec(statement)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func getData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)
	w.Write([]byte("{message: 'app running'}"))
}

func retrieveData(db *sql.DB) {
	err := DownloadFile("data.zip", "http://187.191.75.115/gobmx/salud/datos_abiertos/datos_abiertos_covid19.zip")
	if err != nil {
		fmt.Println("ERROR Downloading file", err)
	}
	err = UnzipFile("data.zip", "data")

	if err != nil {
		fmt.Println("ERROR Unzipping file", err)
	}

	files, err := ioutil.ReadDir("data")

	oldPath := fmt.Sprintf("data/%s", files[0].Name())
	os.Rename(oldPath, "data/data.csv")
	writeCSVToDB("./data/data.csv", db)
}
