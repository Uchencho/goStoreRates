package rating

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

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

		if inRedis, averageRating := getFromRedis(company_name, pl.ProductID); inRedis {
			w.WriteHeader(http.StatusOK)
			resp := avgJson{
				ProductID:     pl.ProductID,
				AverageRating: averageRating,
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
		return

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, `{"Message" : "Method not allowed"}`)
		return
	}
}

func AllProductsRating(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	authorized, companyName, err := auth.CheckAuth(req)
	if !authorized {
		auth.UnAuthorizedResponse(w, err)
		return
	}

	switch req.Method {
	case http.MethodGet:
		w.WriteHeader(http.StatusOK)
		if found, allrates := getAllRates(config.Db, companyName); found {
			jsonResp, err := json.Marshal(allrates)
			if err != nil {
				log.Println("Error occured in marshalling json, ", err)
			}
			fmt.Fprint(w, string(jsonResp))
			return
		}
		fmt.Fprint(w, `[]`)
		return

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, `{"Message" : "Method not allowed"}`)
		return
	}
}

func ProductRatingDetail(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	authorized, companyName, err := auth.CheckAuth(req)
	if !authorized {
		auth.UnAuthorizedResponse(w, err)
		return
	}

	productID, err := strconv.Atoi(strings.TrimPrefix(req.URL.Path, "/ratings/"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"detail" : "Invalid url"}`)
		return
	}
	pID := strconv.Itoa(productID)

	switch req.Method {
	case http.MethodGet:

		if found, rate := getSpecificRate(config.Db, companyName, pID); found {
			w.WriteHeader(http.StatusOK)
			jsonResp, err := json.Marshal(rate)
			if err != nil {
				log.Println("Error occured in marshalling json, ", err)
			}
			fmt.Fprint(w, string(jsonResp))
			return
		}
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `{"detail" : "Not Found"}`)
		return
	}
}
