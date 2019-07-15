package packets

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPubcompPacketWrite(t *testing.T) {
	cp := NewPubcompPacket()
	cp.MessageID = 4321

	pubcompPacketBytes := bytes.Buffer{}
	assert.NoError(t, cp.Write(&pubcompPacketBytes))
	assert.Equal(t, []byte{112, 2, 16, 225}, pubcompPacketBytes.Bytes(), "Pubcomp packet write not matched")
	cp.Close()
}

func TestPubcompPacket(t *testing.T) {
	pubcompPacketBytes := bytes.NewBuffer([]byte{112, 2, 4, 210})
	packet, err := ReadPacket(pubcompPacketBytes)
	if err != nil {
		t.Fatalf("Error reading packet: %s", err.Error())
	}
	cp := packet.(*PubcompPacket)

	assert.Equal(t, uint16(1234), cp.MessageID, "Pubcomp messageID")
	packet.Close()
}
