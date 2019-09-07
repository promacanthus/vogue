package main

import (
	"errors"
	"fmt"
)

func f1(arg int) (int, error) { // 按照惯例，errors是最后一个返回的值，error类型(一个内置的接口)
	if arg == 42 {
		return -1, errors.New("can't work with 42") // 使用给定的错误消息来构造一个基本的错误值
	}
	return arg + 3, nil // nil表示没有错误
}

type argError struct { // 自定义错误类型
	arg  int
	prob string
}

func (e *argError) Error() string { // 实现Error方法以实现error接口
	return fmt.Sprintf("%d-%s", e.arg, e.prob)
}

func f2(arg int) (int, error) {
	if arg == 42 {
		return -1, &argError{arg, "can't work with it"} // 使用&argError语法来构建一个新的结构体，同时为字段提供初始化值
	}
	return arg + 3, nil
}

func main() {
	for _, i := range []int{7, 42} {
		if r, e := f1(i); e != nil {
			fmt.Println("f1 failed:", e)
		} else {
			fmt.Println("f1 worked:", r)
		}
	}
	for _, i := range []int{7, 42} {
		if r, e := f2(i); e != nil {
			fmt.Println("f2 failed:", e)
		} else {
			fmt.Println("f2 worked:", r)
		}
	}

	// 如果要以编程方式使用自定义错误中的数据，则需要通过类型断言将错误作为自定义错误类型的实例
	_, e := f2(42)
	if ae, ok := e.(*argError); ok {
		fmt.Println(ae.arg)
		fmt.Println(ae.prob)
	}
}
