package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Uchencho/goStoreRates/api/users"
)

// type SQLDB interface {
// 	Exec(query string, args ...interface{}) (sql.Result, error)
// }

// type MockDB struct{}

// func (mdb *MockDB) Exec(query string, args ...interface{}) (sql.Result, error) {
// 	return nil, nil
// }

func TestRegister(t *testing.T) {

	l := users.Account{
		CompanyName: "Amazon",
		Email:       "ama@gmail.com",
		Password:    "SomeLongString",
	}
	reqBody, err := json.Marshal(l)
	if err != nil {
		t.Errorf("Could not marshal json with error, %v", err)
	}
	req, err := http.NewRequest("POST", "localhost:8000/registers", bytes.NewBuffer(reqBody))

	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rec := httptest.NewRecorder()
	users.RegisterUser(rec, req)

	resp := rec.Result()
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status created; got %v", resp.StatusCode)
	}
	defer resp.Body.Close()
}

func TestAddRate(t *testing.T) {
	// req, err := http.NewRequest("POST", "localhost:8000/create")
}

// t.Fatalf("This test fails and STOPS running")
// t.Errorf("This test fails but continues running")
