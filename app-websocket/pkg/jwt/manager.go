package jwt

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type ErrorBody struct {
	Message string `json:"message"`
}

func ProcessError(w http.ResponseWriter, msg string, code int) {
	body := ErrorBody{
		Message: msg,
	}
	buf, _ := json.Marshal(body)

	w.WriteHeader(code)
	_, _ = w.Write(buf)
}

type AuthClaims struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	jwt.RegisteredClaims
}

func VerifyToken(token, secretKey string) (*jwt.Token, error) {
	// Parse the JWT token
	parsedToken, err := jwt.ParseWithClaims(token, &AuthClaims{}, func(token *jwt.Token) (interface{}, error) {

		// Verify the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method: %v", token.Header["alg"])
		}

		// Return the secret key used for signing
		return []byte(secretKey), nil
	})

	if err != nil {
		// Check if the error is due to token expiration
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("access token expired %s", err)
		}
		return nil, fmt.Errorf("invalid access token %s", err)
	}

	// Extract claims and validate
	claims, ok := parsedToken.Claims.(*AuthClaims)
	if !ok || !parsedToken.Valid {
		return nil, errors.New("invalid token")
	}

	// Validate registered claims
	if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New("token expired")
	}
	if claims.Issuer != "go-server" {
		return nil, errors.New("invalid issuer")
	}

	// Return the parsed token
	return parsedToken, nil
}

func Validate(logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			accessToken := r.Header.Get("Authorization")
			if accessToken == "" {

				accessToken = r.URL.Query().Get("access_token")
				if accessToken == "" {
					ProcessError(w, "empty auth header", http.StatusUnauthorized)
					return
				}
			}

			if strings.Contains(accessToken, "Bearer ") {
				accessToken = accessToken[len("Bearer "):]
			}
			_, err := VerifyToken(accessToken, "my_secret_key")

			if err != nil {
				if err.Error() == "access token expired" {

				}
				ProcessError(w, "invalid access token: ", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
