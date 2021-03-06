package gotcp

import (
	"bufio"
	"errors"
	"io"

	"net"
	"sync"
	"sync/atomic"
	"time"

	"gokvstore/common/logging"
	//"fmt"
	//"gokvstore/tcpserver"

)

// Error type
var (
	ErrConnClosing   = errors.New("use of closed network connection")
	ErrWriteBlocking = errors.New("write packet was blocking")
	ErrReadBlocking  = errors.New("read packet was blocking")
)

var (
	PackErr         = int32(-1)
	PackNeedMore    = int32(0)
)

const defaultBufferSize = 16 * 1024
const defaultOutputBufferTimeout = 250 * time.Millisecond

// Conn exposes a set of callbacks for the various events that occur on a connection
type Conn struct {
	Owner                interface{}
	srv                  *ConnWraper
	conn                 *net.TCPConn  // the raw connection
	extraData            string        // to save extra data
	closeOnce            sync.Once     // close the conn, once, per instance
	closeFlag            int32         // close flag
	closeChan            chan struct{} // close chanel
	packetSendChan       chan Packet   // packet send chanel
	packetReceiveChan    chan Packet   // packeet receive chanel
	tickTime             int64         //上次心跳时间
	needHeartBeat        bool
	hbSendInterval       int64 //每隔多久发一次心跳，同时检测是否超时
	hbTimeout            int64
	UnCompleteReadBuffer []byte
	Reader               *bufio.Reader
	Writer               *bufio.Writer
	OutputBufferTimeout  time.Duration
	LenBuf               [4]byte
	LenSlice             []byte
	RecvBuffer           []byte
	sync.RWMutex
}

// ConnCallback is an interface of methods that are used as callbacks on a connection
type ConnCallback interface {
	// OnConnect is called when the connection was accepted,
	// If the return value of false is closed
	OnConnect(*Conn) bool

	// OnMessage is called when the connection receives a packet,
	// If the return value of false is closed
	OnMessage(*Conn, Packet) bool

	// OnClose is called when the connection closed
	OnClose(*Conn)
}

// newConn returns a wrapper of raw conn
func newConn(conn *net.TCPConn, srv *ConnWraper) *Conn {
	c := &Conn{
		srv:                 srv,
		conn:                conn,
		closeChan:           make(chan struct{}),
		packetSendChan:      make(chan Packet, srv.config.PacketSendChanLimit),
		packetReceiveChan:   make(chan Packet, srv.config.PacketReceiveChanLimit),
		tickTime:            time.Now().Unix(),
		needHeartBeat:       srv.needHeartBeat,
		hbSendInterval:      srv.hbSendInterval,
		hbTimeout:           srv.hbTimeout,
		Reader:              bufio.NewReaderSize(conn, defaultBufferSize),
		Writer:              bufio.NewWriterSize(conn, defaultBufferSize),
		OutputBufferTimeout: defaultOutputBufferTimeout,
		RecvBuffer:          make([]byte,0,0),
	}
	c.LenSlice = c.LenBuf[:]
	return c
}

func (c *Conn) ResetTick() {
	c.tickTime = time.Now().Unix()
}

// GetExtraData gets the extra data from the Conn
func (c *Conn) SetOwner(o interface{}) {
	c.Owner = o
}

// GetExtraData gets the extra data from the Conn
func (c *Conn) GetExtraData() string {
	return c.extraData
}

// PutExtraData puts the extra data with the Conn
func (c *Conn) PutExtraData(data string) {
	c.extraData = data
}

// GetRawConn returns the raw net.TCPConn from the Conn
func (c *Conn) GetRawConn() *net.TCPConn {
	return c.conn
}

// Close closes the connection
func (c *Conn) Close() {
	c.closeOnce.Do(func() {
		atomic.StoreInt32(&c.closeFlag, 1)
		close(c.closeChan)
		close(c.packetSendChan)
		c.conn.Close()
		c.srv.callback.OnClose(c)
	})
}

