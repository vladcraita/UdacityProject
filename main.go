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

	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)

	if customer, ok := db[id]; ok {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(customer)
	} else {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(getErrorMessage("No customer found with id " + strconv.FormatUint(id, 10)))
	}
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

	if newCustomer.isValid() {
		id := uint64(len(db)) + 1
		newCustomer.ID = id
		db[id] = newCustomer

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newCustomer)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(getErrorMessage("Name and Email are mandatory fields"))
	}
}

func updateCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var update Customer
	reqBody, _ := io.ReadAll(r.Body)

	valid := json.Valid(reqBody)
	if !valid {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	json.Unmarshal(reqBody, &update)

	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)

	update.setId(id)

	existingCustomer := db[id]
	copier.CopyWithOption(&existingCustomer, &update, copier.Option{IgnoreEmpty: true, DeepCopy: true})

	if existingCustomer.isValid() {
		db[id] = existingCustomer

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(existingCustomer)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(getErrorMessage("Name and Email are mandatory fields"))
	}
}

func deleteCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	delete(db, id)
	json.NewEncoder(w).Encode(db)
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

	if replaceCustomer.isValid() {
		id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
		replaceCustomer.ID = id
		db[id] = replaceCustomer

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(replaceCustomer)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(getErrorMessage("Name and Email are mandatory fields for replacing a customer at id. If a partial update is desired use PATCH"))
	}
}

func (customer *Customer) setId(id uint64) {
	customer.ID = id
}

func (customer *Customer) isValid() bool {
	return len(customer.NAME) > 0 && len(customer.EMAIL) > 0
}

func showIndexPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}

func getErrorMessage(message string) map[string]string {
	return map[string]string{
		"error": message,
	}
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", showIndexPage)
	router.HandleFunc("/customers", getCustomers).Methods("GET")
	router.HandleFunc("/customers", addCustomer).Methods("POST")
	router.HandleFunc("/customers/{id}", getCustomer).Methods("GET")
	router.HandleFunc("/customers/{id}", deleteCustomer).Methods("DELETE")
	router.HandleFunc("/customers/{id}", replaceCustomer).Methods("PUT")
	router.HandleFunc("/customers/{id}", updateCustomer).Methods("PATCH")

	fmt.Println("Server is starting on port 8080...")
	http.ListenAndServe(":8080", router)

}
