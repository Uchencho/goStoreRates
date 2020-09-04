package users

import (
	"database/sql"
	"log"
	"time"
)

type Account struct {
	ID          uint      `json:"id"`
	CompanyName string    `json:"company_name"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	Token       string    `json:"token"`
	Activated   bool      `json:"activated"`
	CreatedOn   time.Time `json:"created_on"`
	LastLogin   time.Time `json:"last_login"`
}

func checkUser(db *sql.DB, user Account) bool {
	query := `SELECT email from users WHERE email = $1;`

	var email string
	row := db.QueryRow(query, user.Email)
	switch err := row.Scan(&email); err {
	case sql.ErrNoRows:
		return false
	case nil:
		return true
	default:
		log.Println("Uncaught error in checking user, ", err)
		return true
	}
}

func addUser(db *sql.DB, user Account) bool {

	if userExists := checkUser(db, user); userExists {
		return false
	}

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
