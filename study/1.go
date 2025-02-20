package main

import "fmt"

func main() {
	ch := make(chan int)
	go func() {
		ch <- 1
		//close(ch)
	}()
	for v := range ch { // 如果通道未关闭，range 不会结束
		fmt.Println(v)
	}

}
