package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

func main() {
	connStr := "user=postgres password=admin dbname=final_project sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	files := []string{"users.csv", "products.csv", "addresses.csv", "audit_logs.csv", "cache.csv", "cart_items.csv", "categories.csv", "order_items.csv", "orders.csv", "payments.csv", "product_images.csv", "reviews.csv", "roles.csv", "sessions.csv", "shopping_cart.csv"}
	for _, file := range files {
		createTableFromCSV(db, file)
	}
}

func createTableFromCSV(db *sql.DB, filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	headers, err := reader.Read()
	if err != nil {
		panic(err)
	}

	tableName := strings.TrimSuffix(filePath, ".csv")

	for i, header := range headers {
		headers[i] = strings.TrimSpace(strings.ToLower(header))
	}

	columnDefinitions := make([]string, len(headers))
	for i, header := range headers {
		columnDefinitions[i] = fmt.Sprintf(`"%s" TEXT`, header) 
	}
	createTableQuery := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS "%s" (%s);`, tableName, strings.Join(columnDefinitions, ", "))

	_, err = db.Exec(createTableQuery)
	if err != nil {
		panic(fmt.Errorf("FAIL CREATING TABLES: %w", err))
	}

	placeholders := make([]string, len(headers))
	for i := range placeholders {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}
	insertQuery := fmt.Sprintf(`INSERT INTO "%s" VALUES (%s);`, tableName, strings.Join(placeholders, ", "))

	for {
		record, err := reader.Read()
		if err != nil {
			break
		}

		values := make([]any, len(record))
		for i, v := range record {
			values[i] = v
		}

		_, err = db.Exec(insertQuery, values...)
		if err != nil {
			panic(fmt.Errorf("ERROR INSERTING DATA: %w", err))
		}
	}
}
