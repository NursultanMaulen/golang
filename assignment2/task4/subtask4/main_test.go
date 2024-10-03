package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type User struct {
	ID     uint    `gorm:"primaryKey"`
	Name   string  `gorm:"size:255"`
	Age    int
}


var db *gorm.DB
var err error

func connectGorm() {
	dsn := "host=localhost user=postgres password=admin dbname=assignment2 port=5432 sslmode=disable TimeZone=Europe/Moscow"

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("FAIL CONNECTIONG GORM: %v", err)
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	fmt.Println("SUCCESS CONNECTION")
}

func createModel() {
	db.AutoMigrate(&User{})
}

func getUsersHandler(w http.ResponseWriter, r *http.Request) {
	var users []User

	if db == nil {
        http.Error(w, "Database not initialized", http.StatusInternalServerError)
        return
    }

	page := r.URL.Query().Get("page")
	pageSize := r.URL.Query().Get("pageSize")

	if page == "" {
		page = "1"
	}
	if pageSize == "" {
		pageSize = "10"
	}

	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt < 1 {
		http.Error(w, "Invalid page number", http.StatusBadRequest)
		return
	}

	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil || pageSizeInt < 1 {
		http.Error(w, "Invalid page size", http.StatusBadRequest)
		return
	}

	offset := (pageInt - 1) * pageSizeInt

	if err := db.Limit(pageSizeInt).Offset(offset).Find(&users).Error; err != nil {
        http.Error(w, "ERROR GETTING USER: "+err.Error(), http.StatusInternalServerError)
        return
    }

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	var user User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "INCORRECT REQUEST", http.StatusBadRequest)
		return
	}

	if err := db.Create(&user).Error; err != nil {
		http.Error(w, "ERROR CREATING USER", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func updateUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	var user User
	if err := db.First(&user, userID).Error; err != nil {
		http.Error(w, "USER NOT FOUND", http.StatusNotFound)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "INCORRECT REQUEST", http.StatusBadRequest)
		return
	}

	if err := db.Save(&user).Error; err != nil {
		http.Error(w, "ERROR UPDATING USER", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	if err := db.Delete(&User{}, userID).Error; err != nil {
		http.Error(w, "ERROR DELETING USER", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func setupTestDB() {
	var err error
	dsn := "host=localhost user=postgres password=admin dbname=assignment2 port=5432 sslmode=disable"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	db.AutoMigrate(&User{})
}

func TestGetUsers(t *testing.T) {
	setupTestDB() 
	router := mux.NewRouter()
	router.HandleFunc("/users", getUsersHandler).Methods("GET")

	db.Create(&User{Name: "Alice", Age: 30})
	db.Create(&User{Name: "Bob", Age: 25})

	req, err := http.NewRequest("GET", "/users", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var users []User
	if err := json.NewDecoder(rr.Body).Decode(&users); err != nil {
		t.Fatalf("Could not decode response: %v", err)
	}

	expectedUsers := map[string]int{
		"Alice": 30,
		"Bob":   25,
	}

	for _, user := range users {
		if age, ok := expectedUsers[user.Name]; ok {
			if user.Age != age {
				t.Errorf("User %s has age %d, want %d", user.Name, user.Age, age)
			}
			delete(expectedUsers, user.Name)
		} else {
			t.Errorf("Unexpected user: %s", user.Name)
		}
	}

	db.Delete(&User{}, "name = ?", "Alice")
	db.Delete(&User{}, "name = ?", "Bob")
}

func TestCreateUser(t *testing.T) {
	setupTestDB()
	router := mux.NewRouter()
	router.HandleFunc("/users", createUserHandler).Methods("POST")

	user := User{Name: "Charlie", Age: 40}
	body, err := json.Marshal(user)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	var createdUser User
	if err := json.NewDecoder(rr.Body).Decode(&createdUser); err != nil {
		t.Fatalf("Could not decode response: %v", err)
	}

	if createdUser.Name != user.Name || createdUser.Age != user.Age {
		t.Errorf("Created user does not match input data: got %+v, want %+v", createdUser, user)
	}

	if err := db.Delete(&User{}, createdUser.ID).Error; err != nil {
		t.Fatalf("Failed to clean up test data: %v", err)
	}
}

func TestUpdateUser(t *testing.T) {
	setupTestDB()
	router := mux.NewRouter()
	router.HandleFunc("/users/{id}", updateUserHandler).Methods("PUT")

	var testUser User
	db.Create(&User{Name: "David", Age: 20})
	db.First(&testUser)

	updatedData := map[string]interface{}{
		"Name": "David Updated",
		"Age":  35,
	}
	body, err := json.Marshal(updatedData)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("PUT", "/users/"+strconv.Itoa(int(testUser.ID)), bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var retrievedUser User
	db.First(&retrievedUser, testUser.ID)

	if retrievedUser.Name != updatedData["Name"] || retrievedUser.Age != updatedData["Age"] {
		t.Errorf("Updated user does not match input data: got %+v, want %+v", retrievedUser, updatedData)
	}

	if err := db.Delete(&User{}, testUser.ID).Error; err != nil {
		t.Fatalf("Failed to clean up test data: %v", err)
	}
}



func TestDeleteUser(t *testing.T) {
	setupTestDB()
	router := mux.NewRouter()
	router.HandleFunc("/users/{id}", deleteUserHandler).Methods("DELETE")

	var testUser User
	db.Create(&User{Name: "Eve", Age: 28})
	db.First(&testUser)

	req, err := http.NewRequest("DELETE", "/users/"+strconv.Itoa(int(testUser.ID)), nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}

	var deletedUser User
	result := db.First(&deletedUser, testUser.ID)
	if result.Error == nil {
		t.Errorf("Expected user to be deleted, but found: %+v", deletedUser)
	} else if result.Error != gorm.ErrRecordNotFound {
		t.Errorf("Unexpected error when checking deleted user: %v", result.Error)
	}

	if err := db.Unscoped().Delete(&User{}, testUser.ID).Error; err != nil {
		t.Fatalf("Failed to clean up test data: %v", err)
	}
}


func main() {
	connectGorm()
	createModel()

	r := mux.NewRouter()
	
	r.HandleFunc("/users", getUsersHandler).Methods("GET")
	r.HandleFunc("/users", createUserHandler).Methods("POST")
	r.HandleFunc("/users/{id}", updateUserHandler).Methods("PUT")
	r.HandleFunc("/users/{id}", deleteUserHandler).Methods("DELETE")

	log.Println("SERVER STARTED ON :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
