package main

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
)

// InsertStatement for creating a new Covid case
const InsertStatement string = `INSERT IGNORE INTO cases(FECHA_ACTUALIZACION,
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

// UpdateStatement for creating a new Covid case
const UpdateStatement string = `UPDATE IGNORE cases(FECHA_ACTUALIZACION,
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

// InitDB uses OS env vars to connect to Mysql DB
func InitDB() *sqlx.DB {
	mysqlUser := os.Getenv("DB_USER")
	mysqlPwd := os.Getenv("DB_PWD")
	mysqlHost := os.Getenv("DB_HOST")
	mysqlDB := os.Getenv("DB_NAME")

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

	return db
}
