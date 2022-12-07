package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/copier"
	"io"
	"net/http"
	"strconv"
)

type Customer struct {
	ID        uint64 `json:"id"`
	NAME      string `json:"name"`
	ROLE      string `json:"role"`
	EMAIL     string `json:"email"`
	PHONE     string `json:"phone"`
	CONTACTED bool   `json:"contacted"`
}

var db = map[uint64]Customer{
	1: {ID: 1, NAME: "Al Bundy", ROLE: "Shoe Salesman", EMAIL: "al.bundy@garys.com", PHONE: "1078212232", CONTACTED: false},
	2: {ID: 2, NAME: "Bob Rooney", ROLE: "Software Engineer", EMAIL: "bob.rooney@google.com", PHONE: "1098218237", CONTACTED: true},
	3: {ID: 3, NAME: "Jefferson D'Arcy", ROLE: "Unemployed", EMAIL: "jefferson.darcy@garys.com", PHONE: "1048211189", CONTACTED: false},
}

func getCustomers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	var customers []Customer
	for _, customer := range db {
		customers = append(customers, customer)
	}
	json.NewEncoder(w).Encode(customers)
}

func getCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	json.NewEncoder(w).Encode(db[id])
}

func addCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newCustomer Customer
	reqBody, _ := io.ReadAll(r.Body)
	valid := json.Valid(reqBody)
	if !valid {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	json.Unmarshal(reqBody, &newCustomer)

	id := uint64(len(db)) + 1
	newCustomer.ID = id
	db[id] = newCustomer

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newCustomer)
}

func replaceCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	var replaceCustomer Customer
	reqBody, _ := io.ReadAll(r.Body)

	valid := json.Valid(reqBody)
	if !valid {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	json.Unmarshal(reqBody, &replaceCustomer)

	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	replaceCustomer.ID = id
	db[id] = replaceCustomer

	json.NewEncoder(w).Encode(replaceCustomer)
}

func updateCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	var update Customer
	reqBody, _ := io.ReadAll(r.Body)

	valid := json.Valid(reqBody)
	if !valid {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	json.Unmarshal(reqBody, &update)

	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)

	existingCustomer := db[id]
	copier.CopyWithOption(&existingCustomer, &update, copier.Option{IgnoreEmpty: true, DeepCopy: true})
	db[id] = existingCustomer
	json.NewEncoder(w).Encode(existingCustomer)
}

func deleteCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	delete(db, id)
	json.NewEncoder(w).Encode(db)
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/customers", getCustomers).Methods("GET")
	router.HandleFunc("/customers", addCustomer).Methods("POST")
	router.HandleFunc("/customers/{id}", getCustomer).Methods("GET")
	router.HandleFunc("/customers/{id}", deleteCustomer).Methods("DELETE")
	router.HandleFunc("/customers/{id}", replaceCustomer).Methods("PUT")
	router.HandleFunc("/customers/{id}", updateCustomer).Methods("PATCH")

	fmt.Println("Server is starting on port 8080...")
	http.ListenAndServe(":8080", router)

}
