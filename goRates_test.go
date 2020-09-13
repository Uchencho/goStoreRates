package main

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/Uchencho/goStoreRates/config"
	_ "github.com/mattn/go-sqlite3"
)

func createTestDb() *sql.DB {

	_, err := os.Create("data.db")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	db, err := sql.Open("sqlite3", "data.db")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	return db
}

func TestCreateUserTable(t *testing.T) {
	dB := createTestDb()

	defer func() {
		err := os.Remove("data.db")
		if err != nil {
			log.Println("Error occured in removing test db, ", err)
		}
	}()

	defer dB.Close()

	config.CreateUsersTable(dB)
}

func TestCreateRatingTable(t *testing.T) {
	dB := createTestDb()

	defer func() {
		err := os.Remove("data.db")
		if err != nil {
			log.Println("Error occured in removing test db, ", err)
		}
	}()

	defer dB.Close()

	config.CreateRatingTable(dB)
}

// t.Fatalf("This test fails and STOPS running")
// t.Errorf("This test fails but continues running")
