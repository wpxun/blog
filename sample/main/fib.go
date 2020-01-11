package main

import (
	"fmt"
	"math/rand"
	"time"
)

func fib(n int) <-chan int {
	result := make(chan int)

	go func() {
		defer close(result)

		if n <= 2 {
			result <- 1
			return
		}

		result <- <-fib(n-1) + <-fib(n-2)
	}()

	return result
}

func main()  {

	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(20)
	fmt.Println(n)

	if n % 2 == 0 {
		race()
	}

	fmt.Printf("fib = %d\n", <-fib(4))
	
}

func race (){
	var data int
	go func() {
		data++
	}()

	if data == 0 {
		fmt.Printf("the value is %v.\n", data)
	}
}
