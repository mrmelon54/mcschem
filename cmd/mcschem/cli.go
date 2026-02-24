package main

import (
	"flag"
	"fmt"
	"github.com/mrmelon54/mcschem/modules/single-blocks"
	"os"
)

type cliFlags struct {
	input        string
	output       string
	singleBlocks bool
}

func main() {
	var f cliFlags
	flag.StringVar(&f.input, "i", "", "Input file")
	flag.StringVar(&f.output, "o", "out", "Output file")
	flag.BoolVar(&f.singleBlocks, "single-blocks", false, "Single blocks")
	flag.Parse()

	if f.singleBlocks {
		if f.input == "" {
			fmt.Println("No input file specified")
			os.Exit(1)
		}
		single_blocks.Run(f.input, f.output)
		return
	}

}
