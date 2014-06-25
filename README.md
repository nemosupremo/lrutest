lrutest
=======

Simple test to compare implementations of LRU Caching.

Results
======


PASS
BenchmarkGetPLRU	   50000	     48464 ns/op
BenchmarkGetLLRU	   50000	     63148 ns/op
BenchmarkNotSequentialPLRU	   50000	     44317 ns/op
BenchmarkNotSequentialLLRU	   50000	     60072 ns/op
BenchmarkSetPLRU	   50000	     53234 ns/op
BenchmarkSetLLRU	   50000	     65030 ns/op

