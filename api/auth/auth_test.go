package auth

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
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

func TestCheckAuth(t *testing.T) {
	req, err := http.NewRequest("GET", "localhost:8000", nil)
	if err != nil {
		t.Fatalf("Could not create request with error %s", err)
	}
	if authorized, _, _ := CheckAuth(req); authorized {
		t.Fatalf("Authorizing when no authentication was passed")
	}

	req.Header.Add("Authorization", "Bearer Token")
	if authorized, _, _ := CheckAuth(req); authorized {
		t.Fatalf("Authorizing when no user has been created")
	}
}

func TestCheckAuthNoToken(t *testing.T) {
	req, err := http.NewRequest("GET", "localhost:8000", nil)
	if err != nil {
		t.Fatalf("Could not create request with error %s", err)
	}

	req.Header.Add("Authorization", "Bearer")
	if authorized, _, _ := CheckAuth(req); authorized {
		t.Fatalf("Authorizing when no user has been created")
	}
}

func TestUnauthorizedResponse(t *testing.T) {

	rec := httptest.NewRecorder()
	err := errors.New("Test error created")
	UnAuthorizedResponse(rec, err)

	resp := rec.Result()
	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("Unauthorized response was not returned")
	}
}
