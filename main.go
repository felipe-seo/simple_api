package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

//Account Struct(model)
type Account struct {
	Id         string    `json:"id"`
	Name       string    `json:"name"`
	Cpf        int       `json:"cpf"`
	Secret     string    `json:"secret"`
	Balance    int64     `json:"balance"`
	Created_at time.Time `json:"created_at"`
}

//init accounts variable as a slice of accounts struct
var accounts []Account

//get all accounts
//w variable, r variable --> respons and request
func getAccounts(w http.ResponseWriter, r *http.Request) {
	//set header value of content type to application/json or else it's just gonna be served as text
	w.Header().Set("Content-Type", "application/json")

	//encode accounts variable as json in the response
	json.NewEncoder(w).Encode(accounts)
}

//get a single account
func getAccount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//we need the id, so we set the params variable
	params := mux.Vars(r) //get params

	//loop through accounts and find the one with a matching id
	for _, item := range accounts {
		if item.Id == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Account{})
}

//create an account
func createAccount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var account Account
	_ = json.NewDecoder(r.Body).Decode(&account)
	account.Id = strconv.Itoa(rand.Intn(10-3+1) + 3)
	accounts = append(accounts, account)
	json.NewEncoder(w).Encode(account)
}

//update an account
func updateAccount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range accounts {
		if item.Id == params["id"] {
			accounts = append(accounts[:index], accounts[index+1:]...)
			w.Header().Set("Content-Type", "application/json")
			var account Account
			_ = json.NewDecoder(r.Body).Decode(&account)
			accounts = append(accounts, account)
		}
	}
	json.NewEncoder(w).Encode(accounts)
}

//delete an account
func deleteAccount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range accounts {
		if item.Id == params["id"] {
			accounts = append(accounts[:index], accounts[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(accounts)
}

func main() {
	//initialize the router
	r := mux.NewRouter()

	//data, having trouble with the date mock
	accounts = append(accounts, Account{Id: "1", Name: "Boris Fausto", Cpf: 12355567812, Secret: "LesPassword", Balance: 100000, Created_at: time.Now()})

	//call router handlers(the functions) and establish endpoints(the route)
	r.HandleFunc("/accounts", getAccounts).Methods("GET")
	r.HandleFunc("/accounts/{id}", getAccount).Methods("GET")
	r.HandleFunc("/accounts", createAccount).Methods("POST")
	r.HandleFunc("/accounts/{id}", updateAccount).Methods("PUT")
	r.HandleFunc("/accounts/{id}", deleteAccount).Methods("DELETE")
	//Wrapped with log, throws error in case it fails					receives a port and a router
	log.Fatal(http.ListenAndServe(":8000", r))
}
