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

// PublishPacket is an internal representation of the fields of the
// Publish MQTT packet
type PublishPacket struct {
	*FixedHeader
	TopicName string
	MessageID uint16
	Payload   []byte
	TraceID   string
}

// NewPublishPacket return the ping request packet
func NewPublishPacket() *PublishPacket {
	return _publishPacketPool.Get().(*PublishPacket)
}

// SetFixedHeader will set fh for our header
func (p *PublishPacket) SetFixedHeader(fh *FixedHeader) {
	p.FixedHeader = fh
}

// SetTraceID will set traceid for tracing
func (p *PublishPacket) SetTraceID(id string) { p.TraceID = id }

// Verify packet availability
func (p *PublishPacket) Verify() bool { return true }

// Type return the packet type
func (p *PublishPacket) Type() byte {
	return p.FixedHeader.MessageType
}

// Reset will initialize the fields in control packet
func (p *PublishPacket) Reset() {
	p.FixedHeader.Dup = false
	p.FixedHeader.QoS = byte(0)
	p.FixedHeader.RemainingLength = 0
	p.FixedHeader.Retain = false
	p.TopicName = ""
	p.MessageID = 0
	p.Payload = []byte{}
}

// Close reset the packet field put the control packet back to pool
func (p *PublishPacket) Close() {
	p.Reset()
	_publishPacketPool.Put(p)
}

// String export the packet of publish information
func (p *PublishPacket) String() string {
	return fmt.Sprintf("%s topicName: %s MessageID: %d payload: %s traceID: %s", p.FixedHeader,
		p.TopicName, p.MessageID, string(p.Payload), p.TraceID)
}

// Write will write the packets mostly into a net.Conn
func (p *PublishPacket) Write(w io.Writer) (err error) {
	b := Getbuf()
	defer Putbuf(b)

	n := len(p.TopicName) + 2
	pb := make([]byte, n+2)
	if err = encodeString(p.TopicName, pb[:n]); err != nil {
		return err
	}

	if p.QoS > 0 {
		if err = encodeUint16(p.MessageID, pb[n:]); err != nil {
			return err
		}
		n += 2
	}
	p.FixedHeader.RemainingLength = n + len(p.Payload)
	m := p.FixedHeader.pack(b.b[:5])

	if _, err = w.Write(b.b[5-m : 5]); err != nil {
		return err
	}

	if _, err = w.Write(pb[:n]); err != nil {
		return err
	}

	// TODO FIXME XXX should be bufferd?
	// should this be in one Write?
	_, err = w.Write(p.Payload)
	return
}

// Unpack decodes the details of a ControlPacket after the fixed
// header has been read
func (p *PublishPacket) Unpack(b []byte) error {
	var (
		n   int
		err error
	)
	var payloadLength = p.FixedHeader.RemainingLength

	p.TopicName, n, err = decodeString(b)
	if err != nil {
		return err
	}
	if p.QoS > 0 {
		p.MessageID, err = decodeUint16(b[n:])
		n += 2
		payloadLength -= len(p.TopicName) + 4
	} else {
		payloadLength -= len(p.TopicName) + 2
	}

	if err != nil {
		return err
	}

	if payloadLength < 0 {
		return fmt.Errorf("Error upacking publish, payload length < 0")
	}

	p.Payload = make([]byte, payloadLength)
	copy(p.Payload, b[n:])
	return nil
}

// Copy creates a new PublishPacket with the same topic and payload
// but an empty fixed header, useful for when you want to deliver
// a message with different properties such as QoS but the same
// content
// XXX need to check if put back a packet will cause noting wrong
func (p *PublishPacket) Copy() *PublishPacket {
	newP := NewControlPacket(Publish).(*PublishPacket)
	newP.TopicName = p.TopicName
	newP.Payload = p.Payload

	return newP
}

// Details returns a Details struct containing the QoS and
// MessageID of this ControlPacket
func (p *PublishPacket) Details() Details {
	return Details{QoS: p.QoS, MessageID: p.MessageID}
}
