package versions

import (
	"os"
	"path"
	"strconv"
	"github.com/shota-makino/diffblogs/persist"
	"strings"
	"errors"
	"fmt"
)

var (
	ErrNoVersions = errors.New("No versions exist for this filename.")
)

type VerError struct {
	Op string
	Err error
}

func (e *VerError) Error() string {
	return e.Op + ": " + e.Err.Error()
}

type Config struct {
	start uint
	latest uint
	result string
	ignore string
	diffs string
	base FileInfo
}

type FileInfo struct {
	Name string
	Ext string
	Dir string
	orig string
}

func NewFileInfo(b string) FileInfo {
	if b == "" {
		return FileInfo {
			"",
			"",
			"",
			"",
		}
	}

	dir := path.Dir(b)
	ext := path.Ext(b)

	var fn string
	if ext == "" {
		fn = b
	} else {
		df := strings.SplitN(b, ext, 2)
		fn = df[0]
	}

	fn = path.Base(fn)

	return FileInfo {
		Name: fn,
		Ext: ext,
		Dir: dir,
		orig: b,
	}
}

func Configure() Config {
	c := Config{
		start: 1,
		latest: 0,
		result: ".html",
		ignore: ".tmp",
		diffs: ".diffs",
		base: FileInfo{},
	}

	return c
}

// Sets the base to a new value. This is useful when you want to keep
// most of the configuration constant but simply change the base file info.
//
// In the diffblogs package, this function is used every time a file
// is uploaded and the configuration stays the same to keep a consistent
// nomenclature across file versioning.
func (c *Config) SetNewBase(b string) Config {
	fi := NewFileInfo(b)

	c.base = fi
	c.latest = 0
	return *c
}

// Gets the latest version of the file with the name
// constructed via the configuration struct.
//
// If not found, then the file with the filename stored in
// Config.base.orig is returned. If the file is found
// then the file name is found, not the file itself.
func (c Config) GetLatestVer() (string, error) {
	v, err := c.GetLatestVersionNumber()
	if err != nil {
		panic(err)
	}

	if v == 0 {
		err = &VerError{Op: "Latest", Err: ErrNoVersions}
		return c.base.orig, err
	}

	fn := c.createVersionFilePath(v)
	return fn, nil
}

func (c Config) MakeLatestVer() (string, error) {
	v, err := c.GetLatestVersionNumber()
	if err != nil {
		panic(err)
	}

	fn := c.createVersionFilePath(v+1)

	fmt.Print("Renaming inputted file...")
	err = os.Rename(c.base.orig, fn)
	if err != nil {
		panic(err)
	}

	fmt.Println("Done.")
	return fn, err
}

func (c Config) GetLatestVersionNumber() (uint, error) {
	if c.latest != 0 {
		return c.latest, nil
	}

	s := c.start

	for cur := s; ; cur++ {
		_, err := os.Stat(c.createVersionFilePath(cur))
		if os.IsNotExist(err) {
			if cur == s {
				return 0, nil
			}
			c.latest = cur - 1
			return cur - 1, nil
		}

		if err != nil {
			return 0, err
		}
	}
}

func (c Config) createVersionFilePath(v uint) string {
	return path.Join(c.base.Dir, c.createVersionFileName(v))
}

func (c Config) createVersionFileName(v uint) string {
	ver := "_" + strconv.Itoa(int(v))
	ext := c.base.Ext
	tmp := c.ignore
	return c.base.Name + ver + tmp + ext
}

// Should only call when it is known that Diffs do
// exist in a file given file name constructed
// from the config file. Otherwise, this will panic.
func (c Config) GetDiffs() persist.Diffs {
	fn := c.createDiffsFileName()

	dfs, err := persist.OpenFile(fn)
	if err != nil {
		panic(err)
	}

	return dfs
}

func (c Config) SaveDiff(d persist.Diffs) bool {
	fn := c.createDiffsFileName()

	err := d.SaveDiffsAs(fn)
	if err != nil {
		panic(err)
	}

	return true
}

func (c Config) createDiffsFileName() string {
	fn := c.base.Name + c.diffs
	return path.Join(c.base.Dir, fn)
}

// Remember to close the file when done
func (c Config) GetResultFile() (*os.File, error) {
	fn := c.createResultFileName()
	return os.Create(fn)
}

func (c Config) createResultFileName() string {
	fn := c.base.Name + c.result
	return path.Join(c.base.Dir, fn)
}