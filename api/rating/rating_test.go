package rating

import (
	"database/sql"
	"log"
	"net/http"
	"net/http/httptest"
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

func TestAddRateToDB(t *testing.T) {
	dB := createTestDb()

	defer func() {
		err := os.Remove("data.db")
		if err != nil {
			log.Println("Error occured in removing test db, ", err)
		}
	}()

	defer dB.Close()
	config.CreateRatingTable(dB)

	if added, _ := AddRateToDB(dB, rateJson{}); added {
		t.Fatalf("Adding empty rates to database")
	}
}

func TestDeleteRate(t *testing.T) {
	dB := createTestDb()

	defer func() {
		err := os.Remove("data.db")
		if err != nil {
			log.Println("Error occured in removing test db, ", err)
		}
	}()

	defer dB.Close()
	config.CreateRatingTable(dB)

	if deleted := deleteRate(dB, ""); !deleted {
		t.Fatalf("Deleting rate that does not exist doesn't throw error")
	}
}

func TestGetSpecificRate(t *testing.T) {
	dB := createTestDb()

	defer func() {
		err := os.Remove("data.db")
		if err != nil {
			log.Println("Error occured in removing test db, ", err)
		}
	}()

	defer dB.Close()
	config.CreateRatingTable(dB)

	if found, _ := getSpecificRate(dB, "Google", "2"); found {
		t.Fatalf("Finding what has not been created")
	}

	rj := rateJson{
		BusinessName: "Google",
		UserID:       "1",
		ProductID:    "1",
		ProductName:  "Mac Book Pro",
		Rating:       4,
	}

	if added, _ := AddRateToDB(dB, rj); added {
		t.Fatalf("Adding rates without user foreignkey constraint to database")
	}
}

func TestGetFromRedis(t *testing.T) {
	retrieved, avgRate := getFromRedis("Google", "1")
	if !retrieved {
		t.Fatalf("Could not retrieve avg rating stored above")
	}
	if avgRate == 0 {
		t.Fatalf("Stored the wrong value in redis")
	}

}

func TestGetAllRates(t *testing.T) {
	dB := createTestDb()

	defer func() {
		err := os.Remove("data.db")
		if err != nil {
			log.Println("Error occured in removing test db, ", err)
		}
	}()

	defer dB.Close()
	config.CreateRatingTable(dB)

	if found, _ := getAllRates(dB, "Google"); found {
		t.Fatalf("finding what does not exist")
	}
}

func TestUpdateRedis(t *testing.T) {

	dB := createTestDb()

	defer func() {
		err := os.Remove("data.db")
		if err != nil {
			log.Println("Error occured in removing test db, ", err)
		}
	}()

	defer dB.Close()
	config.CreateRatingTable(dB)

	if updated := updateRedis(dB, rateJson{}, "1"); updated {
		t.Fatalf("Updating what was not stored on the db")
	}
}

func TestAddRate(t *testing.T) {
	req, err := http.NewRequest("POST", "", nil)
	if err != nil {
		t.Fatalf("Could not create request for addrate, %s", err)
	}
	rec := httptest.NewRecorder()
	AddRate(rec, req)

	resp := rec.Result()
	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("Authentication was not triggered poperly")
	}
}

func TestProductRatingDetail(t *testing.T) {
	req, err := http.NewRequest("POST", "", nil)
	if err != nil {
		t.Fatalf("Could not create request for addrate, %s", err)
	}
	rec := httptest.NewRecorder()
	ProductRatingDetail(rec, req)

	resp := rec.Result()
	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("Authentication was not triggered poperly")
	}
}
