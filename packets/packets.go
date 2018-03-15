package packets

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

const (
	// LeakyBufSize give out the byte slice size
	// according to mqtt 3.1, remain leanth should be lower than 256M
	// mostly MQTT packet should always lower than 500B
	// 1k*10k: in this condition, leakybuf will capture 10m RAM or less
	LeakyBufSize = 102400
	// MaxNBuf give out leaky buffer size
	MaxNBuf = 1024
)

var _leakyBuf = NewLeakyBuf(MaxNBuf, LeakyBufSize)

// ErrInternal will occored on internal error, which should not happen.
var ErrInternal = errors.New("Internal error occured")

//ControlPacket defines the interface for structs intended to hold
//decoded MQTT packets, either from being read or before being
//written
type ControlPacket interface {
	//SetFixedHeader will set fh for our header
	SetFixedHeader(*FixedHeader)
	//Type return the packet type
	Type() byte
	//Write the packet into Writer
	Write(io.Writer) error
	//Unpack the given byte and fill in the control packet object fields
	Unpack([]byte) error
	//Reset will initialize the fields in control packet
	Reset()
	//Close reset the packet field put the control packet back to pool
	Close()
	//Details return packet QoS and MessageID
	Details() Details
	String() string
}

//PacketNames maps the constants for each of the MQTT packet types
//to a string representation of their name.
var PacketNames = []string{
	0:  "Reserved",
	1:  "CONNECT",
	2:  "CONNACK",
	3:  "PUBLISH",
	4:  "PUBACK",
	5:  "PUBREC",
	6:  "PUBREL",
	7:  "PUBCOMP",
	8:  "SUBSCRIBE",
	9:  "SUBACK",
	10: "UNSUBSCRIBE",
	11: "UNSUBACK",
	12: "PINGREQ",
	13: "PINGRESP",
	14: "DISCONNECT",
	15: "Reserved",
}

//Below are the constants assigned to each of the MQTT packet types
const (
	Connect     = 1
	Connack     = 2
	Publish     = 3
	Puback      = 4
	Pubrec      = 5
	Pubrel      = 6
	Pubcomp     = 7
	Subscribe   = 8
	Suback      = 9
	Unsubscribe = 10
	Unsuback    = 11
	Pingreq     = 12
	Pingresp    = 13
	Disconnect  = 14
)

//Below are the const definitions for error codes returned by Connect()
const (
	Accepted                        = 0x00
	ErrRefusedBadProtocolVersion    = 0x01
	ErrRefusedIDRejected            = 0x02
	ErrRefusedServerUnavailable     = 0x03
	ErrRefusedBadUsernameOrPassword = 0x04
	ErrRefusedNotAuthorised         = 0x05
	ErrNetworkError                 = 0xFE
	ErrProtocolViolation            = 0xFF
)

//ConnackReturnCodes is a map of the error codes constants for Connect()
//to a string representation of the error
var ConnackReturnCodes = map[uint8]string{
	0:   "Connection Accepted",
	1:   "Connection Refused: Bad Protocol Version",
	2:   "Connection Refused: Client Identifier Rejected",
	3:   "Connection Refused: Server Unavailable",
	4:   "Connection Refused: Username or Password in unknown format",
	5:   "Connection Refused: Not Authorised",
	254: "Connection Error",
	255: "Connection Refused: Protocol Violation",
}

//Failure defined error codes returned by Connect()
const Failure = 0x80

//SubackReturnCodes is a map of the error codes constants for Subscribe()
//to a string representation of the error
var SubackReturnCodes = map[uint8]string{
	0:   "Subscribe succeed, Maximum QoS 0",
	1:   "Subscribe succeed, Maximum QoS 1",
	2:   "Subscribe succeed, Maximum QoS 2",
	128: "Subscribe failed",
}

//ConnErrors is a map of the errors codes constants for Connect()
//to a Go error
var ConnErrors = map[byte]error{
	Accepted:                        nil,
	ErrRefusedBadProtocolVersion:    errors.New("Unnacceptable protocol version"),
	ErrRefusedIDRejected:            errors.New("Identifier rejected"),
	ErrRefusedServerUnavailable:     errors.New("Server Unavailable"),
	ErrRefusedBadUsernameOrPassword: errors.New("Bad user name or password"),
	ErrRefusedNotAuthorised:         errors.New("Not Authorized"),
	ErrNetworkError:                 errors.New("Network Error"),
	ErrProtocolViolation:            errors.New("Protocol Violation"),
}

