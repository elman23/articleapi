package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/elman23/articleapi/pkg/db"
	"github.com/golang-jwt/jwt/v4"
)

var singleton db.DbServiceSingleton

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

// var users = map[string]string{
// 	"user1": "password1",
// 	"user2": "password2",
// }

// Create a struct to read the username and password from the request body
type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

// Create a struct that will be encoded to a JWT.
// We add jwt.RegisteredClaims as an embedded type, to provide fields like expiry time
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// Create the Signin handler
func Signin(w http.ResponseWriter, r *http.Request) {

	log.Println("Singing in...")

	var creds Credentials
	// Get the JSON body and decode into credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := singleton.GetService().GetUser(creds.Username)
	if err != nil {
		log.Println("Error getting user!")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Get the expected password from our in memory map
	// expectedPassword, ok := users[creds.Username]
	expectedPassword := user.Password

	// If a password exists for the given user
	// AND, if it is the same as the password we received, the we can move ahead
	// if NOT, then we return an "Unauthorized" status
	if expectedPassword != creds.Password {
		log.Println("Wrong password!")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	log.Println("Password control correct.")

	// Declare the expiration time of the token
	// here, we have kept it as 5 minutes
	expirationTime := time.Now().Add(5 * time.Minute)
	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		Username: creds.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Finally, we set the client cookie for "token" as the JWT we just generated
	// we also set an expiry time which is the same as the token itself
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
}

func IsAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		log.Println("Checking authorization...")

		// We can obtain the session token from the requests cookies, which come with every request
		c, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				// If the cookie is not set, return an unauthorized status
				log.Println("Unauthorized.")
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprintf(w, "Not authorized!")
				return
			}
			// For any other type of error, return a bad request status
			log.Println("Bad request.")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Bad request!")
			return
		}

		// Get the JWT string from the cookie
		tknStr := c.Value

		// Initialize a new instance of `Claims`
		claims := &Claims{}

		// Parse the JWT string and store the result in `claims`.
		// Note that we are passing the key in this method as well. This method will return an error
		// if the token is invalid (if it has expired according to the expiry time we set on sign in),
		// or if the signature does not match
		tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				log.Println("Unauthorized.")
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprintf(w, "Not authorized!")
				return
			}
			log.Println("Bad request.")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Bad request!")
			return
		}
		if !tkn.Valid {
			log.Println("Invalid authorization token.")
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "Not authorized!")
			return
		}
		log.Println("Authorization verified.")
		endpoint(w, r)
	})
}

func Refresh(w http.ResponseWriter, r *http.Request) {

	log.Println("Trying to refresh token...")

	// (BEGIN) Same code as in `IsAuthorized`
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			log.Println("Unauthorized.")
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "Not authorized!")
			return
		}
		// For any other type of error, return a bad request status
		log.Println("Bad request.")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Bad request!")
		return
	}

	// Get the JWT string from the cookie
	tknStr := c.Value

	// Initialize a new instance of `Claims`
	claims := &Claims{}

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			log.Println("Unauthorized.")
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "Not authorized!")
			return
		}
		log.Println("Bad request.")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Bad request!")
		return
	}
	if !tkn.Valid {
		log.Println("Invalid authorization token.")
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Not authorized!")
		return
	}
	// (END) Same code as in `IsAuthorized`

	// We ensure that a new token is not issued until enough time has elapsed
	// In this case, a new token will only be issued if the old token is within
	// 30 seconds of expiry. Otherwise, return a bad request status
	if time.Until(claims.ExpiresAt.Time) > 30*time.Second {
		log.Println("Old authorization token still valid for more than 30 seconds.")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Old authorization token still valid for more than 30 seconds.")
		return
	}

	// Now, create a new token for the current use, with a renewed expiration time
	expirationTime := time.Now().Add(5 * time.Minute)
	claims.ExpiresAt = jwt.NewNumericDate(expirationTime)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		log.Println("Internal server error.")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Internal server error.")
		return
	}

	// Set the new token as the users `token` cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
}

func Logout(w http.ResponseWriter, r *http.Request) {
	// immediately clear the token cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Expires: time.Now(),
	})
}
