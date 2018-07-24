package packets

import (
	"fmt"
	"io"
	"sync"
)

var _pubrecPacketPool = sync.Pool{
	New: func() interface{} {
		return &PubrecPacket{FixedHeader: &FixedHeader{MessageType: Pubrec}}
	},
}

//PubrecPacket is an internal representation of the fields of the
//Pubrec MQTT packet
type PubrecPacket struct {
	*FixedHeader
	MessageID uint16
}

// NewPubrecPacket return the ping request packet
func NewPubrecPacket() *PubrecPacket {
	return _pubrecPacketPool.Get().(*PubrecPacket)
}

//SetFixedHeader will set fh for our header
func (pr *PubrecPacket) SetFixedHeader(fh *FixedHeader) {
	pr.FixedHeader = fh
}

//Type return the packet type
func (pr *PubrecPacket) Type() byte {
	return pr.FixedHeader.MessageType
}

//Reset will initialize the fields in control packet
func (pr *PubrecPacket) Reset() {
	pr.FixedHeader.Dup = false
	pr.FixedHeader.QoS = byte(0)
	pr.FixedHeader.RemainingLength = 0
	pr.FixedHeader.Retain = false
	pr.MessageID = 0
}

//Close reset the packet field put the control packet back to pool
func (pr *PubrecPacket) Close() {
	pr.Reset()
	_pubrecPacketPool.Put(pr)
}

func (pr *PubrecPacket) String() string {
	return fmt.Sprintf("%s MessageID: %d", pr.FixedHeader, pr.MessageID)
}

func (pr *PubrecPacket) Write(w io.Writer) (err error) {
	b := Getbuf()
	pr.FixedHeader.RemainingLength = 2
	pr.FixedHeader.pack(b[:5])

	encodeUint16(pr.MessageID, b[5:])
	_, err = w.Write(b[3:7])
	Putbuf(b)
	return
}

//Unpack decodes the details of a ControlPacket after the fixed
//header has been read
func (pr *PubrecPacket) Unpack(b []byte) error {
	pr.MessageID = decodeUint16(b)

	return nil
}

//Details returns a Details struct containing the QoS and
//MessageID of this ControlPacket
func (pr *PubrecPacket) Details() Details {
	return Details{QoS: pr.QoS, MessageID: pr.MessageID}
}
