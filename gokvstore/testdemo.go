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
)


func GeneralTest(){
	//TestAction()
	//Test50Shades()
	//TestChannel()
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















