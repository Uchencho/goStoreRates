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

type allrateJson struct {
	ID          int    `json:"id"`
	UserID      string `json:"user_id"`
	ProductID   string `json:"product_id"`
	ProductName string `json:"product_name"`
	Rating      int    `json:"rating"`
}

type averageJson struct {
	ProductID     string `json:"product_id"`
	ProductName   string `json:"product_name"`
	AverageRating int    `json:"average_rating"`
}

type avgJson struct {
	ProductID     string `json:"product_id"`
	AverageRating int    `json:"average_rating"`
}

func sendtoRedis(businessName, productID string, averageRating int) bool {

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       1,
	})

	// Check if the business has products there already
	val, err := rdb.Get(businessName).Result()
	if err == redis.Nil {

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
	} else if err != nil {
		log.Println("Apparently there is a big problem, ", err)
		return false
	} else {

		// Key actually exists
		availableRates := map[string]int{}

		err := json.Unmarshal([]byte(val), &availableRates)
		if err != nil {
			log.Println("Error in getting map from redis, ", err)
			return false
		}
		availableRates[productID] = averageRating

		jsonString, err := json.Marshal(availableRates)
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

}

// Retrieves the average rate stored for a particular product
func getFromRedis(businessName, productID string) (bool, int) {

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       1,
	})

	val, err := rdb.Get(businessName).Result()
	if err == redis.Nil {
		log.Printf("%s does not exist ", businessName)
		return false, 0
	} else if err != nil {
		log.Println("Apparently there is a big problem, ", err)
		return false, 0
	} else {
		availableRates := map[string]int{}

		err := json.Unmarshal([]byte(val), &availableRates)
		if err != nil {
			log.Println("Error in getting map from redis, ", err)
			return false, 0
		}
		if value, ok := availableRates[productID]; ok {
			return true, value
		}
		return false, 0
	}
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

// func currentRating(dB *sql.DB, business_name, product_id string) (found bool, rating float32) {
// 	query := `SELECT average_rating FROM rates WHERE business_name = $1 and product_id = $2
// 				ORDER BY id DESC LIMIT 1;`

// 	row := dB.QueryRow(query, business_name, product_id)
// 	switch err := row.Scan(&rating); err {
// 	case sql.ErrNoRows:
// 		return false, 0
// 	case nil:
// 		return true, rating
// 	default:
// 		log.Println("Uncaught error in getting rating from rates table, ", err)
// 		return false, 0
// 	}
// }

func getAllRates(dB *sql.DB, businessName string) (found bool, allRates []allrateJson) {
	query := `SELECT id, user_id, product_id, product_name, rating FROM 
				rates WHERE business_name = $1;`

	rows, err := dB.Query(query, businessName)
	if err != nil {
		log.Println("Error occured in retrieving all ratings")
		return false, []allrateJson{}
	}
	defer rows.Close()

	var rate allrateJson
	for rows.Next() {
		err = rows.Scan(&rate.ID, &rate.UserID, &rate.ProductID,
			&rate.ProductName, &rate.Rating)
		if err != nil {
			log.Println(err)
			return false, []allrateJson{}
		}
		allRates = append(allRates, rate)
	}

	err = rows.Err()
	if err != nil {
		log.Println(err)
		return false, []allrateJson{}
	}
	if len(allRates) > 0 {
		return true, allRates

	}
	return false, []allrateJson{}
}

func getSpecificRate(dB *sql.DB, businessName, ratingID string) (found bool, rate allrateJson) {
	query := `SELECT id, user_id, product_id, product_name, rating FROM 
				rates WHERE business_name = $1 AND id = $2;`

	row := dB.QueryRow(query, businessName, ratingID)
	switch err := row.Scan(&rate.ID, &rate.UserID, &rate.ProductID,
		&rate.ProductName, &rate.Rating); err {
	case sql.ErrNoRows:
		return false, allrateJson{}
	case nil:
		return true, rate
	default:
		log.Println("Uncaught error in getting rating from rates table, ", err)
		return false, allrateJson{}
	}
}

func deleteRate(dB *sql.DB, ratingID string) bool {
	query := `DELETE FROM rates WHERE id = $1`
	_, err := dB.Exec(query, ratingID)
	if err != nil {
		log.Println("Error returned in deleting rating, ", err)
		return false
	}
	return true
}

func updateRedis(dB *sql.DB, r rateJson, ratingID string) bool {

	query := `SELECT product_id FROM rates WHERE id = $1;`

	row := dB.QueryRow(query, ratingID)
	switch err := row.Scan(&r.ProductID); err {
	case sql.ErrNoRows:
		log.Println("No rows were found")
		return false
	case nil:
		ratings := getCurrentAverage(dB, r.BusinessName, r.ProductID)

		var theSum int
		for i := 0; i < len(ratings); i++ {
			theSum += ratings[i]
		}
		r.AverageRating = theSum / len(ratings)

		sentToRedis := sendtoRedis(r.BusinessName, r.ProductID, r.AverageRating)
		if !sentToRedis {
			log.Println("Could not send value to redis")
			return false
		}
		return true
	default:
		log.Println("Uncaught error in getting rating from rates table, ", err)
		return false
	}
}
