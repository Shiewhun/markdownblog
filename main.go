package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "<h1>Hello world</h1>")
}

func main() {
	port := os.Getenv("PORT")
	r := mux.NewRouter()
	r.HandleFunc("/", home).Methods("GET")
	fmt.Println("Starting application on port :4000")
	http.ListenAndServe(":"+port, r)
}
