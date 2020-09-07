package rating

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Uchencho/goStoreRates/api/auth"
	"github.com/Uchencho/goStoreRates/config"
)

func AddRate(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	authorized, company_name, err := auth.CheckAuth(req)
	if !authorized {
		auth.UnAuthorizedResponse(w, err)
		return
	}

	switch req.Method {
	case http.MethodPost:
		var pl rateJson
		_ = json.NewDecoder(req.Body).Decode(&pl)
		if pl.UserID == "" || pl.Rating == 0 || pl.ProductName == "" || pl.ProductID == "" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, `{"Message" : "Invalid Payload"}`)
			return
		}

		pl.BusinessName = company_name
		if savedToDB, details := AddRateToDB(config.Db, pl); savedToDB {
			w.WriteHeader(http.StatusCreated)
			resp := averageJson{
				ProductID:     details.ProductID,
				ProductName:   details.ProductName,
				AverageRating: details.AverageRating,
			}
			jsonresp, err := json.Marshal(resp)
			if err != nil {
				log.Println(err)
				fmt.Fprint(w, `{"Message" : "Saved successfully"}`)
				return
			}
			fmt.Fprint(w, string(jsonresp))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, `{"Message" : "Something went wrong"}`)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, `{"Message" : "Method not allowed"}`)
		return
	}

}

func CurrentAverageRating(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	authorized, company_name, err := auth.CheckAuth(req)
	if !authorized {
		auth.UnAuthorizedResponse(w, err)
		return
	}

	switch req.Method {
	case http.MethodPost:
		var pl rateJson
		_ = json.NewDecoder(req.Body).Decode(&pl)
		if pl.ProductID == "" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, `{"product_id" : "Product ID is compulsory"}`)
			return
		}
		pl.BusinessName = company_name

		if ratingAvailable, avergaeRating := currentRating(config.Db, pl.BusinessName, pl.ProductID); ratingAvailable {
			w.WriteHeader(http.StatusOK)
			resp := avgJson{
				ProductID:     pl.ProductID,
				AverageRating: avergaeRating,
			}
			jsonresp, err := json.Marshal(resp)
			if err != nil {
				log.Println(err)
			}
			fmt.Fprint(w, string(jsonresp))
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"Message" : "Product has not been rated"}`)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, `{"Message" : "Method not allowed"}`)
		return
	}
}
