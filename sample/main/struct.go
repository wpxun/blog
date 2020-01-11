package main

type noName struct {
	A int
	B int
}

func (noName)noTest()  {
	
}

type hasName struct {
	D int
	C int
}

func (hasName)hasTest()  {

}

type compose struct {
	noName
	name hasName

}

func main()  {

	myCompose := new(compose)
	_ = myCompose.A
	myCompose.noTest()
	_ = myCompose.name.C
	_ = myCompose.name.hasTest
	//_ = myCompose.C  // 有名称的不能隐匿调用其属性
	//myCompose.hasTest() // 有名称的也不能隐匿调用其方法
	
}
