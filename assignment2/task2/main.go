package main

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct{
	ID uint `gorm:"PrimaryKey"`
	Name string `gorm:"not null"`
	Age uint `gorm:"not null"`
}

var db *gorm.DB;
var err error;

func connectGorm(){
	dsn := "host=localhost user=postgres password=admin dbname=assignment2 port=5432 sslmode=disable TimeZone=Europe/Moscow"

    db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        fmt.Println("Ошибка подключения к базе данных:", err)
        return
    }
	fmt.Println("SUCCESS CONNECTION")
}

func createModel(){
	db.AutoMigrate(&User{})

    user := User{Name: "Nursultan1", Age: 12}
    db.Create(&user)

    fmt.Println("SUCCESS ADDING USER")

	var checkUser User

	res := db.Where("age = ?", 12).First(&checkUser)

	if(res.Error != nil){
		fmt.Println(res.Error)
		return
	}

	fmt.Println("USER READ:", checkUser.ID, "smt", checkUser.Name, checkUser.Age)
}

func main() {
	connectGorm()
	createModel()
}