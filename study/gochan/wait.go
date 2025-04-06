package main

import (
	"fmt"
	"sync"
	//"time"
)

const N = 10

var wg = &sync.WaitGroup{}

func main() {
	//wggoroutine()
	//seqshow()
	DeferFunc4()
}
func DeferFunc4() (t int) {
	defer func(i int) {
		fmt.Println(i)
		fmt.Println(t)
	}(t)
	t = 1
	return 2
}

func seqshow() {
	ch := make(chan int)
	go func() {
		for i := 1; i <= N; i++ {
			ch <- i
		}
		close(ch)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for v := range ch {
			fmt.Println(v)
		}
	}()
	wg.Wait()

}
func wggoroutine() {
	for i := 0; i < N; i++ {
		wg.Add(1)
		go func(i int) {
			println(i) //无序输出
			defer wg.Done()
		}(i)
	}

	wg.Wait()

}
