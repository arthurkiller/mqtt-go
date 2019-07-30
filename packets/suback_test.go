package packets

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubackPacketWrite(t *testing.T) {
	cp := NewSubackPacket()
	cp.MessageID = 1234
	cp.ReturnCodes = []byte{0, 1, 2, 128}

	subackPacketBytes := bytes.Buffer{}
	ll, err := cp.Write(&subackPacketBytes)
	assert.NoError(t, err)
	assert.Equal(t, ll, subackPacketBytes.Len())
	assert.Equal(t, []byte{144, 6, 4, 210, 0, 1, 2, 128}, subackPacketBytes.Bytes(), "Suback packet write not matched")

	cp.Close()
}

func TestSubackPacket(t *testing.T) {
	subackPacketBytes := bytes.NewBuffer([]byte{144, 6, 4, 210, 0, 1, 2, 128})

	packet, ll, err := ReadPacket(subackPacketBytes)
	if err != nil {
		t.Fatalf("Error reading packet: %s", err.Error())
	}
	cp := packet.(*SubackPacket)

	assert.Equal(t, ll, subackPacketBytes.Cap())
	assert.Equal(t, uint16(1234), cp.MessageID, "Suback messageID not matched")
	assert.Equal(t, []byte{0, 1, 2, 128}, cp.ReturnCodes, "Suback QoSs not matched")

	packet.Close()
}
