package auth

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/Uchencho/goStoreRates/config"
	_ "github.com/mattn/go-sqlite3"
)

func TestHashPassword(t *testing.T) {
	_, err := HashPassword("PasswordToHash")
	if err != nil {
		t.Fatalf("Could not hash the password %s", err)
	}
}

func TestGenerateToken(t *testing.T) {
	_, err := GenerateToken("Uchencho")
	if err != nil {
		t.Fatalf("Could not generate token %s", err)
	}
}

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

func TestGetUser(t *testing.T) {

	dB := createTestDb()
	config.CreateUsersTable(dB)
	defer func() {
		err := os.Remove("data.db")
		if err != nil {
			log.Println("Error occured in removing test db, ", err)
		}
	}()

	defer dB.Close()

	retrieved, _ := getUser(dB, "token")
	if retrieved {
		t.Fatalf("Found a user that has not been created")
	}
}
