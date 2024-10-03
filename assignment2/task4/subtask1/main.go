package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var db *sql.DB;
var err error;

const (
	host = "localhost"
	port = 5432
	user = "postgres"
	password = "admin"
	dbname = "assignment2"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func connectSQL(){
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
        host, port, user, password, dbname)
	db, err = sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Fatalf("FAIL CONNECTION DB: %v", err)
	}

	// POOL PARAMS
	db.SetMaxOpenConns(10) // MAX 10 OPEN CONNS
	db.SetMaxIdleConns(5)  // MAX 5 NON ACTIVE CONNS
	db.SetConnMaxLifetime(30 * time.Minute)

	err = db.Ping()
	if err != nil {
		log.Fatalf("FAIL PINGING DB: %v", err)
	}

	fmt.Println("SUCCESS")
}


func getUsers(w http.ResponseWriter, r *http.Request) {
	query := `SELECT id, name, age FROM users`
	var args []interface{}
	var whereClauses []string

	ageFilter := r.URL.Query().Get("age")
	sortBy := r.URL.Query().Get("sort")

	if ageFilter != "" {
		whereClauses = append(whereClauses, "age = $1")
		args = append(args, ageFilter)
	}

	if len(whereClauses) > 0 {
		query += " WHERE " + whereClauses[0]
	}

	if sortBy == "name" {
		query += " ORDER BY name"
	} else {
		query += " ORDER BY id"
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		http.Error(w, fmt.Sprintf("ERROR QUERYING USERS: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Age)
		if err != nil {
			http.Error(w, fmt.Sprintf("ERROR SCANNING ROW: %v", err), http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func isNameUnique(name string) (bool, error) {
	query := `SELECT COUNT(*) FROM users WHERE name = $1`
	var count int
	err := db.QueryRow(query, name).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("ERROR CHECKING UNIQUENESS: %v", err)
	}
	return count == 0, nil
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "INVALID INPUT", http.StatusBadRequest)
		return
	}

	unique, err := isNameUnique(user.Name)
	if err != nil {
		http.Error(w, "FAIL CHECK UNIQUENESS", http.StatusInternalServerError)
		return
	}
	if !unique {
		http.Error(w, "NAME IS ALREADY USED", http.StatusConflict)
		return
	}

	query := `INSERT INTO users (name, age) VALUES ($1, $2) RETURNING id`
	err = db.QueryRow(query, user.Name, user.Age).Scan(&user.ID)
	if err != nil {
		http.Error(w, fmt.Sprintf("ERROR CREATING USER: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "INVALID USER ID", http.StatusBadRequest)
		return
	}

	var user User
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "INVALID INPUT", http.StatusBadRequest)
		return
	}

	unique, err := isNameUnique(user.Name)
	if err != nil {
		http.Error(w, "FAIL CHECK UNIQUENESS", http.StatusInternalServerError)
		return
	}
	if !unique {
		http.Error(w, "NAME IS ALREADY USED", http.StatusConflict)
		return
	}

	query := `UPDATE users SET name = $1, age = $2 WHERE id = $3`
	_, err = db.Exec(query, user.Name, user.Age, id)
	if err != nil {
		http.Error(w, fmt.Sprintf("ERROR UPDATING USER: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "INVALID USER ID", http.StatusBadRequest)
		return
	}

	query := `DELETE FROM users WHERE id = $1`
	result, err := db.Exec(query, id)
	if err != nil {
		http.Error(w, fmt.Sprintf("ERROR DELETING USER: %v", err), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, fmt.Sprintf("ERROR CHECKING DELETION: %v", err), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "USER NOT FOUND", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func main() {
	connectSQL()
	defer db.Close()

	r := mux.NewRouter()
	// r.HandleFunc("/users", getUsers).Methods("GET")
	// r.HandleFunc("/users", createUser).Methods("POST")
	// r.HandleFunc("/users/{id}", updateUser).Methods("PUT")
	r.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", r))
}