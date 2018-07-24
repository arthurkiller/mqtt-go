package packets

import (
	"fmt"
	"io"
	"sync"
)

var _connectPacketPool = sync.Pool{
	New: func() interface{} {
		return &ConnectPacket{FixedHeader: &FixedHeader{MessageType: Connect}}
	},
}

//ConnectPacket is an internal representation of the fields of the
//Connect MQTT packet
type ConnectPacket struct {
	*FixedHeader
	ProtocolName    string
	ProtocolVersion byte
	CleanSession    bool
	WillFlag        bool
	WillQoS         byte
	WillRetain      bool
	UsernameFlag    bool
	PasswordFlag    bool
	// TODO as protocol saying:
	// > The Server MUST validate that the reserved flag in the CONNECT Control Packet is set to zero and
	// > disconnect the Client if it is not zero
	// so if we need to handle the reservedbit?
	ReservedBit byte
	Keepalive   uint16

	ClientIdentifier string
	WillTopic        string
	WillMessage      []byte
	Username         string
	Password         []byte
}

// NewConnectPacket return the connect packet
func NewConnectPacket() *ConnectPacket {
	return _connectPacketPool.Get().(*ConnectPacket)
}

//Reset will initialize the fields in control packet
func (c *ConnectPacket) Reset() {
	c.FixedHeader.Dup = false
	c.FixedHeader.QoS = byte(0)
	c.FixedHeader.RemainingLength = 0
	c.FixedHeader.Retain = false

	c.ProtocolName = ""
	c.ProtocolVersion = byte(0)
	c.CleanSession = false
	c.WillFlag = false
	c.WillQoS = byte(0)
	c.WillRetain = false
	c.UsernameFlag = false
	c.PasswordFlag = false
	c.ReservedBit = byte(0)
	c.Keepalive = 0

	c.ClientIdentifier = ""
	c.WillTopic = ""
	c.WillMessage = []byte{}
	c.Username = ""
	c.Password = []byte{}
}

//Close reset the packet field put the control packet back to pool
func (c *ConnectPacket) Close() {
	c.Reset()
	_connectPacketPool.Put(c)
}

//SetFixedHeader will set fh for our header
func (c *ConnectPacket) SetFixedHeader(fh *FixedHeader) {
	c.FixedHeader = fh
}

//Type return the packet type
func (c *ConnectPacket) Type() byte {
	return c.FixedHeader.MessageType
}

func (c *ConnectPacket) Write(w io.Writer) (err error) {
	var n = 5
	b := Getbuf()

	encodeString(c.ProtocolName, b[n:])
	n += len(c.ProtocolName) + 2 // 8byte

	b[n] = c.ProtocolVersion
	n++

	b[n] = byte(boolToByte(c.CleanSession)<<1 | boolToByte(c.WillFlag)<<2 | c.WillQoS<<3 |
		boolToByte(c.WillRetain)<<5 | boolToByte(c.PasswordFlag)<<6 | boolToByte(c.UsernameFlag)<<7)
	n++

	encodeUint16(c.Keepalive, b[n:])
	n += 2 // 12byte

	encodeString(c.ClientIdentifier, b[n:])
	n += len(c.ClientIdentifier) + 2

	if c.WillFlag {
		encodeString(c.WillTopic, b[n:])
		n += len(c.WillTopic) + 2
		encodeBytes(c.WillMessage, b[n:])
		n += len(c.WillMessage) + 2
	}
	if c.UsernameFlag {
		encodeString(c.Username, b[n:])
		n += len(c.Username) + 2
	}
	if c.PasswordFlag {
		encodeBytes(c.Password, b[n:])
		n += len(c.Password) + 2
	}
	c.FixedHeader.RemainingLength = n - 5

	m := c.FixedHeader.pack(b[:5])
	_, err = w.Write(b[5-m : n])
	Putbuf(b)
	return err
}

//Unpack decodes the details of a ControlPacket after the fixed
//header has been read
func (c *ConnectPacket) Unpack(b []byte) error {
	var n, m int
	c.ProtocolName, n = decodeString(b)
	c.ProtocolVersion = b[n]
	n++

	options := b[n]
	n++
	c.ReservedBit = 1 & options
	c.CleanSession = 1&(options>>1) > 0
	c.WillFlag = 1&(options>>2) > 0
	c.WillQoS = 3 & (options >> 3)
	c.WillRetain = 1&(options>>5) > 0
	c.PasswordFlag = 1&(options>>6) > 0
	c.UsernameFlag = 1&(options>>7) > 0

	c.Keepalive = decodeUint16(b[n:])
	n += 2

	c.ClientIdentifier, m = decodeString(b[n:])
	n += m

	if c.WillFlag {
		c.WillTopic, m = decodeString(b[n:])
		n += m
		c.WillMessage, m = decodeBytes(b[n:])
		n += m
	}
	if c.UsernameFlag {
		c.Username, m = decodeString(b[n:])
		n += m
	}
	if c.PasswordFlag {
		c.Password, _ = decodeBytes(b[n:])
	}

	return nil
}

//Validate performs validation of the fields of a Connect packet
func (c *ConnectPacket) Validate() byte {
	if c.PasswordFlag && !c.UsernameFlag {
		return ErrRefusedBadUsernameOrPassword
	}
	if c.ReservedBit != 0 {
		//Bad reserved bit
		return ErrProtocolViolation
	}
	if (c.ProtocolName == "MQIsdp" && c.ProtocolVersion != 3) || (c.ProtocolName == "MQTT" && c.ProtocolVersion != 4) {
		//Mismatched or unsupported protocol version
		return ErrRefusedBadProtocolVersion
	}
	if c.ProtocolName != "MQIsdp" && c.ProtocolName != "MQTT" {
		//Bad protocol name
		return ErrProtocolViolation
	}
	if len(c.ClientIdentifier) > 65535 || len(c.Username) > 65535 || len(c.Password) > 65535 {
		//Bad size field
		return ErrProtocolViolation
	}
	if len(c.ClientIdentifier) == 0 && !c.CleanSession {
		//Bad client identifier
		return ErrRefusedIDRejected
	}
	return Accepted
}

//Details returns a Details struct containing the QoS and
//MessageID of this ControlPacket
func (c *ConnectPacket) Details() Details {
	return Details{QoS: 0, MessageID: 0}
}

func (c *ConnectPacket) String() string {
	return fmt.Sprintf("%s protocolversion: %d protocolname: %s cleansession: %t willflag: %t WillQoS: %d WillRetain: %t Usernameflag: %t Passwordflag: %t keepalive: %d clientId: %s willtopic: %s willmessage: %s Username: %s Password: %s",
		c.FixedHeader, c.ProtocolVersion, c.ProtocolName, c.CleanSession, c.WillFlag, c.WillQoS, c.WillRetain,
		c.UsernameFlag, c.PasswordFlag, c.Keepalive, c.ClientIdentifier, c.WillTopic, c.WillMessage, c.Username, c.Password)
}
