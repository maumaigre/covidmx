// * Descargar .zip de gob.mx
// * Guardar .zip en servidor
// * Extraer .zip con sus contenidos (.csv)
// * Parsear .csv y guardar a DB MySQL -----
// * Crear REST API con rutas para busqueda y filtrado de datos

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
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

type CovidCase struct {
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
	port := os.Getenv("PORT")

	mysqlUser := os.Getenv("DB_USER")
	mysqlPwd := os.Getenv("DB_PWD")
	mysqlHost := os.Getenv("DB_HOST")
	// mysqlPort := os.Getenv("DB_PORT")
	mysqlDB := os.Getenv("DB_NAME")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	conn := fmt.Sprintf(
		"%s:%s@(%s)/%s?parseTime=true",
		mysqlUser,
		mysqlPwd,
		mysqlHost,
		mysqlDB,
	)

	var err error
	db, err = sqlx.Connect("mysql", conn)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// If cron job is due (present or past due_Date), retrieve data
	// go retrieveData()

	router := mux.NewRouter()
	router.HandleFunc("/", getData).Methods("GET")

	log.Println("App running at port 5000")
	http.ListenAndServe(":"+port, router)
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

func writeCSVToDB(inputCsvFile string) {
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
	var covidCases []CovidCase
	w.Header().Set("Content-Type", "application/json")
	rows, err := db.Queryx("SELECT * FROM cases WHERE resultado = 1;")

	if err != nil {
		w.Write([]byte(`{"error": "ERROR Querying DB"}`))
	}

	for rows.Next() {
		var covidCase CovidCase
		err = rows.StructScan(&covidCase)

		covidCases = append(covidCases, covidCase)
	}

	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprintf(`{"confirmed": "%d"}`, len(covidCases))))
}

func retrieveData() {
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
	writeCSVToDB("./data/data.csv")
}
