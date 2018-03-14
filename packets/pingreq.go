package packets

import (
	"fmt"
	"io"
	"sync"
)

var _pingreqPacketPool = sync.Pool{
	New: func() interface{} {
		return &PingreqPacket{FixedHeader: &FixedHeader{MessageType: Pingreq}}
	},
}

//PingreqPacket is an internal representation of the fields of the
//Pingreq MQTT packet
type PingreqPacket struct {
	*FixedHeader
}

// NewPingreqPacket return the ping request packet
func NewPingreqPacket() *PingreqPacket {
	return _pingreqPacketPool.Get().(*PingreqPacket)
}

func (pr *PingreqPacket) Reset() {
	pr.FixedHeader.Dup = false
	pr.FixedHeader.QoS = byte(0)
	pr.FixedHeader.RemainingLength = 0
	pr.FixedHeader.Retain = false
}

func (pr *PingreqPacket) Close() {
	pr.Reset()
	_pingreqPacketPool.Put(pr)
}

func (pr *PingreqPacket) SetFixedHeader(fh *FixedHeader) {
	pr.FixedHeader = fh
}

func (pr *PingreqPacket) Type() byte {
	return pr.FixedHeader.MessageType
}

func (pr *PingreqPacket) String() string {
	str := fmt.Sprintf("%s", pr.FixedHeader)
	return str
}

func (pr *PingreqPacket) Write(w io.Writer) error {
	b := _leakyBuf.Get()
	pr.FixedHeader.pack(b[:5])
	_, err := w.Write(b[3:5])
	_leakyBuf.Put(b)
	return err
}

//Unpack decodes the details of a ControlPacket after the fixed
//header has been read
func (pr *PingreqPacket) Unpack([]byte) error {
	return nil
}

//Details returns a Details struct containing the QoS and
//MessageID of this ControlPacket
func (pr *PingreqPacket) Details() Details {
	return Details{QoS: 0, MessageID: 0}
}
