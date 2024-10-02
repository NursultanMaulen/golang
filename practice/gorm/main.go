package main

import (
	"database/sql"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct{
	ID uint `gorm:"PrimaryKey"`
	Name string `gorm:"size:100"`
	Email string `gorm:"unique"`
}

func useSql(){
	dsn := "postgres://postgres:admin@localhost:5432/example_go"
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	err = db.Ping()
	if err != nil{
		fmt.Println(err)
		return
	}

	fmt.Println("SUCCESS")

	// sqlIns := `INSERT INTO users (name, email) VALUES ($1, $2)`

	// _, err = db.Exec(sqlIns, "Example_", "example@mail.ru")

	// if err != nil{
	// 	fmt.Println(err)
	// 	return
	// }

	res, err := db.Query("select * from users")
	if err != nil{
		fmt.Println(err)
		return
	}
	defer res.Close()

	for res.Next() {
		var id int
		var name, email string

		err := res.Scan(&id, &name, &email)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(id, name, email)
	}

}

func useGorm(){
	dsn := "host=localhost user=postgres password=admin dbname=example_go port=5432 sslmode=disable TimeZone=Europe/Moscow"

    // Подключение к базе данных
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        fmt.Println("Ошибка подключения к базе данных:", err)
        return
    }

    // Автоматическое создание таблицы для модели User
    db.AutoMigrate(&User{})

    // Создание новой записи в базе данных
    user := User{Name: "Nursultan1", Email: "nursultan1@example.com"}
    db.Create(&user)

    // fmt.Println("Успешно подключено к базе данных и добавлен пользователь")

	var checkUser User

	res := db.Where("email = ?", "nursultan1@example.com").First(&checkUser)

	if(res.Error != nil){
		fmt.Println(res.Error)
		return
	}

	fmt.Println("USER READ:", checkUser.ID, "smt", checkUser.Name, checkUser.Email)
}

func main() {
	useGorm()
}