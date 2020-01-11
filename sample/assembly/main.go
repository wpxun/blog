package main

import (
	"fmt"
)

func add(a, b int) int // 汇编函数声明

func main() {
	fmt.Println(add(10, 11))
}