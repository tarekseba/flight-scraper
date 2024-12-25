package db

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func panicIfEmpty(arg, val string) {
	if val == "" {
		panic(fmt.Sprintf("%s missing or empty in .env file", arg))
	}
}

/*
Panics if connection fails or missing env params
*/
func InitDB() *sqlx.DB {
	user := os.Getenv(DB_USER)
	panicIfEmpty(DB_USER, user)
	password := os.Getenv(DB_PASSWORD)
	panicIfEmpty(DB_PASSWORD, password)
	ssl := os.Getenv(DB_SSL)
	panicIfEmpty(DB_SSL, ssl)
	dbName := os.Getenv(DB_NAME)
	panicIfEmpty(DB_NAME, dbName)
	driver := os.Getenv(DB_DRIVER)
	panicIfEmpty(DB_DRIVER, driver)

	var dsn string = fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s", user, password, dbName, ssl)
	return sqlx.MustConnect(driver, dsn)
}
