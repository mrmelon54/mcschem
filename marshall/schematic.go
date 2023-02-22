package marshall

import (
	"fmt"
	"github.com/MrMelon54/mcschem/region"
	"github.com/MrMelon54/mcschem/structure"
	"github.com/MrMelon54/mcschem/xyz"
	"github.com/Tnze/go-mc/nbt"
	"github.com/Tnze/go-mc/save"
	"io"
	"regexp"
)

var blockNameParsing = regexp.MustCompile(`(?mi)^([a-z0-9_]+?):([a-z0-9_]+?)(\[([a-z0-9_=,]+?)])?$`)

type SchematicFormat struct {
	name string
	data structure.SchematicData
}

func yzxToIndex(xyz xyz.XYZ, size xyz.XYZ) int32 {
	return xyz.Y*size.X*size.Z + xyz.Z*size.X + xyz.X
}

func ParseSchematicToLitematic(reader io.Reader) (*LitematicFormat, error) {
	name, data, err := decodeGeneric[structure.SchematicData](reader)
	if err != nil {
		return nil, err
	}
	size := xyz.XYZ{X: data.Width, Y: data.Height, Z: data.Length}

	empty := region.Empty()
	empty.SetSize(size)
	empty.ResetToMatchSize()

	// Convert palette
	palette := data.Palette
	pOut := make([]save.BlockState, len(palette))
	for s, i := range palette {
		pOut[i] = parseStringToBlockState(s)
	}

	// Convert block data
	blockData := data.BlockData
	size.LoopAllBlocks(func(xyz xyz.XYZ) {
		b := blockData[yzxToIndex(xyz, size)]
		empty.SetBlock(xyz, pOut[b])
	})

	r, err := empty.ToLitematicRegion()
	if err != nil {
		return nil, err
	}

	return &LitematicFormat{name: name, Data: structure.LitematicData{
		Version:              data.Version,
		MinecraftDataVersion: data.DataVersion,
		Metadata: structure.LitematicMetadata{
			EnclosingSize: size,
			Author:        data.Metadata.Author,
			Description:   data.Metadata.Generator,
			Name:          data.Metadata.Name,
			RegionCount:   1,
			TimeCreated:   data.Metadata.Date,
			TimeModified:  data.Metadata.Date,
		},
		Regions: map[string]structure.LitematicRegion{
			"Unnamed": r,
		},
	}}, nil
}

func parseStringToBlockState(s string) save.BlockState {
	match := blockNameParsing.FindStringSubmatch(s)
	if match == nil {
		fmt.Printf("Failed to parse '%s' as a block state\n", s)
		return save.BlockState{Name: "minecraft:fire"}
	}
	aNamespace := match[1]
	aBlock := match[2]
	aProp := match[4]
	return save.BlockState{
		Name:       fmt.Sprintf("%s:%s", aNamespace, aBlock),
		Properties: parseStringToBlockStateProperties(aProp),
	}
}

func parseStringToBlockStateProperties(s string) nbt.RawMessage {
	a := make(map[string]string)
	var name, value string
	for _, b := range s {
		switch b {
		case '=':
			name = value
			value = ""
		case ',':
			a[name] = value
			name = ""
			value = ""
		default:
			value += string(b)
		}
	}
	if len(s) > 0 {
		a[name] = value
	}
	marshal, err := nbt.Marshal(a)
	if err != nil {
		return nbt.RawMessage{}
	}
	var c nbt.RawMessage
	if err = nbt.Unmarshal(marshal, &c); err != nil {
		return nbt.RawMessage{}
	}
	return c
}
