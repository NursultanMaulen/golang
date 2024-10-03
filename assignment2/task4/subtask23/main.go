package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	// _ "subtask23/docs"
)

// @title User API
// @version 1.0
// @description API for managing users.
// @host localhost:8080
// @BasePath /

type User struct {
	ID     uint    `gorm:"primaryKey"`
	Name   string  `gorm:"size:255"`
	Age    int
	Profile Profile `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type Profile struct {
	ID               uint   `gorm:"primaryKey"`
	UserID           uint   `gorm:"unique;not null;constraint:OnDelete:CASCADE;"`
	Bio              string `gorm:"size:1024"`
	ProfilePictureURL string `gorm:"size:255"`
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
	db.AutoMigrate(&User{}, &Profile{})
}


// @Summary Create a new user
// @Description Creates a new user in the database
// @Tags Users
// @Accept  json
// @Produce  json
// @Param   user  body   User  true  "User Data"
// @Success 201 {object} User
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 500 {object} map[string]string "Server error"
// @Router /users [post]
func getUsersHandler(w http.ResponseWriter, r *http.Request) {
	var users []User

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

	query := db.Preload("Profile").Limit(pageSizeInt).Offset(offset).Find(&users)
	if err := query.Find(&users).Error; err != nil {
		http.Error(w, "ERROR GETTING USER", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// @Summary Create a new user
// @Description Creates a new user in the database
// @Tags Users
// @Accept  json
// @Produce  json
// @Param   user  body   User  true  "User Data"
// @Success 201 {object} User
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 500 {object} map[string]string "Server error"
// @Router /users [post]
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

func main() {
	connectGorm()
	createModel()

	r := mux.NewRouter()
	fs := http.FileServer(http.Dir("./docs"))
	r.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger/", fs))

	r.HandleFunc("/users", getUsersHandler).Methods("GET")
	r.HandleFunc("/users", createUserHandler).Methods("POST")
	r.HandleFunc("/users/{id}", updateUserHandler).Methods("PUT")
	r.HandleFunc("/users/{id}", deleteUserHandler).Methods("DELETE")

	log.Println("SERVER STARTED ON :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
