package lru

import (
	"strconv"
	"testing"
)

func BenchmarkGetPLRU(b *testing.B) {
	plru := NewPreparedLRU()
	plru.setMaxStmts(1000)
	for i := 0; i < 1000; i++ {
		plru.set("A", strconv.FormatInt(int64(i), 10), "A")
	}
	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		plru.get("A", strconv.FormatInt(int64(n), 10))
	}
}

func BenchmarkGetLLRU(b *testing.B) {
	plru := NewListLRU()
	plru.setMaxStmts(1000)
	for i := 0; i < 1000; i++ {
		plru.set("A", strconv.FormatInt(int64(i), 10), "A")
	}
	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		plru.get("A", strconv.FormatInt(int64(n), 10))
	}
}

func BenchmarkNotSequentialPLRU(b *testing.B) {
	plru := NewPreparedLRU()
	plru.setMaxStmts(1000)
	for i := 0; i < 1000; i++ {
		plru.set("A", strconv.FormatInt(int64(i), 10), "A")
	}
	// run the Fib function b.N times
	for n := 1; n <= b.N; n++ {
		plru.get("A", strconv.FormatInt(int64(1000%n), 10))
	}
}

func BenchmarkNotSequentialLLRU(b *testing.B) {
	plru := NewListLRU()
	plru.setMaxStmts(1000)
	for i := 0; i < 1000; i++ {
		plru.set("A", strconv.FormatInt(int64(i), 10), "A")
	}
	// run the Fib function b.N times
	for n := 1; n <= b.N; n++ {
		plru.get("A", strconv.FormatInt(int64(1000%n), 10))
	}
}

func BenchmarkSetPLRU(b *testing.B) {
	plru := NewPreparedLRU()
	plru.setMaxStmts(1000)
	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		plru.set("A", strconv.FormatInt(int64(n), 10), "A")
	}
}

func BenchmarkSetLLRU(b *testing.B) {
	plru := NewListLRU()
	plru.setMaxStmts(1000)
	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		plru.set("A", strconv.FormatInt(int64(n), 10), "A")
	}
}
