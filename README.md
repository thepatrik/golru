# golru

golru is a simple LRU cache, configurable with a max number of entries. If the cache is full and another entry is added, the least recently used entry is evicted. The cache is safe for concurrent access.

The implementation is based on a doubly linked list and a hash table (map). The map makes the time complexity of `get()` to O(1). The list makes the time complexity of `put()` to O(1). The space complexity is O(n).

```go
import (
	"fmt"

	"github.com/thepatrik/golru"
)

func main() {
	lru, _ := golru.New[string, any](golru.WithMaxEntries(100))

	lru.SetOnEvicted(func(k string, v any) {
		fmt.Printf("evicted key %v with value %v", k, v)
	})

	lru.Put("abracadabra", "Puff, the Magic Dragon")

	val, ok := lru.Get("abracadabra")
	if ok {
		fmt.Printf("found value \"%v\"\n", val)
	}
}
```