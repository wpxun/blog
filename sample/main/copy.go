// 接口可以接受值或者指针,指针接口作用不大
// 函数是传值的，传值就是拷贝，结构体值拷贝是浅拷贝，浅拷贝复制是开辟一个新的空间，而值完全复制原来的值，如果原来的值有地址，则会相互影响
// 数组不是指针，切片是指向数组的指针
// map 不能对元素的值取地址，整个map可以取址，而数组是可以对元素取址的

package main

import (
	"fmt"
	"time"
)

type Locker interface {
	test()
}

type LPoint struct {
	n int
	m *int
	o []byte
}

func(l *LPoint)test(){
	fmt.Printf("%p,%p\n", l, &l)
	*(l.m) = 170
	l.n = 555
	fmt.Printf("%p,%p\n", l, &l)
}

type LValue struct {
	n int
	m *int
	o []byte
}

func(l LValue)test(){
	fmt.Printf("%p,%p,%p,%v\n", l, &l, l.m,*l.m)
	*l.m = 169
	l.n = 666
	fmt.Printf("%p,%p,%p,%v\n", l, &l, &l.m,*l.m)
}


func change(l Locker){
	fmt.Printf("%p,%p\n", l, &l)
	l.test()
	fmt.Printf("%p,%p\n", l, &l)
}

type LCopy struct {
	n int
	m *int
	o []byte
	q [4]byte
	p map[int]string
}


func main(){
	copyChannel()
	return


	m0 := 158
	m := 159
	var a LPoint = LPoint{888, &m0, []byte("abc")}
	var b LValue = LValue{999, &m, []byte("def")}

	c := &a

	fmt.Printf("%p,%p\n", b, &b)
	change(b)
	fmt.Printf("%p,%p\n", b, &b)

	fmt.Printf("%p,%p\n", c,&c)
	change(c)
	fmt.Printf("%p,%p\n", c,&c)


	fmt.Println(a, b)
	fmt.Printf("%#v, %#v\n", a, b)

}

// 赋值后要注意引用类型，结构体要看里面的属性是否是引用类型
func copyMap()  {
	sr := map[int]string{1:"hh",5:"xx"}
	cp := sr
	cp[5] = "change value"
	fmt.Println(sr, cp)
}

// 相互转换类型
func copyInt()  {
	type my int
	my1 := 45
	var my2 my
	my2 = my(my1)

	my3 := int(my2)

	fmt.Println(my2,my3)
}


func copyStruct(){
	m := 159
	var b LCopy = LCopy{999, &m, []byte("def"), [4]byte{'u','e','r','y'}, map[int]string{1:"hh",5:"xx"}}
	d := b
	d.n = 123
	*d.m = 888
	d.o[1] = 'y'
	d.q[1] = 'm'
	d.p[5] = "xxhh"

	fmt.Println(d,b,m, &m)
	fmt.Printf("%p\n",&b.q[1])
}

// 变量本地化就是新开辟的内存，要注意新开辟的内存存的是地址还是值，如果是地址就一定会影响原来的值
// 新开辟的内存后面还可以再重新赋值，这时就和第一次赋值的内容不相关了
func outCopy()  {
	out := []int{45}
	out2 := 45

	go func() {
		for i:=0; i<1; i++ {
			var out,out2 = out,out2  // 本地out存放了外部out的地址（相关），本地out2存放了外部out的值（不相关）
			out[0] = 78
			out2 = 78
			out = []int{89}  // 本地out存放了新的内容，这里本地out和外部out不相关了
			time.Sleep(3e9)
			fmt.Println(out, out2)
		}
	}()

	time.Sleep(4e9)
	fmt.Println(out, out2)

}

func copyChannel()  {
	outCopy()
	return

	in := make(chan int,0)
	copyIn := in
	fmt.Printf("%p, %p\n", in, copyIn)

	go func() {
		in<-123
		fmt.Println(len(in),len(copyIn))
		close(in)
	}()


	fmt.Println(<-copyIn)

	fmt.Println(len(in),len(copyIn))

	time.Sleep(13e9)
}