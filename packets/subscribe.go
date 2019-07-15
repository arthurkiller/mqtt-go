package packets

import (
	"fmt"
	"io"
	"sync"
)

var _subscribePacketPool = sync.Pool{
	New: func() interface{} {
		return &SubscribePacket{FixedHeader: &FixedHeader{MessageType: Subscribe, QoS: 1}}
	},
}

// SubscribePacket is an internal representation of the fields of the
// Subscribe MQTT packet
type SubscribePacket struct {
	*FixedHeader
	MessageID uint16
	Topics    []string
	QoSs      []byte
	TraceID   string
}

// NewSubscribePacket return the ping request packet
func NewSubscribePacket() *SubscribePacket {
	return _subscribePacketPool.Get().(*SubscribePacket)
}

// SetFixedHeader will set fh for our header
func (s *SubscribePacket) SetFixedHeader(fh *FixedHeader) {
	s.FixedHeader = fh
}

// SetTraceID will set traceid for tracing
func (s *SubscribePacket) SetTraceID(id string) { s.TraceID = id }

// Verify packet availability
func (s *SubscribePacket) Verify() bool { return true }

// Type return the packet type
func (s *SubscribePacket) Type() byte {
	return s.FixedHeader.MessageType
}

// Reset will initialize the fields in control packet
func (s *SubscribePacket) Reset() {
	s.FixedHeader.Dup = false
	s.FixedHeader.QoS = byte(1)
	s.FixedHeader.RemainingLength = 0
	s.FixedHeader.Retain = false
	s.MessageID = 0
	s.Topics = []string{}
	s.QoSs = []byte{}
}

// Close reset the packet field put the control packet back to pool
func (s *SubscribePacket) Close() {
	s.Reset()
	_subscribePacketPool.Put(s)
}

// String export the packet of connack information
func (s *SubscribePacket) String() string {
	return fmt.Sprintf("%s MessageID: %d topics: %s traceID: %s", s.FixedHeader, s.MessageID, s.Topics, s.TraceID)
}

// Write will write the packets mostly into a net.Conn
func (s *SubscribePacket) Write(w io.Writer) (err error) {
	b := Getbuf()
	defer Putbuf(b)
	if err = encodeUint16(s.MessageID, b.b[5:]); err != nil {
		return err
	}

	n := 0
	sb := make([]byte, n)
	for i, topic := range s.Topics {
		sb = append(sb, make([]byte, len(topic)+3)...)
		if err = encodeString(topic, sb[n:]); err != nil {
			return err
		}
		n += len(topic) + 2
		sb[n] = s.QoSs[i]
		n++
	}
	s.FixedHeader.RemainingLength = n + 2

	m := s.FixedHeader.pack(b.b[:5])
	if _, err = w.Write(b.b[5-m : 7]); err != nil {
		return err
	}

	_, err = w.Write(sb[:n])
	return err
}

// Unpack decodes the details of a ControlPacket after the fixed
// header has been read
func (s *SubscribePacket) Unpack(b []byte) error {
	var (
		n, m int
		err  error
	)
	var topic string
	s.MessageID, err = decodeUint16(b)
	if err != nil {
		return err
	}
	n = 2
	payloadLength := s.FixedHeader.RemainingLength - 2
	for payloadLength > 0 {
		// cut off buf n
		if len(b) <= n {
			return io.ErrShortBuffer
		}
		topic, m, err = decodeString(b[n:])
		if err != nil {
			return err
		}
		n += m
		s.Topics = append(s.Topics, topic)

		// get qos 1 byte
		if len(b) <= n {
			return io.ErrShortBuffer
		}
		qos := b[n]
		n++
		s.QoSs = append(s.QoSs, qos)

		payloadLength -= m + 1 // m give the offset of decode topic string, plus 1 byte for QoS
	}
	return nil
}

// Details returns a Details struct containing the QoS and
// MessageID of this ControlPacket
func (s *SubscribePacket) Details() Details {
	return Details{QoS: 1, MessageID: s.MessageID}
}
