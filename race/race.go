package main 

import (
	"fmt"
	"sync"
)

func main() {
	var mu sync.Mutex
	var counter = make(map[string]int)
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			counter["hits"]++
			mu.Unlock()
		}()

	}
	wg.Wait()
	fmt.Println("Total: ", counter["hits"])
}