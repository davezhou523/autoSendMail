package main

import "fmt"

type Person struct {
	name string
	age  int
}

var list map[string]*Person

func main() {
	map1()
	map2()
}

func map1() {

	var student Person
	student.name = "jack"
	fmt.Println(student)
	list = make(map[string]*Person)
	list["jack"] = &student
	list["jack"].name = "ad"
	fmt.Println(list)
}
func map2() {
	m := make(map[string]*Person)
	student := []Person{
		{name: "jack", age: 18},
		{name: "david", age: 19},
		{name: "tom", age: 20},
	}
	for _, item := range student {
		//m[item.name] = &student[index]
		m[item.name] = &item
	}
	for key, value := range m {
		fmt.Println(key, value)
	}

}
