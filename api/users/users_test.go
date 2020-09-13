package users

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

type SQLDB interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}

type MockDB struct{}

func (mdb *MockDB) Exec(query string, args ...interface{}) (sql.Result, error) {

	return nil, nil
}

func createTestDb() *sql.DB {

	err := os.MkdirAll(".../TestDB", 0755)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	_, err = os.Create(".../TestDB/data.db")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	db, err := sql.Open("sqlite3", ".../TestDB/data.db")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	return db
}

func addUserT(db SQLDB, user Account) bool {

	query := `INSERT INTO users (
		company_name, email, password, token, created_on, last_login
	) VALUES (
		$1, $2, $3, $4, $5, $6
	) RETURNING id;`

	_, err := db.Exec(query, user.CompanyName, user.Email, user.Password,
		user.Token, user.CreatedOn, user.LastLogin)
	if err != nil {
		log.Println("Error adding user to the User's table, ", err)
		return false
	}
	return true
}

func TestAddUser(t *testing.T) {

	createTestDb()

	mockDB := new(MockDB)

	if added := addUserT(mockDB, Account{}); !added {
		t.Fatalf("Expected true in creating user with test db")
	}

}
