package main

import (
	"fmt"
)

func main() {
	var a int = 1
	var b float64 = 2.12313
	c := "ok"
	var d = true
	fmt.Printf("%v has a type %T \n", a, a)
	fmt.Printf("%v has a type %T \n", b, b)
	fmt.Printf("%v has a type %T \n", c, c)
	fmt.Printf("%v has a type %T \n", d, d)
}
