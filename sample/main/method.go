//Go 函数本质上都是值传递，但可以传地址值
//是值传递还是地址传递，取决于调用的函数实现，而不是调用者
//比如以下，不是取决于testInt或MyInt，而是取决于 string 的实现方式

package main

import "fmt"

type (
	Stringer interface {
		mv()
	}
	Inter interface {
		mp()
	}
	MyInt int
)

// 这种方式会自动生成一个新的 string 方法
//func (I *MyInt) string(){
//	(*I).string()
//	return
//}

func (I MyInt) mv() {
	fmt.Println("value receiver")
	I += 10
}

func (I *MyInt) mp() {
	fmt.Println("pointer receiver, change it.")
	*I += 10
}

func main() {

	var testInt MyInt = 123

	// 方法 -------------------------------------------------
	//int 方法是指针接收者，每次调用后testInt都会加10
	(testInt).mp()
	(&testInt).mp()
	(*MyInt).mp(&testInt)
	//MyInt.mp(testInt) //不允许这样调用

	//string 方法是值接收者，每次调用后testInt不会变
	(testInt).mv()
	(&testInt).mv()
	(*MyInt).mv(&testInt)
	MyInt.mv(testInt)

	// 接口 -------------------------------------------------
	var testPi Inter = &testInt
	//var testVi Inter = testInt //不允许这样调用
	var testPs Stringer = &testInt
	var testVs Stringer = testInt
	(testPi).mp()
	(testPs).mv()
	(testVs).mv()

	//类型断言 -------------------------------------------------
	//要断言的变量必须是接口，如果不是接口可以用 interface{}() 转换
	if xx, ok := testPi.(*MyInt); ok {
		xx.mp()
	}
	if xx, ok := testPs.(*MyInt); ok {
		xx.mv()
	}
	if xx, ok := testVs.(MyInt); ok {
		xx.mv()
	}

}
