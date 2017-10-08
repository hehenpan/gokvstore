package gotcp

import (
	"net"
)

type Packet interface {
	Serialize() []byte
	SetConn(c *Conn)
	GetConn() *Conn
}

type Protocol interface {
	ReadPacket(conn *net.TCPConn) (Packet, error)
	Unpack(c *Conn, readerChannel chan Packet) error
	GetHeatBeatData() Packet

	/*added by hehenpan
	Description: use to parse a whole package on tcp socket stream,
	this method must be implemented, and the parse result must be
	returned. the parse result value indicate the following status:
	-1: data error
	0:  need more data to compose the whole package
	n:  n>0 one package recv finished, the len of the package is n

	*/

	ParsePacket(bufferRecieved []byte) (Packet, int32)
}