// IsClosed indicates whether or not the connection is closed
func (c *Conn) IsClosed() bool {
	return atomic.LoadInt32(&c.closeFlag) == 1
}
/*
// AsyncReadPacket async reads a packet, this method will never block
func (c *Conn) AsyncReadPacket(timeout time.Duration) (Packet, error) {
	if c.IsClosed() {
		return nil, ErrConnClosing
	}

	if timeout == 0 {
		select {
		case p := <-c.packetReceiveChan:
			return p, nil

		default:
			return nil, ErrReadBlocking
		}

	} else {
		select {
		case p := <-c.packetReceiveChan:
			return p, nil

		case <-c.closeChan:
			return nil, ErrConnClosing

		case <-time.After(timeout):
			return nil, ErrReadBlocking
		}
	}
}
*/
// 发送报文
func (c *Conn) WritePacket(p Packet) error{
	if c.IsClosed()==true {
		logging.Debug("conn already closed, drop send")
		return nil
	}
	c.packetSendChan<-p
	return nil
}
/*
//同步发送，异步发送太慢
func (c *Conn) syncWritePacket(p Packet) error {
	if c.IsClosed() {
		return ErrConnClosing
	}
	c.conn.SetWriteDeadline(time.Now().Add(time.Second * 20))
	packetstr := p.Serialize()

	c.Lock()
	_, err := c.Writer.Write(packetstr)
	c.Writer.Flush()
	c.Unlock()
	if err != nil {
		logging.Error("con  SyncWritePacket write found a error: %v", err)
		return err
	}
	return nil
}
*/
/*
// AsyncWritePacket async writes a packet, this method will never block
func (c *Conn) asyncWritePacket(p Packet, timeout time.Duration) error {
	if c.IsClosed() {
		return ErrConnClosing
	}

	c.conn.SetWriteDeadline(time.Now().Add(time.Second * 20))
	packetstr := p.Serialize()

	//写到缓冲区而已
	c.Lock()
	_, err := c.Writer.Write(packetstr)
	c.Unlock()

	if err != nil {
		logging.Error("con  AsyncWritePacket write found a error: %v", err)
		return err
	}

	return nil
	/*
		if timeout == 0 {
			select {
			case c.packetSendChan <- p:
				return nil

				//default:
				//return ErrWriteBlocking
			}

		} else {
			select {
			case c.packetSendChan <- p:
				return nil

			case <-c.closeChan:
				return ErrConnClosing

			case <-time.After(timeout):
				return ErrWriteBlocking
			}
		}
	*/
//}

// Do it
func (c *Conn) Do() {
	if !c.srv.callback.OnConnect(c) {
		return
	}
	//c.conn.SetDeadline(time.Now().Add(time.Second * 30))
	go c.handleLoop()
	go c.readStickPackLoop()
	//go c.readLoop()
	//go c.writeStickPacketLoop()
	//go c.heartbeatLoop()
	go c.writeLoop()
}

