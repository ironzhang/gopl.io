package memo_test

import (
	"testing"

	"github.com/ironzhang/gopl.io/ch9/memo2"
	"github.com/ironzhang/gopl.io/ch9/memotest"
)

func TestSequential(t *testing.T) {
	m := memo.New(memotest.HTTPGetBody)
	memotest.Sequential(t, m)
}

func TestConcurrent(t *testing.T) {
	m := memo.New(memotest.HTTPGetBody)
	memotest.Concurrent(t, m)
}
