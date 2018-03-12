package main

import (
	"io/ioutil"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/shota-makino/diffblogs/persist"
	"github.com/shota-makino/diffblogs/versions"
	"unicode/utf8"
	"fmt"
)

// Runs the engine for
func run(a, b string, v versions.Config) {
	// Used to hold the string values of the files we are comparing
	var at, bt string

	if a != "" {
		aa, err := ioutil.ReadFile(a)
		if err != nil {
			panic(err)
		}
		at = string(aa)
	}

	if b != "" {
		bb, err := ioutil.ReadFile(b)
		if err != nil {
			panic(err)
		}
		bt = string(bb)
	}

	dmp := diffmatchpatch.New()
	// Diff between latest file and the current file
	dmpD := dmp.DiffMain(at, bt, false)

	// Convert []dmp.Diff -> persist.Diffs
	vr, err := v.GetLatestVersionNumber()
	if err != nil {
		panic(err)
	}
	d := persist.Dtod(vr, dmpD)

	if a == "" || b == "" {
		if ok := v.SaveDiff(d); !ok {
			panic("Could not save diffs")
		}

		v.DiffsToHTMLFile(d)
		//fmt.Printf("^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^\n")
		return
	}

	// Main diff file
	bd := v.GetDiffs()
	fin := mergeDiffs(bd, d)

	if ok := v.SaveDiff(fin); !ok {
		panic("Could not save aggregated diff")
	}

	v.DiffsToHTMLFile(fin)
	//fmt.Printf("^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^\n")
}

func mergeDiffs(base, cur persist.Diffs) persist.Diffs {
	fmt.Print("Merging diffs...")
	var (
		// Counters for indexing base and cur
		indBase = 0
		indCur = 0
		lenBase = len(base)
		lenCur = len(cur)
		// Keeps track of current Diff.Text
		txtBase string
		txtCur string
		capR int
	)

	//printPersistDiffs(base)
	//printPersistDiffs(cur)
	merged := make(persist.Diffs, 0, lenBase + lenCur + 2)

	for indBase < lenBase && indCur < lenCur {
		baseD := base[indBase].Diff
		curD := cur[indCur].Diff

		if baseD.Type == diffmatchpatch.DiffDelete {
			merged = append(merged, base[indBase])
			indBase++
			continue
		}

		if curD.Type == diffmatchpatch.DiffInsert {
			merged = append(merged, cur[indCur])
			indCur++
			continue
		}

		// Four combinations of (base, cur):
		// 1. eq, eq
		// 2. eq, del
		// 3. ins, eq
		// 4. ins, del
		if len(txtBase) == 0 {
			txtBase = baseD.Text
		}
		if len(txtCur) == 0 {
			txtCur = curD.Text
		}

		nrBase := utf8.RuneCountInString(txtBase)
		nrCur := utf8.RuneCountInString(txtCur)

		if nrBase > nrCur {
			capR = nrCur
		} else {
			capR = nrBase
		}
		aggR := make([]rune, 0, capR)

		for len(txtBase) > 0 && len(txtCur) > 0 {
			cBase, wdBase := utf8.DecodeRuneInString(txtBase)
			cCur, wdCur := utf8.DecodeRuneInString(txtCur)

			if wdBase != wdCur || cBase != cCur {
				panic("Could not match diffs")
			}

			aggR = append(aggR, cBase)
			txtBase = txtBase[wdBase:]
			txtCur = txtCur[wdCur:]
		}

		var diff = diffmatchpatch.Diff {
			Text: string(aggR),
			Type: curD.Type,
		}

		var pdiff = persist.Diff {
			Diff: diff,
			VS: base[indBase].VS,
			VE: cur[indCur].VE,
		}

		merged = append(merged, pdiff)

		if len(txtBase) == 0 {
			indBase++
		}

		if len(txtCur) == 0 {
			indCur++
		}
	}

	if indCur < lenCur {
		for i := indCur; i < lenCur; i++ {
			diff := cur[i].Diff

			switch diff.Type {
			case diffmatchpatch.DiffInsert:
				merged = append(merged, cur[i])
			case diffmatchpatch.DiffDelete:
				fmt.Println("\n\tCUR: Did not match Delete")
				merged = append(merged, cur[i])
			case diffmatchpatch.DiffEqual:
				fmt.Println("\n\tCUR: Did not match Equal")
				merged = append(merged, cur[i])
			}
		}
	}

	if indBase < lenBase {
		for i := indBase; i < lenBase; i++ {
			diff := base[i].Diff

			switch diff.Type {
			case diffmatchpatch.DiffDelete:
				merged = append(merged, base[i])
			case diffmatchpatch.DiffInsert:
				fmt.Println("\n\tBASE: Did not match an insert")
				base[i].Diff.Type = diffmatchpatch.DiffEqual
				base[i].VE = 0
				merged = append(merged, base[i])
			case diffmatchpatch.DiffEqual:
				fmt.Println("\n\tBASE: Did not match an equal")
				base[i].VE = 0
				merged = append(merged, base[i])
			}
		}
	}
	fmt.Println("Done.")

	printPersistDiffs(merged)
	return merged
}