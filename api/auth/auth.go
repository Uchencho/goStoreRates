package auth

import (
	"encoding/base64"
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

func CheckAuth(r *http.Request) (bool, string) {
	if r.Header["Authorization"] != nil {

		if len(strings.Split(r.Header["Authorization"][0], " ")) < 2 {
			return false, ""
		}

		accesstoken := strings.Split(r.Header["Authorization"][0], " ")[1]
		found, company_name := getUser(config.Db, accesstoken)
		return found, company_name

	}
	return false, ""
}

func UnAuthorizedResponse(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusForbidden)
	fmt.Fprint(w, err.Error())
}
