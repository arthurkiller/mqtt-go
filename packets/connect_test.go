package packets

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConectPacketWrite(t *testing.T) {
	cp := NewConnectPacket()
	cp.ProtocolName = "MQTT"
	cp.ProtocolVersion = 4
	cp.CleanSession = true
	cp.WillFlag = true
	cp.WillQoS = 2
	cp.WillRetain = true

	cp.UsernameFlag = true
	cp.PasswordFlag = true
	cp.ReservedBit = 0
	cp.Keepalive = 300

	cp.ClientIdentifier = "test"
	cp.WillTopic = "/test"
	cp.WillMessage = []byte("hello,world")
	cp.Username = "tom"
	cp.Password = []byte("mqtt")

	connectPacketBytes := bytes.Buffer{}
	cp.Write(&connectPacketBytes)
	assert.Equal(t, []byte{16, 47, 0, 4, 77, 81, 84, 84, 4, 246, 1, 44, 0, 4, 116, 101, 115, 116, 0, 5, 47,
		116, 101, 115, 116, 0, 11, 104, 101, 108, 108, 111, 44, 119, 111, 114, 108, 100, 0, 3, 116, 111, 109,
		0, 4, 109, 113, 116, 116}, connectPacketBytes.Bytes(), "Connect packet write not matched")
	cp.Close()
}

func TestConnectPacket(t *testing.T) {
	connectPacketBytes := bytes.NewBuffer([]byte{16, 52, 0, 4, 77, 81, 84, 84, 4, 204, 0, 0, 0, 0, 0, 4, 116,
		101, 115, 116, 0, 12, 84, 101, 115, 116, 32, 80, 97, 121, 108, 111, 97, 100, 0, 8, 116, 101, 115, 116,
		117, 115, 101, 114, 0, 8, 116, 101, 115, 116, 112, 97, 115, 115})
	packet, err := ReadPacket(connectPacketBytes)
	if err != nil {
		t.Fatalf("Error reading packet: %s", err.Error())
	}
	cp := packet.(*ConnectPacket)
	if cp.ProtocolName != "MQTT" {
		t.Errorf("Connect Packet ProtocolName is %s, should be %s", cp.ProtocolName, "MQTT")
	}
	if cp.ProtocolVersion != 4 {
		t.Errorf("Connect Packet ProtocolVersion is %d, should be %d", cp.ProtocolVersion, 4)
	}
	if cp.UsernameFlag != true {
		t.Errorf("Connect Packet UsernameFlag is %t, should be %t", cp.UsernameFlag, true)
	}
	if cp.Username != "testuser" {
		t.Errorf("Connect Packet Username is %s, should be %s", cp.Username, "testuser")
	}
	if cp.PasswordFlag != true {
		t.Errorf("Connect Packet PasswordFlag is %t, should be %t", cp.PasswordFlag, true)
	}
	if string(cp.Password) != "testpass" {
		t.Errorf("Connect Packet Password is %s, should be %s", string(cp.Password), "testpass")
	}
	if cp.WillFlag != true {
		t.Errorf("Connect Packet WillFlag is %t, should be %t", cp.WillFlag, true)
	}
	if cp.WillTopic != "test" {
		t.Errorf("Connect Packet WillTopic is %s, should be %s", cp.WillTopic, "test")
	}
	if cp.WillQoS != 1 {
		t.Errorf("Connect Packet WillQoS is %d, should be %d", cp.WillQoS, 1)
	}
	if cp.WillRetain != false {
		t.Errorf("Connect Packet WillRetain is %t, should be %t", cp.WillRetain, false)
	}
	if string(cp.WillMessage) != "Test Payload" {
		t.Errorf("Connect Packet WillMessage is %s, should be %s", string(cp.WillMessage), "Test Payload")
	}
	packet.Close()
}
