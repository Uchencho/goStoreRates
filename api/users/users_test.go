package users

import (
	"database/sql"
	"log"
	"os"
	"testing"
	"time"

	"github.com/Uchencho/goStoreRates/config"
	_ "github.com/mattn/go-sqlite3"
)

func createTestDb() *sql.DB {

	err := os.MkdirAll("../TestDB", 0755)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	_, err = os.Create("../TestDB/data.db")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	db, err := sql.Open("sqlite3", "../TestDB/data.db")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	return db
}

func TestAddUser(t *testing.T) {

	dB := createTestDb()
	config.CreateUsersTable(dB)
	defer func() {
		err := os.RemoveAll("../TestDB/")
		if err != nil {
			log.Println("Error occured in removing test db, ", err)
		}
	}()

	defer dB.Close()

	testUser := Account{
		Email:       "alozyuche@gmail.com",
		Password:    "strongPassword",
		CompanyName: "Google",
		Token:       "random string",
		CreatedOn:   time.Now(),
		LastLogin:   time.Now(),
	}

	if created := addUser(dB, testUser); !created {
		t.Fatalf("Could not add user to database")
	}

	if added := addUser(dB, testUser); added {
		t.Fatalf("Added duplicated data to the db")
	}

}
