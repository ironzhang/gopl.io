package memotest

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"testing"
	"time"
)

func httpGetBody(url string) (interface{}, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func mockGetBody(url string) (interface{}, error) {
	time.Sleep(2 * time.Second)
	return []byte(url), nil
}

var HTTPGetBody = mockGetBody

func incomingURLs() <-chan string {
	ch := make(chan string)
	go func() {
		for _, url := range []string{
			"https://github.com/adonovan/gopl.io/blob/master/ch9/memo1/memo.go",
			"https://github.com/adonovan/gopl.io/blob/master/ch9/memo1/memo.go",
			"https://github.com/adonovan/gopl.io/blob/master/ch9/memo1/memo.go",
			"https://github.com/adonovan/gopl.io/blob/master/ch9/memo1/memo.go",
			"https://github.com/adonovan/gopl.io/blob/master/ch9/memo1/memo.go",
			"https://github.com/adonovan/gopl.io/blob/master/ch9/memo1/memo.go",
			"https://github.com/adonovan/gopl.io/blob/master/ch9/memo1/memo.go",
			"https://github.com/adonovan/gopl.io/blob/master/ch9/memo1/memo.go",
			"https://github.com/adonovan/gopl.io/blob/master/ch9/memo1/memo.go",
			"https://github.com/adonovan/gopl.io/blob/master/ch9/memo1/memo.go",
			"https://github.com/adonovan/gopl.io/blob/master/ch9/memo1/memo.go",
			"https://github.com/adonovan/gopl.io/blob/master/ch9/memo1/memo.go",
			"https://github.com/adonovan/gopl.io/blob/master/ch9/memo1/memo.go",
			"https://github.com/adonovan/gopl.io/blob/master/ch9/memo1/memo.go",
			"https://github.com/adonovan/gopl.io/blob/master/ch9/memo1/memo.go",
			"https://github.com/adonovan/gopl.io/blob/master/ch9/memo1/memo.go",
			"https://github.com/adonovan/gopl.io/blob/master/ch9/memo1/memo.go",
			"https://github.com/adonovan/gopl.io/blob/master/ch9/memo1/memo.go",
			"https://github.com/adonovan/gopl.io/blob/master/ch9/memo1/memo.go",
			"https://github.com/adonovan/gopl.io/blob/master/ch9/memo1/memo.go",
			"https://github.com/adonovan/gopl.io/blob/master/ch9/memo1/memo.go",
			"https://github.com/adonovan/gopl.io/blob/master/ch9/memo1/memo.go",
			"https://github.com/adonovan/gopl.io/blob/master/ch9/memo1/memo.go",
			"https://github.com/adonovan/gopl.io/blob/master/ch9/memo1/memo.go",
			"https://github.com/adonovan/gopl.io/blob/master/ch9/memo1/memo.go",
			"https://github.com/adonovan/gopl.io/blob/master/ch9/memo1/memo.go",
			"https://github.com/adonovan/gopl.io/blob/master/ch9/memo1/memo.go",
			"https://github.com/adonovan/gopl.io/blob/master/ch9/memo1/memo.go",
			"https://github.com/adonovan/gopl.io/blob/master/ch9/memo1/memo.go",
			"https://github.com/adonovan/gopl.io/blob/master/ch9/memo1/memo.go",
			"https://github.com/adonovan/gopl.io/blob/master/ch9/memo1/memo.go",
			"https://github.com/adonovan/gopl.io/blob/master/ch9/memo1/memo.go",
			"https://github.com/adonovan/gopl.io/blob/master/ch9/memo1/memo.go",
			"https://github.com/adonovan/gopl.io/blob/master/ch9/memo2/memo.go",
			"https://github.com/adonovan/gopl.io/blob/master/ch9/memo3/memo.go",
			"https://github.com/adonovan/gopl.io/blob/master/ch9/memo2/memo.go",
			"https://github.com/adonovan/gopl.io/blob/master/ch9/memo4/memo.go",
			"https://github.com/adonovan/gopl.io/blob/master/ch9/memo5/memo.go",
			"https://github.com/adonovan/gopl.io/blob/master/ch9/memo5/memo.go",
			"https://github.com/adonovan/gopl.io/blob/master/ch9/memo1/memo.go",
			"https://github.com/adonovan/gopl.io/blob/master/ch9/memo1/memo.go",
			"https://github.com/adonovan/gopl.io/blob/master/ch9/memo1/memo.go",
			"https://github.com/adonovan/gopl.io/blob/master/ch9/memo2/memo.go",
			"https://github.com/adonovan/gopl.io/blob/master/ch9/memo3/memo.go",
			"https://github.com/adonovan/gopl.io/blob/master/ch9/memo2/memo.go",
			"https://github.com/adonovan/gopl.io/blob/master/ch9/memo4/memo.go",
			"https://github.com/adonovan/gopl.io/blob/master/ch9/memo5/memo.go",
			"https://github.com/adonovan/gopl.io/blob/master/ch9/memo5/memo.go",
			"https://github.com/adonovan/gopl.io/blob/master/ch9/memo1/memo.go",
			"https://github.com/adonovan/gopl.io/blob/master/ch9/memo1/memo.go",
		} {
			ch <- url
		}
		close(ch)
	}()
	return ch
}

type M interface {
	Get(key string) (interface{}, error)
}

func Sequential(t *testing.T, m M) {
	for url := range incomingURLs() {
		start := time.Now()
		value, err := m.Get(url)
		if err != nil {
			log.Print(err)
			continue
		}
		fmt.Printf("%s, %s, %d bytes\n", url, time.Since(start), len(value.([]byte)))
	}
}

func Concurrent(t *testing.T, m M) {
	t1 := time.Now()
	var wg sync.WaitGroup
	for url := range incomingURLs() {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			start := time.Now()
			value, err := m.Get(url)
			if err != nil {
				log.Print(err)
				return
			}
			fmt.Printf("%s, %s, %d bytes\n", url, time.Since(start), len(value.([]byte)))
		}(url)
	}
	wg.Wait()
	fmt.Printf("total: %s\n", time.Since(t1))
}
