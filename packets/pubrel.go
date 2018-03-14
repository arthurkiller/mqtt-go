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

//PubrelPacket is an internal representation of the fields of the
//Pubrel MQTT packet
type PubrelPacket struct {
	*FixedHeader
	MessageID uint16
}

// NewPubrelPacket return the ping request packet
func NewPubrelPacket() *PubrelPacket {
	return _pubrelPacketPool.Get().(*PubrelPacket)
}

//SetFixedHeader will set fh for our header
func (pr *PubrelPacket) SetFixedHeader(fh *FixedHeader) {
	pr.FixedHeader = fh
}

//Type return the packet type
func (pr *PubrelPacket) Type() byte {
	return pr.FixedHeader.MessageType
}

//Reset will initialize the fields in control packet
func (pr *PubrelPacket) Reset() {
	pr.FixedHeader.Dup = false
	pr.FixedHeader.QoS = byte(0)
	pr.FixedHeader.RemainingLength = 0
	pr.FixedHeader.Retain = false
	pr.MessageID = 0
}

//Close reset the packet field put the control packet back to pool
func (pr *PubrelPacket) Close() {
	pr.Reset()
	_pubrelPacketPool.Put(pr)
}

func (pr *PubrelPacket) String() string {
	return fmt.Sprintf("%s MessageID: %d", pr.FixedHeader, pr.MessageID)
}

func (pr *PubrelPacket) Write(w io.Writer) (err error) {
	b := _leakyBuf.Get()

	pr.FixedHeader.RemainingLength = 2
	pr.FixedHeader.pack(b[:5])
	encodeUint16(pr.MessageID, b[5:])

	_, err = w.Write(b[3:7])
	_leakyBuf.Put(b)
	return
}

//Unpack decodes the details of a ControlPacket after the fixed
//header has been read
func (pr *PubrelPacket) Unpack(b []byte) error {
	pr.MessageID = decodeUint16(b)
	return nil
}

//Details returns a Details struct containing the QoS and
//MessageID of this ControlPacket
func (pr *PubrelPacket) Details() Details {
	return Details{QoS: pr.QoS, MessageID: pr.MessageID}
}
