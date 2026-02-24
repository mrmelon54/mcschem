package region

import (
	"fmt"
	"github.com/mrmelon54/mcschem/structure"
	"github.com/mrmelon54/mcschem/xyz"
	"github.com/Tnze/go-mc/level"
	"github.com/Tnze/go-mc/nbt"
	"github.com/Tnze/go-mc/save"
	"reflect"
)

type Access struct {
	region structure.LitematicRegion
	data   *PaletteContainer[BlocksState]
}

func Empty() *Access {
	return &Access{
		region: structure.LitematicRegion{
			Size:              xyz.XYZ{},
			Entities:          []nbt.RawMessage{},
			BlockStatePalette: []save.BlockState{},
			BlockStates:       []uint64{},
			Position:          xyz.XYZ{},
			PendingFluidTicks: []nbt.RawMessage{},
			TileEntities:      []nbt.RawMessage{},
			PendingBlockTicks: []nbt.RawMessage{},
		},
		data: NewStatesPaletteContainer(0, level.BlocksState(0)),
	}
}

func (a *Access) parseInternal(size xyz.XYZ, blockPalette []save.BlockState, states []uint64) error {
	_, palette, err := readStatesPalette(size, blockPalette, states)
	if err != nil {
		return err
	}
	a.data = palette
	return err
}

func (a *Access) writeInternal() ([]save.BlockState, []uint64, error) {
	palette, data, err := writeStatesPalette(a.data)
	if err != nil {
		return nil, nil, err
	}
	return palette, data, nil
}

func FromLitematicRegion(region structure.LitematicRegion) (*Access, error) {
	a := &Access{region: region}
	return a, a.parseInternal(region.Size, region.BlockStatePalette, region.BlockStates)
}

func (a *Access) ToLitematicRegion() (structure.LitematicRegion, error) {
	fmt.Println("Generating litematic region raw data")
	palette, data, err := a.writeInternal()
	if err != nil {
		return structure.LitematicRegion{}, err
	}
	a.region.BlockStatePalette = palette
	a.region.BlockStates = data
	fmt.Printf("%#v\n", a.region.BlockStatePalette)
	fmt.Printf("%#v\n", a.region.BlockStates)
	return a.region, nil
}

func (a *Access) ResetToMatchSize() {
	b, _ := blockStateToId(save.BlockState{Name: "minecraft:air"})
	a.data = NewStatesPaletteContainer(int(a.Size().Count()), b)
}

func (a *Access) SetPosition(pos xyz.XYZ) {
	a.region.Position = pos
}

func (a *Access) Position() xyz.XYZ {
	return a.region.Position
}

func (a *Access) SetSize(size xyz.XYZ) {
	a.region.Size = size
}

func (a *Access) Size() xyz.XYZ {
	return a.region.Size
}

func (a *Access) AbsSize() xyz.XYZ {
	return a.region.Size.Abs()
}

func (a *Access) Count() int32 {
	return a.region.Size.Count()
}

func (a *Access) yzxToIndex(xyz xyz.XYZ) int32 {
	p := xyz.Abs()
	size := a.AbsSize()
	return p.Y*size.X*size.Z + p.Z*size.X + p.X
}

func (a *Access) checkPalette(state save.BlockState) bool {
	for _, i := range a.region.BlockStatePalette {
		if i.Name == state.Name && reflect.DeepEqual(i.Properties, state.Properties) {
			return true
		}
	}
	return false
}

func (a *Access) SetBlock(xyz xyz.XYZ, state save.BlockState) {
	stateId, err := blockStateToId(state)
	if err != nil {
		return
	}
	a.data.Set(int(a.yzxToIndex(xyz)), stateId)
}

func (a *Access) GetBlock(xyz xyz.XYZ) save.BlockState {
	b := a.data.Get(int(a.yzxToIndex(xyz)))
	return idToBlockState(b)
}

func (a *Access) Palette() []save.BlockState {
	return a.region.BlockStatePalette
}

func (a *Access) CloneEmpty() (*Access, error) {
	oPalette, oData, err := a.writeInternal()
	if err != nil {
		return nil, err
	}
	b := &Access{
		region: structure.LitematicRegion{
			Size:              a.region.Size,
			Entities:          a.region.Entities,
			BlockStatePalette: oPalette,
			BlockStates:       oData,
			Position:          a.region.Position,
			PendingFluidTicks: a.region.PendingFluidTicks,
			TileEntities:      a.region.TileEntities,
			PendingBlockTicks: a.region.PendingBlockTicks,
		},
	}
	b.ResetToMatchSize()
	return b, nil
}
