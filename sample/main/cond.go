//条件变量，Broadcast 要在 Wait 之后调用才能通知到。
package main

import (
	"fmt"
	"sync"
	"time"
)

func main(){
	subscribe := func(c *sync.Cond, fn func()){
		//var gor sync.WaitGroup
		//gor.Add(1)
		go func() {
			//gor.Done()
			c.L.Lock()
			defer c.L.Unlock()
			fmt.Println("one", time.Now().UnixNano())
			c.Wait()
			fn()
		}()
		//gor.Wait()
	}
	Clicked := sync.NewCond(&sync.Mutex{})

	var wg sync.WaitGroup
	wg.Add(1)
	subscribe(Clicked, func(){
		fmt.Println("tow", time.Now().UnixNano())
		wg.Done()
	})
	time.Sleep(1e7)
	Clicked.Broadcast()
	fmt.Println("thr", time.Now().UnixNano())
	wg.Wait()
}