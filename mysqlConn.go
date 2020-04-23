package main

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
)

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
