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

func queryUsers(ageFilter int, page, pageSize int) {
	query := `SELECT id, name, age FROM users`
	var args []interface{}

	query += ` WHERE age <= $1`
	args = append(args, ageFilter)

	query += ` ORDER BY id LIMIT $2 OFFSET $3`
	args = append(args, pageSize, (page-1)*pageSize)

	rows, err := db.Query(query, args...)
	if err != nil {
		log.Fatalf("FAIL QUERY REQUEST: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		var age int
		err := rows.Scan(&id, &name, &age)
		if err != nil {
			log.Fatalf("FAIL SCAN ROW: %v", err)
		}
		fmt.Printf("ID: %d, Name: %s, Age: %d\n", id, name, age)
	}

	if err = rows.Err(); err != nil {
		log.Fatalf("FAIL ITERATING: %v", err)
	}
}

func deleteUser(id int) error {
	query := `DELETE FROM users WHERE id = $1`
	result, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("ERROR DELETING USER: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ERROR ROWS AFFECTED: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("USER ID %d NOT FOUND", id)
	}

	fmt.Printf("USER ID %d DELETED SUCCESSFULLY\n", id)
	return nil
}

func main() {
	connectSQL()
	// dropTable()
	// createTable()
	// insertTransaction([]map[string]interface{}{
	// 	{"name":"Clar","age":41},
	// 	{"name":"Walsh","age":49},
	// })

	// ageFilter := 25
	// page := 1
	// pageSize := 10

	// queryUsers(ageFilter, page, pageSize)

	// deleteUser(1)

	defer db.Close()
}