package main

import (
	"github.com/sergi/go-diff/diffmatchpatch"
	"fmt"
	"github.com/shota-makino/diffblogs/versions"
	"flag"
	"github.com/shota-makino/diffblogs/persist"
//	"path"
)

var WATCHDIR *string

func init() {
	WATCHDIR = flag.String("watch", ".", "Set the directory that you would like to process docs.")
}

func main() {
	flag.Parse()

	v := versions.Configure()
	watch(*WATCHDIR, v)

	// Testing purposes
	//v.SetNewBase("/Users/shota/Documents/example/a.txt")
	//
	//run("/Users/shota/Documents/example/a_1.tmp.txt", "/Users/shota/Documents/example/a_2.tmp.txt", v)
}


// Prints the diff array in a pretty format
func printDiffs(diffs []diffmatchpatch.Diff) {
	fmt.Println("\nDIFF -------------------------------")
	for _, aDiff := range diffs {
		switch aDiff.Type {
		case diffmatchpatch.DiffInsert:
			fmt.Printf("Insert: %s\n", aDiff.Text)
		case diffmatchpatch.DiffEqual:
			fmt.Printf("Equal: %s\n", aDiff.Text)
		case diffmatchpatch.DiffDelete:
			fmt.Printf("Delete: %s\n", aDiff.Text)
		}
	}
}

func printPersistDiffs(diffs persist.Diffs) {
	fmt.Println("\nPERSIST DIFF -------------------------------")
	for _, aDiff := range diffs {
		switch aDiff.Diff.Type {
		case diffmatchpatch.DiffInsert:
			fmt.Printf("Insert: %s\n", aDiff.Diff.Text)
		case diffmatchpatch.DiffEqual:
			fmt.Printf("Equal: %s\n", aDiff.Diff.Text)
		case diffmatchpatch.DiffDelete:
			fmt.Printf("Delete: %s\n", aDiff.Diff.Text)
		}
		fmt.Printf("\tVS: %d \t VE: %d\n", aDiff.VS, aDiff.VE)
	}
	fmt.Printf("\n")
}