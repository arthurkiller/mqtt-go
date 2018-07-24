package packets

import (
	"fmt"
	"io"
	"sync"
)

var _connackPacketPool = sync.Pool{
	New: func() interface{} {
		return &ConnackPacket{FixedHeader: &FixedHeader{MessageType: Connack}}
	},
}

//ConnackPacket is an internal representation of the fields of the
//Connack MQTT packet
type ConnackPacket struct {
	*FixedHeader
	SessionPresent bool
	ReturnCode     byte
}

// NewConnackPacket return the connect ack packet
func NewConnackPacket() *ConnackPacket {
	return _connackPacketPool.Get().(*ConnackPacket)
}

//Reset will initialize the fields in control packet
func (ca *ConnackPacket) Reset() {
	ca.FixedHeader.Dup = false
	ca.FixedHeader.QoS = byte(0)
	ca.FixedHeader.RemainingLength = 0
	ca.FixedHeader.Retain = false

	ca.SessionPresent = false
	ca.ReturnCode = byte(0)
}

//Close reset the packet field put the control packet back to pool
func (ca *ConnackPacket) Close() {
	ca.Reset()
	_connackPacketPool.Put(ca)
}

//SetFixedHeader will set fh for our header
func (ca *ConnackPacket) SetFixedHeader(fh *FixedHeader) {
	ca.FixedHeader = fh
}

//Type return the packet type
func (ca *ConnackPacket) Type() byte {
	return ca.FixedHeader.MessageType
}

// Write will write the packets mostly into a net.Conn
// | ---  Fixed Header  --- |
// | --- SessionPresent --- |
// | ---   ReturnCode   --- |
// Connack is a specify length packet. So we only use one bytes
func (ca *ConnackPacket) Write(w io.Writer) (err error) {
	b := Getbuf()
	ca.FixedHeader.RemainingLength = 2
	ca.FixedHeader.pack(b[:5])

	b[5], b[6] = boolToByte(ca.SessionPresent), ca.ReturnCode
	_, err = w.Write(b[3:7])
	Putbuf(b)
	return
}

//Unpack decodes the details of a ControlPacket after the fixed
//header has been read
func (ca *ConnackPacket) Unpack(b []byte) error {
	ca.SessionPresent = 0x01&b[0] > 0
	ca.ReturnCode = b[1]
	return nil
}

//Details returns a Details struct containing the QoS and
//MessageID of this ControlPacket
func (ca *ConnackPacket) Details() Details {
	return Details{QoS: 0, MessageID: 0}
}

func (ca *ConnackPacket) String() string {
	return fmt.Sprintf("%s sessionpresent: %t returncode: %d", ca.FixedHeader, ca.SessionPresent, ca.ReturnCode)
}
