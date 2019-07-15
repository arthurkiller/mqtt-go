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

// UnsubackPacket is an internal representation of the fields of the
// Unsuback MQTT packet
type UnsubackPacket struct {
	*FixedHeader
	MessageID uint16
	TraceID   string
}

// NewUnsubackPacket return the ping request packet
func NewUnsubackPacket() *UnsubackPacket {
	return _unsubackPacketPool.Get().(*UnsubackPacket)
}

// SetTraceID will set traceid for tracing
func (ua *UnsubackPacket) SetTraceID(id string) { ua.TraceID = id }

// Verify packet availability
func (ua *UnsubackPacket) Verify() bool { return true }

// Type return the packet type
func (ua *UnsubackPacket) Type() byte {
	return ua.FixedHeader.MessageType
}

// Reset will initialize the fields in control packet
func (ua *UnsubackPacket) Reset() {
	ua.FixedHeader.Dup = false
	ua.FixedHeader.QoS = byte(0)
	ua.FixedHeader.RemainingLength = 0
	ua.FixedHeader.Retain = false
	ua.MessageID = 0
}

// SetFixedHeader will set fh for our header
func (ua *UnsubackPacket) SetFixedHeader(fh *FixedHeader) {
	ua.FixedHeader = fh
}

// Close reset the packet field put the control packet back to pool
func (ua *UnsubackPacket) Close() {
	ua.Reset()
	_unsubackPacketPool.Put(ua)
}

// String export the packet of unsuback information
func (ua *UnsubackPacket) String() string {
	return fmt.Sprintf("%s MessageID: %d traceID: %s", ua.FixedHeader, ua.MessageID, ua.TraceID)
}

// Write will write the packets mostly into a net.Conn
func (ua *UnsubackPacket) Write(w io.Writer) (err error) {
	b := Getbuf()
	defer Putbuf(b)
	ua.FixedHeader.RemainingLength = 2
	ua.FixedHeader.pack(b.b[:5])
	if err = encodeUint16(ua.MessageID, b.b[5:]); err != nil {
		return err
	}
	_, err = w.Write(b.b[3:7])
	return
}

// Unpack decodes the details of a ControlPacket after the fixed
// header has been read
func (ua *UnsubackPacket) Unpack(b []byte) (err error) {
	if len(b) < 2 {
		return io.ErrShortBuffer
	}
	ua.MessageID, err = decodeUint16(b)
	return err
}

// Details returns a Details struct containing the QoS and
// MessageID of this ControlPacket
func (ua *UnsubackPacket) Details() Details {
	return Details{QoS: 0, MessageID: ua.MessageID}
}
