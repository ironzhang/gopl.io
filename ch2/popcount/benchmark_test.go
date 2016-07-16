package popcount_test

import (
	"testing"

	"github.com/ironzhang/gopl.io/ch2/popcount"
)

func BenchmarkPopCount(b *testing.B) {
	for i := 0; i < b.N; i++ {
		popcount.PopCount(0x1234567890ABCDEF)
	}
}

//go test -cpu=4 -bench=. github.com/ironzhang/gopl.io/ch2/popcount
