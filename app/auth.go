package app

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/nishant01/mybookstore_items-api/utils/http_utils"
	"github.com/nishant01/procard-go-api/models/accounts"
	"github.com/nishant01/procard-go-api/utils/rest_errors"
	"net/http"
	"strings"
)

const (
	tokenKey = "thisIsTheJwtSecretPassword" //os.Getenv("token_key")
)
var tokenHeader2 string = ""

var jwtAuthentication = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		noAuth := []string{"/register", "/login", "/ping"} // List of endpoints that doesn't require auth
		requestPath := r.URL.Path

		//check if request does not need authentication, serve the request if it doesn't need it
		for _, value := range noAuth {
			if value == requestPath {
				next.ServeHTTP(w,r)
				return
			}
		}

		// We can obtain the session token from the requests cookies, which come with every request
		c, err := r.Cookie("procard_token")
		if err != nil {
			tokenHeader2 = ""
			// if err == http.ErrNoCookie {
				// If the cookie is not set, return an unauthorized status
				//w.WriteHeader(http.StatusUnauthorized)
				//return
			// }
			// For any other type of error, return a bad request status
			//w.WriteHeader(http.StatusBadRequest)
			//return
		}

		// Get the JWT string from the cookie
		tokenHeader2 = c.Value

		// response := make(map[string] interfase{})
		tokenHeader := r.Header.Get("Authorization") //Grab the token from the header

		if tokenHeader == "" || tokenHeader2 == "" { //Token is missing, returns with error code 403 Unauthorized
			respErr := rest_errors.NewUnauthorizedError("Missing auth token")
			http_utils.RespondError(w, respErr)
			return
		}

		splitToken := strings.Split(tokenHeader, " ") //The token normally comes in format `Bearer {token-body}`, we check if the retrieved token matched this requirement

		if len(splitToken) != 2 {
			respErr := rest_errors.NewUnauthorizedError("invalid/Malformed auth token")
			http_utils.RespondError(w, respErr)
			return
		}

		tokenPart := splitToken[1] //Grab the token part, what we are truly interested in
		tk := &accounts.Token{}

		token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(tokenKey), nil
		})

		if err != nil { //Malformed token, returns with http code 403 as usual
			respErr := rest_errors.NewUnauthorizedError("Malformed authentication token")
			http_utils.RespondError(w, respErr)
			return
		}

		if !token.Valid { //Token is invalid, maybe not signed on this server
			respErr := rest_errors.NewUnauthorizedError("Token is not valid.")
			http_utils.RespondError(w, respErr)
			return
		}

		//Everything went well, proceed with the request and set the caller to the user retrieved from the parsed token
		fmt.Sprintf("User %", tk.Email) //Useful for monitoring
		ctx := context.WithValue(r.Context(), "user", tk.UserID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r) //proceed in the middleware chain!
	})
}
