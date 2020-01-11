package main

import (
	"fmt"
	"time"
)


// 通道的数据流向通道一定要注意用 <-<- 而不是单个 <-
func main2()  {
	valueStream := make(chan interface{})
	takeStream := make(chan interface{})
	var oh = 123
	go func() {
		valueStream<- oh
	}()

	go func() {
		//takeStream <-valueStream  // 传递地址，比如 0xc000086000
		takeStream<- <-valueStream  // 传递数据
	}()

	//getValue,ifOk := <-takeStream
	//fmt.Println(getValue, ifOk)

	done := time.After(3e9)

	Loop:
		for {
			select {
			case e := <-takeStream:
				fmt.Println(e)
			case <-done:
				break Loop
			}
		}

	addrSelect()
	fmt.Println("")
	//intSelect()
}


func repeat(done <-chan interface{}, values ...interface{}) <-chan interface{}{
	valueStream := make(chan interface{})
	go func() {
		defer close(valueStream)

		for {
			for _, v := range values {
				select {
				case <-done:
					return
				case valueStream <- v:
				}
			}
		}
	}()

	return valueStream
}

// 特别要注意通道到通道的转换
func take(done <-chan interface{}, valuseStream <-chan interface{}, num int) <-chan interface{}{
	takeStream := make(chan interface{})
	go func() {
		defer close(takeStream)
		for i:=0; i< num ; i++ {
			select {
			case <-done:
				return
			case takeStream <-<- valuseStream:  // 输出数值
			//case takeStream <- valuseStream: // 输出通道地址，后期得把 interface{} 地址断言成为 <-chan interface{}，再输出通道的内容
			}
		}
	}()

	return takeStream
}

func toInt(done <-chan interface{}, valuseStream <-chan interface{}) <-chan int{
	intStream := make(chan int)
	go func() {
		defer close(intStream)
		for v := range valuseStream {
			select {
			case <-done:
				return
			case intStream <- v.(int):
			}
		}
	}()

	return intStream
}

// 输出interface{} 的通道地址，后期得把 interface{} 地址断言成为 <-chan interface{}，再输出通道的内容
func addrSelect()  {
	done := make(chan interface{})
	defer close(done)

	for num := range take(done, repeat(done, 1), 10) {
		fmt.Printf("%v ", num) // 输出通道地址
		fmt.Printf("%v ", <-(num.(<-chan interface{}))) //把通道地址转为通道，再输出通道的内容
	}
}

// 输出的是 interface{} 的 int 数据，转为 int 才能进行运算
func testSelect()  {
	done := make(chan interface{})
	defer close(done)

	for num := range take(done, repeat(done, 1), 10) {
		//fmt.Printf("%v ", num.(int)) // panic: interface conversion: interface {} is <-chan interface {}, not int
		fmt.Printf("%v ", num.(int)+1) // num 是 interface{}，必须先转型
	}
}

// 调用了 toInt
func intSelect()  {
	done := make(chan interface{})
	defer close(done)

	for num := range toInt(done, take(done, repeat(done, 1), 10)) {
		//fmt.Printf("%v ", num.(int)) // panic: interface conversion: interface {} is <-chan interface {}, not int
		fmt.Printf("%v ", num+1) // 采用了 toInt state，已经转型为 int
	}
}

func main()  {
	done := make(chan interface{})
	defer close(done)

	out1, out2 := tee(done, take(done, repeat(done, 1, 2),1))

	for val1 := range out1{
		fmt.Println("out1:%v, out2:%v\n", val1, <-out2)
	}
}

func orDone(done, c <-chan interface{}) <-chan interface{} {
	valStream := make(chan interface{})

	go func() {
		defer close(valStream)

		for {
			select {
			case <-done:
				return
			case v, ok := <-c:
				if ok == false {
					return
				}

				select {
				case valStream <- v:
				case <-done:
				}
			}
		}
	}()

	return valStream
}

func tee(done <-chan interface{}, in <-chan interface{}) (_,_ <-chan interface{}) {
	out1 := make(chan interface{})
	out2 := make(chan interface{})

	go func() {
		defer close(out1)
		defer close(out2)

		for val := range orDone(done, in)  {
			var out1, out2 = out1, out2

			for i := 0; i < 2; i++ {
				select {
				case <-done:
				case out1 <- val:
					out1 = nil
				case out2 <- val:
					out2 = nil
				}
			}
		}
	}()

	return out1, out2
}