package tcpserver

import(
	"gokvstore/common/gotcp" 
	"net"
	"gokvstore/common/logging"

	_ "gokvstore/core"
	"time"
	_ "reflect"

	"encoding/json"
	"gokvstore/core"


)

var PACK_HEAD_LEN=uint32(8)

// 需要实现gotcp.ConnCallback接口
type ConnCallbackImp struct {
	// OnConnect is called when the connection was accepted,
	// If the return value of false is closed
	//OnConnect(*Conn) bool

	// OnMessage is called when the connection receives a packet,
	// If the return value of false is closed
	//OnMessage(*Conn, Packet) bool

	// OnClose is called when the connection closed
	//OnClose(*Conn)
	//a
	//a string
}

func (callback *ConnCallbackImp) OnConnect(conn *gotcp.Conn) bool{
	logging.Debug("OnConnect new conn remote addr:%s",
	conn.GetRawConn().RemoteAddr().String())
	logging.Debug("%s", conn.GetRawConn().RemoteAddr().String())
	return true
}
func (callback *ConnCallbackImp) OnMessage(c *gotcp.Conn, p gotcp.Packet) bool{
	//sendPack:=&PacketImp{
	//	Data:       p.Serialize(),
	//}
	//c.WritePacket(sendPack)


	//go func() {
	//	time.Sleep(time.Duration(time.Second)*10)
	//	c.Close()
	//}()

	processMessage(p,c)

	return true
}

func (callback *ConnCallbackImp) OnClose(*gotcp.Conn){
	
}

func processMessage(p gotcp.Packet, c *gotcp.Conn)  {

	pImp,ok:= p.(*PacketImp)
	if ok!=true{
		logging.Error("invalid Packet type")
		return
	}
	logging.Debug("get body:%s len:%d", pImp.CmdInfo, len(pImp.CmdInfo))
	switch pImp.CmdType {
	case core.CMDTYPE_GET:{
		req:=&core.CmdGetReq{}
		err:=json.Unmarshal(pImp.CmdInfo, req)
		//logging.Debug("CmdGetReq: %#v ",req)
		if err!=nil {
			logging.Error("cmd get info json unmarshal failed, info:%v",
				pImp.CmdInfo)
			c.Close()
			return
		}
		reply, err:=core.ProcessCmdGet(req)
		if err!=nil{
			logging.Error("ProcessCmdGet failed, err:%s",err.Error())
			c.Close()
			return
		}
		//logging.Debug("ProcessCmdGet reply:%s",reply)
		protocolImp:=&ProtocolImp{}
		pack,err:=protocolImp.ProducePacket(reply, core.CMDTYPE_GET_ACK)
		if err!=nil{
			logging.Error("ProducePacket failed, err:%s",err.Error())
			c.Close()
			return
		}
		c.WritePacket(pack)
		return

	}
	case core.CMDTYPE_SET:{
		req:=&core.CmdSetReq{}
		err:=json.Unmarshal(pImp.CmdInfo, req)
		//logging.Debug("CmdSetReq: %#v ",req)
		if err!=nil {
			logging.Error("cmd set info json unmarshal failed, info:%v",
				pImp.CmdInfo)
			c.Close()
			return
		}
		reply, err:=core.ProcessCmdSet(req)
		if err!=nil {
			logging.Error("ProcessCmdSet failed, err:%s",err.Error())
			c.Close()
			return
		}
		//logging.Debug("ProcessCmdSet %s",reply)
		protocolImp:=&ProtocolImp{}
		pack,_:=protocolImp.ProducePacket(reply, core.CMDTYPE_SET_ACK)
		c.WritePacket(pack)

		return
	}
	default:
		logging.Error("invalid cmdtype:%d, close socket", pImp.CmdType)
		c.Close()
		return
	}


	//msg:=pImp.GetBody()
	//cmdCommon:=&core.CmdCommonReq{}
	//err:=json.Unmarshal(msg, cmdCommon)
	//if err!=nil {
	//	logging.Error("invalid msg,err:%s %s",err.Error(), msg)
	//	return
	//}
	//logging.Debug("cmdcommon:%s", cmdCommon)
	//logging.Debug("info:%s",cmdCommon.Info)
	//if cmdCommon.Cmd==core.CmdTypeGet{
	//	reply, err:=core.ProcessCmdGet(cmdCommon)
	//	if err!=nil {
	//		logging.Error("ProcessCmdGet failed, cmd:%s",pImp.GetBody())
	//		return
	//	}

		//protocolImp:=&ProtocolImp{}
		//pack,err:=protocolImp.ProducePacket(reply)
		//if err!=nil{
		//	logging.Error("ProducePacket failed, err:%s",err.Error())
		//	return
		//}
		
		//c.WritePacket(pack)
	//	return

	//}
	//if cmdCommon.Cmd==core.CmdTypeSet{

	//	return
	//}

	//logging.Debug("invalid cmdtype:%s",cmdCommon.Cmd)

}


// 需要实现gotcp.Packet接口
type PacketImp struct{
	Data []byte
	conn *gotcp.Conn
	HeadBuffer []byte
	TotalLen uint32
	CmdType uint32
	CmdInfo []byte
}

