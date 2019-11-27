package main

import (
	"errors"
	"fmt"
)

type operate func(x, y int) int

func calculate(x int, y int, op operate) (int, error) {
	// 利用卫语句进行参数检查
	if op == nil {
		return 0, errors.New("invalid operation")
	}
	// 参数正确，进行相应的 operate 操作，没有错误发生，第二个返回值返回 nil
	return op(x, y), nil
}

type calculateFunc func(x int, y int) (int, error)

func genCalculator(op operate) calculateFunc {
	return func(x int, y int) (int, error) {
		if op == nil {
			return 0, errors.New("invalid operation")
		}
		return op(x, y), nil
	}
}

func main() {
	// 调用方法一：声明一个函数并赋值给一个变量
	op := func(x, y int) int {
		return x + y
	}
	res, _ := calculate(1, 2, op)
	fmt.Printf("res1: %d\n", res)

	// 调用方法二：实现了 operate 类型的匿名函数
	res2, _ := calculate(2, 3, func(x, y int) int {
		return x * y
	})
	fmt.Printf("res2: %d\n", res2)

	x, y := 56, 78
	add := genCalculator(op)
	res3, _ := add(x, y)
	fmt.Printf("res3: %d\n", res3)
}
