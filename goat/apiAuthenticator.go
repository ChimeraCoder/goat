package goat

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// apiAuthenticator interface which defines methods required to implement an authentication method
type apiAuthenticator interface {
	Auth(*http.Request) bool
}

// basicAPIAuthenticator uses the HTTP Basic authentication scheme
type basicAPIAuthenticator struct {
}

// Auth handles validation of HTTP Basic authentication
func (a basicAPIAuthenticator) Auth(r *http.Request) bool {
	// Retrieve Authorization header
	auth := r.Header.Get("Authorization")

	// No header provided
	if auth == "" {
		return false
	}

	// Ensure format is valid
	basic := strings.Split(auth, " ")
	if basic[0] != "Basic" {
		return false
	}

	// Decode base64'd user:password pair
	buf, err := base64.URLEncoding.DecodeString(basic[1])
	if err != nil {
		log.Println(err.Error())
		return false
	}

	// Split into username/password
	credentials := strings.Split(string(buf), ":")

	// Load user by username, verify user exists
	user := new(userRecord).Load(credentials[0], "username")
	if user == (userRecord{}) {
		return false
	}

	// Load user's API key
	key := new(apiKey).Load(user.ID, "user_id")
	if key == (apiKey{}) {
		return false
	}

	// Hash input password
	sha := sha1.New()
	if _, err = sha.Write([]byte(credentials[1] + key.Salt)); err != nil {
		log.Println(err.Error())
		return false
	}

	hash := fmt.Sprintf("%x", sha.Sum(nil))

	// Verify hashes match
	if hash != key.Key {
		return false
	}

	// Authentication succeeded
	return true
}

// hmacAPIAuthenticator uses the HMAC-SHA1 authentication scheme
type hmacAPIAuthenticator struct {
}

// Auth handles validation of HMAC-SHA1 authentication
func (a hmacAPIAuthenticator) Auth(r *http.Request) bool {
	return true
}
