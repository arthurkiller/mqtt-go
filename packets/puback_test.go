package packets

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPubackPacketWrite(t *testing.T) {
	cp := NewPubackPacket()
	cp.MessageID = 4321

	pubackPacketBytes := bytes.Buffer{}
	ll, err := cp.Write(&pubackPacketBytes)
	assert.NoError(t, err)
	assert.Equal(t, pubackPacketBytes.Len(), ll)
	assert.Equal(t, []byte{64, 2, 16, 225}, pubackPacketBytes.Bytes(), "Puback packet write not matched")
	cp.Close()
}

func TestPubackPacket(t *testing.T) {
	pubackPacketBytes := bytes.NewBuffer([]byte{64, 2, 4, 210})
	packet, ll, err := ReadPacket(pubackPacketBytes)
	if err != nil {
		t.Fatalf("Error reading packet: %s", err.Error())
	}
	cp := packet.(*PubackPacket)

	assert.Equal(t, uint16(1234), cp.MessageID, "Puback messageID")
	assert.Equal(t, pubackPacketBytes.Cap(), ll)
	packet.Close()
}
