// +build ignore

// Run with:
//   $ go run ~/gopath/src/github.com/ironzhang/gopl.io/ch8/thumbnail/main.go
//   foo.jpeg
//   ^D

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/ironzhang/gopl.io/ch8/thumbnail"
)

func main() {
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		thumb, err := thumbnail.ImageFile(input.Text())
		if err != nil {
			log.Print(err)
			continue
		}
		fmt.Println(thumb)
	}
	if err := input.Err(); err != nil {
		log.Fatal(err)
	}
}
