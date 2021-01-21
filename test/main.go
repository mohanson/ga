package main

import (
	"fmt"
)

type Dog struct {
	age  int
	name string
}

func main() {
	roger := Dog{5, "Roger"}
	mydog := roger

	fmt.Printf("%p\n", &roger)
	fmt.Printf("%p\n", &mydog)
}
