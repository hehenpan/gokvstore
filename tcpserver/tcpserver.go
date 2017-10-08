package tcpserver

import(
	"gokvstore/common/gotcp" 
	"net"
	"gokvstore/common/logging"


	"time"
)



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
	a string
}

func (callback *ConnCallbackImp) OnConnect(conn *gotcp.Conn) bool{
	logging.Debug("OnConnect new conn remote addr:%s",
	conn.GetRawConn().RemoteAddr().String())
	logging.Debug("%s", conn.GetRawConn().RemoteAddr().String())
	return true
}
func (callback *ConnCallbackImp) OnMessage(c *gotcp.Conn, p gotcp.Packet) bool{
	sendPack:=&PacketImp{
		Data:       p.Serialize(),
	}
	c.WritePacket(sendPack)


	go func() {
		time.Sleep(time.Duration(time.Second)*10)
		c.Close()
	}()

	return true
}

func (callback *ConnCallbackImp) OnClose(*gotcp.Conn){
	
}

// 需要实现gotcp.Packet接口
type PacketImp struct{
	Data []byte
	conn *gotcp.Conn
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


//type Packet interface {
//	Serialize() []byte
//}

// 需要实现gotcp.Protocol接口
type ProtocolImp struct{
	
}

func (p *ProtocolImp)ParsePacket(bufferRecieved []byte) (gotcp.Packet, int32){
	if len(bufferRecieved) < 7 {
		return nil,gotcp.PackNeedMore
	}
	pack:=&PacketImp{
		Data:       bufferRecieved[0:7],
	}
	return pack, int32(len(pack.Data))
}

func (p *ProtocolImp)ReadPacket(conn *net.TCPConn) (gotcp.Packet, error){
	buffer:=make([]byte, 0, 1000)
	readlen, err:=conn.Read(buffer)
	if err!=nil {
		return nil,err
	}
	logging.Debug("readlen:%d",readlen)
	pack:=PacketImp{
		Data:    buffer,
	}
	//pack.Data=append(pack.Data, buffer)

	return &pack,nil
}

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
	logging.Debug("cfg:%v",cfg)
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






























