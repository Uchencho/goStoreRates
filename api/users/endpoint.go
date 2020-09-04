package users

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Uchencho/goStoreRates/config"

	"github.com/Uchencho/goStoreRates/api/auth"
)

func RegisterUser(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	switch req.Method {
	case http.MethodPost:
		var (
			user Account
			err  error
		)
		_ = json.NewDecoder(req.Body).Decode(&user)
		if user.CompanyName == "" || user.Password == "" || user.Email == "" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, `{"message":"email, company_name and password is necessary")`)
			return
		}

		user.LastLogin = time.Now()
		user.CreatedOn = time.Now()

		user.Password, err = auth.HashPassword(user.Password)
		if err != nil {
			log.Println("Error occurred while hashing password, ", err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `{"message":"Something went wrong"}`)
			return
		}

		user.Token, err = auth.GenerateToken(user.CompanyName)
		if err != nil {
			log.Println("Errror occured while generating token")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `{"message":"Something went wrong"}`)
			return
		}

		if created := addUser(config.Db, user); created {
			w.WriteHeader(http.StatusCreated)
			fmt.Fprint(w, `{"Message" : "Successfully Created"}`)
		} else {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, `{"Message" : "User already exists, please login"}`)
			return
		}

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, `{"Message" : "Method not allowed"}`)
		return
	}
}
