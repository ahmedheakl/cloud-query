package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
	"github.com/gorilla/mux"
)

func notFoundFunc(w http.ResponseWriter, r *http.Request) {
	log.Println("Not found", r.RequestURI)
	http.Error(w, fmt.Sprintf("Not found: %s", r.RequestURI), http.StatusNotFound)
}

func root(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello from root")
}

func route(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(ExecuteQuery("SELECT * FROM data LIMIT 10"))
}

func main() {
	r := mux.NewRouter()

	r.NotFoundHandler = http.HandlerFunc(notFoundFunc)
	r.HandleFunc("/", root)
	r.HandleFunc("/select-all/", route)

	adapter := gorillamux.NewV2(r)
	lambda.Start(adapter.ProxyWithContext)
}
