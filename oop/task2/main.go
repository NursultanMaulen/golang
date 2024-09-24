package main

import "fmt"

type Employee struct{
	Name string
	ID int
}

type Manager struct{
	Employee
	Department string
}

func (e Employee) Work(){
	fmt.Println(e.Name, e.ID)
}

func main(){
	/*Learn how to achieve composition in Go through embedding.*/
	mng := Manager{Employee: Employee{Name: "Nursultan", ID: 123123}, 
	Department: "Big Data"}

	mng.Work()
}