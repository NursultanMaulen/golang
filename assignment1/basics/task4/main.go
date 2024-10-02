package main

import "fmt"

func add(a int, b int) int {
	return a + b
}

func swap(s1 string, s2 string) (string, string){
	return s2, s1
}

func getQuotAndRem(num1 int, num2 int) (quot int, rem int){
	quot = num1 / num2
	rem = num1 % num2
	return
}

func main(){
	// Subtask 1 Write a function add that takes two integers as arguments and returns their sum.
	a := 0
	b := 0
	fmt.Print("Enter a value: ")
	fmt.Scan(&a)
	fmt.Print("Enter b value: ")
	fmt.Scan(&b)
	fmt.Println("a + b = ", add(a, b))

	// Subtask 2 Create a function swap that returns the two input strings in reverse order.
	s1 := ""
	s2 := ""
	fmt.Print("Enter s1 value: ")
	fmt.Scan(&s1)
	fmt.Print("Enter s2 value: ")
	fmt.Scan(&s2)
	fmt.Println("Original s1:", s1, ", Original s2: ", s2)
	s1, s2 = swap(s1, s2)
	fmt.Println("Reversed s1:", s1, ", Reversed s2: ", s2)

	// Subtask 3 Write a function that returns both the quotient and remainder of two integers.
	num1 := 5
	num2 := 4
	fmt.Print("Enter num1 value: ")
	fmt.Scan(&num1)
	fmt.Print("Enter num2 value: ")
	fmt.Scan(&num2)
	fmt.Println(getQuotAndRem(num1, num2))
}