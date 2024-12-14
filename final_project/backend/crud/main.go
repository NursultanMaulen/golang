package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func connectToDB() (*sql.DB, error) {
	connStr := "user=postgres password=admin dbname=final_project sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к базе данных: %w", err)
	}
	return db, nil
}


func createUser(db *sql.DB, username, passwordHash, email, role string) error {
	query := `
        INSERT INTO users (username, password_hash, email, role) 
        VALUES ($1, $2, $3, $4)`
	_, err := db.Exec(query, username, passwordHash, email, role)
	if err != nil {
		return fmt.Errorf("ошибка добавления пользователя: %w", err)
	}
	return nil
}

func getUsers(db *sql.DB, limit int) ([]map[string]interface{}, error) {
	query := `SELECT user_id, username, email, role FROM users LIMIT $1`
	rows, err := db.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("ERROR FETCHING USERS: %w", err)
	}
	defer rows.Close()

	users := []map[string]interface{}{}
	for rows.Next() {
		var id sql.NullInt64
		var username, email, role sql.NullString
		var createdAt sql.NullTime

		err := rows.Scan(&id, &username, &email, &role)
		if err != nil {
			return nil, fmt.Errorf("ERROR SCANNING ROW: %w", err)
		}

		user := map[string]interface{}{
			"user_id":    nullToValue(id),
			"username":   nullToValue(username),
			"email":      nullToValue(email),
			"created_at": nullToValue(createdAt),
			"role":       nullToValue(role),
		}
		users = append(users, user)
	}

	return users, nil
}

func nullToValue(v interface{}) interface{} {
	switch value := v.(type) {
	case sql.NullInt64:
		if value.Valid {
			return value.Int64
		}
	case sql.NullString:
		if value.Valid {
			return value.String
		}
	case sql.NullTime:
		if value.Valid {
			return value.Time
		}
	}
	return nil
}


func updateUser(db *sql.DB, userID int, email string, role string) error {
	query := `UPDATE users SET email = $1, role = $2 WHERE user_id = $3`
	_, err := db.Exec(query, email, role, userID)
	if err != nil {
		return fmt.Errorf("ошибка обновления пользователя: %w", err)
	}
	return nil
}

func deleteUser(db *sql.DB, userID int) error {
	query := `DELETE FROM users WHERE user_id = $1`
	_, err := db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("ошибка удаления пользователя: %w", err)
	}
	return nil
}

func getTotalRevenue(db *sql.DB, startDate, endDate string) (float64, error) {
	var totalRevenue float64
	query := `
		SELECT COALESCE(SUM(order_items.quantity::NUMERIC * order_items.price::NUMERIC), 0) AS total_revenue
		FROM orders
		JOIN order_items ON orders.order_id = order_items.order_id
		WHERE orders.order_date BETWEEN $1 AND $2;
	`
	err := db.QueryRow(query, startDate, endDate).Scan(&totalRevenue)
	if err != nil {
		return 0, fmt.Errorf("ERROR FETCHING TOTAL REVENUE: %w", err)
	}
	return totalRevenue, nil
}



func main() {
	db, err := connectToDB()
	if err != nil {
		log.Fatalf("Ошибка подключения: %v", err)
	}
	defer db.Close()

	// err = createUser(db, "nurs", "hashed_password", "nurs@example.com", "admin")
	// if err != nil {
	// 	log.Printf("Ошибка добавления пользователя: %v", err)
	// } else {
	// 	log.Println("Пользователь успешно добавлен")
	// }

	// users, err := getUsers(db, 2)
	// if err != nil {
	// 	log.Printf("Ошибка получения пользователей: %v", err)
	// } else {
	// 	log.Println("Пользователи:", users)
	// }

	// err = updateUser(db, 100, "new_email1111@example.com", "user")
	// if err != nil {
	// 	log.Printf("Ошибка обновления пользователя: %v", err)
	// } else {
	// 	log.Println("Пользователь успешно обновлен")
	// }

	// err = deleteUser(db, 101)
	// if err != nil {
	// 	log.Printf("Ошибка удаления пользователя: %v", err)
	// } else {
	// 	log.Println("Пользователь успешно удалён")
	// }

	startDate := "2022-12-01"
	endDate := "2023-12-31"

	totalRevenue, err := getTotalRevenue(db, startDate, endDate)
	if err != nil {
		log.Fatalf("ERROR FETCHING TOTAL REVENUE: %v", err)
	}

	log.Printf("TOTAL REVENUE BETWEEN %s AND %s: $%.2f", startDate, endDate, totalRevenue)
}


