package security

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

type Authentication struct {
	encryptionKey []byte
}

type QuestClaims struct {
	IsAdmin  bool   `json:"isadmin"`
	AuthType string `json:"authtype"`
	jwt.StandardClaims
}

func NewAuthentication() *Authentication {

	return &Authentication{
		encryptionKey: []byte("abxx002151!"),
	}
}

func (a *Authentication) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow OPTIONS passthrough
		if (r.Method == "OPTIONS") {
			next.ServeHTTP(w, r)
			return
		}

		var token string

		// Get token from the Authorization header
		// format: Authorization: Bearer <token>
		tokens, ok := r.Header["Authorization"]
		if ok && len(tokens) >= 1 {
			token = tokens[0]
			token = strings.TrimPrefix(token, "Bearer ")
		}

		// If the token is empty...
		if token == "" {
			// If we get here, the required token is missing
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		// Now parse the token
		parsedToken, err := jwt.ParseWithClaims(token, &QuestClaims{}, func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				msg := fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				return nil, msg
			}
			return a.encryptionKey, nil
		})
		if err != nil {
			http.Error(w, "Error parsing token", http.StatusUnauthorized)
			return
		}

		// Check token is valid
		if claims, ok := parsedToken.Claims.(*QuestClaims); ok && parsedToken != nil && parsedToken.Valid {

			// Everything worked! Set the user in the context.
			fmt.Println("User authenticated")
			// TODO: "user" or "device", pull from token
			r.Header.Add("QUEST_AUTH_TYPE", claims.AuthType)
			r.Header.Add("QUEST_IS_ADMIN", strconv.FormatBool(claims.IsAdmin)) // "true"/"false"
			next.ServeHTTP(w, r)
			return
		}

		// Token is invalid
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	})
}
