package core

import (
	"os"
	"gokvstore/common/logging"
	"errors"
	"io"
	//"bufio"
)

type DbInfo struct{
	Filename string
	Filepath string
	FileWholename string
	FileObj	 *os.File
}

type KVInfo struct{
	ValueType uint8
	Key []byte
	Value []byte
	IsExpire uint8
	ExpireTs uint64
}

func NewDbInfo(filepath,filename  string) (*DbInfo, error){
	dbinfo:=&DbInfo{
		Filename:	filename,
		Filepath:	filepath,
		FileWholename:   filepath+filename,
	}
	exist:=dbinfo.IsExist()
	if exist==false {
		logging.Debug("file not exist")
		return nil, errors.New("file not exist")
	}
	fileobj,err:=os.Open(dbinfo.FileWholename)
	if err!=nil {
		logging.Error("open file err:%s",err.Error())
		return nil, err
	}
	dbinfo.FileObj=fileobj
	return dbinfo,nil
}

func (db *DbInfo)IsExist()bool{
	wholepath:=db.Filepath+"/"+db.Filename
	fileinfo,err:=os.Stat(wholepath)
	if err != nil{
		logging.Error("file not exist, path:%s err:%s",wholepath,err.Error())
		return false
	}
	logging.Debug("fileinfo:%v",fileinfo)
	return true
}

func (db *DbInfo)LoadPrefix()error{
	headRedis:=make([]byte, 5)
	headDbVer:=make([]byte, 4)
	readLen:=0
	var err error
	readLen, err =io.ReadAtLeast(db.FileObj,headRedis, len(headRedis))
	if err!=nil{
		logging.Error("read redis head failed:%s",err.Error())
		return err
	}
	logging.Debug("load head redis:%d  %s",readLen,headRedis)
	readLen, err = io.ReadAtLeast(db.FileObj, headDbVer, len(headDbVer))
	if err != nil{
		logging.Error("read redis ver failed:%s",err.Error())
		return err
	}
	logging.Debug("load ver:%d  %s",readLen, headDbVer)
	return nil
}

func (db *DbInfo)LoadDbs()error{
	selectDbTag:=make([]byte,1)
	readLen:=0
	var err error
	readLen,err=io.ReadAtLeast(db.FileObj, selectDbTag, 1)
	if err != nil {
		logging.Error("read selectdbid failed, err:%s", err.Error())
		return err
	}
	logging.Debug("read selectDbTag len:%d   value:%d", readLen,selectDbTag[0])

	selectDbId := make([]byte,1)
	readLen,err=io.ReadAtLeast(db.FileObj,selectDbId,1)
	if err != nil{
		logging.Error("read selectdbid failed, err:%s",err.Error())
		return err
	}
	logging.Debug("read selectDbId len:%d value:%d",readLen, selectDbId[0])

	return nil


}

func (db *DbInfo)LoadKVItems()(*KVInfo, error){
	kvType:=make([]byte,1)
	readLen:=0
	var err error
	readLen,err=io.ReadAtLeast(db.FileObj, kvType, 1)
	if err!=nil{
		logging.Debug("read kv type failed,:%s",err.Error())
		return nil,err
	}
	logging.Debug("readLen  %d",readLen)

	return nil,nil






}



































