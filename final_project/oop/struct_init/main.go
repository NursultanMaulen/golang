package main

import "fmt"

type Product struct {
    ID    int
    Name  string
    Price float64
}

func (p Product) Display() {
    fmt.Printf("Product ID: %d, Name: %s, Price: %.2f\n", p.ID, p.Name, p.Price)
}

func main() {
    p := Product{ID: 1, Name: "Laptop", Price: 1200.50}
    p.Display()
}
