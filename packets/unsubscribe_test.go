package packets

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnsubscribePacketWrite(t *testing.T) {
	cp := NewUnsubscribePacket()
	cp.MessageID = 1234
	cp.Topics = []string{"hello", "world", "/this/is/mqtt"}

	unsubscribePacketBytes := bytes.Buffer{}
	ll, err := cp.Write(&unsubscribePacketBytes)
	assert.NoError(t, err)
	assert.Equal(t, ll, unsubscribePacketBytes.Len())
	assert.Equal(t, []byte{162, 31, 4, 210, 0, 5, 104, 101, 108, 108, 111, 0, 5, 119, 111,
		114, 108, 100, 0, 13, 47, 116, 104, 105, 115, 47, 105, 115, 47, 109, 113, 116, 116},
		unsubscribePacketBytes.Bytes(), "Ununsubscribe packet write not matched")

	cp.Close()
}

func TestUnsubscribePacket(t *testing.T) {
	unsubscribePacketBytes := bytes.NewBuffer([]byte{162, 31, 4, 210, 0, 5,
		104, 101, 108, 108, 111, 0, 5, 119, 111, 114, 108, 100, 0, 13, 47,
		116, 104, 105, 115, 47, 105, 115, 47, 109, 113, 116, 116})

	packet, ll, err := ReadPacket(unsubscribePacketBytes)
	if err != nil {
		t.Fatalf("Error reading packet: %s", err.Error())
	}
	cp := packet.(*UnsubscribePacket)

	assert.Equal(t, unsubscribePacketBytes.Cap(), ll)
	assert.Equal(t, uint16(1234), cp.MessageID, "Ununsubscribe messageID not matched")
	assert.Equal(t, []string{"hello", "world", "/this/is/mqtt"}, cp.Topics, "Subscribe topics not matched")
	packet.Close()
}
