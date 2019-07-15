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

// SubackPacket is an internal representation of the fields of the
// Suback MQTT packet
type SubackPacket struct {
	*FixedHeader
	MessageID   uint16
	ReturnCodes []byte
	TraceID     string
}

// NewSubackPacket return the ping request packet
func NewSubackPacket() *SubackPacket {
	return _subackPacketPool.Get().(*SubackPacket)
}

// SetFixedHeader will set fh for our header
func (sa *SubackPacket) SetFixedHeader(fh *FixedHeader) {
	sa.FixedHeader = fh
}

// SetTraceID will set traceid for tracing
func (sa *SubackPacket) SetTraceID(id string) { sa.TraceID = id }

// Verify packet availability
func (sa *SubackPacket) Verify() bool { return true }

// Type return the packet type
func (sa *SubackPacket) Type() byte {
	return sa.FixedHeader.MessageType
}

// Reset will initialize the fields in control packet
func (sa *SubackPacket) Reset() {
	sa.FixedHeader.Dup = false
	sa.FixedHeader.QoS = byte(0)
	sa.FixedHeader.RemainingLength = 0
	sa.FixedHeader.Retain = false
	sa.MessageID = 0
	sa.ReturnCodes = []byte{}
}

// Close reset the packet field put the control packet back to pool
func (sa *SubackPacket) Close() {
	sa.Reset()
	_subackPacketPool.Put(sa)
}

func (sa *SubackPacket) String() string {
	return fmt.Sprintf("%s MessageID: %d traceID: %s", sa.FixedHeader, sa.MessageID, sa.TraceID)
}

// Write will write the packets mostly into a net.Conn
func (sa *SubackPacket) Write(w io.Writer) (err error) {
	b := Getbuf()
	defer Putbuf(b)
	if err = encodeUint16(sa.MessageID, b.b[5:]); err != nil {
		return err
	}
	sa.FixedHeader.RemainingLength = 2 + len(sa.ReturnCodes)

	m := sa.FixedHeader.pack(b.b[:5])
	if _, err = w.Write(b.b[5-m : 7]); err != nil {
		return err
	}

	_, err = w.Write(sa.ReturnCodes)
	return err
}

// Unpack decodes the details of a ControlPacket after the fixed
// header has been read
func (sa *SubackPacket) Unpack(b []byte) (err error) {
	sa.MessageID, err = decodeUint16(b)
	if err != nil {
		return err
	}
	sa.ReturnCodes = make([]byte, len(b)-2)
	copy(sa.ReturnCodes, b[2:])
	return err
}

// Details returns a Details struct containing the QoS and
// MessageID of this ControlPacket
func (sa *SubackPacket) Details() Details {
	return Details{QoS: 0, MessageID: sa.MessageID}
}
