package lru

import (
	"github.com/golang/groupcache/lru"
	"strconv"
	"testing"
)

const (
	lruSize = 10000
)

func BenchmarkGetPLRU(b *testing.B) {
	plru := NewPreparedLRU()
	plru.setMaxStmts(lruSize)
	for i := 0; i < lruSize; i++ {
		plru.set("A", strconv.FormatInt(int64(i), 10), "A")
	}
	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		plru.get("A", strconv.FormatInt(int64(n), 10))
	}
}

func BenchmarkGetLLRU(b *testing.B) {
	plru := NewListLRU()
	plru.setMaxStmts(lruSize)
	for i := 0; i < lruSize; i++ {
		plru.set("A", strconv.FormatInt(int64(i), 10), "A")
	}
	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		plru.get("A", strconv.FormatInt(int64(n), 10))
	}
}

func BenchmarkGetGPLRU(b *testing.B) {
	plru := lru.New(lruSize)
	for i := 0; i < lruSize; i++ {
		plru.Add("A"+strconv.FormatInt(int64(i), 10), "A")
	}
	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		plru.Get("A" + strconv.FormatInt(int64(n), 10))
	}
}

func BenchmarkNotSequentialPLRU(b *testing.B) {
	plru := NewPreparedLRU()
	plru.setMaxStmts(lruSize)
	for i := 0; i < lruSize; i++ {
		plru.set("A", strconv.FormatInt(int64(i), 10), "A")
	}
	// run the Fib function b.N times
	for n := 1; n <= b.N; n++ {
		plru.get("A", strconv.FormatInt(int64(lruSize%n), 10))
	}
}

func BenchmarkNotSequentialLLRU(b *testing.B) {
	plru := NewListLRU()
	plru.setMaxStmts(lruSize)
	for i := 0; i < lruSize; i++ {
		plru.set("A", strconv.FormatInt(int64(i), 10), "A")
	}
	// run the Fib function b.N times
	for n := 1; n <= b.N; n++ {
		plru.get("A", strconv.FormatInt(int64(lruSize%n), 10))
	}
}

func BenchmarkNotSequentialGPLRU(b *testing.B) {
	plru := lru.New(lruSize)
	for i := 0; i < lruSize; i++ {
		plru.Add("A"+strconv.FormatInt(int64(i), 10), "A")
	}
	// run the Fib function b.N times
	for n := 1; n <= b.N; n++ {
		plru.Get("A" + strconv.FormatInt(int64(lruSize%n), 10))
	}
}

func BenchmarkSetPLRU(b *testing.B) {
	plru := NewPreparedLRU()
	plru.setMaxStmts(lruSize)
	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		plru.set("A", strconv.FormatInt(int64(n), 10), "A")
	}
}

func BenchmarkSetLLRU(b *testing.B) {
	plru := NewListLRU()
	plru.setMaxStmts(lruSize)
	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		plru.set("A", strconv.FormatInt(int64(n), 10), "A")
	}
}

func BenchmarkSetGPLRU(b *testing.B) {
	plru := lru.New(lruSize)
	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		plru.Add("A"+strconv.FormatInt(int64(n), 10), "A")
	}
}
