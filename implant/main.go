package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type command struct {
	Name      string `json:"name"`
	Script    string `json:"script"`
	Arguments string `json:"arguments"`
}

type response struct {
	Status     bool   `json:"status"` // true if ok, false otherwise
	StatusCode int    `json:"statuscode"`
	Message    string `json:"message"`
}

var notFoundResponse response = response{ // can't const a struct yay
	Status:     false,
	Message:    "URL or Method invalid",
	StatusCode: http.StatusNotFound,
}

var pongResponse response = response{
	Status:     true,
	Message:    "pong",
	StatusCode: http.StatusOK,
}

func serverError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{}")) // be polite and still send actual JSON back
}

func sendResponse(res response, w http.ResponseWriter) {
	out, err := json.Marshal(res)
	if err != nil {
		log.Printf("Error serializing json!  %+v\n", res)
		serverError(w)
		return
	}
	w.WriteHeader(res.StatusCode)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(out))
}

func startCommand(r *http.Request) response {
	return response{
		Status:     false,
		StatusCode: http.StatusNotImplemented,
		Message:    "Function not yet implemented",
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s\n", r.Method, r.URL)
	// handle requests on a one-by-one basis, there really arent that many...
	// refactor to a switch or two for method and URL if this needs more than the actual endpoint and a test one
	if r.Method == "GET" && r.URL.Path == "/ping" {
		sendResponse(pongResponse, w)
	} else if r.Method == "POST" && r.URL.Path == "/run" {
		sendResponse(startCommand(r), w)
	} else {
		sendResponse(notFoundResponse, w)
	}
}

func main() {
	log.Print("Starting server\n")
	http.HandleFunc("/", handleRequest)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
