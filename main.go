package main

import (
	"log"
	"net/http"

	"github.com/Uchencho/goStoreRates/api/utils"

	"github.com/Uchencho/goStoreRates/api/users"
	"github.com/Uchencho/goStoreRates/config"
)

func main() {

	defer config.Db.Close()
	config.CreateUsersTable(config.Db)

	http.HandleFunc("/register", users.RegisterUser)
	if err := http.ListenAndServe(utils.GetServerddress(), nil); err != http.ErrServerClosed {
		log.Println("Error occured in listen and serve ", err)
	}
}
