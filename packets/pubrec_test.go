package packets

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPubrecPacketWrite(t *testing.T) {
	cp := NewPubrecPacket()
	cp.MessageID = 4321

	pubrecPacketBytes := bytes.Buffer{}
	assert.NoError(t, cp.Write(&pubrecPacketBytes))
	assert.Equal(t, []byte{80, 2, 16, 225}, pubrecPacketBytes.Bytes(), "Pubrec packet write not matched")
	cp.Close()
}

func TestPubrecPacket(t *testing.T) {
	pubrecPacketBytes := bytes.NewBuffer([]byte{80, 2, 4, 210})
	packet, err := ReadPacket(pubrecPacketBytes)
	if err != nil {
		t.Fatalf("Error reading packet: %s", err.Error())
	}
	cp := packet.(*PubrecPacket)

	assert.Equal(t, uint16(1234), cp.MessageID, "Pubrec messageID")
	packet.Close()
}
