package packets

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubscribePacketWrite(t *testing.T) {
	cp := NewSubscribePacket()
	cp.MessageID = 1234
	cp.Topics = []string{"hello", "world", "/this/is/mqtt"}
	cp.QoSs = []byte{0, 1, 2}

	subscribePacketBytes := bytes.Buffer{}
	assert.NoError(t, cp.Write(&subscribePacketBytes))
	assert.Equal(t, []byte{130, 34, 4, 210, 0, 5, 104, 101, 108, 108, 111, 0, 0, 5, 119, 111,
		114, 108, 100, 1, 0, 13, 47, 116, 104, 105, 115, 47, 105, 115, 47, 109, 113, 116, 116, 2},
		subscribePacketBytes.Bytes(), "Subscribe packet write not matched")

	cp.Close()
}

func TestSubscribePacket(t *testing.T) {
	subscribePacketBytes := bytes.NewBuffer([]byte{130, 34, 4, 210, 0, 5, 104, 101, 108, 108, 111, 0, 0, 5,
		119, 111, 114, 108, 100, 1, 0, 13, 47, 116, 104, 105, 115, 47, 105, 115, 47, 109, 113, 116, 116, 2})

	packet, err := ReadPacket(subscribePacketBytes)
	if err != nil {
		t.Fatalf("Error reading packet: %s", err.Error())
	}
	cp := packet.(*SubscribePacket)

	assert.Equal(t, uint16(1234), cp.MessageID, "Subscribe messageID not matched")
	assert.Equal(t, []string{"hello", "world", "/this/is/mqtt"}, cp.Topics, "Subscribe topics not matched")
	assert.Equal(t, []byte{0, 1, 2}, cp.QoSs, "Subscribe QoSs not matched")

	packet.Close()
}
