package main

import "fmt"

func main() {
	fmt.Println(comma("123"))
	fmt.Println(comma("12345"))
	fmt.Println(comma("123456"))
	fmt.Println(comma("123456789"))
	fmt.Println(comma("1234567890"))
	fmt.Println(comma("abcef"))
}
