package main

import (
	"gokvstore/common/action"
	"gokvstore/common/logging"
	"time"
	_ "gokvstore/core"
)


func GeneralTest(){
	//TestAction()
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
		actor.StopuActor()
	}()
}















