package packets

import (
	"io"
	"sync"
)

var _pingrespPacketPool = sync.Pool{
	New: func() interface{} {
		return &PingrespPacket{FixedHeader: &FixedHeader{MessageType: Pingresp}}
	},
}

// PingrespPacket is an internal representation of the fields of the
// Pingresp MQTT packet
type PingrespPacket struct {
	*FixedHeader
}

// NewPingrespPacket return the ping response packet
func NewPingrespPacket() *PingrespPacket {
	return _pingrespPacketPool.Get().(*PingrespPacket)
}

// Type return the packet type
func (pr *PingrespPacket) Type() byte {
	return pr.FixedHeader.MessageType
}

// Reset will initialize the fields in control packet
func (pr *PingrespPacket) Reset() {
	pr.FixedHeader.Dup = false
	pr.FixedHeader.QoS = byte(0)
	pr.FixedHeader.RemainingLength = 0
	pr.FixedHeader.Retain = false
}

// SetTraceID will set traceid for tracing
func (pr *PingrespPacket) SetTraceID(string) {}

// Verify packet availability
func (pr *PingrespPacket) Verify() bool { return true }

// SetFixedHeader will set fh for our header
func (pr *PingrespPacket) SetFixedHeader(fh *FixedHeader) {
	pr.FixedHeader = fh
}

// Close reset the packet field put the control packet back to pool
func (pr *PingrespPacket) Close() {
	pr.Reset()
	_pingrespPacketPool.Put(pr)
}

// String export the packet of connack info
func (pr *PingrespPacket) String() string {
	return pr.FixedHeader.String()
}

// Write will write the packets mostly into a net.Conn
func (pr *PingrespPacket) Write(w io.Writer) error {
	b := Getbuf()
	defer Putbuf(b)
	pr.FixedHeader.pack(b.b[:5])
	_, err := w.Write(b.b[3:5])
	return err
}

// Unpack decodes the details of a ControlPacket after the fixed
// header has been read
func (pr *PingrespPacket) Unpack([]byte) error {
	return nil
}

// Details returns a Details struct containing the QoS and
// MessageID of this ControlPacket
func (pr *PingrespPacket) Details() Details {
	return Details{QoS: 0, MessageID: 0}
}
