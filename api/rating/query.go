package rating

import (
	"database/sql"
	"encoding/json"
	"log"

	"github.com/go-redis/redis"
)

type rateJson struct {
	BusinessName  string `json:"company_name"`
	UserID        string `json:"user_id"`
	ProductID     string `json:"product_id"`
	ProductName   string `json:"product_name"`
	Rating        int    `json:"rating"`
	AverageRating int    `json:"average_rating"`
}

type averageJson struct {
	ProductID     string `json:"product_id"`
	ProductName   string `json:"product_name"`
	AverageRating int    `json:"average_rating"`
}

type avgJson struct {
	ProductID     string  `json:"product_id"`
	AverageRating float32 `json:"average_rating"`
}

func sendtoRedis(businessName, productID string, averageRating int) bool {

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       1,
	})

	value := map[string]int{
		productID: averageRating,
	}
	jsonString, err := json.Marshal(value)
	if err != nil {
		log.Println("Could not marshal the json ", err)
		return false
	}
	err = rdb.Set(businessName, string(jsonString), 0).Err()
	if err != nil {
		log.Println("Could not save to redis db, ", err)
		return false
	}
	return true
}

func getCurrentAverage(dB *sql.DB, businessName, productID string) (ratings []int) {
	query := `SELECT rating FROM rates WHERE business_name = $1 and product_id = $2;`

	rows, err := dB.Query(query, businessName, productID)
	if err != nil {
		log.Println(err)
		return []int{}
	}
	defer rows.Close()

	var rat int

	for rows.Next() {
		err := rows.Scan(&rat)
		if err != nil {
			log.Println(err)
			return []int{}
		}
		ratings = append(ratings, rat)
	}

	err = rows.Err()
	if err != nil {
		log.Println(err)
		return []int{}
	}
	return
}

func AddRateToDB(dB *sql.DB, r rateJson) (bool, rateJson) {

	query := `INSERT INTO rates (
		business_name, user_id, product_id,
		product_name, rating, average_rating
		) VALUES (
			$1, $2, $3, $4, $5, $6
		) RETURNING id;`

	ratings := getCurrentAverage(dB, r.BusinessName, r.ProductID)
	ratings = append(ratings, r.Rating)

	var theSum int
	for i := 0; i < len(ratings); i++ {
		theSum += ratings[i]
	}
	r.AverageRating = theSum / len(ratings)

	sentToRedis := sendtoRedis(r.BusinessName, r.ProductID, r.AverageRating)
	if !sentToRedis {
		log.Println("Could not send value to redis")
	}

	_, err := dB.Exec(query, r.BusinessName, r.UserID, r.ProductID,
		r.ProductName, r.Rating, r.AverageRating)
	if err != nil {
		log.Println("Error in saving details to db, ", err)
		return false, rateJson{}
	}
	return true, r
}

func currentRating(dB *sql.DB, business_name, product_id string) (found bool, rating float32) {
	query := `SELECT average_rating FROM rates WHERE business_name = $1 and product_id = $2 
				ORDER BY id DESC LIMIT 1;`

	row := dB.QueryRow(query, business_name, product_id)
	switch err := row.Scan(&rating); err {
	case sql.ErrNoRows:
		return false, 0
	case nil:
		return true, rating
	default:
		log.Println("Uncaught error in getting rating from rates table, ", err)
		return false, 0
	}
}
