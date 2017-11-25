package persistence

import (
	//"errors"
	"os"
	"gokvstore/common/logging"

	"sync"
	"gokvstore/common/gotcp"
	"errors"
	"encoding/json"
)

type JournalItem struct {
	CmdType     uint32
	Info        interface{}
}

type JournalDb struct{
	FileName    string
	FilePath    string
	OffSet      int64
	finfo       *os.FileInfo
	fileobj     *os.File
	runOnce     *sync.Once

}


func NewJournalDb(filename, filepath string)(*JournalDb, error){
	jdb:=&JournalDb{
		FileName:       filename,
		FilePath:       filepath,
		OffSet:         0,
		finfo:          nil,
		fileobj:        nil,
	}
	WholePath:=filepath+"/"+filename
	fileobj,err:=os.Open(WholePath)
	if err!=nil {
		logging.Error("open journal file failed, path;%s",WholePath)
		return nil, err
	}
	jdb.fileobj=fileobj
	fileinfo, err:=fileobj.Stat()
	if err!=nil {
		logging.Error("file stat error %v", err.Error())
		return nil,err
	}
	jdb.finfo=&fileinfo

	Offset,err:=fileobj.Seek(fileinfo.Size(),0)
	if err!=nil {
		logging.Error("Seek failed, path:%s, err:%s",WholePath, err.Error())
		return nil,err
	}
	jdb.OffSet=Offset
	return jdb,nil
}

func (j *JournalDb) Close() {
	if j.fileobj==nil {
		return
	}
	j.runOnce.Do(
		func () {
			j.fileobj.Close()
		})
}

func (j * JournalDb) WriteJournal(msg []byte) error {
	sizeNeedSave:=len(msg)
	OffsetSaved:=0
	for{
		writeSlice:=msg[OffsetSaved:]
		writesize,err:=j.fileobj.Write(writeSlice)
		if err!=nil {
			logging.Error("Write failed, err:%s",err.Error())
			return err
		}
		OffsetSaved=OffsetSaved+writesize
		if OffsetSaved<sizeNeedSave {
			logging.Debug("write journal %d %d %d",sizeNeedSave,writesize, OffsetSaved)
			continue
		}else {
			break
		}
	}
	//length,err:=j.fileobj.Write(msg)
	j.OffSet=j.OffSet+int64(sizeNeedSave)
	return nil
}

func (j *JournalDb) ReadJournal() (*JournalItem, error){
	headSlice:=make([]byte,0)
	readLen, err:=j.fileobj.ReadAt(headSlice, 8)
	if err != nil {
		logging.Error("read Journal failed, err:%v",err.Error())
		return nil ,err
	}
	if readLen!=8 {
		logging.Error("read journal head failed, get headlen:%d",readLen)
		return nil, errors.New("invalid heead, file error")
	}
	lengthSlice:=headSlice[0:4]
	length:=gotcp.BytesToUInt32BigEndian(lengthSlice)
	cmdTypeSlice:=headSlice[4:]
	cmdType:=gotcp.BytesToUInt32BigEndian(cmdTypeSlice)

	infoSlice:=make([]byte,0)
	readLen,err = j.fileobj.ReadAt(infoSlice, int64(length)-int64(8))
	if err!=nil {
		logging.Error("read journal body failed, err:%s", err.Error())
		return nil, err
	}
	if int64(readLen) != int64(length) - int64(8) {
		logging.Error("read journal body failed, file err readlen:%d expect:%d",readLen,
			length-8)
		return nil,errors.New("file error")
	}
	var info interface{}
	err = json.Unmarshal(infoSlice, info)
	if err != nil {
		logging.Error("info json unmarshal failed, err:%s info:%s",err.Error(),
			infoSlice)
		return nil, err
	}

	journalItem:=&JournalItem{
		CmdType:    cmdType,
		Info:       info,
	}
	return journalItem, nil

}























