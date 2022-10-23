package main

import (
	"log"
	"net/http"
	"redirector/handlers"

	"github.com/gorilla/mux"
)

func handleRoutes(router *mux.Router) {
	router.PathPrefix("/").HandlerFunc(handlers.RedirectorGETHandler).Methods("GET")
	router.PathPrefix("/").HandlerFunc(handlers.RedirectorHEADHandler).Methods("HEAD")
}

func createServer() {
	muxRoute := mux.NewRouter().StrictSlash(false)
	log.Println("starting server")
	handleRoutes(muxRoute)
	log.Fatal(http.ListenAndServe("0.0.0.0:10010", muxRoute))
}

func main() {
	createServer()
}
