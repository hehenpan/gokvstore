package action

import (
	"gokvstore/common/logging"
	"sync"
)

type ActorMotion interface {
	Action(msg *Message)
}


type Actor struct{
	queueMsg chan *Message
	queueExit chan uint32
	actorMotion ActorMotion
	runonce sync.Once
}


type Message struct{
	Data interface{}
	replyActor *Actor
	//MsgType int32
}

func NewMessage(data interface{}, replyActor *Actor) *Message {
	msg:=&Message{
		Data:       data,
		replyActor: replyActor,
	}

	return msg
}

func (a *Actor) PostMsg(msg *Message){
	a.queueMsg<-msg
}

func NewActor(queuesize uint32, actorMotion ActorMotion) *Actor{
	actor:=&Actor{
		queueMsg:      make(chan *Message, queuesize),
		queueExit:     make(chan uint32, 10),
		actorMotion:   actorMotion,
	}


	return actor
}

func (a *Actor)StopuActor(){
	a.queueExit<-uint32(0)
}

func (a *Actor)StartActor(){
	a.runonce.Do( func() {
		go a.ActorRoutine()
				})

}

func (a *Actor)ActorRoutine(){
	defer func(){
		close(a.queueExit)
		close(a.queueMsg)
	}()

	for true  {
		select {
			case <-a.queueExit:
				logging.Debug("ActorRoutine exit")
				return
			case msg:=<-a.queueMsg:
				// 处理
				a.actorMotion.Action(msg)
				continue
		}
	}
}

















