package packets

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnsubackPacketWrite(t *testing.T) {
	cp := NewUnsubackPacket()
	cp.MessageID = 1234

	unsubackPacketBytes := bytes.Buffer{}
	cp.Write(&unsubackPacketBytes)
	assert.Equal(t, []byte{176, 2, 4, 210}, unsubackPacketBytes.Bytes(), "Ununsuback packet write not matched")

	cp.Close()
}

func TestUnsubackPacket(t *testing.T) {
	unsubackPacketBytes := bytes.NewBuffer([]byte{176, 2, 4, 210})

	packet, err := ReadPacket(unsubackPacketBytes)
	if err != nil {
		t.Fatalf("Error reading packet: %s", err.Error())
	}
	cp := packet.(*UnsubackPacket)

	assert.Equal(t, uint16(1234), cp.MessageID, "Ununsuback messageID not matched")
	packet.Close()
}
