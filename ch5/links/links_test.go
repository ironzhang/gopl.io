package links

import (
	"fmt"
	"testing"
)

func TestExtract(t *testing.T) {
	url := "https://golang.org"
	links, err := Extract(url)
	if err != nil {
		t.Fatalf("Extract: %v", err)
	}
	for _, link := range links {
		fmt.Println(link)
	}
}
