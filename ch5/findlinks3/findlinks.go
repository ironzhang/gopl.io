package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ironzhang/gopl.io/ch5/links"
)

func breadthFirst(f func(item string) []string, worklist []string) {
	seen := make(map[string]bool)
	for len(worklist) > 0 {
		var list []string
		for _, item := range worklist {
			if !seen[item] {
				seen[item] = true
				list = append(list, f(item)...)
			}
		}
		worklist = list
	}
}

func crawl(url string) []string {
	fmt.Println(url)
	list, err := links.Extract(url)
	if err != nil {
		log.Print(err)
	}
	return list
}

func main() {
	breadthFirst(crawl, os.Args[1:])
}
