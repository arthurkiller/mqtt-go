package packets

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPacketNames(t *testing.T) {
	if PacketNames[1] != "CONNECT" {
		t.Errorf("PacketNames[1] is %s, should be %s", PacketNames[1], "CONNECT")
	}
	if PacketNames[2] != "CONNACK" {
		t.Errorf("PacketNames[2] is %s, should be %s", PacketNames[2], "CONNACK")
	}
	if PacketNames[3] != "PUBLISH" {
		t.Errorf("PacketNames[3] is %s, should be %s", PacketNames[3], "PUBLISH")
	}
	if PacketNames[4] != "PUBACK" {
		t.Errorf("PacketNames[4] is %s, should be %s", PacketNames[4], "PUBACK")
	}
	if PacketNames[5] != "PUBREC" {
		t.Errorf("PacketNames[5] is %s, should be %s", PacketNames[5], "PUBREC")
	}
	if PacketNames[6] != "PUBREL" {
		t.Errorf("PacketNames[6] is %s, should be %s", PacketNames[6], "PUBREL")
	}
	if PacketNames[7] != "PUBCOMP" {
		t.Errorf("PacketNames[7] is %s, should be %s", PacketNames[7], "PUBCOMP")
	}
	if PacketNames[8] != "SUBSCRIBE" {
		t.Errorf("PacketNames[8] is %s, should be %s", PacketNames[8], "SUBSCRIBE")
	}
	if PacketNames[9] != "SUBACK" {
		t.Errorf("PacketNames[9] is %s, should be %s", PacketNames[9], "SUBACK")
	}
	if PacketNames[10] != "UNSUBSCRIBE" {
		t.Errorf("PacketNames[10] is %s, should be %s", PacketNames[10], "UNSUBSCRIBE")
	}
	if PacketNames[11] != "UNSUBACK" {
		t.Errorf("PacketNames[11] is %s, should be %s", PacketNames[11], "UNSUBACK")
	}
	if PacketNames[12] != "PINGREQ" {
		t.Errorf("PacketNames[12] is %s, should be %s", PacketNames[12], "PINGREQ")
	}
	if PacketNames[13] != "PINGRESP" {
		t.Errorf("PacketNames[13] is %s, should be %s", PacketNames[13], "PINGRESP")
	}
	if PacketNames[14] != "DISCONNECT" {
		t.Errorf("PacketNames[14] is %s, should be %s", PacketNames[14], "DISCONNECT")
	}
}

func TestPacketConsts(t *testing.T) {
	if Connect != 1 {
		t.Errorf("Const for Connect is %d, should be %d", Connect, 1)
	}
	if Connack != 2 {
		t.Errorf("Const for Connack is %d, should be %d", Connack, 2)
	}
	if Publish != 3 {
		t.Errorf("Const for Publish is %d, should be %d", Publish, 3)
	}
	if Puback != 4 {
		t.Errorf("Const for Puback is %d, should be %d", Puback, 4)
	}
	if Pubrec != 5 {
		t.Errorf("Const for Pubrec is %d, should be %d", Pubrec, 5)
	}
	if Pubrel != 6 {
		t.Errorf("Const for Pubrel is %d, should be %d", Pubrel, 6)
	}
	if Pubcomp != 7 {
		t.Errorf("Const for Pubcomp is %d, should be %d", Pubcomp, 7)
	}
	if Subscribe != 8 {
		t.Errorf("Const for Subscribe is %d, should be %d", Subscribe, 8)
	}
	if Suback != 9 {
		t.Errorf("Const for Suback is %d, should be %d", Suback, 9)
	}
	if Unsubscribe != 10 {
		t.Errorf("Const for Unsubscribe is %d, should be %d", Unsubscribe, 10)
	}
	if Unsuback != 11 {
		t.Errorf("Const for Unsuback is %d, should be %d", Unsuback, 11)
	}
	if Pingreq != 12 {
		t.Errorf("Const for Pingreq is %d, should be %d", Pingreq, 12)
	}
	if Pingresp != 13 {
		t.Errorf("Const for Pingresp is %d, should be %d", Pingresp, 13)
	}
	if Disconnect != 14 {
		t.Errorf("Const for Disconnect is %d, should be %d", Disconnect, 14)
	}
}

func TestConnackConsts(t *testing.T) {
	if Accepted != 0x00 {
		t.Errorf("Const for Accepted is %d, should be %d", Accepted, 0)
	}
	if ErrRefusedBadProtocolVersion != 0x01 {
		t.Errorf("Const for RefusedBadProtocolVersion is %d, should be %d", ErrRefusedBadProtocolVersion, 1)
	}
	if ErrRefusedIDRejected != 0x02 {
		t.Errorf("Const for RefusedIDRejected is %d, should be %d", ErrRefusedIDRejected, 2)
	}
	if ErrRefusedServerUnavailable != 0x03 {
		t.Errorf("Const for RefusedServerUnavailable is %d, should be %d", ErrRefusedServerUnavailable, 3)
	}
	if ErrRefusedBadUsernameOrPassword != 0x04 {
		t.Errorf("Const for RefusedBadUsernameOrPassword is %d, should be %d", ErrRefusedBadUsernameOrPassword, 4)
	}
	if ErrRefusedNotAuthorised != 0x05 {
		t.Errorf("Const for RefusedNotAuthorised is %d, should be %d", ErrRefusedNotAuthorised, 5)
	}
}

