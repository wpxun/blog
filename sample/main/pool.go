package main

import (
	"fmt"
	"sync"
)

type abc struct{
	N int
	M string
}


func main(){

	myPool := &sync.Pool{
		//New: func() interface{} {
		//	fmt.Println("creating new instance.")
		//	return abc{0,"abc"}
		//},
	}

	myPool.Put(abc{0,"abc"})
	myPool.Put(abc{0,"abc"})
	myPool.Put(abc{0,"abc"})
	myPool.Put(abc{0,"abc"})
	myPool.Put(abc{0,"abc"})

	g1 := myPool.Get().(abc)
	g2 := myPool.Get().(abc)
	g3 := myPool.Get().(abc)
	g4 := myPool.Get().(abc)
	g5 := myPool.Get().(abc)
	g1.N = 1
	g2.N = 2
	g3.N = 3
	g4.N = 4
	//g5.N = 5
	myPool.Put(g1)
	myPool.Put(g2)
	myPool.Put(g3)
	myPool.Put(g5)
	myPool.Put(g4)

	fmt.Println(myPool.Get())
	fmt.Println(myPool.Get())
	fmt.Println(myPool.Get())
}
