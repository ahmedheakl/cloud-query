package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func notFoundFunc(w http.ResponseWriter, r *http.Request) {
	http.Error(w, fmt.Sprintf("Not found: %s", r.RequestURI), http.StatusNotFound)
	log.Println("Not found", r.RequestURI)
}
func root(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello from root")
}

func all(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	resp, err := executeQuery("SELECT * FROM test LIMIT 10")
	if err != nil {
		w.Write([]byte(err.Error()))
		log.Println(err.Error())
	}
	w.Write(resp)
}

func customQuery(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	query, err := io.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(err.Error()))
		log.Println(err.Error())
	}
	log.Printf("Executing query %s", query)

	resp, err := executeQuery(string(query))
	if err != nil {
		w.Write([]byte(err.Error()))
		log.Println(err.Error())
	}
	w.Write(resp)
}

func main() {
	PORT := ":8080"
	r := mux.NewRouter()

	_, err := getOrCreate()
	if err != nil {
		log.Fatalf(err.Error())
	}

	r.NotFoundHandler = http.HandlerFunc(notFoundFunc)
	r.HandleFunc("/", root)
	r.HandleFunc("/select-all/", all)
	r.HandleFunc("/query/", customQuery)
	log.Println("Started server on port", PORT)
	http.ListenAndServe(PORT, r)
}
