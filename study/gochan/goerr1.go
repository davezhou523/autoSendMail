package main

import "fmt"

func main() {
	//chanNotClose()
	//defer1()
	//fmt.Println(defer2())
	//fmt.Println("res", defer3())
	defer4()
}
func chanNotClose() {
	ch := make(chan int)
	go func() {
		for i := 0; i < 10; i++ {
			ch <- i
		}

		//close(ch)
	}()
	//未关闭通道导致 range 或读操作永远阻塞。
	for v := range ch {
		println(v)
	}
}

// defer 先进后出 栈存储
func defer1() {
	defer func() { fmt.Println("defer1") }()
	defer func() { fmt.Println("defer2") }()
}

// 如果返回值是匿名的，defer 无法直接影响它。
// 如果返回值是命名的，defer 可以通过作用域修改它。返回变量名地址相同
func defer2() (res int) {
	defer func() {
		res += 1
	}()
	res = 100
	fmt.Println(res)
	return res
}
func defer3() int {
	res := 100
	defer func() {
		res += 1
		fmt.Println(res) //100
	}()
	return res //执行return 值时把结果赋值给匿名变量
}
func defer4() {
	defer func() { fmt.Println("clean") }()
	defer func() {
		fmt.Println("recover:", recover())
	}()
	panic("panic in defer4")
}
