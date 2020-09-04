package config

import (
	"database/sql"
	"log"
)

func CreateUsersTable(db *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS users (
		id serial PRIMARY KEY,
		company_name VARCHAR (100) UNIQUE NOT NULL,
		email VARCHAR (100) UNIQUE NOT NULL,
		password VARCHAR (200) NOT NULL,
		token VARCHAR(200),
		activated BOOL,
		created_on TIMESTAMP,
		last_login TIMESTAMP
	);`

	_, err := db.Exec(query)
	if err != nil {
		log.Println(err)
	}
}
