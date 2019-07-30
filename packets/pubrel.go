package packets

import (
	"fmt"
	"io"
	"sync"
)

var _pubrelPacketPool = sync.Pool{
	New: func() interface{} {
		return &PubrelPacket{FixedHeader: &FixedHeader{MessageType: Pubrel, QoS: 1}}
	},
}

// PubrelPacket is an internal representation of the fields of the
// Pubrel MQTT packet
type PubrelPacket struct {
	*FixedHeader
	MessageID uint16
	TraceID   string
}

// NewPubrelPacket return the ping request packet
func NewPubrelPacket() *PubrelPacket {
	return _pubrelPacketPool.Get().(*PubrelPacket)
}

// SetFixedHeader will set fh for our header
func (pr *PubrelPacket) SetFixedHeader(fh *FixedHeader) {
	pr.FixedHeader = fh
}

// SetTraceID will set traceid for tracing
func (pr *PubrelPacket) SetTraceID(id string) { pr.TraceID = id }

// Verify packet availability
func (pr *PubrelPacket) Verify() bool { return true }

// Type return the packet type
func (pr *PubrelPacket) Type() byte {
	return pr.FixedHeader.MessageType
}

// Reset will initialize the fields in control packet
func (pr *PubrelPacket) Reset() {
	pr.FixedHeader.Dup = false
	pr.FixedHeader.QoS = byte(1)
	pr.FixedHeader.RemainingLength = 0
	pr.FixedHeader.Retain = false
	pr.MessageID = 0
}

// Close reset the packet field put the control packet back to pool
func (pr *PubrelPacket) Close() {
	pr.Reset()
	_pubrelPacketPool.Put(pr)
}

// String export the packet of pubrel information
func (pr *PubrelPacket) String() string {
	return fmt.Sprintf("%s MessageID: %d traceID: %s", pr.FixedHeader, pr.MessageID, pr.TraceID)
}

// Write will write the packets mostly into a net.Conn
func (pr *PubrelPacket) Write(w io.Writer) (n int, err error) {
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
func (pr *PubrelPacket) Unpack(b []byte) (err error) {
	if len(b) < 2 {
		return io.ErrShortBuffer
	}
	pr.MessageID, err = decodeUint16(b)
	return err
}

// Details returns a Details struct containing the QoS and
// MessageID of this ControlPacket
func (pr *PubrelPacket) Details() Details {
	return Details{QoS: pr.QoS, MessageID: pr.MessageID}
}
