package packets

import (
	"fmt"
	"io"
	"sync"
)

var _pingrespPacketPool = sync.Pool{
	New: func() interface{} {
		return &PingrespPacket{FixedHeader: &FixedHeader{MessageType: Pingresp}}
	},
}

//PingrespPacket is an internal representation of the fields of the
//Pingresp MQTT packet
type PingrespPacket struct {
	*FixedHeader
}

// NewPingrespPacket return the ping response packet
func NewPingrespPacket() *PingrespPacket {
	return _pingrespPacketPool.Get().(*PingrespPacket)
}

//Type return the packet type
func (pr *PingrespPacket) Type() byte {
	return pr.FixedHeader.MessageType
}

//Reset will initialize the fields in control packet
func (pr *PingrespPacket) Reset() {
	pr.FixedHeader.Dup = false
	pr.FixedHeader.QoS = byte(0)
	pr.FixedHeader.RemainingLength = 0
	pr.FixedHeader.Retain = false
}

//SetFixedHeader will set fh for our header
func (pr *PingrespPacket) SetFixedHeader(fh *FixedHeader) {
	pr.FixedHeader = fh
}

//Close reset the packet field put the control packet back to pool
func (pr *PingrespPacket) Close() {
	pr.Reset()
	_pingrespPacketPool.Put(pr)
}

func (pr *PingrespPacket) String() string {
	str := fmt.Sprintf("%s", pr.FixedHeader)
	return str
}

func (pr *PingrespPacket) Write(w io.Writer) error {
	b := Getbuf()
	pr.FixedHeader.pack(b[:5])
	_, err := w.Write(b[3:5])
	Putbuf(b)
	return err
}

//Unpack decodes the details of a ControlPacket after the fixed
//header has been read
func (pr *PingrespPacket) Unpack([]byte) error {
	return nil
}

//Details returns a Details struct containing the QoS and
//MessageID of this ControlPacket
func (pr *PingrespPacket) Details() Details {
	return Details{QoS: 0, MessageID: 0}
}
