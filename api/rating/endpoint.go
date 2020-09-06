package rating

import (
	"encoding/json"
	"fmt"
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
		if savedToDB := AddRateToDB(config.Db, pl); savedToDB {
			w.WriteHeader(http.StatusCreated)
			fmt.Fprint(w, `{"Message" : "Successfully saved"}`)
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
