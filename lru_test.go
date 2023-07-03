package golru_test

import (
	"fmt"
	"testing"

	"github.com/thepatrik/golru"
)

func TestSouce(t *testing.T) {
	lru, _ := golru.New[string, string](golru.WithMaxEntries(1000))

	lru.Put("abracadabra", "magic dragon")

	val, ok := lru.Get("abracadabra")
	if ok {
		fmt.Printf("found value \"%v\"\n", val)
	}
}

func BenchmarkLRU_Rand(b *testing.B) {
	l, err := golru.New[int64, int64](golru.WithMaxEntries(8192))
	if err != nil {
		b.Fatalf("err: %v", err)
	}

	trace := make([]int64, b.N*2)
	for i := 0; i < b.N*2; i++ {
		trace[i] = getRand(b) % 32768
	}

	b.ResetTimer()

	var hit, miss int
	for i := 0; i < 2*b.N; i++ {
		if i%2 == 0 {
			l.Put(trace[i], trace[i])
		} else {
			if _, ok := l.Get(trace[i]); ok {
				hit++
			} else {
				miss++
			}
		}
	}
	b.Logf("hit: %d miss: %d ratio: %f", hit, miss, float64(hit)/float64(hit+miss))
}

func BenchmarkLRU_Freq(b *testing.B) {
	l, err := golru.New[int64, int64](golru.WithMaxEntries(8192))
	if err != nil {
		b.Fatalf("err: %v", err)
	}

	trace := make([]int64, b.N*2)
	for i := 0; i < b.N*2; i++ {
		if i%2 == 0 {
			trace[i] = getRand(b) % 16384
		} else {
			trace[i] = getRand(b) % 32768
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		l.Put(trace[i], trace[i])
	}
	var hit, miss int
	for i := 0; i < b.N; i++ {
		if _, ok := l.Get(trace[i]); ok {
			hit++
		} else {
			miss++
		}
	}
	b.Logf("hit: %d miss: %d ratio: %f", hit, miss, float64(hit)/float64(hit+miss))
}

func TestLRU(t *testing.T) {
	evictCounter := 0
	l, err := golru.New[int, int](golru.WithMaxEntries(128))
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	l.SetOnEvicted(func(k int, v int) {
		if k != v {
			t.Fatalf("Evict values not equal (%v!=%v)", k, v)
		}
		evictCounter++
	})

	for i := 0; i < 256; i++ {
		l.Put(i, i)
	}
	if l.Len() != 128 {
		t.Fatalf("bad len: %v", l.Len())
	}

	if evictCounter != 128 {
		t.Fatalf("bad evict count: %v", evictCounter)
	}

	for i, k := range l.Keys() {
		if v, ok := l.Get(k); !ok || v != k || v != i+128 {
			t.Fatalf("bad key: %v", k)
		}
	}
	for i, v := range l.Values() {
		if v != i+128 {
			t.Fatalf("bad value: %v", v)
		}
	}
	for i := 0; i < 128; i++ {
		if _, ok := l.Get(i); ok {
			t.Fatalf("should be evicted")
		}
	}
	for i := 128; i < 256; i++ {
		if _, ok := l.Get(i); !ok {
			t.Fatalf("should not be evicted")
		}
	}
	for i := 128; i < 192; i++ {
		l.Remove(i)
		if _, ok := l.Get(i); ok {
			t.Fatalf("should be deleted")
		}
	}

	l.Get(192) // expect 192 to be last key in l.Keys()

	for i, k := range l.Keys() {
		if (i < 63 && k != i+193) || (i == 63 && k != 192) {
			t.Fatalf("out of order key: %v", k)
		}
	}
}

// test that Add returns true/false if an eviction occurred
func TestLRUPut(t *testing.T) {
	l, err := golru.New[int, int](golru.WithMaxEntries(1))
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	evictCounter := 0

	l.SetOnEvicted(func(k int, v int) {
		evictCounter++
	})

	l.Put(1, 1)
	if evictCounter != 0 {
		t.Errorf("should not have an eviction")
	}

	l.Put(2, 2)
	if evictCounter != 1 {
		t.Errorf("should have an eviction")
	}
}

// test that Contains doesn't update recent-ness
func TestLRUContains(t *testing.T) {
	l, err := golru.New[int, int](golru.WithMaxEntries(2))
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	l.Put(1, 1)
	l.Put(2, 2)
	if !l.Contains(1) {
		t.Errorf("1 should be contained")
	}

	l.Put(3, 3)
	if l.Contains(1) {
		t.Errorf("Contains should not have updated recent-ness of 1")
	}
}

// test that Peek doesn't update recent-ness
func TestLRUPeek(t *testing.T) {
	l, err := golru.New[int, int](golru.WithMaxEntries(2))
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	l.Put(1, 1)
	l.Put(2, 2)
	if v, ok := l.Peek(1); !ok || v != 1 {
		t.Errorf("1 should be set to 1: %v, %v", v, ok)
	}

	l.Put(3, 3)
	if l.Contains(1) {
		t.Errorf("should not have updated recent-ness of 1")
	}

	if v, ok := l.Peek(1); ok || v != 0 {
		t.Errorf("0 should not exist: %v, %v", v, ok)
	}
}
