package main

import (
	"gokvstore/common/action"
	"gokvstore/common/logging"
	"time"
	_ "gokvstore/core"
	_"fmt"
	"reflect"
	_"fmt"
	"fmt"
	"gokvstore/persistence"
	"sync"
	"runtime"
	"gokvstore/core"
	"os"
)


func GeneralTest(){
	//TestAction()
	//Test50Shades()
	//TestChannel()
	//TestPersistence()
	//TestMethod()
	//TestInterview()
	TestKvDb()

	os.Exit(0)
}

func TestKvDb(){
	fileinfo,_:=core.NewDbInfo("/root/winshare/winshare/sz/","dump.rdb")
	exist:=fileinfo.IsExist()
	logging.Debug("exist: %v",exist)

	fileinfo.LoadPrefix()
	fileinfo.LoadDbs()

}



func TestChannel(){
	//ch1:=make(chan int, 100)
	var i int
	for i =0; i<10000000; i++ {
		ch:=make(chan int,100)
		ch <- 11
		ch <- 12
		if i%1000==0{
			time.Sleep(time.Duration(time.Second))
			logging.Debug("%d",i)
		}
		close(ch)
		select{
		case item,ok:=<-ch:
			logging.Debug("item:%v  %v",item,ok)

		}
		select{
		case item,ok:=<-ch:
			logging.Debug("item:%v  %v",item,ok)
			return
		}
	}


}

type ActorTest struct{

}

func (a *ActorTest)Action(msg *action.Message){
	logging.Debug("ActorTest Action triggered")

}











func TestAction(){
	actortest:=&ActorTest{}
	actor:=action.NewActor(1000,actortest)
	actor.StartActor()

	go func(){
		time.Sleep(time.Duration(time.Second*10))
		msg:=action.NewMessage("testdata", nil)
		actor.PostMsg(msg)
	}()

	go func(){
		time.Sleep(time.Duration(time.Second*12))
		actor.StopActor()
	}()
}

type info struct {
	result int
	testmap map[string]string
}

func work() (info,error) {
	//return 13,nil
	return info{13,nil},nil
}

type data1 struct {
	num int
	fp float32
	complex complex64
	str string
	char rune
	yes bool
	events <-chan string
	handler interface{}
	ref *byte
	raw [10]byte
}

func Test50Shades(){
	//var three int //error, even though it's assigned 3 on the next line
	//three = 3

	//var data info
	//data.result:=11
	//data.testmap:=make(map[string]string)
	//data:=&info{}

	//data, err := work() //error
	//fmt.Printf("info: %+v\n",data,err.Error())

	//var one int   //error, unused variable
	//two := 2      //error, unused variable
	//var three int //error, even though it's assigned 3 on the next line
	//three = 3

	v1 := data1{}
	v2 := data1{}

	if v1==v2{
		logging.Debug("equal")
	}else {
		logging.Debug("not equal")
	}

	if reflect.DeepEqual(v1,v2){
		logging.Debug("equal")
	}else {
		logging.Debug("not equal")
	}

	dataint1:=1
	dataint2:=2
	logging.Debug("dataint: %v %v",
		dataint1==dataint2,reflect.DeepEqual(dataint1,dataint2))

	datamap1:=make(map[string]string)
	datamap2:=make(map[string]string)
	logging.Debug("datamap:   %v",
		reflect.DeepEqual(datamap1,datamap2))
	datamap1["test1"]="test1"
	logging.Debug("datamap:   %v",
		reflect.DeepEqual(datamap1,datamap2))

	dataslice1:=[]byte{'1'}
	dataslice2:=[]byte{'2'}
	logging.Debug("dataslice: %v", reflect.DeepEqual(dataslice1, dataslice2))

	datastr1:="123"
	datastr2:="124"
	logging.Debug("datastr: %v  %v",
		datastr1==datastr2, reflect.DeepEqual(datastr1,datastr2))


	m := map[string]*field {"x":{"one"}}
	m["x"].name = "two" //error
	//m["x"]
	logging.Debug("%v",m)

	data := []*field{{"one"},{"two"},{"three"}}

	for _,v := range data {
		//v := v
		go v.print()
	}
	ii:=100
	for ;ii<10000; ii++{
		go func(){
			for true{
				time.Sleep(time.Duration(time.Millisecond*5))
			}

	}()
	}
}


