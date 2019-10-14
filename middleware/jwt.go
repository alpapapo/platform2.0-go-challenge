package middleware

import (
	"context"
	"fmt"
	c "github.com/GlobalWebIndex/platform2.0-go-challenge/common"
	"github.com/GlobalWebIndex/platform2.0-go-challenge/models"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"os"
	"strings"
)

var JwtAuthentication = func(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		notAuth := []string{"/api/user/register", "/api/user/login", "/populate/assets"} //List of endpoints that doesn't require auth
		requestPath := r.URL.Path                               //current request path

		//check if request does not need authentication, serve the request if it doesn't need it
		for _, value := range notAuth {

			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}

		response := make(map[string]interface{})
		tokenHeader := r.Header.Get("Authorization") //Grab the token from the header

		if tokenHeader == "" { //Token is missing, returns with error code 403 Unauthorized
			response = c.Message(false, "Missing auth token")
			c.Respond(w, http.StatusForbidden, response)
			return
		}

		splitted := strings.Split(tokenHeader, " ") //The token normally comes in format `Bearer {token-body}`, we check if the retrieved token matched this requirement
		if len(splitted) != 2 {
			response = c.Message(false, "Invalid/Malformed auth token")
			c.Respond(w, http.StatusForbidden, response)
			return
		}

		tokenPart := splitted[1] //Grab the token part, what we are truly interested in
		tk := &models.Token{}

		token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_TOKEN_PASSWORD")), nil
		})

		if err != nil { //Malformed token, returns with http code 403 as usual
			response = c.Message(false, "Malformed authentication token")
			c.Respond(w, http.StatusForbidden, response)
			return
		}

		if !token.Valid { //Token is invalid, maybe not signed on this server
			response = c.Message(false, "Token is not valid.")
			c.Respond(w, http.StatusForbidden, response)
			return
		}

		//Everything went well, proceed with the request and set the caller to the user retrieved from the parsed token
		fmt.Sprintf("User %d", tk.UserId) //Useful for monitoring
		ctx := context.WithValue(r.Context(), "user", tk.UserId)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r) //proceed in the middleware chain!
	});
}
