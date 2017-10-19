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

/* ----------------------------
	扔过去就不管了,全部采用异步的方式，同步的方式意义不大
	至于需不需要回复，由
		Action(msg *Message)
	接口的实现自行处理，在msg中会挂上需要回复的actor的指针，
    便于实现过程回扔消息

*/
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

func (a *Actor)StopActor(){
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

















