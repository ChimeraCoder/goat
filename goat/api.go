package goat

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// apiError represents an error response from the API
type apiError struct {
	Error string `json:"error"`
}

// apiRouter handles the routing of HTTP API requests
func apiRouter(w http.ResponseWriter, r *http.Request) {
	// API is read-only, at least for the time being
	if r.Method != "GET" {
		http.Error(w, string(apiErrorResponse("Method not allowed")), 405)
		return
	}

	// Log API calls
	log.Printf("API: [http %s] %s %s\n", r.RemoteAddr, r.Method, r.URL.Path)

	// Split request path
	urlArr := strings.Split(r.URL.Path, "/")

	// Verify API method set
	if len(urlArr) < 3 {
		http.Error(w, string(apiErrorResponse("No API call")), 404)
		return
	}

	// Check for an ID
	ID := -1
	if len(urlArr) == 4 {
		i, err := strconv.Atoi(urlArr[3])
		if err != nil || i < 1 {
			http.Error(w, string(apiErrorResponse("Invalid integer ID")), 400)
			return
		} else {
			ID = i
		}
	}

	// API response chan
	apiChan := make(chan []byte)

	// Choose API method
	switch urlArr[2] {
	// Files on tracker
	case "files":
		go getFilesJSON(ID, apiChan)
	// Server status
	case "status":
		go getStatusJSON(apiChan)
	// Return error response
	default:
		http.Error(w, string(apiErrorResponse("Undefined API call")), 404)
		close(apiChan)
		return
	}

	w.Write(<-apiChan)
	close(apiChan)
	return
}

// apiErrorResponse returns an apiError as JSON
func apiErrorResponse(msg string) []byte {
	res := apiError{
		msg,
	}

	out, err := json.Marshal(res)
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	return out
}
