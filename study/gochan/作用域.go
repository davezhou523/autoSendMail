package main

import (
	"fmt"
	"github.com/pkg/errors"
)

var ErrDidNotWork = errors.New("did not work")

func do(r bool) (err error) {
	var res string
	if r {
		//res, err := thing() //变量作用域 if 语句块内的 err 变量会遮罩函数作用域内的 err 变量
		res, err = thing() //调整
		if err != nil || res != "it worked" {
			err = ErrDidNotWork
		}
	}
	return err
}
func thing() (string, error) {
	return "", ErrDidNotWork
}
func main() {
	fmt.Println(do(true))
	fmt.Println(do(false))
}