func (p *PacketImp)Serialize() []byte{
	return p.Data
}

func (p *PacketImp)SetConn(c *gotcp.Conn){
	p.conn=c
}
func (p *PacketImp)GetConn() *gotcp.Conn{
	return p.conn
}

func (p *PacketImp)GetBody() []byte{
	return p.Data[4:]
}



//type Packet interface {
//	Serialize() []byte
//}

// 需要实现gotcp.Protocol接口
type ProtocolImp struct{

}

func (p *ProtocolImp)ProducePacket(bufferSend []byte,cmdType uint32) (gotcp.Packet, error){
	packetImp:=&PacketImp{
		Data:       gotcp.UInt32ToBytesEndian(uint32( len(bufferSend))+PACK_HEAD_LEN),
	}
	cmdTypeSlice:=gotcp.UInt32ToBytesEndian(cmdType)
	packetImp.Data=append(packetImp.Data, cmdTypeSlice...)
	packetImp.Data=append(packetImp.Data, bufferSend...)
	return packetImp,nil

	//buffer:=gotcp.UInt32ToBytes(len(bufferSend))
}

func (p *ProtocolImp)ParsePacket(bufferRecieved []byte) (gotcp.Packet, int32){
	/*
	if len(bufferRecieved) < 7 {
		return nil,gotcp.PackNeedMore
	}
	pack:=&PacketImp{
		Data:       bufferRecieved[0:7],
	}
	return pack, int32(len(pack.Data))

	*/
	if len(bufferRecieved)<=int(PACK_HEAD_LEN){

		return nil, gotcp.PackNeedMore
	}
	headbytes:=bufferRecieved[0:4]

	//logging.Debug("%d %d  %d  %d",headbytes[0],headbytes[1],headbytes[2],headbytes[3])

	length:=gotcp.BytesToUInt32BigEndian(headbytes)
	//logging.Debug("pack len:%d",length)
	if length>uint32( len(bufferRecieved) ) {
		return nil,gotcp.PackNeedMore
	}
	cmdType:=gotcp.BytesToUInt32BigEndian(bufferRecieved[4:PACK_HEAD_LEN])
	pack:=&PacketImp{
		Data:       bufferRecieved[0:length],
		HeadBuffer: bufferRecieved[0:PACK_HEAD_LEN],
		TotalLen:   length,
		CmdType:    cmdType,
		CmdInfo:    bufferRecieved[PACK_HEAD_LEN:],

	}
	return pack, int32(length)

	/*
	HeadBuffer []byte
	TotalLen uint32
	CmdType uint32
	CmdInfo []byte
	*/
}


func (p *ProtocolImp)ReadPacket(conn *net.TCPConn) (gotcp.Packet, error){
	//buffer:=make([]byte, 0, 1000)
	//readlen, err:=conn.Read(buffer)
	//if err!=nil {
	//	return nil,err
	//}
	//logging.Debug("readlen:%d",readlen)
	//pack:=PacketImp{
	//	Data:    buffer,
	//}
	//pack.Data=append(pack.Data, buffer)

	//return &pack,nil
	return nil,nil
}

//func (p *ProtocolImp)ParseMessage(pack gotcp.Packet) string{
//	return string( pack.Serialize()[4:])
//}

func (p *ProtocolImp)Unpack(c *gotcp.Conn, readerChannel chan gotcp.Packet) error{
	logging.Debug("unpack called")
	return  nil
}

func (p *ProtocolImp)GetHeatBeatData() gotcp.Packet{
	logging.Debug("GetHeatBeatData called")
	return  nil
}

//type Protocol interface {
//	ReadPacket(conn *net.TCPConn) (Packet, error)
//	Unpack(c *Conn, readerChannel chan Packet) error
//	GetHeatBeatData() Packet
//}








func StartTcpServer(PortInfo string, SendChanLimit uint32, RecvChanLimit uint32,
					SendTimeoutSec uint32,RecvTimeoutSec uint32)error{
	
	cfg:=&gotcp.Config{
		PacketSendChanLimit: SendChanLimit,
		PacketReceiveChanLimit: RecvChanLimit,
		ReadTimeOut: RecvTimeoutSec,
		WriteTimeOut: SendTimeoutSec,
	}
	cbk:=&ConnCallbackImp{}
	protoc:=&ProtocolImp{}
	//logging.Debug("cfg:%v",cfg)
	svr:=gotcp.NewServer(cfg,cbk,protoc,100,50)
	///netListener,err:=net.TCPListener{}("tcp", PortInfo)
	tcpaddr,err:=net.ResolveTCPAddr("tcp",PortInfo)
	if err!=nil {
		logging.Debug("inalid portinfo:%s",PortInfo)
		return err
	}
	netListener,err:=net.ListenTCP("tcp", tcpaddr)
	if err!=nil{
		logging.Debug("listen failed protinfo:%s err:%s",PortInfo, err.Error())
		return err
	}
	svr.Start(netListener, time.Duration(time.Second))
	return nil
}

//PacketSendChanLimit    uint32 // the limit of packet send channel
//	PacketReceiveChanLimit uint32 // the limit of packet receive channel
//	ReadTimeOut            uint32
//	WriteTimeOut           uint32






























