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

func AddRateToDB(dB *sql.DB, r rateJson) bool {
	query := `INSERT INTO rates (
		business_name, user_id, product_id,
		product_name, rating, average_rating
		) VALUES (
			$1, $2, $3, $4, $5, $6
		) RETURNING id`

	_, err := dB.Exec(query, r.BusinessName, r.UserID, r.ProductID,
		r.ProductName, r.Rating, r.AverageRating)
	if err != nil {
		log.Println("Error in saving details to db, ", err)
		return false
	}
	return true
}
