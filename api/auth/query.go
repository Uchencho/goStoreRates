package auth

import (
	"database/sql"
	"log"
)

func getUser(dB *sql.DB, token string) (bool, string) {
	query := `SELECT company_name FROM users WHERE token = $1;`

	var company_name string
	row := dB.QueryRow(query, token)
	switch err := row.Scan(&company_name); err {
	case sql.ErrNoRows:
		return false, ""
	case nil:
		return true, company_name
	default:
		log.Println("Uncaught error in checking user, ", err)
		return false, ""
	}
}