type field struct {
	name string
}

func (p *field) print() {
	//logging.Debug(p.name)
	fmt.Println(p.name)
}



func TestPersistence(){
	persistence.NewJournalDb("abc","def")
}


type Object struct{
	Item int
}

func (o *Object) add(a int){
	o.Item=o.Item+a
}

func TestMethod(){

	var o Object
	o.Item=0
	o.add(1)
	logging.Debug("item %d",o.Item)

	var i1 interface{}
	var i2 interface{}
	logging.Debug("%v",i1==i2)
	i1=&Object{
		Item:   10,
	}
	logging.Debug("%v",i1==i2)
	i2=&Object{
		Item:   10,
	}
	logging.Debug("%v",i1==i2)

	logging.Debug("%T",i1)

}

/*----------------------------------------------------*/
func TestInterview(){
	logging.Debug("TestInterview")
	//defer_call()
	//Gofunc()
	//t := Teacher{}
	//logging.Debug("t: %#v, %T  %p",t,t,&t)
	//t.ShowA()
	//ChanPanic()

	//a := 1
	//b := 2
	//defer calc("1", a, calc("10", a, b))
	//a = 0
	//defer calc("2", a, calc("20", a, b))
	//b = 1

	//s := make([]int, 5)
	//s = append(s, 1, 2, 3)
	//fmt.Println(s)
	//ch := make(chan interface{})
	//ch <- 10
	//ch <- 10
	//logging.Debug("chan finish")

	var peo People1 = &Stduent{}
	think := "bitch"
	fmt.Println(peo.Speak(think))

	//peoType:=reflect.TypeOf(peo)

	//logging.Debug("%v  %v",peoType,  peo.(type)  )
	//peoType==Type(People1)


	switch peo.(type) {
	//case int:
	//	fmt.Println("int")
	//case string:
	//	fmt.Println("string")
	//case People1:
	//	fmt.Println("People1")
	case interface{}:
		fmt.Println("interface")
	default:
		fmt.Println("unknown")
	}


	//var x *int = nil
	//Foo(x)

	fmt.Println(x,y,z,k,p)

}
const (
	x = iota
	y
	z = "zz"
	k
	p = iota
)

func Foo(x interface{}) {
	if x == nil {
		fmt.Println("empty interface")
		return
	}
	fmt.Println("non-empty interface")
}


func DeferFunc1(i int) (t int) {
	t = i
	defer func() {
		t += 3
	}()
	return 100
}


func defer_call() {
	defer func() { logging.Debug("打印前") }()
	defer func() { logging.Debug("打印中") }()
	defer func() { logging.Debug("打印后") }()

	//panic("触发异常")
}

func Gofunc(){
	runtime.GOMAXPROCS(1)
	wg := sync.WaitGroup{}
	wg.Add(20)
	for i := 0; i < 10; i++ {
		go func() {
			fmt.Println("ai: ", i)
			wg.Done()
		}()
	}
	for i := 0; i < 10; i++ {
		go func(i int) {
			fmt.Println("bi: ", i)
			wg.Done()
		}(i)
	}
	wg.Wait()
}

type People struct{}

func (p *People) ShowA() {
	logging.Debug("showA, %#v  %T    %p",p,p,p)
	p.ShowB()
}
func (p *People) ShowB() {
	logging.Debug("showB, %#v  %T  %p",p,p,p)
}

type Teacher struct {
	People
}

func (t *Teacher) ShowB() {
	logging.Debug("teacher showB")
}



func ChanPanic(){
	runtime.GOMAXPROCS(1)
	int_chan := make(chan int, 1)
	string_chan := make(chan string, 1)
	int_chan <- 1
	string_chan <- "hello"
	select {
	case value := <-int_chan:
		fmt.Println(value)
	case value := <-string_chan:
		panic(value)
	}
	logging.Debug("ChanPanic finish")
}

func calc(index string, a, b int) int {
	ret := a + b
	fmt.Println(index, a, b, ret)
	return ret
}

type People1 interface {
	Speak(string) string
}

type Stduent struct{}

func (stu *Stduent) Speak(think string) (talk string) {
	if think == "bitch" {
		talk = "You are a good boy"
	} else {
		talk = "hi"
	}
	return
}




