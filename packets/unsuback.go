package packets

import (
	"fmt"
	"io"
	"sync"
)

var _unsubackPacketPool = sync.Pool{
	New: func() interface{} {
		return &UnsubackPacket{FixedHeader: &FixedHeader{MessageType: Unsuback}}
	},
}

//UnsubackPacket is an internal representation of the fields of the
//Unsuback MQTT packet
type UnsubackPacket struct {
	*FixedHeader
	MessageID uint16
}

// NewUnsubackPacket return the ping request packet
func NewUnsubackPacket() *UnsubackPacket {
	return _unsubackPacketPool.Get().(*UnsubackPacket)
}

//Type return the packet type
func (ua *UnsubackPacket) Type() byte {
	return ua.FixedHeader.MessageType
}

//Reset will initialize the fields in control packet
func (ua *UnsubackPacket) Reset() {
	ua.FixedHeader.Dup = false
	ua.FixedHeader.QoS = byte(0)
	ua.FixedHeader.RemainingLength = 0
	ua.FixedHeader.Retain = false
	ua.MessageID = 0
}

//SetFixedHeader will set fh for our header
func (ua *UnsubackPacket) SetFixedHeader(fh *FixedHeader) {
	ua.FixedHeader = fh
}

//Close reset the packet field put the control packet back to pool
func (ua *UnsubackPacket) Close() {
	ua.Reset()
	_unsubackPacketPool.Put(ua)
}

func (ua *UnsubackPacket) String() string {
	return fmt.Sprintf("%s MessageID: %d", ua.FixedHeader, ua.MessageID)
}

func (ua *UnsubackPacket) Write(w io.Writer) (err error) {
	b := Getbuf()
	ua.FixedHeader.RemainingLength = 2
	ua.FixedHeader.pack(b[:5])
	encodeUint16(ua.MessageID, b[5:])
	_, err = w.Write(b[3:7])
	Putbuf(b)
	return
}

//Unpack decodes the details of a ControlPacket after the fixed
//header has been read
func (ua *UnsubackPacket) Unpack(b []byte) error {
	ua.MessageID = decodeUint16(b)
	return nil
}

//Details returns a Details struct containing the QoS and
//MessageID of this ControlPacket
func (ua *UnsubackPacket) Details() Details {
	return Details{QoS: 0, MessageID: ua.MessageID}
}
