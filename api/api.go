package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/MohamedAbdeen21/cloud-store/cmd"
	"github.com/gorilla/mux"
)

func notFoundFunc(w http.ResponseWriter, r *http.Request) {
	http.Error(w, fmt.Sprintf("Not found: %s", r.RequestURI), http.StatusNotFound)
	log.Println("Not found", r.RequestURI)
}

func root(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello from root")
}

func main() {
	r := mux.NewRouter()

	r.NotFoundHandler = http.HandlerFunc(notFoundFunc)
	r.HandleFunc("/", root).Methods("GET")
	r.HandleFunc("/signin/", cmd.Signin).Methods("POST")
	r.HandleFunc("/signup/", cmd.Signup).Methods("POST")
	r.HandleFunc("/query/", cmd.CustomQuery).Methods("POST")

	r.HandleFunc("/add/", cmd.AddItem).Methods("PUT")
	r.HandleFunc("/remove/", cmd.RemoveItem).Methods("PUT")
	r.HandleFunc("/items/", cmd.CheckItems).Methods("GET")
	r.HandleFunc("/checkout/", cmd.Checkout).Methods("GET")

	PORT := ":8080"
	log.Println("Started server on port", PORT)
	http.ListenAndServe(PORT, r)
}
