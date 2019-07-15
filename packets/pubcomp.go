package packets

import (
	"fmt"
	"io"
	"sync"
)

var _pubcompPacketPool = sync.Pool{
	New: func() interface{} {
		return &PubcompPacket{FixedHeader: &FixedHeader{MessageType: Pubcomp}}
	},
}

// PubcompPacket is an internal representation of the fields of the
// Pubcomp MQTT packet
type PubcompPacket struct {
	*FixedHeader
	MessageID uint16
	TraceID   string
}

// NewPubcompPacket return the ping request packet
func NewPubcompPacket() *PubcompPacket {
	return _pubcompPacketPool.Get().(*PubcompPacket)
}

// SetFixedHeader will set fh for our header
func (pc *PubcompPacket) SetFixedHeader(fh *FixedHeader) {
	pc.FixedHeader = fh
}

// SetTraceID will set traceid for tracing
func (pc *PubcompPacket) SetTraceID(id string) { pc.TraceID = id }

// Verify packet availability
func (pc *PubcompPacket) Verify() bool { return true }

// Type return the packet type
func (pc *PubcompPacket) Type() byte {
	return pc.FixedHeader.MessageType
}

// Reset will initialize the fields in control packet
func (pc *PubcompPacket) Reset() {
	pc.FixedHeader.Dup = false
	pc.FixedHeader.QoS = byte(0)
	pc.FixedHeader.RemainingLength = 0
	pc.FixedHeader.Retain = false
	pc.MessageID = 0
}

// Close reset the packet field put the control packet back to pool
func (pc *PubcompPacket) Close() {
	pc.Reset()
	_pubcompPacketPool.Put(pc)
}

// String export the packet of connack info
func (pc *PubcompPacket) String() string {
	return fmt.Sprintf("%s MessageID: %d traceID: %s", pc.FixedHeader, pc.MessageID, pc.TraceID)
}

// Write will write the packets mostly into a net.Conn
func (pc *PubcompPacket) Write(w io.Writer) (err error) {
	b := Getbuf()
	defer Putbuf(b)
	pc.FixedHeader.RemainingLength = 2
	pc.FixedHeader.pack(b.b[:5])
	if err = encodeUint16(pc.MessageID, b.b[5:]); err != nil {
		return err
	}
	_, err = w.Write(b.b[3:7])
	return err
}

// Unpack decodes the details of a ControlPacket after the fixed
// header has been read
func (pc *PubcompPacket) Unpack(b []byte) (err error) {
	if len(b) < 2 {
		return io.ErrShortBuffer
	}
	pc.MessageID, err = decodeUint16(b)
	return err
}

// Details returns a Details struct containing the QoS and
// MessageID of this ControlPacket
func (pc *PubcompPacket) Details() Details {
	return Details{QoS: pc.QoS, MessageID: pc.MessageID}
}
