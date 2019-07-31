package packets

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

const (
	// BufferSize give out the byte slice size in pool
	bufferSize = 7

	// MaxRemainingLength the max length of remain
	MaxRemainingLength = 268435455 // 0xFF, 0xFF, 0xFF, 0x7F
)

func init() {
	NewBufferPool(bufferSize)
}

var (
	// ErrInternal will occored on internal error, which should not happen.
	ErrInternal = errors.New("internal error occurred")
	// ErrReadPacketLimitation return on packet remain length too large
	ErrReadPacketLimitation = errors.New("packet size exceeded limitation")
	// ErrMQTTPacketLimitation return on the mqtt packet remain length too large
	ErrMQTTPacketLimitation = errors.New("the mqtt packet size exceeded limitation")
)

// func MarshalPacket(b []byte) (p ControlPacket, err error) { return nil, nil }
// func UnmarshalPacket(ControlPacket) ([]byte, error)       { return nil, nil }

// WritePacket write a packet to w
func WritePacket(w io.Writer, p ControlPacket) (int, error) { return p.Write(w) }

// ControlPacket defines the interface for structs intended to hold
// decoded MQTT packets, either from being read or before being
// written
type ControlPacket interface {
	// SetFixedHeader will set fh for our header
	SetFixedHeader(*FixedHeader)
	// SetTraceID set the traceable label into packets. whild the packet was delievering.
	SetTraceID(id string)
	// Verify will check the package if validate
	Verify() bool
	// Type return the packet type
	Type() byte
	// Write the packet into Writer
	Write(io.Writer) (int, error)
	// Unpack the given byte and fill in the control packet object fields
	Unpack([]byte) error
	// Reset will initialize the fields in control packet
	Reset()
	// Close reset the packet field put the control packet back to pool
	Close()
	// Details return packet QoS and MessageID
	Details() Details

	// String export the packet of connack info
	String() string
}

// PacketNames maps the constants for each of the MQTT packet types
// to a string representation of their name.
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

// Below are the constants assigned to each of the MQTT packet types
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

// Below are the const definitions for error codes returned by Connect()
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

// ConnackReturnCodes is a map of the error codes constants for Connect()
// to a string representation of the error
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

// Failure defined error codes returned by Connect()
const Failure = 0x80

// SubackReturnCodes is a map of the error codes constants for Subscribe()
// to a string representation of the error
var SubackReturnCodes = map[uint8]string{
	0:   "Subscribe succeed, Maximum QoS 0",
	1:   "Subscribe succeed, Maximum QoS 1",
	2:   "Subscribe succeed, Maximum QoS 2",
	128: "Subscribe failed",
}

// ConnErrors is a map of the errors codes constants for Connect()
// to a Go error
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

// ReadPacket takes an instance of an io.Reader (such as net.Conn) and attempts
// to read an MQTT packet from the stream. It returns a ControlPacket
// representing the decoded MQTT packet and an error. One of these returns will
// always be nil, a nil ControlPacket indicating an error occurred.
// [0][1] [2 3 4] ... [m+1] ... [n]
// 	 --- fh --- | vh |  datagram  |
func ReadPacket(r io.Reader) (cp ControlPacket, length int, err error) {
	var fh FixedHeader
	b := Getbuf()
	defer Putbuf(b)

	n, err := io.ReadFull(r, b.b[0:1])
	if err != nil {
		return nil, n, err
	}
	if n != 1 {
		return nil, n, io.ErrUnexpectedEOF
	}
	if n = fh.unpack(r, b.b); n <= 1 {
		return nil, n, io.ErrUnexpectedEOF
	}

	// 5 bit the packets header
	length = n

	if fh.RemainingLength > MaxRemainingLength {
		return nil, n, ErrMQTTPacketLimitation
	}

	rb := make([]byte, fh.RemainingLength)
	n, err = io.ReadFull(r, rb)
	// add remain  length
	length += n
	if err != nil {
		return nil, length, err
	}

	if fh.RemainingLength != n {
		return nil, length, errors.New("Failed to read expected data")
	}

	cp = NewControlPacketWithHeader(&fh)
	if cp == nil {
		return nil, length, errors.New("Bad data from client")
	}

	err = cp.Unpack(rb)
	return cp, length, err
}

