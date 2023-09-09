package snowflake

import (
	"fmt"
	"sync"
	"testing"
)

func TestConcurrentGenerateId(t *testing.T) {
	sf, err := NewSnowflake(1)
	if err != nil {
		t.Fatalf("Error initializing snowflake: %v", err)
	}

	const count = 10000
	var wg sync.WaitGroup
	ch := make(chan int64, count)
	m := make(map[int64]byte)
	wg.Add(count)
	defer close(ch)

	for i := 0; i < count; i++ {
		go func() {
			defer wg.Done()
			id, _ := sf.GenerateId()
			ch <- id
		}()
	}

	wg.Wait()

	for i := 0; i < count; i++ {
		id := <-ch
		if _, exists := m[id]; exists {
			t.Fatalf("Duplicate ID: %d", id)
		}
		m[id] = 1
	}
	fmt.Println("All", count, "IDs uniquely generated!")
}
