package packets

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPublishPacketWrite(t *testing.T) {
	cp := NewPublishPacket()
	cp.FixedHeader.QoS = 1
	cp.TopicName = "/welcome/to/use/mqtt/protocol"
	cp.MessageID = 4321
	cp.Payload = []byte("this message is send via mqtt protocol 3.1")

	publishPacketBytes := bytes.Buffer{}
	assert.NoError(t, cp.Write(&publishPacketBytes))
	assert.Equal(t, []byte{50, 75, 0, 29, 47, 119, 101, 108, 99, 111, 109, 101, 47, 116, 111, 47, 117, 115,
		101, 47, 109, 113, 116, 116, 47, 112, 114, 111, 116, 111, 99, 111, 108, 16, 225, 116, 104, 105, 115,
		32, 109, 101, 115, 115, 97, 103, 101, 32, 105, 115, 32, 115, 101, 110, 100, 32, 118, 105, 97, 32, 109,
		113, 116, 116, 32, 112, 114, 111, 116, 111, 99, 111, 108, 32, 51, 46, 49}, publishPacketBytes.Bytes(), "Publish packet write not matched")

	cp.Close()

	cp = NewPublishPacket()
	cp.FixedHeader.QoS = 0
	cp.TopicName = "/welcome/to/use/mqtt/protocol"
	cp.MessageID = 4321
	cp.Payload = []byte("this message is send via mqtt protocol 3.1")

	publishPacketBytes.Reset()
	assert.NoError(t, cp.Write(&publishPacketBytes))
	assert.Equal(t, []byte{48, 73, 0, 29, 47, 119, 101, 108, 99, 111, 109, 101, 47, 116, 111, 47, 117, 115,
		101, 47, 109, 113, 116, 116, 47, 112, 114, 111, 116, 111, 99, 111, 108, 116, 104, 105, 115, 32, 109,
		101, 115, 115, 97, 103, 101, 32, 105, 115, 32, 115, 101, 110, 100, 32, 118, 105, 97, 32, 109, 113,
		116, 116, 32, 112, 114, 111, 116, 111, 99, 111, 108, 32, 51, 46, 49}, publishPacketBytes.Bytes(), "Publish packet write not matched")
	cp.Close()
}

func TestPublishPacket(t *testing.T) {
	publishPacketBytes := bytes.NewBuffer([]byte{50, 75, 0, 29, 47, 119, 101, 108, 99, 111, 109, 101, 47,
		116, 111, 47, 117, 115, 101, 47, 109, 113, 116, 116, 47, 112, 114, 111, 116, 111, 99, 111, 108, 4,
		210, 116, 104, 105, 115, 32, 109, 101, 115, 115, 97, 103, 101, 32, 105, 115, 32, 115, 101, 110, 100,
		32, 118, 105, 97, 32, 109, 113, 116, 116, 32, 112, 114, 111, 116, 111, 99, 111, 108, 32, 51, 46, 49})

	packet, err := ReadPacket(publishPacketBytes)
	if err != nil {
		t.Fatalf("Error reading packet: %s", err.Error())
	}
	cp := packet.(*PublishPacket)

	assert.Equal(t, uint16(1234), cp.MessageID, "Publish messageID not matched")
	assert.Equal(t, "/welcome/to/use/mqtt/protocol", cp.TopicName, "Publish topic not matched")
	assert.Equal(t, []byte("this message is send via mqtt protocol 3.1"), cp.Payload, "Publish payload not matched")

	packet.Close()

	publishPacketBytes = bytes.NewBuffer([]byte{48, 73, 0, 29, 47, 119, 101, 108, 99, 111, 109, 101, 47,
		116, 111, 47, 117, 115, 101, 47, 109, 113, 116, 116, 47, 112, 114, 111, 116, 111, 99, 111, 108, 116,
		104, 105, 115, 32, 109, 101, 115, 115, 97, 103, 101, 32, 105, 115, 32, 115, 101, 110, 100,
		32, 118, 105, 97, 32, 109, 113, 116, 116, 32, 112, 114, 111, 116, 111, 99, 111, 108, 32, 51, 46, 49})
	packet, err = ReadPacket(publishPacketBytes)
	if err != nil {
		t.Fatalf("Error reading packet: %s", err.Error())
	}
	cp = packet.(*PublishPacket)
	assert.Equal(t, uint16(0), cp.MessageID, "Publish messageID not matched")
	assert.Equal(t, "/welcome/to/use/mqtt/protocol", cp.TopicName, "Publish topic not matched")
	assert.Equal(t, []byte("this message is send via mqtt protocol 3.1"), cp.Payload, "Publish payload not matched")

	packet.Close()
}
