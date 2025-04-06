package main

import (
	"fmt"
	"time"
)

func main() {
	//cloneRange()
	slice1()
}

func cloneRange() {
	var a = [5]int{1, 2, 3, 4, 5}
	var r [5]int
	for i, v := range a {
		if i == 0 {
			a[1] = 12
			a[2] = 13
		}
		r[i] = v
	}
	fmt.Println("r = ", r)
	fmt.Println("a = ", a)
}

func slice1() {
	slice := []int{1, 2, 3}
	m := make(map[int]int)
	var slice2 [3]int
	for index, value := range slice {
		slice = append(slice, value)
		go func() {
			fmt.Println("in goroutine: ", index, value)
		}()
		//time.Sleep(time.Second * 1)
		m[index] = value
		if index == 0 {
			slice[1] = 11
			slice[2] = 22
		}
		slice2[index] = value
	}
	fmt.Println("slice: ", slice)
	for key, value := range m {
		fmt.Println("in map: ", key, "->", value)
	}
	fmt.Println("slice2: ", slice2)
	time.Sleep(time.Second * 10)
}