func (c *Conn) readStickPackLoop() {
	c.srv.waitGroup.Add(1)
	defer func() {
		//recover()
		c.Close()
		c.srv.waitGroup.Done()
	}()

	//reader := bufio.NewReader(c.conn)


	for {

		select {
		case <-c.srv.exitChan:
			return

		case <-c.closeChan:
			return

		default:
		}
		c.conn.SetReadDeadline(time.Now().Add(time.Second * 20))

		//err := c.srv.protocol.Unpack(c, c.packetReceiveChan)

		//if e, ok := err.(net.Error); ok && e.Timeout() {
			//l4g.Info("con read found a timeout error, i can do")

		//	continue
			// This was a timeout
		//}
		//if err != nil {
		//	if err == io.EOF {
		//		logging.Info("close by peer")
		//		return
		//	}
		//	logging.Info("con read found a error: %v", err)
		//	return
		//}


			buffer := make([]byte, 1024*2)
			readlen, err := c.Reader.Read(buffer)
			if e, ok := err.(net.Error); ok && e.Timeout() {

				logging.Debug("read timeout")
				continue
				// This was a timeout
			}
			if err != nil {
				if err == io.EOF {
					logging.Debug("close by peer")
					return
				}
				logging.Info("conn read found a error: %s", err.Error())
				return
			}

			if readlen > 0 {
				//logging.Debug("readlen is %d, data:%s", readlen, buffer[0:readlen])
				c.RecvBuffer=append(c.RecvBuffer, buffer[0:readlen] ...)
				//logging.Debug("recvbuffer:%v",c.RecvBuffer)
			}
			for true{
				//logging.Debug("ready to parse packet")
				pack, result:=c.srv.protocol.ParsePacket(c.RecvBuffer)
				if result==-1 {
					logging.Error("invalid data, ready to active close socket")
					return
				}
				if result==0 {
					//logging.Debug("buffer size:%d package need more, continue read",
					//		len(c.RecvBuffer))
					break
				}
				if result>0 {
					//logging.Debug("fetch one packege, len:%d c.recvbuffer:%d",
					//			result,len(c.RecvBuffer))
					// 这个地方处理一个报文
					pack.SetConn(c)
					c.packetReceiveChan<-pack
					if result==int32(len(c.RecvBuffer)) {
						c.RecvBuffer=make([]byte,0,0)   //
						//logging.Debug("no data left in buffer, continue read")
						break
					}
					c.RecvBuffer=c.RecvBuffer[result:]
					//logging.Debug("buffer left size:%d",len(c.RecvBuffer))
					continue
				}
			}

	}
}
/*
func (c *Conn) readLoop() {
	c.srv.waitGroup.Add(1)
	defer func() {
		//recover()
		c.Close()
		c.srv.waitGroup.Done()
	}()

	for {
		select {
		case <-c.srv.exitChan:
			return

		case <-c.closeChan:
			return

		default:
		}

		p, err := c.srv.protocol.ReadPacket(c.conn)
		if err != nil {
			logging.Info("con ReadPacket found a error: %v", err)
			return
		}

		c.packetReceiveChan <- p
	}
}
*/
func (c *Conn) writeLoop() {
	c.srv.waitGroup.Add(1)
	defer func() {
		//recover()
		c.Close()
		c.srv.waitGroup.Done()
	}()

	for {
		if c.IsClosed()==true {
			return
		}
		select {
		case <-c.srv.exitChan:
			return

		case <-c.closeChan:
			return

		case p := <-c.packetSendChan:
			if c.IsClosed()==true {
				//logging.Debug("socket already closed, drop send in writeLoop")
				return
			}
			if _, err := c.conn.Write(p.Serialize()); err != nil {
				logging.Info("con write found a error: %v", err)
				return
			}else{
				c.Writer.Flush()
				//logging.Debug("packetSendChan get one packet, write success")

			}
		}
	}
}
/*
func (c *Conn) writeStickPacketLoop() {
	c.srv.waitGroup.Add(1)
	defer func() {
		//recover()
		c.Close()
		c.srv.waitGroup.Done()
	}()

	outputBufferTicker := time.NewTicker(c.OutputBufferTimeout)
	for {
		select {
		case <-c.srv.exitChan:
			return

		case <-c.closeChan:
			return
		case <-outputBufferTicker.C: //每隔一段时间写入对端
			c.Lock()
			err := c.Writer.Flush()
			c.Unlock()
			if err != nil {
				logging.Error("conn writeStickPacketLoop fFlush failed,err=%s", err.Error())
				return
			}

			/*
				case p := <-c.packetSendChan:
					c.conn.SetWriteDeadline(time.Now().Add(time.Second * 180))
					packetstr := p.Serialize()

					//写到缓冲区而已
					c.Lock()
					_, err := c.Writer.Write(packetstr)
					c.Unlock()
					if e, ok := err.(net.Error); ok && e.Timeout() {
						//l4g.Info("con read found a timeout error, i can do")
						c.packetSendChan <- p //写回去
						continue
					}
					// This was a timeout

					if err != nil {
						logging.Info("con write found a error: %v", err)
						return
					}
			*/
/*		}
	}
}
*/
/*
func (c *Conn) heartbeatLoop() {
	c.srv.waitGroup.Add(1)
	defer func() {
		//recover()
		c.srv.waitGroup.Done()
	}()

	if !c.needHeartBeat {
		return
	}

	timercheck := time.NewTicker(time.Duration(c.hbSendInterval) * time.Second)
	for {
		select {
		case <-c.srv.exitChan:
			return

		case <-c.closeChan:
			return
		case <-timercheck.C:
			curTime := time.Now().Unix()
			//logging.Debug("conn %s timeout,curTime=%d,tickTime=%d,hbTimeout=%d", c.GetExtraData(), int(curTime), int(c.tickTime), c.hbTimeout)
			if curTime >= c.tickTime+c.hbTimeout {
				logging.Error("conn %s timeout,curTime=%d,tickTime=%d,hbTimeout=%d,", c.GetExtraData(), int(curTime), int(c.tickTime), c.hbTimeout)
				c.Close()
				return
			}
			c.SyncWritePacket(c.srv.protocol.GetHeatBeatData())
		}
	}
}
*/
func (c *Conn) handleLoop() {
	c.srv.waitGroup.Add(1)
	defer func() {
		//recover()
		c.Close()
		c.srv.waitGroup.Done()
	}()

	for {
		if c.IsClosed()==true {
			return
		}
		select {
		case <-c.srv.exitChan:
			return

		case <-c.closeChan:
			return

		case p ,ok:= <-c.packetReceiveChan:
			//logging.Debug("receive msg:%s", string(p.Serialize()))
			if ok==false{
				return
			}
			if !c.srv.callback.OnMessage(c, p) {
				return
			}
		}
	}
}
