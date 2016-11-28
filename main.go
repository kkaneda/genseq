package main

import (
	"fmt"
	"log"
	"os"

	"github.com/kkaneda/genseq/fasta"
	"github.com/kkaneda/genseq/seq"
)

// main defines the main method for the goseq program.
func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s <input_filename>\n", os.Args[0])
		os.Exit(1)
	}
	filename := os.Args[1]
	seqSet, err := fasta.LoadFromFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	output, err := seq.Run(seqSet)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(output)
}
