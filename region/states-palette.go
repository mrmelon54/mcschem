package region

import (
	"fmt"
	"github.com/MrMelon54/mcschem/xyz"
	"github.com/Tnze/go-mc/level/block"
	"github.com/Tnze/go-mc/save"
	"github.com/fatih/color"
)

func writeStatesPalette(paletteData *PaletteContainer[BlocksState]) (palette []save.BlockState, data []uint64, err error) {
	rawPalette := paletteData.palette.export()
	palette = make([]save.BlockState, len(rawPalette))
	for i, v := range rawPalette {
		palette[i] = idToBlockState(v)
	}
	fmt.Printf("%s %v\n", color.RedString("[writeStatesPalette()]"), palette)
	data = paletteData.data.Raw()
	fmt.Printf("%s Data length: %d\n", color.RedString("[writeStatesPalette()]"), len(data))
	return
}

func readStatesPalette(xyz xyz.XYZ, palette []save.BlockState, data []uint64) (blockCount int16, paletteData *PaletteContainer[BlocksState], err error) {
	statePalette := make([]BlocksState, len(palette))
	for i, v := range palette {
		s, err := blockStateToId(v)
		if err != nil {
			return 0, nil, err
		}
		if !block.IsAir(s) {
			blockCount++
		}
		statePalette[i] = s
	}
	paletteData = NewStatesPaletteContainerWithData(int(xyz.Count()), data, statePalette)
	return
}
