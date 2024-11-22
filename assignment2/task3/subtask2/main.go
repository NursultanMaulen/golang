package main

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

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

var db *gorm.DB;
var err error;

func connectGorm(){
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

func createModel(){
	db.AutoMigrate(&User{}, &Profile{})

	user := User{
		Name: "John Doe",
		Age:  30,
		Profile: Profile{
			Bio:              "Software Engineer",
			ProfilePictureURL: "https://example.com/profile.jpg",
		},
	}

	result := db.Create(&user)
	if result.Error != nil {
		log.Fatalf("FAIL INSERT USER: %v", result.Error)
	}

	fmt.Printf("NEW USER WITH ID: %d\n", user.ID)

	var retrievedUser User
	db.Preload("Profile").First(&retrievedUser, user.ID)
	fmt.Printf("RETRIEVED USER: %v\n", retrievedUser)
	fmt.Printf("USER PROFILE: %v\n", retrievedUser.Profile)
}

func insertUserWithProfile() error {
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	user := User{
		Name: "Jane Doe",
		Age:  28,
		Profile: Profile{
			Bio:              "Designer and Artist",
			ProfilePictureURL: "https://example.com/jane.jpg",
		},
	}

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	fmt.Println("USER AND PROFILE ADDED!")
	return nil
}

func queryUsersWithProfiles() {
	var users []User

	result := db.Preload("Profile").Find(&users)
	if result.Error != nil {
		log.Fatalf("FAIL SELECTING USERS: %v", result.Error)
	}

	for _, user := range users {
		fmt.Printf("USER: ID=%d, Name=%s, Age=%d\n", user.ID, user.Name, user.Age)
		fmt.Printf("PROFILE: Bio=%s, ProfilePictureURL=%s\n", user.Profile.Bio, user.Profile.ProfilePictureURL)
	}
}

func deleteUser(userID uint) error {
	var user User

	result := db.Preload("Profile").First(&user, userID)
	if result.Error != nil {
		return result.Error
	}

	if err := db.Delete(&user).Error; err != nil {
		return err
	}

	fmt.Printf("USER ID %d AND PROFILE DELETED!\n", userID)
	return nil
}

func updateUserProfile(db *gorm.DB, userID uint, newBio string, newProfilePictureURL string) error {
	var profile Profile

	result := db.Where("user_id = ?", userID).First(&profile)
	if result.Error != nil {
		return result.Error
	}

	profile.Bio = newBio
	profile.ProfilePictureURL = newProfilePictureURL

	if err := db.Save(&profile).Error; err != nil {
		return err
	}

	fmt.Printf("USER ID %d PROFILE UPDATED!\n", userID)
	return nil
}

func main() {
	connectGorm()
	// createModel()
	// insertUserWithProfile()
	queryUsersWithProfiles()
	// deleteUser(1)
// 	err := updateUserProfile(db, 2, "NEW BIO", "https://example.com/new-profile-picture.jpg")
// 	if err != nil {
// 		log.Fatalf("Ошибка при обновлении профиля: %v", err)
// 	}
}