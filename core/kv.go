package core

import (
	"sync"


)

type CmdCommonReq struct {
	Cmd   string `json:"cmd,omitempty"`
	//Info  interface{}  `json:"info,omitempty"`
	Info  interface{}  `json:"info"`
}

type CmdGetReq struct {
	Key   string `json:"key,omitempty"`
}
type CmdGetRsp struct{
	Value   string `json:"value,omitempty"`
}


var CoreMemMap *CoreMem

var CmdTypeGet = "get"
var CmdTypeSet = "set"

type CoreMem struct {
	storeMap map[string]interface{}
	lock sync.RWMutex
}


func NewCoreMem() *CoreMem{

	coremem:=&CoreMem{
		storeMap:        make(map[string]interface{}),
	}


	return coremem
}

func (c *CoreMem)Get(key string) interface{}{
	c.lock.Lock()
	defer func() {c.lock.Unlock()}()
	v,ok:=c.storeMap[key]
	if ok {
		return v
	}
	return nil
}

func (c *CoreMem)Set(key string, v interface{})  {
	c.lock.Lock()
	defer func() {c.lock.Unlock()}()

	delete(c.storeMap, key)
	c.storeMap[key]=v
}

func ProcessCmdGet(cmdCommon *CmdCommonReq) ([]byte, error){
	return make([]byte,0,0), nil
}






































