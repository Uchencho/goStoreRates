package rating

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Uchencho/goStoreRates/config"

	"github.com/Uchencho/goStoreRates/api/auth"
)

func AddRate(w http.ResponseWriter, req *http.Request) {
	// Check for authentication
	authorized, company_name := auth.CheckAuth(req)
	if !authorized {
		err := errors.New("Check the auth function")
		auth.UnAuthorizedResponse(w, err)
	}

	// Use business_name returned to add rate to db
	var pl rateJson
	_ = json.NewDecoder(req.Body).Decode(&pl)
	if pl.UserID == "" || pl.Rating == 0 {
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

}
