package rating

import (
	"database/sql"
	"log"
)

type rateJson struct {
	BusinessName  string `json:"company_name"`
	UserID        string `json:"user_id"`
	ProductID     string `json:"product_id"`
	ProductName   string `json:"product_name"`
	Rating        int    `json:"rating"`
	AverageRating int    `json:"average_rating"`
}

func getCurrentAverage(dB *sql.DB, business_name, product_id string) (ratings []int) {
	query := `SELECT rating FROM rates WHERE business_name = $1 and product_id = $2;`

	rows, err := dB.Query(query, business_name, product_id)
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

func AddRateToDB(dB *sql.DB, r rateJson) bool {

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

	_, err := dB.Exec(query, r.BusinessName, r.UserID, r.ProductID,
		r.ProductName, r.Rating, r.AverageRating)
	if err != nil {
		log.Println("Error in saving details to db, ", err)
		return false
	}
	return true
}
