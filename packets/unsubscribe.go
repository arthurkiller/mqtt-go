package packets

import (
	"fmt"
	"io"
	"sync"
)

var _unsubscribePacketPool = sync.Pool{
	New: func() interface{} {
		return &UnsubscribePacket{FixedHeader: &FixedHeader{MessageType: Unsubscribe, QoS: 1}}
	},
}

// UnsubscribePacket is an internal representation of the fields of the
// Unsubscribe MQTT packet
type UnsubscribePacket struct {
	*FixedHeader
	MessageID uint16
	Topics    []string
	TraceID   string
}

// NewUnsubscribePacket return the ping request packet
func NewUnsubscribePacket() *UnsubscribePacket {
	return _unsubscribePacketPool.Get().(*UnsubscribePacket)
}

// SetFixedHeader will set fh for our header
func (u *UnsubscribePacket) SetFixedHeader(fh *FixedHeader) {
	u.FixedHeader = fh
}

// SetTraceID will set traceid for tracing
func (u *UnsubscribePacket) SetTraceID(id string) { u.TraceID = id }

// Verify packet availability
func (u *UnsubscribePacket) Verify() bool { return true }

// Type return the packet type
func (u *UnsubscribePacket) Type() byte {
	return u.FixedHeader.MessageType
}

// Reset will initialize the fields in control packet
func (u *UnsubscribePacket) Reset() {
	u.FixedHeader.Dup = false
	u.FixedHeader.QoS = byte(1)
	u.FixedHeader.RemainingLength = 0
	u.FixedHeader.Retain = false
	u.MessageID = 0
	u.Topics = []string{}
}

// Close reset the packet field put the control packet back to pool
func (u *UnsubscribePacket) Close() {
	u.Reset()
	_unsubscribePacketPool.Put(u)
}

// String export the packet of unsubscribe information
func (u *UnsubscribePacket) String() string {
	return fmt.Sprintf("%s MessageID: %d traceID: %s", u.FixedHeader, u.MessageID, u.TraceID)
}

// Write will write the packets mostly into a net.Conn
func (u *UnsubscribePacket) Write(w io.Writer) (n int, err error) {
	b := Getbuf()
	defer Putbuf(b)
	if err = encodeUint16(u.MessageID, b.b[5:]); err != nil {
		return 0, err
	}

	ub := make([]byte, n)
	for _, topic := range u.Topics {
		ub = append(ub, make([]byte, len(topic)+2)...)
		if err = encodeString(topic, ub[n:]); err != nil {
			return 0, err
		}
		n += len(topic) + 2
	}
	u.FixedHeader.RemainingLength = n + 2

	m := u.FixedHeader.pack(b.b[:5])
	if m, err = w.Write(b.b[5-m : 7]); err != nil {
		return m, err
	}

	n, err = w.Write(ub[:n])
	return n + m, err
}

// Unpack decodes the details of a ControlPacket after the fixed
// header has been read
func (u *UnsubscribePacket) Unpack(b []byte) error {
	var (
		n, m int
		err  error
	)
	u.MessageID, err = decodeUint16(b)
	if err != nil {
		return err
	}
	n = 2

	var topic string
	payloadLength := u.FixedHeader.RemainingLength - 2
	for payloadLength > 0 {
		topic, m, err = decodeString(b[n:])
		if err != nil {
			return err
		}
		n += m
		u.Topics = append(u.Topics, topic)
		payloadLength -= m
	}
	return nil
}

// Details returns a Details struct containing the QoS and
// MessageID of this ControlPacket
func (u *UnsubscribePacket) Details() Details {
	return Details{QoS: 1, MessageID: u.MessageID}
}
