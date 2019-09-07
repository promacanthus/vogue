package main

import "fmt"

// person ：结构体包含name和age两个字段
type person struct {
	name string
	age  int
}

// NewPerson ：根据给定的name创建一个person结构体
func NewPerson(name string) *person {
	p := person{name: name}
	p.age = 42
	return &p // 可以安全的返回指向局部变量的指针，因为局部变量将在函数范围内存活
}

func main() {
	fmt.Println(person{"Bob", 20})              //这种语法创建一个新的结构体
	fmt.Println(person{name: "alice", age: 30}) // 可以在初始化结构体时命名字段
	fmt.Println(person{name: "Fred"})           // 忽略的字段将会被命名为该字段的零值
	fmt.Println(&person{name: "Ann", age: 40})  // &前缀将产生一个指向结构体的指针
	fmt.Println(NewPerson("Jon"))               // 在构造函数中封装新结构的创建不太好

	s := person{name: "Sean", age: 50}
	fmt.Println(s.name) // 使用点访问结构体的字段

	sp := &s
	fmt.Println(sp.age) // 在结构体指针上也可以使用点，此时指针会自动接触引用

	sp.age = 51 // 结构体是可变的
	fmt.Println(sp.age)
}
