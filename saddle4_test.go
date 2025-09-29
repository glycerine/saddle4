package main

import (
	"testing"
)

func Test(t *testing.T) {
	//path := "fifo.pony"
	//by, err := os.ReadFile(path)
	//panicOn(err)
	if want, got := 3, firstNonSpace([]byte("   a")); want != got {
		panicf("want %v, got %v", want, got)
	}
	if want, got := 0, firstNonSpace([]byte("a")); want != got {
		panicf("want %v, got %v", want, got)
	}
	if want, got := 1, firstNonSpace([]byte(" a")); want != got {
		panicf("want %v, got %v", want, got)
	}
	if want, got := 5, firstNonSpace([]byte("     a")); want != got {
		panicf("want %v, got %v", want, got)
	}

	if want, got := 1, closestStopDelta(firstNonSpace([]byte("   a")), 4); want != got {
		panicf("want %v, got %v", want, got)
	}
	if want, got := 0, closestStopDelta(firstNonSpace([]byte("a")), 4); want != got {
		panicf("want %v, got %v", want, got)
	}
	if want, got := 3, closestStopDelta(firstNonSpace([]byte(" a")), 4); want != got {
		panicf("want %v, got %v", want, got)
	}
	if want, got := -1, closestStopDelta(firstNonSpace([]byte("     a")), 4); want != got {
		panicf("want %v, got %v", want, got)
	}

}
