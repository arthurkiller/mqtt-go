package packets

import (
	"fmt"
	"io"
	"sync"
)

var _pubrecPacketPool = sync.Pool{
	New: func() interface{} {
		return &PubrecPacket{FixedHeader: &FixedHeader{MessageType: Pubrec}}
	},
}

// PubrecPacket is an internal representation of the fields of the
// Pubrec MQTT packet
type PubrecPacket struct {
	*FixedHeader
	MessageID uint16
	TraceID   string
}

// NewPubrecPacket return the ping request packet
func NewPubrecPacket() *PubrecPacket {
	return _pubrecPacketPool.Get().(*PubrecPacket)
}

// SetFixedHeader will set fh for our header
func (pr *PubrecPacket) SetFixedHeader(fh *FixedHeader) {
	pr.FixedHeader = fh
}

// SetTraceID will set traceid for tracing
func (pr *PubrecPacket) SetTraceID(id string) { pr.TraceID = id }

// Verify packet availability
func (pr *PubrecPacket) Verify() bool { return true }

// Type return the packet type
func (pr *PubrecPacket) Type() byte {
	return pr.FixedHeader.MessageType
}

// Reset will initialize the fields in control packet
func (pr *PubrecPacket) Reset() {
	pr.FixedHeader.Dup = false
	pr.FixedHeader.QoS = byte(0)
	pr.FixedHeader.RemainingLength = 0
	pr.FixedHeader.Retain = false
	pr.MessageID = 0
}

// Close reset the packet field put the control packet back to pool
func (pr *PubrecPacket) Close() {
	pr.Reset()
	_pubrecPacketPool.Put(pr)
}

// String export the packet of pubrec information
func (pr *PubrecPacket) String() string {
	return fmt.Sprintf("%s MessageID: %d traceID: %s", pr.FixedHeader, pr.MessageID, pr.TraceID)
}

// Write will write the packets mostly into a net.Conn
func (pr *PubrecPacket) Write(w io.Writer) (n int, err error) {
	b := Getbuf()
	defer Putbuf(b)
	pr.FixedHeader.RemainingLength = 2
	pr.FixedHeader.pack(b.b[:5])
	if err = encodeUint16(pr.MessageID, b.b[5:]); err != nil {
		return 0, err
	}
	return w.Write(b.b[3:7])
}

// Unpack decodes the details of a ControlPacket after the fixed
// header has been read
func (pr *PubrecPacket) Unpack(b []byte) (err error) {
	if len(b) < 2 {
		return io.ErrShortBuffer
	}
	pr.MessageID, err = decodeUint16(b)
	return err
}

// Details returns a Details struct containing the QoS and
// MessageID of this ControlPacket
func (pr *PubrecPacket) Details() Details {
	return Details{QoS: pr.QoS, MessageID: pr.MessageID}
}
