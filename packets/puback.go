package packets

import (
	"fmt"
	"io"
	"sync"
)

var _pubackPacketPool = sync.Pool{
	New: func() interface{} {
		return &PubackPacket{FixedHeader: &FixedHeader{MessageType: Puback}}
	},
}

// PubackPacket is an internal representation of the fields of the
// Puback MQTT packet
type PubackPacket struct {
	*FixedHeader
	MessageID uint16
	TraceID   string
}

// NewPubackPacket return the puback packet
func NewPubackPacket() *PubackPacket {
	return _pubackPacketPool.Get().(*PubackPacket)
}

// Type return the packet type
func (pa *PubackPacket) Type() byte {
	return pa.FixedHeader.MessageType
}

// Reset will initialize the fields in control packet
func (pa *PubackPacket) Reset() {
	pa.FixedHeader.Dup = false
	pa.FixedHeader.QoS = byte(0)
	pa.FixedHeader.RemainingLength = 0
	pa.FixedHeader.Retain = false
	pa.MessageID = 0
}

// SetTraceID will set traceid for tracing
func (pa *PubackPacket) SetTraceID(id string) { pa.TraceID = id }

// Verify packet availability
func (pa *PubackPacket) Verify() bool { return true }

// SetFixedHeader will set fh for our header
func (pa *PubackPacket) SetFixedHeader(fh *FixedHeader) {
	pa.FixedHeader = fh
}

// Close reset the packet field put the control packet back to pool
func (pa *PubackPacket) Close() {
	pa.Reset()
	_pubackPacketPool.Put(pa)
}

// String export the packet of connack info
func (pa *PubackPacket) String() string {
	return fmt.Sprintf("%s MessageID: %d traceID: %s", pa.FixedHeader, pa.MessageID, pa.TraceID)
}

// Write will write the packets mostly into a net.Conn
func (pa *PubackPacket) Write(w io.Writer) (err error) {
	b := Getbuf()
	defer Putbuf(b)
	pa.FixedHeader.RemainingLength = 2
	pa.FixedHeader.pack(b.b[:5])

	if err = encodeUint16(pa.MessageID, b.b[5:]); err != nil {
		return err
	}
	_, err = w.Write(b.b[3:7])
	return err
}

// Unpack decodes the details of a ControlPacket after the fixed
// header has been read
func (pa *PubackPacket) Unpack(b []byte) (err error) {
	if len(b) < 2 {
		return io.ErrShortBuffer
	}
	pa.MessageID, err = decodeUint16(b)
	return err
}

// Details returns a Details struct containing the QoS and
// MessageID of this ControlPacket
func (pa *PubackPacket) Details() Details {
	return Details{QoS: pa.QoS, MessageID: pa.MessageID}
}
