package single_blocks

import (
	"compress/gzip"
	"fmt"
	structure2 "github.com/MrMelon54/mcschem/structure"
	"github.com/Tnze/go-mc/nbt"
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
	reader, err := gzip.NewReader(open)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	decoder := nbt.NewDecoder(reader)

	var data structure2.SchematicData
	decode, err := decoder.Decode(&data)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Tag name: %s\n", decode)
	fmt.Printf("Size: %d %d %d\n", data.Width, data.Height, data.Length)

	fmt.Println("Palette:")
	for i := range data.Palette {
		fmt.Printf("  - %s\n", i)
	}
	fmt.Println()
	fmt.Println("Generating a separate schematic for each palette block")

	blob := structure2.NewBlockDataBlob(data.Palette, data.BlockData, data.Width, data.Height, data.Length)
	for i := range data.Palette {
		oPalette, oBlockData, oWidth, oHeight, oLength := blob.SingleBlockData(i).Output()
		out := data
		out.Palette = oPalette
		out.BlockData = oBlockData
		out.Width = oWidth
		out.Height = oHeight
		out.Length = oLength
		out.Metadata.Name = fmt.Sprintf("%s_%s", data.Metadata.Name, i)
		baseName := path.Base(input)
		blankName := strings.TrimSuffix(baseName, path.Ext(baseName))
		safeName := fmt.Sprintf("%s-%s.schem", blankName, strings.ReplaceAll(i, ":", "_"))
		create, err := os.Create(path.Join(output, safeName))
		if err != nil {
			fmt.Printf("Failed to output '%s': %s\n", i, err)
			continue
		}
		writer := gzip.NewWriter(create)
		encoder := nbt.NewEncoder(writer)
		err = encoder.Encode(out, decode)
		if err != nil {
			fmt.Printf("Failed to encode '%s': %s\n", i, err)
		}
		err = writer.Close()
		if err != nil {
			fmt.Println("Failed to close gzip:", err)
		}
	}

	fmt.Println("")
	fmt.Println("Finished outputting schematics")
}