func TestPackUnpackControlPackets(t *testing.T) {
	packets := []ControlPacket{
		NewControlPacket(Connect).(*ConnectPacket),
		NewControlPacket(Connack).(*ConnackPacket),
		NewControlPacket(Publish).(*PublishPacket),
		NewControlPacket(Puback).(*PubackPacket),
		NewControlPacket(Pubrec).(*PubrecPacket),
		NewControlPacket(Pubrel).(*PubrelPacket),
		NewControlPacket(Pubcomp).(*PubcompPacket),
		NewControlPacket(Subscribe).(*SubscribePacket),
		NewControlPacket(Suback).(*SubackPacket),
		NewControlPacket(Unsubscribe).(*UnsubscribePacket),
		NewControlPacket(Unsuback).(*UnsubackPacket),
		NewControlPacket(Pingreq).(*PingreqPacket),
		NewControlPacket(Pingresp).(*PingrespPacket),
		NewControlPacket(Disconnect).(*DisconnectPacket),
	}
	buf := new(bytes.Buffer)
	for _, packet := range packets {
		buf.Reset()
		if err := packet.Write(buf); err != nil {
			t.Errorf("Write of %T returned error: %s", packet, err)
		}
		read, err := ReadPacket(buf)
		if err != nil {
			t.Errorf("Read of packed %T returned error: %s", packet, err)
		}
		if read.String() != packet.String() {
			t.Errorf("Read of packed %T did not equal original.\nExpected: %v\n     Got: %v", packet, packet, read)
		}
	}
}

func TestType(t *testing.T) {
	packets := []ControlPacket{
		NewControlPacket(Connect).(*ConnectPacket),
		NewControlPacket(Connack).(*ConnackPacket),
		NewControlPacket(Publish).(*PublishPacket),
		NewControlPacket(Puback).(*PubackPacket),
		NewControlPacket(Pubrec).(*PubrecPacket),
		NewControlPacket(Pubrel).(*PubrelPacket),
		NewControlPacket(Pubcomp).(*PubcompPacket),
		NewControlPacket(Subscribe).(*SubscribePacket),
		NewControlPacket(Suback).(*SubackPacket),
		NewControlPacket(Unsubscribe).(*UnsubscribePacket),
		NewControlPacket(Unsuback).(*UnsubackPacket),
		NewControlPacket(Pingreq).(*PingreqPacket),
		NewControlPacket(Pingresp).(*PingrespPacket),
		NewControlPacket(Disconnect).(*DisconnectPacket),
	}
	for i, packet := range packets {
		assert.Equal(t, byte(i+1), packet.Type())
	}
}

func TestReset(t *testing.T) {
	packets := []ControlPacket{
		NewControlPacket(Connect).(*ConnectPacket),
		NewControlPacket(Connack).(*ConnackPacket),
		NewControlPacket(Publish).(*PublishPacket),
		NewControlPacket(Puback).(*PubackPacket),
		NewControlPacket(Pubrec).(*PubrecPacket),
		NewControlPacket(Pubrel).(*PubrelPacket),
		NewControlPacket(Pubcomp).(*PubcompPacket),
		NewControlPacket(Subscribe).(*SubscribePacket),
		NewControlPacket(Suback).(*SubackPacket),
		NewControlPacket(Unsubscribe).(*UnsubscribePacket),
		NewControlPacket(Unsuback).(*UnsubackPacket),
		NewControlPacket(Pingreq).(*PingreqPacket),
		NewControlPacket(Pingresp).(*PingrespPacket),
		NewControlPacket(Disconnect).(*DisconnectPacket),
	}
	for _, packet := range packets {
		packet.Reset()
		packet.Details()
		packet.Close()
	}
}

func TestEncoding(t *testing.T) {
	if res := decodeUint16([]byte{0x56, 0x78}); res != 22136 {
		t.Errorf("decodeUint16([0x5678]) did not return 22136 but %d", res)
	}
	res := make([]byte, 64)
	if encodeUint16(22136, res); !bytes.Equal(res[:2], []byte{0x56, 0x78}) {
		t.Errorf("encodeUint16(22136) did not return [0x5678] but [0x%X]", res)
	}

	strings := map[string][]byte{
		"foo":         {0x00, 0x03, 'f', 'o', 'o'},
		"\U0000FEFF":  {0x00, 0x03, 0xEF, 0xBB, 0xBF},
		"A\U0002A6D4": {0x00, 0x05, 'A', 0xF0, 0xAA, 0x9B, 0x94},
	}
	for str, encoded := range strings {
		if res, _ := decodeString(encoded); res != str {
			t.Errorf(`decodeString(%v) did not return "%s", but "%s"`, encoded, str, res)
		}
		res := make([]byte, 64)
		if encodeString(str, res); !bytes.Equal(res[:len(str)+2], encoded) {
			t.Errorf(`encodeString("%s") did not return [0x%X], but [0x%X]`, str, encoded, res)
		}
	}

	lengths := map[int][]byte{
		0:         {0x00},
		127:       {0x7F},
		128:       {0x80, 0x01},
		16383:     {0xFF, 0x7F},
		16384:     {0x80, 0x80, 0x01},
		2097151:   {0xFF, 0xFF, 0x7F},
		2097152:   {0x80, 0x80, 0x80, 0x01},
		268435455: {0xFF, 0xFF, 0xFF, 0x7F},
	}
	buf := make([]byte, 5)
	for length, encoded := range lengths {
		if res, _ := decodeLength(bytes.NewBuffer(encoded), buf); res != length {
			t.Errorf("decodeLength([0x%X]) did not return %d, but %d", encoded, length, res)
		}
		res := make([]byte, 64)
		if n := encodeLength(length, res); !bytes.Equal(res[5-n:5], encoded) {
			t.Errorf("encodeLength(%d) did not return [0x%X], but [0x%X]", length, encoded, res[:n])
		}
	}
}
