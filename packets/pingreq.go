package packets

import (
	"io"
	"sync"
)

var _pingreqPacketPool = sync.Pool{
	New: func() interface{} {
		return &PingreqPacket{FixedHeader: &FixedHeader{MessageType: Pingreq}}
	},
}

// PingreqPacket is an internal representation of the fields of the
// Pingreq MQTT packet
type PingreqPacket struct {
	*FixedHeader
}

// NewPingreqPacket return the ping request packet
func NewPingreqPacket() *PingreqPacket {
	return _pingreqPacketPool.Get().(*PingreqPacket)
}

// Reset will initialize the fields in control packet
func (pr *PingreqPacket) Reset() {
	pr.FixedHeader.Dup = false
	pr.FixedHeader.QoS = byte(0)
	pr.FixedHeader.RemainingLength = 0
	pr.FixedHeader.Retain = false
}

// Close reset the packet field put the control packet back to pool
func (pr *PingreqPacket) Close() {
	pr.Reset()
	_pingreqPacketPool.Put(pr)
}

// SetTraceID will set traceid for tracing
func (pr *PingreqPacket) SetTraceID(string) {}

// Verify packet availability
func (pr *PingreqPacket) Verify() bool { return true }

// SetFixedHeader will set fh for our header
func (pr *PingreqPacket) SetFixedHeader(fh *FixedHeader) {
	pr.FixedHeader = fh
}

// Type return the packet type
func (pr *PingreqPacket) Type() byte {
	return pr.FixedHeader.MessageType
}

// String export the packet of ping information
func (pr *PingreqPacket) String() string {
	return pr.FixedHeader.String()
}

// Write will write the packets mostly into a net.Conn
func (pr *PingreqPacket) Write(w io.Writer) (int, error) {
	b := Getbuf()
	defer Putbuf(b)
	pr.FixedHeader.pack(b.b[:5])
	return w.Write(b.b[3:5])
}

// Unpack decodes the details of a ControlPacket after the fixed
// header has been read
func (pr *PingreqPacket) Unpack([]byte) error {
	return nil
}

// Details returns a Details struct containing the QoS and
// MessageID of this ControlPacket
func (pr *PingreqPacket) Details() Details {
	return Details{QoS: 0, MessageID: 0}
}
