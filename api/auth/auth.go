package auth

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Uchencho/goStoreRates/config"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 4)
	return string(bytes), err
}

func GenerateToken(username string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(username), 5)
	return base64.StdEncoding.EncodeToString(bytes), err
}

func CheckAuth(r *http.Request) (bool, string, error) {
	if r.Header["Authorization"] != nil {

		if len(strings.Split(r.Header["Authorization"][0], " ")) < 2 {
			return false, "", errors.New("Authentication was not provided")
		}

		accesstoken := strings.Split(r.Header["Authorization"][0], " ")[1]
		found, company_name := getUser(config.Db, accesstoken)
		if found {
			return found, company_name, nil
		}
		return found, company_name, errors.New("Invalid credentials")

	}
	return false, "", errors.New("Authentication was not provided")
}

func UnAuthorizedResponse(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusForbidden)
	fmt.Fprint(w, err.Error())
}
