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

//SubscribePacket is an internal representation of the fields of the
//Subscribe MQTT packet
type SubscribePacket struct {
	*FixedHeader
	MessageID uint16
	Topics    []string
	QoSs      []byte
}

// NewSubscribePacket return the ping request packet
func NewSubscribePacket() *SubscribePacket {
	return _subscribePacketPool.Get().(*SubscribePacket)
}

//SetFixedHeader will set fh for our header
func (s *SubscribePacket) SetFixedHeader(fh *FixedHeader) {
	s.FixedHeader = fh
}

//Type return the packet type
func (s *SubscribePacket) Type() byte {
	return s.FixedHeader.MessageType
}

//Reset will initialize the fields in control packet
func (s *SubscribePacket) Reset() {
	s.FixedHeader.Dup = false
	s.FixedHeader.QoS = byte(1)
	s.FixedHeader.RemainingLength = 0
	s.FixedHeader.Retain = false
	s.MessageID = 0
	s.Topics = []string{}
	s.QoSs = []byte{}
}

//Close reset the packet field put the control packet back to pool
func (s *SubscribePacket) Close() {
	s.Reset()
	_subscribePacketPool.Put(s)
}

func (s *SubscribePacket) String() string {
	return fmt.Sprintf("%s MessageID: %d topics: %s", s.FixedHeader, s.MessageID, s.Topics)
}

func (s *SubscribePacket) Write(w io.Writer) (err error) {
	var n, m int
	b := Getbuf()
	encodeUint16(s.MessageID, b[5:])
	n = 7

	for i, topic := range s.Topics {
		encodeString(topic, b[n:])
		n += len(topic) + 2
		b[n] = s.QoSs[i]
		n++
	}
	s.FixedHeader.RemainingLength = n - 5

	m = s.FixedHeader.pack(b[:5])
	_, err = w.Write(b[5-m : n])
	Putbuf(b)
	return err
}

//Unpack decodes the details of a ControlPacket after the fixed
//header has been read
func (s *SubscribePacket) Unpack(b []byte) error {
	var n, m int
	var topic string
	s.MessageID = decodeUint16(b)
	n = 2
	payloadLength := s.FixedHeader.RemainingLength - 2
	for payloadLength > 0 {
		topic, m = decodeString(b[n:])
		n += m
		s.Topics = append(s.Topics, topic)
		qos := b[n]
		n++
		s.QoSs = append(s.QoSs, qos)

		payloadLength -= m + 1 // m give the offset of decode topic string, plus 1 byte for QoS
	}
	return nil
}

//Details returns a Details struct containing the QoS and
//MessageID of this ControlPacket
func (s *SubscribePacket) Details() Details {
	return Details{QoS: 1, MessageID: s.MessageID}
}