//ReadPacket takes an instance of an io.Reader (such as net.Conn) and attempts
//to read an MQTT packet from the stream. It returns a ControlPacket
//representing the decoded MQTT packet and an error. One of these returns will
//always be nil, a nil ControlPacket indicating an error occurred.
// [0][1] [2 3 4] ... [m+1] ... [n]
// 	 --- fh --- | vh |  datagram  |
func ReadPacket(r io.Reader) (cp ControlPacket, err error) {
	var fh FixedHeader
	b := _leakyBuf.Get()

	n, err := io.ReadFull(r, b[0:1])
	if err != nil {
		return nil, err
	}

	offset := fh.unpack(r, b)
	// FIXME cause we are on a buffer pool, the

	n, err = io.ReadFull(r, b[offset:offset+fh.RemainingLength])
	if err != nil {
		return nil, err
	}
	if fh.RemainingLength != n {
		return nil, errors.New("Failed to read expected data")
	}

	cp = NewControlPacketWithHeader(&fh)
	if cp == nil {
		return nil, errors.New("Bad data from client")
	}

	err = cp.Unpack(b[offset : offset+fh.RemainingLength])
	_leakyBuf.Put(b)
	return cp, err
}

//ReadPacketLimitSize limit MQTT packet size. Other is same to ReadPacket.
func ReadPacketLimitSize(r io.Reader, maxSize int) (cp ControlPacket, err error) {
	var fh FixedHeader
	b := _leakyBuf.Get()

	n, err := io.ReadFull(r, b[0:0])
	if err != nil {
		return nil, err
	}

	offset := fh.unpack(r, b)
	n, err = io.ReadFull(r, b[offset:offset+fh.RemainingLength])
	if err != nil {
		return nil, err
	}
	if fh.RemainingLength != n {
		return nil, errors.New("Failed to read expected data")
	}

	// FIXME max size only caculate varable header + payload size.
	if fh.RemainingLength > maxSize {
		return nil, fmt.Errorf("ReminingLength(%d) > maxSize(%d)", fh.RemainingLength, maxSize)
	}

	cp = NewControlPacketWithHeader(&fh)
	if cp == nil {
		return nil, errors.New("Bad data from client")
	}

	err = cp.Unpack(b[offset : offset+fh.RemainingLength])
	_leakyBuf.Put(b)
	return cp, err
}

//NewControlPacket is used to create a new ControlPacket of the type specified
//by packetType, this is usually done by reference to the packet type constants
//defined in packets.go. The newly created ControlPacket is empty and a pointer
//is returned.
func NewControlPacket(packetType byte) (cp ControlPacket) {
	switch packetType {
	case Connect:
		cp = NewConnectPacket()
	case Connack:
		cp = NewConnackPacket()
	case Disconnect:
		cp = NewDisconnectPacket()
	case Publish:
		cp = NewPublishPacket()
	case Puback:
		cp = NewPubackPacket()
	case Pubrec:
		cp = NewPubrecPacket()
	case Pubrel:
		cp = NewPubrelPacket()
	case Pubcomp:
		cp = NewPubcompPacket()
	case Subscribe:
		cp = NewSubscribePacket()
	case Suback:
		cp = NewSubackPacket()
	case Unsubscribe:
		cp = NewUnsubscribePacket()
	case Unsuback:
		cp = NewUnsubackPacket()
	case Pingreq:
		cp = NewPingreqPacket()
	case Pingresp:
		cp = NewPingrespPacket()
	default:
		return nil
	}
	return cp
}

//NewControlPacketWithHeader is used to create a new ControlPacket of the type
//specified within the FixedHeader that is passed to the function.
//The newly created ControlPacket is empty and a pointer is returned.
func NewControlPacketWithHeader(fh *FixedHeader) (cp ControlPacket) {
	switch fh.MessageType {
	case Connect:
		cp = NewConnectPacket()
	case Connack:
		cp = NewConnackPacket()
	case Disconnect:
		cp = NewDisconnectPacket()
	case Publish:
		cp = NewPublishPacket()
	case Puback:
		cp = NewPubackPacket()
	case Pubrec:
		cp = NewPubrecPacket()
	case Pubrel:
		cp = NewPubrelPacket()
	case Pubcomp:
		cp = NewPubcompPacket()
	case Subscribe:
		cp = NewSubscribePacket()
	case Suback:
		cp = NewSubackPacket()
	case Unsubscribe:
		cp = NewUnsubscribePacket()
	case Unsuback:
		cp = NewUnsubackPacket()
	case Pingreq:
		cp = NewPingreqPacket()
	case Pingresp:
		cp = NewPingrespPacket()
	default:
		return nil
	}
	cp.SetFixedHeader(fh)
	return cp
}

//Details struct returned by the Details() function called on
//ControlPackets to present details of the QoS and MessageID
//of the ControlPacket
type Details struct {
	QoS       byte
	MessageID uint16
}

