package main

import "fmt"

type MyImplement struct {}

func (m *MyImplement) String() string  {
	return "hi"
}

func GetStringer() fmt.Stringer  {
	var s *MyImplement = nil
	//var s fmt.Stringer = nil

	if  s == nil {
		fmt.Println("s is nil")  // s 不是接口，所以不像接口那样判断
	} else {
		fmt.Println("s no nil")
	}

	return s
}

func main()  {
	if GetStringer() == nil {  // 接口会判断 类型和值
		fmt.Println("is nil")
	} else {
		fmt.Println("no nil")
	}

	testi()
}






type testI interface {
	test()
}

type valueT int

func (m valueT)test()  {
}

type valueS struct {
	A int
}

func (m valueS)test()  {
}

func (m valueS)selfTest()  {
}

func testi()  {
	var source valueT = 3456
	var i interface{}
	i = source
	source.test()
	//v, ok := i.(valueT)  //true，可以调用 valueT 的方法，如果是结构体，还可以调用属性等
	//v, ok := i.(testI)  //true，只能调用 testI 接口的方法，如果是结构体，不用调用其属性，必须得断言
	v, ok := i.(int)  // false，只能断言到变量定义的类型，不能找到类型的原始类型
	fmt.Println(v, ok)

	sourceS := valueS{
		A: 789,
	}
	i = sourceS
	v2, ok2 := i.(valueS)
	fmt.Println(v2, ok2)
	fmt.Println(v2.A) // 这是可行的
	v2.selfTest()  // 这是可行的

	v3, ok3 := i.(testI)
	fmt.Println(v3, ok3)
	//fmt.Println(v3.A) // 这是不行的
	//v3.selfTest() //这是不行的
	return
}