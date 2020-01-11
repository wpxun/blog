package main

import (
	"fmt"
	"runtime"
)

func consumer(ch chan int)  {
	data := <-ch
	fmt.Println(data)
}

func main()  {
	ch := make(chan int)

	for {
		var dummy string
		fmt.Scan(&dummy)

		go consumer(ch)

		//启动
		//ch<- 3

		fmt.Println("goroutine:", runtime.NumGoroutine())
	}
}
