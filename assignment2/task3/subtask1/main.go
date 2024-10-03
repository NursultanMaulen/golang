package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

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

func dropTable(){
	dropTableSQL := `DROP TABLE IF EXISTS users;`
	_, err := db.Exec(dropTableSQL)
	if err != nil {
		log.Fatalf("FAIL DELETING TABLE: %v", err)
	}

	fmt.Println("SUCCESS DROPPING TABLE.")
}

func createTable() {
	createTableSQL := `
	CREATE TABLE users (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL UNIQUE,
		age INT NOT NULL
	);
	`

	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("FAIL CREATING TABLE: %v", err)
	}

	fmt.Println("SUCCESS TABLE CREATE!")
}

func insertTransaction(users []map[string]interface{}){
	tx, err := db.Begin()

	if err != nil {
		log.Fatalf("FAIL STARTING TRANSACTION: %v", err)
	}

	for _, user := range users {
		_, err := tx.Exec("INSERT INTO users (name, age) VALUES ($1, $2)", user["name"], user["age"])
		if err != nil {
			tx.Rollback()
			log.Fatalf("FAIL INSERT TRANSACTION: %v", err)
			return
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Fatalf("FAIL COMMIT TRANSACTION: %v", err)
	}

	fmt.Println("SUCCESSFUL TRANSACTION")
}

func queryUsers(ageFilter *int, page, pageSize int) {

}

func main() {
	connectSQL()
	// dropTable()
	// createTable()
	// insertTransaction([]map[string]interface{}{
	// 	{"name":"Clarey","age":41},
	// 	{"name":"Shurlock","age":121},
	// 	{"name":"Amanda","age":35},
	// 	{"name":"Irita","age":7},
	// 	{"name":"Billy","age":96},
	// 	{"name":"Lorine","age":5},
	// 	{"name":"Goran","age":85},
	// 	{"name":"Brigham","age":105},
	// 	{"name":"Nixie","age":34},
	// 	{"name":"Walsh","age":49},
	// })

	// var ageFilter *int = &25
	// page := 1
	// pageSize := 10

	// queryUsers(ageFilter, page, pageSize)

	defer db.Close()
}