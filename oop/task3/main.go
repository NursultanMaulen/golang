package main

import (
	"fmt"
	"math"
)


type Shape interface {
	Area() float64
}

type Circle struct {
	Radius float64
}

type Rectangle struct {
	Height float64
	Width  float64
}

func (c Circle) Area() float64 {
	return float64(math.Pi * c.Radius * c.Radius)
}
func (r Rectangle) Area() float64 {
	return float64(r.Height * r.Width)
}

func PrintArea(s1 Shape, s2 Shape){
	fmt.Println("Area of Circle:", s1.Area())
	fmt.Println("Area of Rectangle:", s2.Area())
}

func main() {
	var c Shape
	var r Shape
	c = Circle{Radius: 3}
	r = Rectangle{Height: 3, Width: 4}

	PrintArea(c, r)
	
}