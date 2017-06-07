package security

import (
  "fmt"
  "net/http"
  "strings"

	"github.com/dgrijalva/jwt-go"
)

type Authentication struct {
  encryptionKey []byte
}

func NewAuthentication() *Authentication {

  return &Authentication{
    encryptionKey: []byte("abxx002151!"),
  }
}

func (a *Authentication) Authenticate(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
		parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
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
    		if parsedToken != nil && parsedToken.Valid {
    			// Everything worked! Set the user in the context.
          fmt.Println("User authenticated")
    			next.ServeHTTP(w, r)
          return
    		}

    		// Token is invalid
    		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
    		return
  })
}
