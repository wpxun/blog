package main

import "fmt"

func main()  {
	times := finds(50)
	fmt.Println("找到了：", times)
}


func finds(n int) (times int) {

	fmt.Println("1元%t2元%t5元%t10元")
	for i:=0; i <= n; i++ {
		for j:=0; j <= n/2; j++ {
			for z:=0; z <= n/5; z++ {
				for x:=0; x <= n/10; x++ {
					if i*1 + j*2 + z*5 + x*10 == n {
						fmt.Printf("%d\t%d\t%d\t%d\n", i, j, z, x)
						times++
					}
				}
			}
		}
	}

	return
}