package main

import (
	"context"

	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis"
)

const (
	host = "localhost"
	port = 5432
	user = "postgres"
	password = "admin"
	dbname = "usejoins"
)

var db *sql.DB
var err error

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

func createTables() error {
    userTable := `
    CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		email VARCHAR(100) NOT NULL UNIQUE
	);
	`

    orderTable := `
    CREATE TABLE IF NOT EXISTS orders (
    	id SERIAL PRIMARY KEY,
    	user_id INT NOT NULL,
    	product VARCHAR(100) NOT NULL,
    	quantity INT NOT NULL,
    	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);`

    _, err := db.Exec(userTable)
    if err != nil {
        log.Fatalf("Error creating users table: %v", err)
    }

    _, err = db.Exec(orderTable)
    if err != nil {
        log.Fatalf("Error creating orders table: %v", err)
    }

    fmt.Println("Tables created successfully!")
    return nil
}


type User struct {
    ID    int
    Name  string
    Email string
    Orders []Order 
}

type Order struct {
    ID       int
    UserID   int
    Product  string
    Quantity int
}

type UserWithOrders struct {
    User
    Orders []Order
}

type userRepository struct {
    db *sql.DB // assuming you're using the database/sql package
}

type UserRepository interface {
    CreateUser(ctx context.Context, name string, email string) (int, error)
    CreateOrder(ctx context.Context, userID int, product string, quantity int) (int, error)
    GetUserWithOrders(ctx context.Context, userID int) (*User, error)
}

// createUser inserts a new user into the users table
func createUser(ctx context.Context, db *sql.DB, name, email string) (int, error) {
    query := `INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id`
    var userID int
    err := db.QueryRowContext(ctx, query, name, email).Scan(&userID)
    if err != nil {
        return 0, err
    }
    return userID, nil
}

// createOrder inserts a new order into the orders table for a user
func createOrder(ctx context.Context, db *sql.DB, userID int, product string, quantity int) (int, error) {
    query := `INSERT INTO orders (user_id, product, quantity) VALUES ($1, $2, $3) RETURNING id`
    var orderID int
    err := db.QueryRowContext(ctx, query, userID, product, quantity).Scan(&orderID)
    if err != nil {
        return 0, err
    }
    return orderID, nil
}

// getUserWithOrders retrieves a user and their orders using a JOIN query
func getUserWithOrders(ctx context.Context, db *sql.DB, userID int) (*User, error) {
    // Check if the result is in the cache
    cacheKey := fmt.Sprintf("user:%d", userID)
    cachedUser, err := redisClient.Get(ctx, cacheKey).Result()
    if err == nil {
        // If found in cache, unmarshal it and return
        var user User
        // Unmarshal the cached user data (assuming you stored it as a JSON string)
        if err := json.Unmarshal([]byte(cachedUser), &user); err != nil {
            return nil, err
        }
        return &user, nil // Return cached user
    } else if err != redis.Nil {
        return nil, err // Return error if it's not a cache miss
    }

    // Cache miss, so query the database
    query := `
    SELECT u.id, u.name, u.email, o.id, o.product, o.quantity
    FROM users u
    LEFT JOIN orders o ON u.id = o.user_id
    WHERE u.id = $1`

    rows, err := db.QueryContext(ctx, query, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var user User
    orders := []Order{}

    // Iterate through the result rows
    for rows.Next() {
        var order Order
        if err := rows.Scan(&user.ID, &user.Name, &user.Email, &order.ID, &order.Product, &order.Quantity); err != nil {
            return nil, err
        }
        if order.ID != 0 { // Only append if the order exists
            orders = append(orders, order)
        }
    }

    user.Orders = orders

    // Cache the result in Redis for future requests
    userData, err := json.Marshal(user)
    if err != nil {
        return nil, err
    }
    err = redisClient.Set(ctx, cacheKey, userData, 10*time.Minute).Err() // Set cache with 10 minutes expiration
    if err != nil {
        return nil, err
    }

    return &user, nil
}



var redisClient *redis.Client

func connectRedis() {
	redisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6380", // Redis server address
	})

	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("FAIL CONNECTING TO REDIS: %v", err)
	}

	fmt.Println("SUCCESS: Connected to Redis")
}

func main(){
	fmt.Println("123")
	connectSQL()
    connectRedis()
	// createTables()

	// Use a context with a timeout
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Step 1: Insert a mock user
    userID, err := createUser(ctx, db, "John Doe", "john@example.com")
    if err != nil {
        log.Fatal("Error creating user:", err)
    }

    // Step 2: Insert a mock order for the user
    orderID, err := createOrder(ctx, db, userID, "Smartphone", 2)
    if err != nil {
        log.Fatal("Error creating order:", err)
    }
    fmt.Printf("User created with ID: %d, Order created with ID: %d\n", userID, orderID)

    // Step 3: Get user with orders using a JOIN query
    user, err := getUserWithOrders(ctx, db, userID)
    if err != nil {
        log.Fatal("Error getting user with orders:", err)
    }

    fmt.Printf("User: %+v\n", user)
    for _, order := range user.Orders {
        fmt.Printf("Order: %+v\n", order)
    }
}