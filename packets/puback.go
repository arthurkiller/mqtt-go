package packets

import (
	"fmt"
	"io"
	"sync"
)

var _pubackPacketPool = sync.Pool{
	New: func() interface{} {
		return &PubackPacket{FixedHeader: &FixedHeader{MessageType: Puback}}
	},
}

//PubackPacket is an internal representation of the fields of the
//Puback MQTT packet
type PubackPacket struct {
	*FixedHeader
	MessageID uint16
}

// NewPubackPacket return the puback packet
func NewPubackPacket() *PubackPacket {
	return _pubackPacketPool.Get().(*PubackPacket)
}

//Type return the packet type
func (pa *PubackPacket) Type() byte {
	return pa.FixedHeader.MessageType
}

//Reset will initialize the fields in control packet
func (pa *PubackPacket) Reset() {
	pa.FixedHeader.Dup = false
	pa.FixedHeader.QoS = byte(0)
	pa.FixedHeader.RemainingLength = 0
	pa.FixedHeader.Retain = false
	pa.MessageID = 0
}

//SetFixedHeader will set fh for our header
func (pa *PubackPacket) SetFixedHeader(fh *FixedHeader) {
	pa.FixedHeader = fh
}

//Close reset the packet field put the control packet back to pool
func (pa *PubackPacket) Close() {
	pa.Reset()
	_pubackPacketPool.Put(pa)
}

func (pa *PubackPacket) String() string {
	return fmt.Sprintf("%s MessageID: %d", pa.FixedHeader, pa.MessageID)
}

func (pa *PubackPacket) Write(w io.Writer) (err error) {
	b := _leakyBuf.Get()
	pa.FixedHeader.RemainingLength = 2
	pa.FixedHeader.pack(b[:5])
	encodeUint16(pa.MessageID, b[5:])
	_, err = w.Write(b[3:7])
	_leakyBuf.Put(b)
	return err
}

//Unpack decodes the details of a ControlPacket after the fixed
//header has been read
func (pa *PubackPacket) Unpack(b []byte) error {
	pa.MessageID = decodeUint16(b)
	return nil
}

//Details returns a Details struct containing the QoS and
//MessageID of this ControlPacket
func (pa *PubackPacket) Details() Details {
	return Details{QoS: pa.QoS, MessageID: pa.MessageID}
}