//FixedHeader is a struct to hold the decoded information from
//the fixed header of an MQTT ControlPacket
// |  7  6  5  4 |  3  |  2   1  |  0   |
// | MessageType | Dup |   QoS   |Retain|
// | byte        | bool|  byte   | bool |
//
//FixedHeader only contains 5 bytes at most. And it can align to memory boundry
//based on 64bit system and hardware
// | MsgType | dup | qos | retain |
// |        remain length         |
// |           ... ...            |
// |           ... ...            |
// |           ... ...            |
// |    ---    4 bytes     ---    |
//
// FixedHeader memory layout
// |    MsgType    |-----
// |      dup      |    |
// |      qos      |    |
// |    retain     | 8 bytes
// | remain length0|    |
// | remain length1|    |
// | remain length2|    |
// | remain length3|-----
type FixedHeader struct {
	MessageType     byte
	Dup             bool
	QoS             byte
	Retain          bool
	RemainingLength int
}

func (fh *FixedHeader) String() string {
	return fmt.Sprintf("%s: dup: %t qos: %d retain: %t rLength: %d", PacketNames[int(fh.MessageType)],
		fh.Dup, fh.QoS, fh.Retain, fh.RemainingLength)
}

func (fh *FixedHeader) pack(dst []byte) (n int) {
	if len(dst) != 5 {
		// this should never happen or
		// someone write an internal bug
		panic(ErrInternal)
	}
	n = encodeLength(fh.RemainingLength, dst)

	// | Message Type |DUP flg|  QoS  |Retain|
	// |      4       |   1   |   2   |   1  |
	dst[4-n] = byte(fh.MessageType<<4 | boolToByte(fh.Dup)<<3 | fh.QoS<<1 | boolToByte(fh.Retain))
	return n + 1
}

func (fh *FixedHeader) unpack(r io.Reader, b []byte) int {
	var n int
	fh.MessageType = b[0] >> 4
	fh.Dup = (b[0]>>3)&0x01 > 0
	fh.QoS = (b[0] >> 1) & 0x03
	fh.Retain = b[0]&0x01 > 0
	fh.RemainingLength, n = decodeLength(r, b[1:])
	return n + 1
}

// boolToByte wtire a byte into dst
func boolToByte(b bool) byte {
	if b {
		return 1
	}
	return 0
}

// decodeUint16 will only decode first 2 bytes int num.
func decodeUint16(num []byte) uint16 {
	return binary.BigEndian.Uint16(num)
}

// encodeUint16 will encode the num into dst first 2 bytes
func encodeUint16(num uint16, dst []byte) {
	binary.BigEndian.PutUint16(dst[0:], num)
}

// decodeString decode the given bytes and read a string and return remaining slice
func decodeString(b []byte) (string, int) {
	r, n := decodeBytes(b)
	return string(r), n
}

// encodeString into dst, length will be len(field) + 2
func encodeString(field string, dst []byte) {
	encodeBytes([]byte(field), dst)
}

// decodeBytes read length and return bytes and remaining slice
func decodeBytes(b []byte) ([]byte, int) {
	fieldLength := decodeUint16(b[:2])
	return b[2 : 2+fieldLength], 2 + int(fieldLength)
}

// encodeBytes encode field length and field into dst
func encodeBytes(field []byte, dst []byte) {
	binary.BigEndian.PutUint16(dst[:2], uint16(len(field)))
	copy(dst[2:], field)
}

// follow MQTT v3.1 implements
// buf should be a 5 length slice
func encodeLength(length int, buf []byte) int {
	var n int
	// [ 0   1   2   3   4]
	// | flg | len(max 4) |
	for n = 1; n < 5; n++ {
		// FIXME Concern about BigEndian or LittleEndian crossing different platform
		buf[n] = byte(length % 0x80) // 0x80 = 128 = 1000 0000
		length /= 0x80

		// if there are more data to encode, set the top bit of this byte
		if length > 0 {
			buf[n] |= 0x80
		} else {
			break
		}
	}

	// shrink slice to bottom
	for i := n; i > 0 && n != 4; i-- {
		buf[i+4-n] = buf[i]
	}

	return n
}

// follow MQTT v3.1 implements
func decodeLength(r io.Reader, bs []byte) (int, int) {
	var value, n int
	var multiplier = 1
	for n = 0; ; n++ {
		// FIXME handle ERROR?
		io.ReadFull(r, bs[n:n+1])

		// FIXME Concern about BigEndian or LittleEndian crossing different platform
		value += (int(bs[n]&0x7f) * multiplier)
		multiplier *= 128
		if multiplier > 2097152 { // 128*128*128
			// TODO should panic here?
			return value, n + 1
		}
		if (bs[n] & 0x80) == 0 { // 0x80 = 128 = 1000 0000 ; 0x7f = 127 = 0111 1111
			break
		}
	}
	return value, n + 1
}
