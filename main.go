package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mauriliommachado/go-metric/engine"
)

var wapi *engine.WorkerAPI

func main() {
	fmt.Printf("hello, world\n")
	wapi = engine.New(1)
	wapi.Storage.MockMetrics()
	initAPI()
}

func initAPI() {
	router := httprouter.New()
	router.PUT("/metric", index)
	router.GET("/metric", get)
	log.Println("Starting server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func get(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	j, _ := json.Marshal(wapi.Storage.GetItens())
	w.Header().Add("Content-Type", "application/json")
	fmt.Fprintf(w, "%s", j)
}

func index(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var data map[string]interface{}
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &data)
	wapi.SubmitWork(data)
	writeOKResponse(w)
}

func writeOKResponse(w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// Writes the error response as a Standard API JSON response with a response code
func writeErrorResponse(w http.ResponseWriter, errorCode int, errorMsg string) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(errorCode)
}
