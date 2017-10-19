package core

import (
	"sync"

	"encoding/json"
	"errors"
	"gokvstore/common/logging"
	"reflect"
)

var CMDTYPE_GET                 = uint32(0x00)
var CMDTYPE_GET_ACK             = uint32(0x01)
var CMDTYPE_SET                 = uint32(0x02)
var CMDTYPE_SET_ACK             = uint32(0x03)

var ERR_CODE_OK                 = uint32(0)
var ERR_CODE_INVALID_VALUE_TYPE = uint32(1)
var ERR_CODE_NOT_EXIST          = uint32(2)


var ERR_MSG_EMPTY               = "ok"
var ERR_MSG_INVALID_VALUE_TYPE  = "invalid type"
var ERR_MSG_NOT_EXIST           = "not exist"
//type CmdCommonReq struct {
//	Cmd   string `json:"cmd,omitempty"`
	//Info  interface{}  `json:"info,omitempty"`
//	Info  interface{}  `json:"info"`
//}

//----------------- GET ----------------------------
type CmdGetReq struct {
	//Key   string `json:"key,omitempty"`
	Key   string `json:"key"`
}
type CmdGetRsp struct{
	Value   string `json:"value"`
	Ok      uint32 `json:"ok"`
	Err     string `json:"err"`
}
//---------------- SET ------------------------------
type CmdSetReq struct{
	Key     string `json:"key"`
	Value   string `json:"value"`
}
type CmdSetRsp struct {
	Ok      uint32 `json:"ok"`
	Err     string `json:"err,omitempty"`
}

var CoreMemMap *CoreMem

//var CmdTypeGet = "get"
//var CmdTypeSet = "set"

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

/*
返回值：
result: -1       // 不存在
result: -2       // 存在，但value的类型不是string
result: n>=0     // 存在
*/
func (c *CoreMem)Get(key string) (string,int){
	c.lock.Lock()
	defer func() {c.lock.Unlock()}()
	v,ok:=c.storeMap[key]
	if ok {
		if reflect.TypeOf(v).String()==reflect.String.String() {
			myv:=v.(string)
			return myv,len(myv)
		}
		return "", -2
	}
	return "",-1
}

func (c *CoreMem)Set(key string, v interface{})  {
	c.lock.Lock()
	defer func() {c.lock.Unlock()}()

	delete(c.storeMap, key)
	c.storeMap[key]=v
}

func ProcessCmdGet(cmdGetReq *CmdGetReq) ([]byte, error){
	reply:=&CmdGetRsp{
		Ok:         ERR_CODE_OK,
		Err:        ERR_MSG_EMPTY,
	}
	var result int
	reply.Value,result=CoreMemMap.Get(cmdGetReq.Key)
	if result==-2 {
		reply.Ok=ERR_CODE_INVALID_VALUE_TYPE
		reply.Err=ERR_MSG_INVALID_VALUE_TYPE

	}else if result==-1 {
		reply.Ok=ERR_CODE_NOT_EXIST
		reply.Err=ERR_MSG_NOT_EXIST
	}else {
		reply.Ok=ERR_CODE_OK
		reply.Err=ERR_MSG_EMPTY
	}
	buf,err:=json.Marshal(reply)
	if err!=nil {
		return nil,errors.New("json marshal error, err:"+err.Error())
	}
	return buf, nil
}

func ProcessCmdSet(cmdSetReq *CmdSetReq)([]byte, error){
	reply:=&CmdSetRsp{
		Ok:     ERR_CODE_OK,
		Err:    ERR_MSG_EMPTY,
	}
	//buf1,err:=json.Marshal(reply)
	//logging.Debug("CmdSetRsp %s  %v",buf1,reply)
	//return buf1,nil
	reply.Ok=ERR_CODE_OK
	reply.Err=ERR_MSG_EMPTY
	if len(cmdSetReq.Key)==0 {
		logging.Debug("len(cmdSetReq.Key)==0")
		buffer,err:=json.Marshal(reply)
		return buffer,err
	}
	CoreMemMap.Set(cmdSetReq.Key, cmdSetReq.Value)
	buf,err:=json.Marshal(reply)
	if err!=nil {
		logging.Error("CmdSetRsp json marshal failed")
	}
	//logging.Debug("CmdSetRsp %s  %v",buf,reply)
	return buf,err
}






































