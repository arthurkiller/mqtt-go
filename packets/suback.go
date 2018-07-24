package packets

import (
	"fmt"
	"io"
	"sync"
)

var _subackPacketPool = sync.Pool{
	New: func() interface{} {
		return &SubackPacket{FixedHeader: &FixedHeader{MessageType: Suback}}
	},
}

//SubackPacket is an internal representation of the fields of the
//Suback MQTT packet
type SubackPacket struct {
	*FixedHeader
	MessageID   uint16
	ReturnCodes []byte
}

// NewSubackPacket return the ping request packet
func NewSubackPacket() *SubackPacket {
	return _subackPacketPool.Get().(*SubackPacket)
}

//SetFixedHeader will set fh for our header
func (sa *SubackPacket) SetFixedHeader(fh *FixedHeader) {
	sa.FixedHeader = fh
}

//Type return the packet type
func (sa *SubackPacket) Type() byte {
	return sa.FixedHeader.MessageType
}

//Reset will initialize the fields in control packet
func (sa *SubackPacket) Reset() {
	sa.FixedHeader.Dup = false
	sa.FixedHeader.QoS = byte(0)
	sa.FixedHeader.RemainingLength = 0
	sa.FixedHeader.Retain = false
	sa.MessageID = 0
	sa.ReturnCodes = []byte{}
}

//Close reset the packet field put the control packet back to pool
func (sa *SubackPacket) Close() {
	sa.Reset()
	_subackPacketPool.Put(sa)
}

func (sa *SubackPacket) String() string {
	return fmt.Sprintf("%s MessageID: %d", sa.FixedHeader, sa.MessageID)
}

func (sa *SubackPacket) Write(w io.Writer) (err error) {
	var n, m int
	b := Getbuf()
	encodeUint16(sa.MessageID, b[5:])
	n = 7
	m = copy(b[n:], sa.ReturnCodes)
	n += m
	sa.FixedHeader.RemainingLength = n - 5

	m = sa.FixedHeader.pack(b[:5])
	_, err = w.Write(b[5-m : n])
	Putbuf(b)
	return err
}

//Unpack decodes the details of a ControlPacket after the fixed
//header has been read
func (sa *SubackPacket) Unpack(b []byte) error {
	sa.MessageID = decodeUint16(b)
	sa.ReturnCodes = make([]byte, len(b)-2)
	copy(sa.ReturnCodes, b[2:])
	return nil
}

//Details returns a Details struct containing the QoS and
//MessageID of this ControlPacket
func (sa *SubackPacket) Details() Details {
	return Details{QoS: 0, MessageID: sa.MessageID}
}
