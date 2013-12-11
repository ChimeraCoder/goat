package goat

import (
	"net"
	"net/http"
)

// ConnHandler interface method Handle defines how to handle incoming network connections
type ConnHandler interface {
	Handle(l net.Listener, logChan chan string)
}

// HttpConnHandler handles incoming HTTP (TCP) network connections
type HttpConnHandler struct {
}

// Handle incoming HTTP connections and serve
func (h HttpConnHandler) Handle(l net.Listener, logChan chan string) {
	// Set up HTTP routes for handling functions
	http.HandleFunc("/", parseHttp)
	http.HandleFunc("/announce", parseHttp)

	// Serve HTTP requests
	err := http.Serve(l, nil)
	if err != nil {
		logChan <- err.Error()
	}
}

// Parse incoming HTTP connections before making tracker calls
func parseHttp(w http.ResponseWriter, r *http.Request) {
	// Create channel to return bencoded response to client on
	resChan := make(chan []byte)

	// Handle tracker functions via different URLs
	switch r.URL.Path {
	// Tracker announce
	case "/announce":
		go TrackerAnnounce(resChan)
	// Any undefined handlers
	default:
		go TrackerError(resChan, "Malformed announce")
	}

	// Add header to identify goat
	w.Header().Add("Server", APP+"-git")

	// Wait for response, and send it when ready
	w.Write(<-resChan)
}

// UdpConnHandler handles incoming UDP network connections
type UdpConnHandler struct {
}

// Handle incoming UDP connections and return response
func (u UdpConnHandler) Handle(l net.Listener, logChan chan string) {
}
