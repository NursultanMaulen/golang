package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var db *sql.DB
var rdb *redis.Client

func connectRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr: "redis-container:6379",
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Error connecting to Redis: %v", err)
	}
	log.Println("Connected to Redis!")
}

func connectToDB() {
	var err error
	connStr := "host=host.docker.internal user=postgres password=admin dbname=final_project sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("ERROR CONNECTING TO DB: %v", err)
	}
}

type Product struct {
	ID          int     `json:"product_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	CategoryID  int     `json:"category_id"`
}

type User struct {
	ID       int    `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	PasswordHash string `json:"password_hash"`
	Role     string `json:"role"`
}

func getProducts(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT product_id, name, description, price, stock, category_id FROM products")
	if err != nil {
		http.Error(w, fmt.Sprintf("ERROR FETCHING PRODUCTS: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Stock, &p.CategoryID); err != nil {
			http.Error(w, fmt.Sprintf("ERROR SCANNING PRODUCT: %v", err), http.StatusInternalServerError)
			return
		}
		products = append(products, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func addProduct(w http.ResponseWriter, r *http.Request) {
	var p Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "INVALID REQUEST PAYLOAD", http.StatusBadRequest)
		return
	}

	query := "INSERT INTO products (name, description, price, stock, category_id) VALUES ($1, $2, $3, $4, $5) RETURNING product_id"
	if err := db.QueryRow(query, p.Name, p.Description, p.Price, p.Stock, p.CategoryID).Scan(&p.ID); err != nil {
		http.Error(w, fmt.Sprintf("ERROR INSERTING PRODUCT: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT user_id, username, email, role FROM users")
	if err != nil {
		http.Error(w, fmt.Sprintf("ERROR FETCHING USERS: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.Role); err != nil {
			http.Error(w, fmt.Sprintf("ERROR SCANNING USER: %v", err), http.StatusInternalServerError)
			return
		}
		users = append(users, u)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func addUser(w http.ResponseWriter, r *http.Request) {
	var newUser User
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		http.Error(w, "INVALID REQUEST PAYLOAD", http.StatusBadRequest)
		return
	}

	if newUser.Role == "" {
		newUser.Role = "user"
	}

	query := "INSERT INTO users (username, email, password_hash, role) VALUES ($1, $2, $3, $4) RETURNING user_id"
	err := db.QueryRow(query, newUser.Username, newUser.Email, newUser.PasswordHash, newUser.Role).Scan(&newUser.ID)
	if err != nil {
		http.Error(w, fmt.Sprintf("ERROR INSERTING USER: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}

func getUserByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var user User
	query := "SELECT user_id, username, email, password_hash, role FROM users WHERE user_id = $1"
	err = db.QueryRow(query, userID).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("Error fetching user: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func updateUserByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var updatedUser User
	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	query := "UPDATE users SET username = $1, email = $2, password_hash = $3 WHERE user_id = $4"
	_, err = db.Exec(query, updatedUser.Username, updatedUser.Email, updatedUser.PasswordHash, userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error updating user: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User updated successfully"})
}

func deleteUserByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	query := "DELETE FROM users WHERE user_id = $1"
	_, err = db.Exec(query, userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error deleting user: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User deleted successfully"})
}

func authenticateUser(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "INVALID REQUEST PAYLOAD", http.StatusBadRequest)
		return
	}

	var user User
	query := "SELECT user_id, username, email, password_hash, role FROM users WHERE email = $1"
	err := db.QueryRow(query, credentials.Email).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "INVALID EMAIL OR PASSWORD", http.StatusUnauthorized)
			return
		}
		http.Error(w, fmt.Sprintf("ERROR FETCHING USER: %v", err), http.StatusInternalServerError)
		return
	}

	if user.PasswordHash != credentials.Password {
		http.Error(w, "INVALID EMAIL OR PASSWORD", http.StatusUnauthorized)
		return
	}

	token := "mock-token"
	response := map[string]interface{}{
		"token": token,
		"user": map[string]interface{}{
			"user_id":  user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func addToCart(w http.ResponseWriter, r *http.Request) {
    var cartItem struct {
        UserID    int `json:"user_id"`
        ProductID int `json:"product_id"`
        Quantity  int `json:"quantity"`
    }

    log.Println("Starting addToCart handler")

    // Декодирование JSON запроса
    if err := json.NewDecoder(r.Body).Decode(&cartItem); err != nil {
        log.Printf("ERROR DECODING REQUEST BODY: %v", err)
        http.Error(w, "INVALID REQUEST PAYLOAD", http.StatusBadRequest)
        return
    }
    log.Printf("Decoded cart item: %+v", cartItem)

    // Проверка количества
    if cartItem.Quantity <= 0 {
        log.Printf("INVALID QUANTITY: %d", cartItem.Quantity)
        http.Error(w, "INVALID QUANTITY", http.StatusBadRequest)
        return
    }

	// Check available stock for the product
	var stock int
	err := db.QueryRow(`SELECT stock FROM products WHERE product_id = $1`, cartItem.ProductID).Scan(&stock)
	if err != nil {
		log.Printf("ERROR FETCHING PRODUCT STOCK: %v", err)
		http.Error(w, "ERROR FETCHING PRODUCT STOCK", http.StatusInternalServerError)
		return
	}

	if cartItem.Quantity > stock {
		log.Printf("REQUESTED QUANTITY EXCEEDS STOCK: %d > %d", cartItem.Quantity, stock)
		http.Error(w, fmt.Sprintf("ONLY %d ITEMS IN STOCK", stock), http.StatusBadRequest)
		return
	}

    var cartID sql.NullInt64

	query := `
		WITH existing_cart AS (
			SELECT cart_id FROM shopping_cart WHERE user_id = $1 AND status = 'open' LIMIT 1
		),
		new_cart AS (
			INSERT INTO shopping_cart (user_id, status, created_at) 
			SELECT $1, 'open', CURRENT_TIMESTAMP
			WHERE NOT EXISTS (SELECT 1 FROM existing_cart)
			RETURNING cart_id
		)
		SELECT COALESCE(
			(SELECT cart_id FROM existing_cart),
			(SELECT cart_id FROM new_cart)
		) AS cart_id;
	`

	log.Println("Executing query to find or create cart...")
	err = db.QueryRow(query, cartItem.UserID).Scan(&cartID)
	if err != nil {
		log.Printf("ERROR FINDING OR CREATING CART: %v", err)
		http.Error(w, fmt.Sprintf("ERROR FINDING OR CREATING CART: %v", err), http.StatusInternalServerError)
		return
	}

	if !cartID.Valid {
		log.Println("UNABLE TO CREATE OR FIND CART")
		http.Error(w, "UNABLE TO CREATE OR FIND CART", http.StatusInternalServerError)
		return
	}

	log.Printf("Cart ID: %d", cartID.Int64)

    // SQL-запрос для добавления продукта в корзину
    _, err = db.Exec(`
        INSERT INTO cart_items (cart_id, product_id, quantity, updated_at) 
        VALUES ($1, $2, $3, CURRENT_TIMESTAMP)
        ON CONFLICT (cart_id, product_id) 
        DO UPDATE SET quantity = cart_items.quantity + EXCLUDED.quantity, updated_at = CURRENT_TIMESTAMP;
    `, cartID.Int64, cartItem.ProductID, cartItem.Quantity)

    if err != nil {
        log.Printf("ERROR ADDING ITEM TO CART: %v", err)
        http.Error(w, fmt.Sprintf("ERROR ADDING ITEM TO CART: %v", err), http.StatusInternalServerError)
        return
    }

    log.Printf("Product added to cart: cart_id=%d, product_id=%d, quantity=%d", cartID.Int64, cartItem.ProductID, cartItem.Quantity)

    // Успешный ответ
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "message": "Product added to cart successfully",
        "cart_id": cartID.Int64,
    })
}

func getCartItems(w http.ResponseWriter, r *http.Request) {
    userID := r.URL.Query().Get("user_id")
    if userID == "" {
        http.Error(w, "USER ID IS REQUIRED", http.StatusBadRequest)
        return
    }

    rows, err := db.Query(`
       SELECT 
		ci.product_id, 
		p.name, 
		p.description, 
		p.price, 
		ci.quantity 
	FROM 
		cart_items ci
	JOIN 
		products p ON ci.product_id = p.product_id
	JOIN 
		shopping_cart sc ON ci.cart_id = sc.cart_id
	WHERE 
		sc.user_id = $1 AND sc.status = 'open';
    `, userID)

    if err != nil {
        http.Error(w, fmt.Sprintf("ERROR FETCHING CART ITEMS: %v", err), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var cartItems []map[string]interface{}
    for rows.Next() {
        var item struct {
            ProductID   int     `json:"product_id"`
            Name        string  `json:"name"`
            Description string  `json:"description"`
            Price       float64 `json:"price"`
            Quantity    int     `json:"quantity"`
        }

        if err := rows.Scan(&item.ProductID, &item.Name, &item.Description, &item.Price, &item.Quantity); err != nil {
            http.Error(w, fmt.Sprintf("ERROR SCANNING ROW: %v", err), http.StatusInternalServerError)
            return
        }

        cartItems = append(cartItems, map[string]interface{}{
            "product_id": item.ProductID,
            "name":       item.Name,
            "description": item.Description,
            "price":      item.Price,
            "quantity":   item.Quantity,
        })
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(cartItems)
}


func removeFromCart(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID, err := strconv.Atoi(vars["productId"])
	if err != nil {
		http.Error(w, "INVALID PRODUCT ID", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil || userID <= 0 {
		http.Error(w, "INVALID USER ID", http.StatusBadRequest)
		return
	}

	_, err = db.Exec(`
		DELETE FROM cart_items
		WHERE product_id = $1 
		  AND cart_id = (
			SELECT cart_id 
			FROM shopping_cart 
			WHERE user_id = $2 AND status = 'open'
		)
	`, productID, userID)

	if err != nil {
		http.Error(w, fmt.Sprintf("ERROR REMOVING PRODUCT FROM CART: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func placeOrder(w http.ResponseWriter, r *http.Request) {
    log.Println("Starting placeOrder handler")

    var orderRequest struct {
        UserID int `json:"user_id"`
        Items  []struct {
            ProductID int `json:"product_id"`
            Quantity  int `json:"quantity"`
            Price     float64 `json:"price"` // Ensure price is included in the payload
        } `json:"items"`
    }

    // Decode the JSON payload
    if err := json.NewDecoder(r.Body).Decode(&orderRequest); err != nil {
        log.Printf("Error decoding order request: %v", err)
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    log.Printf("Decoded order request: %+v", orderRequest)

    // Validate the request
    if orderRequest.UserID == 0 || len(orderRequest.Items) == 0 {
        log.Println("Invalid order request: Missing user_id or items")
        http.Error(w, "Invalid order request: Missing user_id or items", http.StatusBadRequest)
        return
    }

    // Calculate the total amount
    totalAmount := 0.0
    for _, item := range orderRequest.Items {
        totalAmount += float64(item.Quantity) * item.Price
    }
    log.Printf("Total order amount: %.2f", totalAmount)

    // Insert the order into the `orders` table
    var orderID string
	err := db.QueryRow(`
		INSERT INTO orders (user_id, order_date, status, total_amount)
		VALUES ($1, CURRENT_TIMESTAMP, 'pending', $2)
		RETURNING order_id
	`, orderRequest.UserID, totalAmount).Scan(&orderID)
	if err != nil {
		log.Printf("Error inserting order: %v", err)
		http.Error(w, fmt.Sprintf("Error inserting order: %v", err), http.StatusInternalServerError)
		return
	}


    log.Printf("Created order with ID: %s", orderID)

    // Insert items into the `order_items` table
    for _, item := range orderRequest.Items {
        _, err = db.Exec(`
            INSERT INTO order_items (order_id, product_id, quantity)
            VALUES ($1, $2, $3)
        `, orderID, item.ProductID, item.Quantity)

        if err != nil {
            log.Printf("Error inserting order item (product_id: %d): %v", item.ProductID, err)
            http.Error(w, "Error creating order items", http.StatusInternalServerError)
            return
        }
    }

    log.Printf("Order items added for order ID: %s", orderID)

    // Respond with the created order ID
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "message":  "Order placed successfully",
        "order_id": orderID,
        "total_amount": totalAmount,
    })
}

func updateProductQuantity(w http.ResponseWriter, r *http.Request) {
    log.Println("Starting updateProductQuantity handler")

    userRole := r.Header.Get("Role")
    log.Printf("Role Header: %s", userRole)

    if userRole != "admin" {
        http.Error(w, "FORBIDDEN: You don't have permission to update product quantities.", http.StatusForbidden)
        return
    }

    // Parse product ID
    vars := mux.Vars(r)
    log.Printf("Mux Vars: %+v", vars)

    productID, err := strconv.Atoi(vars["productID"])
    if err != nil {
        log.Printf("Invalid product ID: %v", err)
        http.Error(w, "INVALID PRODUCT ID", http.StatusBadRequest)
        return
    }

    // Decode request body
    var requestBody struct {
        Quantity int `json:"quantity"`
    }
    if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
        log.Printf("Error decoding request body: %v", err)
        http.Error(w, "INVALID REQUEST PAYLOAD", http.StatusBadRequest)
        return
    }

    log.Printf("Decoded Request Body: %+v", requestBody)

    if requestBody.Quantity <= 0 {
        log.Printf("Invalid quantity: %d", requestBody.Quantity)
        http.Error(w, "INVALID QUANTITY", http.StatusBadRequest)
        return
    }

    query := `UPDATE products SET stock = $1 WHERE product_id = $2`
    log.Printf("Executing query: %s", query)

    _, err = db.Exec(query, requestBody.Quantity, productID)
    if err != nil {
        log.Printf("ERROR UPDATING PRODUCT QUANTITY: %v", err)
        http.Error(w, "ERROR UPDATING PRODUCT QUANTITY", http.StatusInternalServerError)
        return
    }

    log.Println("Product quantity updated successfully")

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Product quantity updated successfully"})
}

func getUserAndCart(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "USER ID IS REQUIRED", http.StatusBadRequest)
		return
	}

	ctx := context.Background()

	cacheKey := fmt.Sprintf("user_cart:%s", userID)
	cachedData, err := rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		log.Println("Cache hit!")
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cachedData))
		return
	}

	log.Println("Cache miss! Fetching data from DB...")

	userChan := make(chan map[string]interface{})
	cartChan := make(chan []map[string]interface{})

	go func() {
		var user struct {
			UserID   int    `json:"user_id"`
			Username string `json:"username"`
			Email    string `json:"email"`
		}

		err := db.QueryRow(`SELECT user_id, username, email FROM users WHERE user_id = $1`, userID).
			Scan(&user.UserID, &user.Username, &user.Email)
		if err != nil {
			userChan <- nil
			return
		}
		userChan <- map[string]interface{}{
			"user_id":  user.UserID,
			"username": user.Username,
			"email":    user.Email,
		}
	}()

	go func() {
		rows, err := db.Query(`SELECT product_id, quantity FROM cart_items WHERE cart_id = (SELECT cart_id FROM shopping_cart WHERE user_id = $1 AND status = 'open')`, userID)
		if err != nil {
			cartChan <- nil
			return
		}
		defer rows.Close()

		var cart []map[string]interface{}
		for rows.Next() {
			var item struct {
				ProductID int `json:"product_id"`
				Quantity  int `json:"quantity"`
			}
			if err := rows.Scan(&item.ProductID, &item.Quantity); err != nil {
				continue
			}
			cart = append(cart, map[string]interface{}{
				"product_id": item.ProductID,
				"quantity":   item.Quantity,
			})
		}
		cartChan <- cart
	}()

	user := <-userChan
	cart := <-cartChan

	if user == nil || cart == nil {
		http.Error(w, "ERROR FETCHING DATA", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"user": user,
		"cart": cart,
	}
	responseJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "ERROR MARSHALING RESPONSE", http.StatusInternalServerError)
		return
	}

	err = rdb.Set(ctx, cacheKey, responseJSON, 10*time.Minute).Err()
	if err != nil {
		log.Printf("Error caching data: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}




func main() {
	connectToDB()
	connectRedis()
	defer db.Close()

	r := mux.NewRouter()
	r.HandleFunc("/api/users/authenticate", authenticateUser).Methods("POST")
	r.HandleFunc("/api/products", getProducts).Methods("GET")
	r.HandleFunc("/api/products", addProduct).Methods("POST")
	r.HandleFunc("/api/products/{productID}", updateProductQuantity).Methods("PUT")
	r.HandleFunc("/api/users", getUsers).Methods("GET")
	r.HandleFunc("/api/users/{id}", getUserByID).Methods("GET")
	r.HandleFunc("/api/users/{id}", updateUserByID).Methods("PUT")
	r.HandleFunc("/api/users/{id}", deleteUserByID).Methods("DELETE")
	r.HandleFunc("/api/users", addUser).Methods("POST")
	r.HandleFunc("/api/cart/add", addToCart).Methods("POST")
	r.HandleFunc("/api/cart/items", getCartItems).Methods("GET")
	r.HandleFunc("/api/cart/{productId}", removeFromCart).Methods("DELETE")
	r.HandleFunc("/api/orders", placeOrder).Methods("POST")
	r.HandleFunc("/api/usersapi", getUserAndCart).Methods("GET")


	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:3000"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization", "Role"}),
	)

	log.Println("SERVER STARTED ON PORT 8080")
	log.Fatal(http.ListenAndServe(":8080", corsHandler(r)))
}
