package main

import "fmt"

func main() {
	/* Subtask 1 Write a program that takes an integer input and 
		prints whether the number is positive, negative, or zero 
		using an if statement.
	*/
	_in := 0
	fmt.Print("Enter a number(default is 0): ")
	fmt.Scanln(&_in)
	if (_in > 0){
		fmt.Println("It is a positive number")
	} else if (_in < 0){
		fmt.Println("It is a negative number")
	} else {
		fmt.Println("It is a 0")
	}

	// Subtask 2 Implement a for loop that calculates the sum of the first 10 natural numbers.
	sumVal := 0
	for i := 1; i <= 10; i++{
		sumVal += i
	}
	fmt.Println("Sum of first 10 natural numbers is:", sumVal)

	/* Subtask 3 Write a switch statement 
	that prints the day of the week based on 
	an integer input (1 for Monday, 2 for Tuesday, etc.).
	*/
	fmt.Print("Input the number beetwen 1 and 7: ")
	day := 1
	fmt.Scanln(&day)
	switch day{
	case 1: fmt.Println("Monday")
	case 2: fmt.Println("Tuesday")
	case 3: fmt.Println("Wednesday")
	case 4: fmt.Println("Thursday")
	case 5: fmt.Println("Friday")
	case 6: fmt.Println("Saturday")
	case 7: fmt.Println("Sunday")
	default: fmt.Println("Not a weekday")
	}

}