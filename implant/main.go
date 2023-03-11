package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
)

type command struct {
	Name      string `json:"name"`
	Script    string `json:"script"`
	Arguments string `json:"arguments"` // base64 encoded to avoid breaking json or cli arguments; plugins should decode the base64
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
	w.Header().Set("Content-Type", "application/json") // seems to not be respected??
	w.Write([]byte(out))
}

func startCommand(r *http.Request) response {
	body, err := ioutil.ReadAll(r.Body) // try to read the request body into a variable
	if err != nil {
		return response{
			Status:     false,
			StatusCode: http.StatusInternalServerError,
			Message:    fmt.Sprintf("Could not read request body: %s", err.Error()),
		}
	}

	var cmd command
	err = json.Unmarshal(body, &cmd) // try to coax the body from JSON into a command object, ignores unknown fields
	if err != nil {
		return response{
			Status:     false,
			StatusCode: http.StatusInternalServerError,
			Message:    fmt.Sprintf("Could not parse request body: %s", err.Error()),
		}
	}

	modpath := fmt.Sprintf("./modules/%s", cmd.Script)
	if _, err := os.Stat(modpath); err != nil {
		log.Printf("Module not found or inaccessible: %s", modpath)
		return response{
			Status:     false,
			StatusCode: http.StatusNotFound,
			Message:    fmt.Sprintf("Module not found or inaccessible: %s (%s)", modpath, err.Error()),
		}
	}

	log.Printf("Running %s", modpath)
	result, err := exec.Command("python3", modpath, cmd.Arguments).Output() // attempt to run the python module
	if err != nil {
		return response{
			Status:     false,
			StatusCode: http.StatusInternalServerError,
			Message:    fmt.Sprintf("Could not run the module: %s", err.Error()),
		}
	}

	enc := base64.StdEncoding.EncodeToString(result) // encode results to base64 to not have any funky chars in HTTP responses
	return response{
		Status:     true,
		StatusCode: http.StatusOK,
		Message:    enc,
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
