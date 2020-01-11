package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main()  {
	teeChannel()
	return

	newRandStream := func() chan int{
		randStream := make(chan int)
		go func() {
			for{
				randStream <- rand.Int()
			}
		}()
		return randStream
	}

	randStream := newRandStream()
	fmt.Println("3 random ints:")
	for i:=1; i<=3 ; i++  {
		fmt.Printf("%d: %d\n", i, <-randStream)
	}

	//close(randStream)

	time.Sleep(10*1e9)
}

func teeChannel()  {

	send := make(chan interface{})

	out1, out2 := teeFun(send)

	time.Sleep(1e9)
	go func() {
		send <- 1
		time.Sleep(3e9)
		close(send)
	}()

	for val1 := range out1 {
		fmt.Printf("out1:%v, out2:%v\n", val1, <-out2)
	}

	time.Sleep(1e9)

}

func teeFun(send <-chan interface{}) (_,_ <-chan interface{}) {
	rev1 := make(chan interface{})
	rev2 := make(chan interface{})

	go func() {
		defer close(rev1)
		defer close(rev2)

		val := <- send

		for i := 0; i < 2 ; i++ {
			var rev1, rev2 = rev1, rev2
			select {
			case rev1 <- val:
				rev1 = nil
			case rev2 <- val:
				rev2 = nil
			}
		}


		//for val := range send  {
		//	for i := 0; i < 2 ; i++ {
		//		var rev1, rev2 = rev1, rev2
		//		select {
		//		case rev1 <- val:
		//			rev1 = nil
		//		case rev2 <- val:
		//			rev2 = nil
		//		}
		//	}
		//}
	}()
	return rev1, rev2
}
