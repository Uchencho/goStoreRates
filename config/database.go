package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

const (
	POSTGRES_USER     = "golang"
	POSTGRES_PASSWORD = "googleGo"
	DB_NAME           = "rateService"
	POSTGRES_HOST     = "localhost"
	POSTGRES_PORT     = 5432
)

func databaseUrl() string {
	dbUrl, present := os.LookupEnv("DATABASE_URL")
	if present {
		return dbUrl
	}

	local_postgress_conn := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable",
		POSTGRES_HOST, POSTGRES_PORT, POSTGRES_USER, POSTGRES_PASSWORD,
		DB_NAME)
	return local_postgress_conn
}

func ConnectDatabase() *sql.DB {

	db, err := sql.Open("postgres", databaseUrl())
	if err != nil {
		log.Println(err)
		panic("Failed to connect to database")
	}

	dbErr := db.Ping()
	if dbErr != nil {
		log.Println("Error occured in pinging database")
		panic(dbErr)
	}
	fmt.Println("\nConnected successfully")
	return db
}

var Db = ConnectDatabase()
