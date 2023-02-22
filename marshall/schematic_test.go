package marshall

import (
	"github.com/Tnze/go-mc/nbt"
	"github.com/Tnze/go-mc/save"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseStringToBlockState(t *testing.T) {
	state := parseStringToBlockState("minecraft:basalt[axis=y]")
	assert.Equal(t, save.BlockState{
		Name: "minecraft:basalt",
		Properties: nbt.RawMessage{
			Type: 0xa,
			Data: []byte{0x8, 0x0, 0x4, 0x61, 0x78, 0x69, 0x73, 0x0, 0x1, 0x79, 0x0},
		},
	}, state)
}
