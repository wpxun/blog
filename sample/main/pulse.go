package main

import (
	"fmt"
	"time"
)

func main()  {
	done := make(chan interface{})
	time.AfterFunc(10e9, func() { close(done) })

	const timeout = 2e9
	heartbeat, results := doWork(done, timeout/2)

	for  {
		select {
		case _, ok := <-heartbeat:
			if ok == false {
				return
			}
		case r, ok := <-results:
			if ok == false {
				return
			}
			fmt.Println("ok", r)
		case <- time.After(timeout):
			fmt.Println("呵呵")
			return
		}
	}
}

func doWork(done <-chan interface{}, pulseInterval time.Duration) (<-chan interface{}, <-chan time.Time)  {
	heartbeat := make(chan interface{})
	results := make(chan time.Time)

	go func() {
		pulse := time.Tick(pulseInterval)
		workGen := time.Tick(2*pulseInterval)

		sendPulse := func() {
			select {
			case heartbeat<-struct {}{}:
			default:
			}
		}

		sendResult := func(r time.Time) {
			for {
				select {
				case <-pulse:
					sendPulse()
				case results <- r:
					return
				}
			}
		}

		for i := 0; i < 2 ; i++ {
			select {
			case <-done:
				return
			case <-pulse:
				sendPulse()
			case r := <-workGen:
				sendResult(r)
			}
		}
	}()
	return heartbeat, results
}
