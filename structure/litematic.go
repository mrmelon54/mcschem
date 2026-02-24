package structure

import (
	"github.com/mrmelon54/mcschem/xyz"
	"github.com/Tnze/go-mc/nbt"
	"github.com/Tnze/go-mc/save"
)

type LitematicData struct {
	Version              int32
	MinecraftDataVersion int32
	Metadata             LitematicMetadata
	Regions              map[string]LitematicRegion
}

type LitematicMetadata struct {
	EnclosingSize xyz.XYZ
	Author        string
	Description   string
	Name          string
	RegionCount   int32
	TimeCreated   int64
	TimeModified  int64
	TotalBlocks   int32
	TotalVolume   int32
}

type LitematicRegion struct {
	Size              xyz.XYZ
	Entities          []nbt.RawMessage
	BlockStatePalette []save.BlockState
	BlockStates       []uint64
	Position          xyz.XYZ
	PendingFluidTicks []nbt.RawMessage
	TileEntities      []nbt.RawMessage
	PendingBlockTicks []nbt.RawMessage
}
