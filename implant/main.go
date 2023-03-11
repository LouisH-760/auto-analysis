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
	Status  bool   `json:"status"` // true if ok, false otherwise
	Message string `json:"message"`
}

func serverError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{}")) // be polite and still send actual JSON back
}

func pong() response {
	return response{
		Status:  true,
		Message: "pong",
	}
}

func notFound() response {
	return response{
		Status:  false,
		Message: "URL or Method invalid",
	}
}

func sendResponse(res response, status int, w http.ResponseWriter) {
	out, err := json.Marshal(res)
	if err != nil {
		log.Printf("Error serializing json!  %+v\n", res)
		serverError(w)
		return
	}
	log.Printf("Sending %s\n", out)
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(out))
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s\n", r.Method, r.URL)
	// handle requests on a one-by-one basis, there really arent that many...
	// refactor to a switch or two for method and URL if this needs more than the actual endpoint and a test one
	if r.Method == "GET" && r.URL.Path == "/ping" {
		sendResponse(pong(), http.StatusOK, w)
	} else {
		sendResponse(notFound(), http.StatusNotFound, w)
	}
}

func main() {
	log.Print("Starting server\n")
	http.HandleFunc("/", handleRequest)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
