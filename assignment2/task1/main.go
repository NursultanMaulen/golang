package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	host = "localhost"
	port = 5432
	user = "postgres"
	password = "admin"
	dbname = "assignment2"
)

var db *sql.DB;
var err error;

func connectDB(){
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
        host, port, user, password, dbname)
	// dsn := "postgres://postgres:admin@localhost:5432/assignment2"
	db, err = sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Не удалось проверить подключение к базе данных: %v", err)
	}

	fmt.Println("SUCCESS")
}

func createTable() {
	createTableSQL := `
	CREATE TABLE users (
		id SERIAL PRIMARY KEY NOT NULL,
		name TEXT NOT NULL,
		age INT NOT NULL
	);
	`

	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Ошибка при создании таблицы: %v", err)
	}

	fmt.Println("Таблица успешно создана!")
}

func insertData(name string, age int) {
	insertUserSQL := `INSERT INTO users (name, age) VALUES ($1, $2) returning id;` 

	userID := 0

	err = db.QueryRow(insertUserSQL, name, age).Scan(&userID)

	if err != nil {
		log.Fatalf("Ошибка при выполнении запроса: %v", err)
	}


	fmt.Println("USER ADDED SUCCESSFULLY")
}

func printUsers(){
	rows, err := db.Query(`SELECT * from users;`)

	if err != nil {
		log.Fatalf("Ошибка при выполнении запроса: %v", err)
	}
	
	defer rows.Close()

	fmt.Println("SUCCESSFULLY READ ALL USERS: ")

	for rows.Next() {
		var id int
		var name string
		var age int

		err = rows.Scan(&id, &name, &age)
		if err != nil {
			log.Fatalf("Ошибка при выполнении запроса: %v", err)
		}

		fmt.Println("ID:", id, "NAME:", name, "AGE:", age)
	}

}

func main() {
	connectDB()
	
	// createTable()

	// insertData("nursultan2", 21)

	printUsers()

	defer db.Close()
}