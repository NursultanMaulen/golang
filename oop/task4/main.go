package main

import (
	"encoding/json"
	"fmt"
)


type Product struct {
	Name string `json: "Name"`
	Price int `json: "Price"`
	Quantity int `json: "Quantity`
}


func (p Product) convertJson() (jsonData []byte){
	jsonData, err := json.Marshal(p)
	if err != nil {
		fmt.Println(err)
		return
	}
	return
}

func main() {

	p := Product{Name: "Tomato", Price: 10, Quantity: 100}
	fmt.Println("Json: ", string(p.convertJson()))

	jsonD := `{"Name": "Carrot", "Price": 30, "Quantity": 1111}`

	var p2 Product
	err := json.Unmarshal([]byte(jsonD), &p2)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Decoded Struct:\n", p2)
}