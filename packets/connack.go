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

// ConnackPacket is an internal representation of the fields of the
// Connack MQTT packet
type ConnackPacket struct {
	*FixedHeader
	SessionPresent bool
	ReturnCode     byte
	TraceID        string
}

// NewConnackPacket return the connect ack packet
func NewConnackPacket() *ConnackPacket {
	return _connackPacketPool.Get().(*ConnackPacket)
}

// Reset will initialize the fields in control packet
func (ca *ConnackPacket) Reset() {
	ca.FixedHeader.Dup = false
	ca.FixedHeader.QoS = byte(0)
	ca.FixedHeader.RemainingLength = 0
	ca.FixedHeader.Retain = false

	ca.SessionPresent = false
	ca.ReturnCode = byte(0)
}

// Close reset the packet field put the control packet back to pool
func (ca *ConnackPacket) Close() {
	ca.Reset()
	_connackPacketPool.Put(ca)
}

// SetFixedHeader will set fh for our header
func (ca *ConnackPacket) SetFixedHeader(fh *FixedHeader) {
	ca.FixedHeader = fh
}

// SetTraceID will set traceid for tracing
func (ca *ConnackPacket) SetTraceID(id string) { ca.TraceID = id }

// Verify packet availability
func (ca *ConnackPacket) Verify() bool { return true }

// Type return the packet type
func (ca *ConnackPacket) Type() byte {
	return ca.FixedHeader.MessageType
}

// Write will write the packets mostly into a net.Conn
// | ---  Fixed Header  --- |
// | --- SessionPresent --- |
// | ---   ReturnCode   --- |
// Connack is a specify length packet. So we only use one bytes
func (ca *ConnackPacket) Write(w io.Writer) (length int, err error) {
	b := Getbuf()
	defer Putbuf(b)
	ca.FixedHeader.RemainingLength = 2
	ca.FixedHeader.pack(b.b[:5])
	b.b[5], b.b[6] = boolToByte(ca.SessionPresent), ca.ReturnCode
	return w.Write(b.b[3:7])
}

// Unpack decodes the details of a ControlPacket after the fixed
// header has been read
func (ca *ConnackPacket) Unpack(b []byte) error {
	if len(b) < 2 {
		return io.ErrShortBuffer
	}
	ca.SessionPresent = 0x01&b[0] > 0
	ca.ReturnCode = b[1]
	return nil
}

// Details returns a Details struct containing the QoS and
// MessageID of this ControlPacket
func (ca *ConnackPacket) Details() Details {
	return Details{QoS: 0, MessageID: 0}
}

// String export the packet of connack information
func (ca *ConnackPacket) String() string {
	return fmt.Sprintf("%s sessionpresent: %t returncode: %d tracdID: %s", ca.FixedHeader, ca.SessionPresent, ca.ReturnCode, ca.TraceID)
}