// ReadPacketLimitSize read a packet with a maximum size of s from net.Conn
func ReadPacketLimitSize(r io.Reader, s int) (cp ControlPacket, length int, err error) {
	var fh FixedHeader
	b := Getbuf()
	defer Putbuf(b)

	n, err := io.ReadFull(r, b.b[0:1])
	if err != nil {
		return nil, n, err
	}
	if n != 1 {
		return nil, n, io.ErrUnexpectedEOF
	}

	if n = fh.unpack(r, b.b); n <= 1 {
		return nil, n, io.ErrUnexpectedEOF
	}

	// 5 bit the header packets
	length = n

	if fh.RemainingLength > MaxRemainingLength {
		return nil, length, ErrMQTTPacketLimitation
	}

	if fh.RemainingLength > s {
		return nil, length, ErrReadPacketLimitation
	}

	rb := make([]byte, fh.RemainingLength)
	n, err = io.ReadFull(r, rb)
	length += n
	if err != nil {
		return nil, length, err
	}

	if fh.RemainingLength != n {
		return nil, length, errors.New("Failed to read expected data")
	}

	cp = NewControlPacketWithHeader(&fh)
	if cp == nil {
		return nil, length, errors.New("Bad data from client")
	}

	err = cp.Unpack(rb)
	return cp, length, err
}

// NewControlPacket is used to create a new ControlPacket of the type specified
// by packetType, this is usually done by reference to the packet type constants
// defined in packets.go. The newly created ControlPacket is empty and a pointer
// is returned.
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

// NewControlPacketWithHeader is used to create a new ControlPacket of the type
// specified within the FixedHeader that is passed to the function.
// The newly created ControlPacket is empty and a pointer is returned.
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

// Details struct returned by the Details() function called on
// ControlPackets to present details of the QoS and MessageID
// of the ControlPacket
type Details struct {
	QoS       byte
	MessageID uint16
}

// FixedHeader is a struct to hold the decoded information from
// the fixed header of an MQTT ControlPacket
//  |  7  6  5  4 |  3  |  2   1  |  0   |
//  | MessageType | Dup |   QoS   |Retain|
//  | byte        | bool|  byte   | bool |
//
// FixedHeader only contains 5 bytes at most. And it can align to memory boundary
// based on 64bit system and hardware
//  | MsgType | dup | qos | retain |
//  |        remain length         |
//  |           ... ...            |
//  |           ... ...            |
//  |           ... ...            |
//  |    ---    4 bytes     ---    |
//
//  FixedHeader memory layout
//  |    MsgType    |-----
//  |      dup      |    |
//  |      qos      |    |
//  |    retain     | 8 bytes
//  | remain length0|    |
//  | remain length1|    |
//  | remain length2|    |
//  | remain length3|-----
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
		return 0
	}
	if n = encodeLength(fh.RemainingLength, dst); n <= 0 {
		return 0
	}

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
func decodeUint16(num []byte) (uint16, error) {
	if len(num) < 2 {
		return 0, io.ErrShortBuffer
	}
	return binary.BigEndian.Uint16(num), nil
}

// encodeUint16 will encode the num into dst first 2 bytes
func encodeUint16(num uint16, dst []byte) error {
	if len(dst) < 0 {
		return ErrInternal
	}
	binary.BigEndian.PutUint16(dst[0:], num)
	return nil
}

// decodeString decode the given bytes and read a string and return remaining slice
func decodeString(b []byte) (string, int, error) {
	r, n, err := decodeBytes(b)
	if err != nil {
		return "", 0, err
	}
	return string(r), n, nil
}

// encodeString into dst, length will be len(field) + 2
func encodeString(field string, dst []byte) error {
	return encodeBytes([]byte(field), dst)
}

// decodeBytes read length and return bytes and remaining slice
func decodeBytes(b []byte) ([]byte, int, error) {
	if len(b) < 2 {
		return nil, 0, io.ErrShortBuffer
	}
	fieldLength, err := decodeUint16(b[:2])
	if err != nil {
		return nil, 0, err
	}
	if len(b) < int(2+fieldLength) {
		return nil, 0, io.ErrShortBuffer
	}
	return b[2 : 2+fieldLength], 2 + int(fieldLength), nil
}

// encodeBytes encode field length and field into dst
func encodeBytes(field []byte, dst []byte) error {
	if cap(dst) < 2+len(field) {
		return ErrInternal
	}
	binary.BigEndian.PutUint16(dst[:2], uint16(len(field)))
	copy(dst[2:], field)
	return nil
}

// follow MQTT v3.1 implements
// buf should be a 5 length slice
func encodeLength(length int, buf []byte) (n int) {
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
func decodeLength(r io.Reader, bs []byte) (value, n int) {
	var multiplier = 1
	// bs max size 4 byte
	for n = 0; n < 4; n++ {
		if nn, err := io.ReadFull(r, bs[n:n+1]); err != nil {
			return 0, -(n + nn)
		}

		// FIXME Concern about BigEndian or LittleEndian crossing different platform
		value += (int(bs[n]&0x7f) * multiplier)
		multiplier *= 128
		if (bs[n] & 0x80) == 0 { // 0x80 = 128 = 1000 0000 ; 0x7f = 127 = 0111 1111
			break
		}
	}
	return value, n + 1
}
