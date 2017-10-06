package udp

import (
	"errors"
	"fmt"
	"net"
	"runtime"
	"sync"
	"time"
)

type UDP_HANDLE_FUNC func(socket *net.UDPConn, remoteAddr *net.UDPAddr, data []byte)

//函数对象
type UdpHandleFuncObj struct {
	handleFunc UDP_HANDLE_FUNC
}

func (o *UdpHandleFuncObj) HandleMsg(socket *net.UDPConn, remoteAddr *net.UDPAddr, data []byte) {
	if o.handleFunc != nil {
		o.handleFunc(socket, remoteAddr, data)
	}
}

//创建函数对象
func NewUdpHandleFuncObj(_func UDP_HANDLE_FUNC) *UdpHandleFuncObj {
	if _func != nil {
		return &UdpHandleFuncObj{
			handleFunc: _func,
		}
	}
	return nil
}

//使用方需要实现的接口
type Udp_Handle_Imp interface {
	HandleMsg(socket *net.UDPConn, remoteAddr *net.UDPAddr, data []byte)
}

type UdpPackage struct {
	Data       []byte
	RemoteAddr *net.UDPAddr
}

type UdpServer struct {
	Addr        *net.UDPAddr
	Conn        *net.UDPConn
	packageChan chan *UdpPackage
	handleObj   Udp_Handle_Imp
	waitGroup   *sync.WaitGroup
	bStop       bool
}

func NewUdpServer(addrstr string, handleObj Udp_Handle_Imp) (*UdpServer, error) {
	if handleObj == nil {
		fmt.Println("callback is nil")
		return nil, errors.New("bad callback")
	}

	addr, err_a := net.ResolveUDPAddr("udp4", addrstr)
	if err_a != nil {
		fmt.Println("bad addr", err_a)
		return nil, err_a
	}

	socket, err := net.ListenUDP("udp4", addr)

	if err != nil {
		fmt.Println("监听失败", err)
		return nil, err
	}
	s := &UdpServer{
		Addr:        addr,
		Conn:        socket,
		packageChan: make(chan *UdpPackage, 100000),
		handleObj:   handleObj,
		waitGroup:   &sync.WaitGroup{},
	}
	return s, nil
}

func (s *UdpServer) Start() {
	if s.Addr == nil || s.Conn == nil || s.handleObj == nil {
		fmt.Println("UdpServer Start failed,not init")
		return
	}

	cpuNum := runtime.NumCPU()
	runtime.GOMAXPROCS(cpuNum)

	//多个接收协程，多个处理协程
	for i := 0; i < cpuNum; i++ {
		go s.handleLoop()
		go s.receiveLoop()
	}
}

func (s *UdpServer) Stop() {
	s.bStop = true
	if s.Conn != nil {
		s.Conn.Close()
	}
	s.waitGroup.Wait()
}

func (s *UdpServer) receiveLoop() {
	s.waitGroup.Add(1)
	for !s.bStop {
		data := make([]byte, 10000)
		n, remoteAddr, err := s.Conn.ReadFromUDP(data)

		if err != nil || remoteAddr == nil {
			fmt.Println("receiveLoop ReadFromUDP failed", err)
			continue
		}
		_package := &UdpPackage{
			Data:       data[:n],
			RemoteAddr: remoteAddr,
		}

		s.packageChan <- _package
	}
	s.waitGroup.Done()
	fmt.Println("receiveLoop exit")
}

func (s *UdpServer) handleLoop() {
	s.waitGroup.Add(1)
	timeout := time.NewTicker(time.Second)
	for !s.bStop {
		select {
		case <-timeout.C: //每隔一秒唤醒一次
			continue
		case _package := <-s.packageChan:
			s.handleObj.HandleMsg(s.Conn, _package.RemoteAddr, _package.Data)
		}
	}
	s.waitGroup.Done()
	fmt.Println("handleLoop exit")
}

/*
使用方法示例

func HandleMsg(socket *net.UDPConn, remoteAddr *net.UDPAddr, data []byte) {
	//fmt.Printf("HandleMsg :%s,remoteAddr=%s\n", string(data), remoteAddr)
	_, err := socket.WriteToUDP([]byte("hehe"), remoteAddr)
	if err != nil {
		fmt.Println("write failed", err)
	}
}


func main() {

	handleObj := NewUdpHandleFuncObj(HandleMsg)
	server, _ := NewUdpServer("0.0.0.0:1987", handleObj)
	server.Start()

	ChanShutdown := make(chan os.Signal)
	signal.Notify(ChanShutdown, syscall.SIGINT)
	signal.Notify(ChanShutdown, syscall.SIGTERM)

	<-ChanShutdown
	server.Stop()
}
*/
