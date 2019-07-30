package packets

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnackPacketWrite(t *testing.T) {
	cp := NewConnackPacket()
	cp.SessionPresent = false
	cp.ReturnCode = 1
	connackPacketBytes := bytes.Buffer{}
	ll, err := cp.Write(&connackPacketBytes)
	assert.NoError(t, err)
	assert.Equal(t, connackPacketBytes.Len(), ll)
	assert.Equal(t, []byte{32, 2, 0, 1}, connackPacketBytes.Bytes(), "Connack packet write not matched")
	cp.Close()
}

func TestConnackPacket(t *testing.T) {
	connackPacketBytes := bytes.NewBuffer([]byte{32, 2, 1, 0})
	packet, ll, err := ReadPacket(connackPacketBytes)
	if err != nil {
		t.Fatalf("Error reading packet: %s", err.Error())
	}
	cp := packet.(*ConnackPacket)

	assert.Equal(t, true, cp.SessionPresent, "Connack packet session present not matched")
	assert.Equal(t, uint8(0), cp.ReturnCode, "Connack packet return code not matched")
	assert.Equal(t, connackPacketBytes.Cap(), ll)

	packet.Close()
}
