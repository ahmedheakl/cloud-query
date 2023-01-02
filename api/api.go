package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MohamedAbdeen21/cloud-store/cmd"
	log "github.com/MohamedAbdeen21/cloud-store/logger"
	"github.com/gorilla/mux"
)

func notFoundFunc(w http.ResponseWriter, r *http.Request) {
	http.Error(w, fmt.Sprintf("Not found: %s", r.RequestURI), http.StatusNotFound)
	log.Info.Println("Not found", r.RequestURI)
}

func root(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	msg, _ := json.Marshal(map[string]string{"message": "Hello from root!"})
	w.Write(msg)
}

func main() {
	r := mux.NewRouter()

	log.Init()
	log.Info.Println("starting server ...")
	err := cmd.InitConnections()
	if err != nil {
		return
	}

	r.NotFoundHandler = http.HandlerFunc(notFoundFunc)
	r.HandleFunc("/", root).Methods("GET")
	r.HandleFunc("/login/", cmd.Login).Methods("POST")
	r.HandleFunc("/signup/", cmd.Signup).Methods("POST")
	r.HandleFunc("/logout/", cmd.Logout).Methods("GET")
	r.HandleFunc("/query/", cmd.CustomQuery).Methods("POST")

	r.HandleFunc("/add/", cmd.AddItem).Methods("POST")
	r.HandleFunc("/remove/", cmd.RemoveItem).Methods("POST")
	r.HandleFunc("/items/", cmd.CheckItems).Methods("GET")
	r.HandleFunc("/checkout/", cmd.Checkout).Methods("GET")
	r.HandleFunc("/cookie/{email}/", cmd.GenerateCookie).Methods("GET")

	PORT := ":8080"
	log.Info.Println("started server on port", PORT)
	fmt.Println("started server on port", PORT)
	err = http.ListenAndServe(PORT, r)
	if err != nil {
		log.Fatal.Printf("can't start server: %s", err.Error())
	}
}
