package lru

import (
	"log"
	"testing"
)

func TestChangeValue(t *testing.T) {
	key := 1
	lru := New[int, int](2)
	lru.Put(key, 15)
	lru.Put(key, 17)

	val, _ := lru.Get(key)

	expect := 17
	if val != expect {
		t.Errorf("got %v, expect %v", val, expect)
	}
}

func TestCapacity(t *testing.T) {
	lru := New[int, int](2)
	lru.Put(1, 1)
	lru.Put(2, 2)
	lru.Put(3, 3)

	_, ok := lru.Get(1)
	if ok {
		t.Errorf("ok should be false")
	}
}

func TestPutAndGet(t *testing.T) {
	lRUCache := New[int, int](2)
	lRUCache.Put(1, 1) // cache is {1=1}
	lRUCache.Put(2, 2) // cache is {1=1, 2=2}

	val, ok := lRUCache.Get(1)
	if !ok {
		t.Errorf("ok should be true")
	}

	expect := 1
	if val != expect {
		t.Errorf("got %v, expect %v", val, expect)
	}

	lRUCache.Put(3, 3) // LRU key was 2, evicts key 2, cache is {1=1, 3=3}
	_, ok = lRUCache.Get(2)
	if ok {
		t.Errorf("ok should be false")
	}

	lRUCache.Put(4, 4) // LRU key was 1, evicts key 1, cache is {4=4, 3=3}

	_, ok = lRUCache.Get(1)
	if ok {
		t.Errorf("ok should be false")
	}

	val, _ = lRUCache.Get(3)
	expect = 3
	if val != expect {
		t.Errorf("got %v, expect %v", val, expect)
	}

	val, _ = lRUCache.Get(4)
	expect = 4
	if val != expect {
		t.Errorf("got %v, expect %v", val, expect)
	}
}

func TestLRU(t *testing.T) {
	evictCounter := 0
	l := New[int, int](128)

	for i := 0; i < 256; i++ {
		ok, _ := l.Put(i, i)
		if ok {
			evictCounter++
		}
	}
	if l.Len() != 128 {
		t.Fatalf("bad len: %v", l.Len())
	}

	if evictCounter != 128 {
		t.Fatalf("bad evict count: %v", evictCounter)
	}

	for i, k := range l.Keys() {
		if v, ok := l.Get(k); !ok || v != k || v != i+128 {
			log.Println(k)
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
		if ok, _ := l.Remove(i); !ok {
			t.Fatalf("should be contained")
		}
		if ok, _ := l.Remove(i); ok {
			t.Fatalf("should not be contained")
		}
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
