package main

import (
	"fmt"
)

func find(arr []int, one int) bool  {

	length := len(arr)

	if length <= 0 {
		return false
	}

	midInt := length/2

	if one < arr[midInt] {
		arr = arr[:midInt]
	} else if one > arr[midInt] {
		arr = arr[midInt+1:]
	} else if one == arr[midInt] {
		return true
	}

	return find(arr, one)
}


func main()  {
	out := find([]int{1,6,9,14,51}, 16)
	fmt.Println(out)
}
