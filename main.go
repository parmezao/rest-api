package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Person struct {
	ID        int    `json:"id,omitempty"`
	FirstName string `json:"firstname,omitempty"`
	LastName  string `json:"lastname,omitempty"`
	Address   string `json:"address,omitempty"`
}

func GetPeople(w http.ResponseWriter, r *http.Request) {
	db := GetDatabase()
	defer db.Close()

	rows, _ := db.Query("SELECT Id, Firstname, Lastname, Address FROM Person")
	defer rows.Close()

	var people []Person
	for rows.Next() {
		var person Person
		rows.Scan(&person.ID, &person.FirstName, &person.LastName, &person.Address)
		people = append(people, person)
	}

	json, _ := json.Marshal(people)

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(json))
}

func GetPerson(w http.ResponseWriter, r *http.Request) {
	db := GetDatabase()
	defer db.Close()

	params := mux.Vars(r)
	id := params["id"]

	var person Person
	db.QueryRow("SELECT Id, Firstname, Lastname, Address FROM Person WHERE Id = ?", id).
		Scan(&person.ID, &person.FirstName, &person.LastName, &person.Address)

	json, _ := json.Marshal(person)

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(json))
}

func CreatePerson(w http.ResponseWriter, r *http.Request) {
	db := GetDatabase()
	defer db.Close()

	var person Person
	_ = json.NewDecoder(r.Body).Decode(&person)

	query := "INSERT INTO Person(Firstname, Lastname, Address) VALUES (?, ?, ?)"
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, person.FirstName, person.LastName, person.Address)
	if err != nil {
		log.Fatal(err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%d person(s) created ", rows)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(person)
}

func DeletePerson(w http.ResponseWriter, r *http.Request) {
	db := GetDatabase()
	defer db.Close()

	params := mux.Vars(r)
	id := params["id"]

	db.Query("DELETE FROM Person WHERE Id = ?", id)
	GetPeople(w, r)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/contato", GetPeople).Methods("GET")
	router.HandleFunc("/contato/{id}", GetPerson).Methods("GET")
	router.HandleFunc("/contato", CreatePerson).Methods("POST")
	router.HandleFunc("/contato/{id}", DeletePerson).Methods("DELETE")

	fmt.Println("Server running in port 8000...")
	log.Fatal(http.ListenAndServe(":8000", router))
}
