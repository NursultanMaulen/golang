package main

import "fmt"

type Displayable interface {
    Display()
}

type User struct {
    Name  string
    Email string
}

func (u User) Display() {
    fmt.Printf("User: %s, Email: %s\n", u.Name, u.Email)
}

func printInfo(d Displayable) {
    d.Display()
}

func main() {
    user := User{Name: "Nursultan Maulen", Email: "maulen@gmail.com"}
    printInfo(user)
}
