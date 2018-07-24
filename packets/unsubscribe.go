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

//UnsubscribePacket is an internal representation of the fields of the
//Unsubscribe MQTT packet
type UnsubscribePacket struct {
	*FixedHeader
	MessageID uint16
	Topics    []string
}

// NewUnsubscribePacket return the ping request packet
func NewUnsubscribePacket() *UnsubscribePacket {
	return _unsubscribePacketPool.Get().(*UnsubscribePacket)
}

//SetFixedHeader will set fh for our header
func (u *UnsubscribePacket) SetFixedHeader(fh *FixedHeader) {
	u.FixedHeader = fh
}

//Type return the packet type
func (u *UnsubscribePacket) Type() byte {
	return u.FixedHeader.MessageType
}

//Reset will initialize the fields in control packet
func (u *UnsubscribePacket) Reset() {
	u.FixedHeader.Dup = false
	u.FixedHeader.QoS = byte(1)
	u.FixedHeader.RemainingLength = 0
	u.FixedHeader.Retain = false
	u.MessageID = 0
	u.Topics = []string{}
}

//Close reset the packet field put the control packet back to pool
func (u *UnsubscribePacket) Close() {
	u.Reset()
	_unsubscribePacketPool.Put(u)
}

func (u *UnsubscribePacket) String() string {
	return fmt.Sprintf("%s MessageID: %d", u.FixedHeader, u.MessageID)
}

func (u *UnsubscribePacket) Write(w io.Writer) (err error) {
	var n, m int
	b := Getbuf()
	encodeUint16(u.MessageID, b[5:])
	n = 7
	for _, topic := range u.Topics {
		encodeString(topic, b[n:])
		n += len(topic) + 2
	}
	u.FixedHeader.RemainingLength = n - 5

	m = u.FixedHeader.pack(b[:5])
	_, err = w.Write(b[5-m : n])
	Putbuf(b)
	return
}

//Unpack decodes the details of a ControlPacket after the fixed
//header has been read
func (u *UnsubscribePacket) Unpack(b []byte) error {
	var n, m int
	u.MessageID = decodeUint16(b)
	n = 2

	var topic string
	payloadLength := u.FixedHeader.RemainingLength - 2
	for payloadLength > 0 {
		topic, m = decodeString(b[n:])
		n += m
		u.Topics = append(u.Topics, topic)
		payloadLength -= m
	}
	return nil
}

//Details returns a Details struct containing the QoS and
//MessageID of this ControlPacket
func (u *UnsubscribePacket) Details() Details {
	return Details{QoS: 1, MessageID: u.MessageID}
}
