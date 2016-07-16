package popcount_test

import (
	"testing"

	"github.com/ironzhang/gopl.io/ch2/popcount"
)

func TestPopCount(t *testing.T) {
	testcases := []struct {
		value uint64
		count int
	}{
		{0x0, 0},
		{0x1, 1},
		{0x0F, 4},
		{0xF0, 4},
		{0x0F0E, 7},
		{0xFF0F, 12},
		{0xFFFFFFFF, 32},
		{0xFFFFFFFFFFFFFFFE, 63},
		{0xFFFFFFFFFFFFFFFF, 64},
	}
	for _, tc := range testcases {
		n := popcount.PopCount(tc.value)
		if n != tc.count {
			t.Errorf("failed, %b %d %d", tc.value, tc.count, n)
		} else {
			t.Logf("%b %d %d", tc.value, tc.count, n)
		}
	}
}
