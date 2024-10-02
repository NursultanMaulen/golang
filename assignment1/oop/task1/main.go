package main

import "fmt"

type Person struct{
	Name string
	Age int
}

func (r Person) Greet(){
	fmt.Println("Hello, ", r.Name, "!")
}

func main(){
	/*Define a struct called Person with fields Name and Age.
		Write a method Greet for the Person struct that prints a greeting message.
		Create an instance of Person, set its fields, and call the Greet method.*/
	pers := Person{Name: "Nursultan", Age: 20}
	pers.Greet()
}