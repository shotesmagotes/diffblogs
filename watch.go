package main

import (
	"regexp"
	"github.com/fsnotify/fsnotify"
	"fmt"
	"log"
	"github.com/shota-makino/diffblogs/versions"
	"github.com/shota-makino/diffblogs/watcher"
	"time"
	"path"
)

var ignoreFiles []*regexp.Regexp = []*regexp.Regexp{
	regexp.MustCompile(`^\.`),
	regexp.MustCompile(`\.html$`),
	regexp.MustCompile(`\.tmp.txt$`),
	regexp.MustCompile(`\.diffs$`),
	regexp.MustCompile(`\.DS_Store$`),
}

var ignoreOps []fsnotify.Op = []fsnotify.Op{
	fsnotify.Write,
	fsnotify.Remove,
	fsnotify.Rename,
	fsnotify.Chmod,
}

func watch(w string, v versions.Config) {
	batcher, err := watcher.New(time.Second)
	defer batcher.Close()
	if err != nil {
		log.Fatal("error creating watcher.")
	}

	act := make(chan fsnotify.Event, 1)

	go func() {
		for {
		OUTER:
			select {
			case events := <-batcher.Events:
				// get only first event
				var event fsnotify.Event
				for _, event = range events {
					if event.Op == fsnotify.Create {
						break
					}
				}
				// does not seem to operate correctly if files are saved in place in directory
				// .DS_Store complicates this because this file is written to every time a file is found
				for _, re := range ignoreFiles {
					if re.MatchString(event.Name) {
						//fmt.Println("Ignored File: ", event.Name)
						break OUTER
					}
				}

				for _, op := range ignoreOps {
					if event.Op&op == op {
						//fmt.Println("Ignored Operation: ", event.Op)
						break OUTER
					}
				}

				act <- event
			case err := <-batcher.Watcher.Errors:
				fmt.Println("Error in watch: ", err)
			}
		}
	}()
	batcher.Watcher.Add(w)

	fmt.Printf("\n\nWatching %s ...\n", w)

	for {
		evt := <-act
		//fmt.Printf("+++++++++++++++++++++++++++++++++++++++++++++++\n")
		//fmt.Printf("RECEIVED: %s; OPERATION: %s\n", evt.Name, evt.Op)

		// Expect small text files because this is intended for publishing blogs. If the file itself is large
		// then consider not using this tool at all in the first place because it will cause a long load
		// time appending multiple versions of large files to itself.

		v.SetNewBase(evt.Name)

		// Get the most recently published article
		fl, err := v.GetLatestVer()
		if err != nil {
			switch err.(type) {
			case *versions.VerError:
				fmt.Printf("First version of %s\n", path.Base(evt.Name))
				fl = ""
			default:
				panic(err)
			}
		}

		// Get the currently processing article
		fc, err := v.MakeLatestVer()
		if err != nil {
			panic(err)
		}

		go run(fl, fc, v)

	}
}

