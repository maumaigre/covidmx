package main

// CovidCase is a structure of Covid-19 cases
type CovidCase struct {
	FechaActualizacion string `db:"FECHA_ACTUALIZACION"`
	IDRegistro         string `db:"ID_REGISTRO"`
	Origen             int    `db:"ORIGEN"`
	Sector             int    `db:"SECTOR"`
	EntidadUm          int    `db:"ENTIDAD_UM"`
	Sexo               int    `db:"SEXO"`
	EntidadNac         int    `db:"ENTIDAD_NAC"`
	EntidadRes         int    `db:"ENTIDAD_RES"`
	MunicipioRes       int    `db:"MUNICIPIO_RES"`
	TipoPaciente       int    `db:"TIPO_PACIENTE"`
	FechaIngreso       string `db:"FECHA_INGRESO"`
	FechaSintomas      string `db:"FECHA_SINTOMAS"`
	FechaDefuncion     string `db:"FECHA_DEF"`
	Intubado           int    `db:"INTUBADO"`
	Neumonia           int    `db:"NEUMONIA"`
	Edad               int    `db:"EDAD"`
	Nacionalidad       int    `db:"NACIONALIDAD"`
	Embarazo           int    `db:"EMBARAZO"`
	HablaLenguaIndig   int    `db:"HABLA_LENGUA_INDIG"`
	Diabetes           int    `db:"DIABETES"`
	Epoc               int    `db:"EPOC"`
	Asma               int    `db:"ASMA"`
	Inmusupr           int    `db:"INMUSUPR"`
	Hipertension       int    `db:"HIPERTENSION"`
	OtraCom            int    `db:"OTRA_COM"`
	Cardiovascular     int    `db:"CARDIOVASCULAR"`
	Obesidad           int    `db:"OBESIDAD"`
	RenalCronica       int    `db:"RENAL_CRONICA"`
	Tabaquismo         int    `db:"TABAQUISMO"`
	OtroCaso           int    `db:"OTRO_CASO"`
	Resultado          int    `db:"RESULTADO"`
	Migrante           int    `db:"MIGRANTE"`
	PaisNacionalidad   string `db:"PAIS_NACIONALIDAD"`
	PaisOrigen         string `db:"PAIS_ORIGEN"`
	Uci                int    `db:"UCI"`
}

type Stats struct {
	Tested    int `db:"total"`
	Confirmed int `db:"confirmed"`
	Dead      int `db:"dead"`
}

type StateStat struct {
	EntidadRes  int `db:"ENTIDAD_RES"`
	Confirmados int `db:"CONFIRMADOS"`
	Fallecidos  int `db:"FALLECIDOS"`
}

type DailyNewStat struct {
	ID           int    `db:"id"`
	FechaIngreso string `db:"fecha_ingreso"`
	NewConfirmed int    `db:"nuevos_confirmados"`
	NewDead      int    `db:"nuevos_fallecidos"`
	NewTested    int    `db:"nuevos_pruebas"`
	Total        int    `db:"total_pruebas"`
	Confirmed    int    `db:"total_confirmados"`
	Dead         int    `db:"total_fallecidos"`
}
