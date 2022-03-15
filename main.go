package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
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

type Transfer struct {
	Id                     string    `json:"id"`
	Account_origin_id      string    `json:"account_origin_id"`
	Account_destination_id string    `json:"account_destination_id"`
	Amount                 int       `json:"amount"`
	Created_at             time.Time `json:"created_at"`
}

//init accounts variable as a slice of accounts struct
var accounts []Account

var transfers []Transfer

//get all accounts
//w variable, r variable --> respons and request
func getAccounts(w http.ResponseWriter, r *http.Request) {
	//set header value of content type to application/json or else it's just gonna be served as text
	w.Header().Set("Content-Type", "application/json")

	//encode accounts variable as json in the response
	json.NewEncoder(w).Encode(accounts)
}

//get a single account(not used in the challenge, change it to specific balance only, also adjust the route)
//handle error on unexistent account
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
	account.Id = uuid.New().String() //insert uuid

	hashedSecret, err := HashPassword(account.Secret)
	//tackle error handling
	if err != nil {
		//fmt.Errorf("error: %s", err)
		fmt.Println(err)
	}

	account.Secret = hashedSecret
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
/*
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
*/
//create a hash for the secret
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

//used on login to check if the informed password translates to the stored hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

//transfers
func getTransfers(w http.ResponseWriter, r *http.Request) {
	//set header value of content type to application/json or else it's just gonna be served as text
	w.Header().Set("Content-Type", "application/json")

	//encode accounts variable as json in the response
	json.NewEncoder(w).Encode(transfers)
}

func createTransfer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var transfer Transfer
	_ = json.NewDecoder(r.Body).Decode(&transfer)
	transfer.Id = uuid.New().String() //insert uuid - does it need to be unique?
	fmt.Println(transfer.Account_origin_id)

	transfers = append(transfers, transfer)

	json.NewEncoder(w).Encode(transfer)
}

//validate transfer
//login after dealing with transfers

//
func main() {
	//initialize the router
	r := mux.NewRouter()

	hash, err := HashPassword("LesPassword")
	if err != nil {
		log.Fatal(err)
	}

	//data
	accounts = append(accounts, Account{Id: "dbd70daa-356c-4bde-be27-a13c3c47d44d", Name: "Boris Fausto", Cpf: 12355567812, Secret: hash, Balance: 100000, Created_at: time.Now()})
	accounts = append(accounts, Account{Id: "dbd74daa-356c-4bde-be27-a13c3c47d49d", Name: "Doriana Yates", Cpf: 11111111506, Secret: hash, Balance: 5000, Created_at: time.Now()})

	transfers = append(transfers, Transfer{Id: "1", Account_origin_id: "dbd70daa-356c-4bde-be27-a13c3c47d44d", Account_destination_id: "dbd74daa-356c-4bde-be27-a13c3c47d49d", Amount: 5000, Created_at: time.Now()})

	//call router handlers(the functions) and establish endpoints(the route)
	r.HandleFunc("/accounts", getAccounts).Methods("GET")
	r.HandleFunc("/accounts/{id}", getAccount).Methods("GET")
	r.HandleFunc("/accounts", createAccount).Methods("POST")
	r.HandleFunc("/accounts/{id}", updateAccount).Methods("PUT")
	//r.HandleFunc("/accounts/{id}", deleteAccount).Methods("DELETE")

	//transfer
	r.HandleFunc("/transfers", getTransfers).Methods("GET")
	r.HandleFunc("/transfers", createTransfer).Methods("POST")

	//Wrapped with log, throws error in case it fails					receives a port and a router
	log.Fatal(http.ListenAndServe(":8000", r))
}
