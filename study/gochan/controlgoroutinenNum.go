package main

import (
	"fmt"
	"log"
	"runtime"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	ch := make(chan struct{}, 3)
	for i := 0; i < 10; i++ {
		ch <- struct{}{} //利用 channel 的缓存区限制并发的协程数量
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			log.Println(i)
			//time.Sleep(time.Second)
			fmt.Println(runtime.NumGoroutine())
			<-ch
		}(i)
	}

	wg.Wait()
}
