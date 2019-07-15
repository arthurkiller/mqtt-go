package packets

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPubrelPacketWrite(t *testing.T) {
	cp := NewPubrelPacket()
	cp.MessageID = 4321

	pubrelPacketBytes := bytes.Buffer{}
	assert.NoError(t, cp.Write(&pubrelPacketBytes))
	assert.Equal(t, []byte{98, 2, 16, 225}, pubrelPacketBytes.Bytes(), "Pubrel packet write not matched")
	cp.Close()
}

func TestPubrelPacket(t *testing.T) {
	pubrelPacketBytes := bytes.NewBuffer([]byte{98, 2, 4, 210})
	packet, err := ReadPacket(pubrelPacketBytes)
	if err != nil {
		t.Fatalf("Error reading packet: %s", err.Error())
	}
	cp := packet.(*PubrelPacket)

	assert.Equal(t, uint16(1234), cp.MessageID, "Pubrel messageID")
	packet.Close()
}
