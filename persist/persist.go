package persist

import (
	"github.com/sergi/go-diff/diffmatchpatch"
	"os"
	"encoding/gob"
	"fmt"
)

type Diff struct {
	Diff diffmatchpatch.Diff
	VS uint
	VE uint
}

type Diffs []Diff

func (dfs Diffs) SaveDiffsAs(filename string) error {
	fmt.Print("Saving diff file...")
	f, err := os.Create(filename)
	defer f.Close()
	if err != nil {
		panic(err)
	}

	enc := gob.NewEncoder(f)
	err = enc.Encode(dfs)
	fmt.Println("Done.")
	return err
}

// Tags all Diff struct with a version and returns the persist.Diffs struct
// as a result. The persist.Diffs struct is the struct we wish to save to the
// file.
func Dtod (vrs uint, dfs []diffmatchpatch.Diff) Diffs {
	diffs := make(Diffs, 0)

	for _, aDiff := range dfs {
		var diff Diff
		switch aDiff.Type {
		case diffmatchpatch.DiffInsert:
			diff = Diff {
				Diff: aDiff,
				VS: vrs,
			}
		case diffmatchpatch.DiffDelete:
			diff = Diff {
				Diff: aDiff,
				VE: vrs,
			}
		case diffmatchpatch.DiffEqual:
			diff = Diff {
				Diff: aDiff,
			}
		}
		diffs = append(diffs, diff)
	}

	return diffs
}

func OpenFile(filename string) (Diffs, error) {
	var dfs Diffs

	fmt.Print("Fetching diffs file...")
	f, err := os.Open(filename)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	dec := gob.NewDecoder(f)
	err = dec.Decode(&dfs)
	if err != nil {
		return nil, err
	}

	fmt.Println("Done.")
	return dfs, nil
}

