package packets

import (
	"fmt"
	"io"
	"sync"
)

var _publishPacketPool = sync.Pool{
	New: func() interface{} {
		return &PublishPacket{FixedHeader: &FixedHeader{MessageType: Publish}}
	},
}

//PublishPacket is an internal representation of the fields of the
//Publish MQTT packet
type PublishPacket struct {
	*FixedHeader
	TopicName string
	MessageID uint16
	Payload   []byte
}

// NewPublishPacket return the ping request packet
func NewPublishPacket() *PublishPacket {
	return _publishPacketPool.Get().(*PublishPacket)
}

func (p *PublishPacket) SetFixedHeader(fh *FixedHeader) {
	p.FixedHeader = fh
}

func (p *PublishPacket) Type() byte {
	return p.FixedHeader.MessageType
}

func (p *PublishPacket) Reset() {
	p.FixedHeader.Dup = false
	p.FixedHeader.QoS = byte(0)
	p.FixedHeader.RemainingLength = 0
	p.FixedHeader.Retain = false
	p.TopicName = ""
	p.MessageID = 0
	p.Payload = []byte{}
}

func (p *PublishPacket) Close() {
	p.Reset()
	_publishPacketPool.Put(p)
}

func (p *PublishPacket) String() string {
	return fmt.Sprintf("%s topicName: %s MessageID: %d payload: %s", p.FixedHeader,
		p.TopicName, p.MessageID, string(p.Payload))
}

func (p *PublishPacket) Write(w io.Writer) (err error) {
	var n, m int
	b := _leakyBuf.Get()

	n = 5
	encodeString(p.TopicName, b[n:])
	n += len(p.TopicName) + 2
	if p.QoS > 0 {
		encodeUint16(p.MessageID, b[n:])
		n += 2
	}
	p.FixedHeader.RemainingLength = n - 5 + len(p.Payload)

	m = p.FixedHeader.pack(b[:5])

	_, err = w.Write(b[5-m : n])
	_leakyBuf.Put(b)
	if err != nil {
		return
	}
	// TODO FIXME XXX should be bufferd?
	// should this be in one Write?
	_, err = w.Write(p.Payload)
	return
}

//Unpack decodes the details of a ControlPacket after the fixed
//header has been read
func (p *PublishPacket) Unpack(b []byte) (err error) {
	var n int
	var payloadLength = p.FixedHeader.RemainingLength

	p.TopicName, n = decodeString(b)
	if p.QoS > 0 {
		p.MessageID = decodeUint16(b[n:])
		n += 2
		payloadLength -= len(p.TopicName) + 4
	} else {
		payloadLength -= len(p.TopicName) + 2
	}

	if payloadLength < 0 {
		return fmt.Errorf("Error upacking publish, payload length < 0")
	}

	p.Payload = make([]byte, payloadLength)
	copy(p.Payload, b[n:])
	return
}

//Copy creates a new PublishPacket with the same topic and payload
//but an empty fixed header, useful for when you want to deliver
//a message with different properties such as QoS but the same
//content
// XXX need to check if put back a packet will cause noting wrong
func (p *PublishPacket) Copy() *PublishPacket {
	newP := NewControlPacket(Publish).(*PublishPacket)
	newP.TopicName = p.TopicName
	newP.Payload = p.Payload

	return newP
}

//Details returns a Details struct containing the QoS and
//MessageID of this ControlPacket
func (p *PublishPacket) Details() Details {
	return Details{QoS: p.QoS, MessageID: p.MessageID}
}
