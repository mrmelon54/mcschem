package single_blocks

import (
	"bytes"
	"fmt"
	"github.com/mrmelon54/mcschem/marshall"
	"github.com/mrmelon54/mcschem/region"
	"github.com/mrmelon54/mcschem/structure"
	"github.com/mrmelon54/mcschem/xyz"
	"github.com/Tnze/go-mc/nbt"
	"github.com/Tnze/go-mc/save"
	"github.com/fatih/color"
	"os"
	"path"
	"strings"
)

func Run(input, output string) {
	err := os.MkdirAll(output, os.ModePerm)
	if err != nil {
		fmt.Println("Failed to make output folder:", err)
		os.Exit(1)
	}

	open, err := os.Open(input)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var l *marshall.LitematicFormat
	switch path.Ext(input) {
	case ".litematic":
		l, err = marshall.ParseLitematic(open)
	case ".schem", ".schematic":
		l, err = marshall.ParseSchematicToLitematic(open)
	}
	if err != nil {
		fmt.Printf("Failed to parse schematic: %s\n", err)
		os.Exit(1)
	}
	if l == nil {
		fmt.Printf("Invalid format: '%s'\n", path.Ext(input))
		os.Exit(1)
	}
	regions := l.ListRegions()
	for _, i := range regions {
		reg, err := l.GetRegion(i)
		if err != nil {
			fmt.Printf("Failed to read region '%s': %s\n", i, err)
			continue
		}
		extractRegionAsSingleBlocks(input, output, i, l.Data, reg)
	}

	fmt.Println("")
	fmt.Println("Finished outputting schematics")
}

func extractRegionAsSingleBlocks(input, output, name string, meta structure.LitematicData, access *region.Access) {
	fmt.Println("Palette:")
	palette := access.Palette()
	for i := range palette {
		if palette[i].Properties.Type == 0 {
			palette[i].Properties = nbt.RawMessage{Type: 0xa, Data: []byte{0}}
		}
		fmt.Printf("  - %s\n", stringifyBlockState(palette[i]))
	}
	fmt.Println()
	fmt.Println("Generating a separate schematic for each palette block")

	out := marshall.NewLitematic(meta.Version, meta.MinecraftDataVersion, meta.Metadata.Name, meta.Metadata.Description, meta.Metadata.Author)

	for _, i := range palette {
		strBlock := stringifyBlockState(i)
		if i.Name == "minecraft:air" {
			continue
		}
		singleBlockRegion, err := extractSinglePaletteBlock(access, i)
		if err != nil {
			fmt.Printf("Failed to extract block '%s': %s\n", strBlock, err)
			continue
		}
		err = out.AddRegion(strBlock, singleBlockRegion)
		if err != nil {
			fmt.Printf("Failed to add region for block '%s': %s\n", strBlock, err)
		}
	}

	baseName := path.Base(input)
	blankName := strings.TrimSuffix(baseName, path.Ext(baseName))
	safeName := fmt.Sprintf("%s-%s.litematic", blankName, name)
	fullName := path.Join(output, safeName)

	create, err := os.Create(fullName)
	if err != nil {
		fmt.Printf("Failed to create file '%s': %s\n", fullName, err)
		return
	}
	out.CalculateMetadata()
	err = out.Save(create)
	if err != nil {
		fmt.Printf("Failed to save litematic '%s': %s\n", fullName, err)
		return
	}
}

func extractSinglePaletteBlock(access *region.Access, i save.BlockState) (*region.Access, error) {
	a, err := access.CloneEmpty()
	if err != nil {
		return nil, err
	}
	size := a.AbsSize()
	a.ResetToMatchSize()
	size.LoopAllBlocks(func(xyz xyz.XYZ) {
		block := access.GetBlock(xyz)
		if blockStatesEqual(block, i) {
			fmt.Printf("%s Setting '%s' block at %v\n", color.HiBlueString("[extractSinglePaletteBlock()]"), i, xyz)
			a.SetBlock(xyz, i)
		}
	})
	return a, nil
}

func stringifyBlockState(i save.BlockState) string {
	return fmt.Sprintf("%s%s", i.Name, i.Properties.String())
}

func blockStatesEqual(a save.BlockState, b save.BlockState) bool {
	fmt.Printf("Comparing %#v\n          %#v\n", a, b)
	if a.Name != b.Name {
		return false
	}
	if a.Properties.Type != b.Properties.Type {
		return false
	}
	return bytes.Equal(a.Properties.Data, a.Properties.Data)
}
