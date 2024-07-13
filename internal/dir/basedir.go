package dir

import (
	"encoding/json"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/tmknom/tfmod/internal/collection"
	"github.com/tmknom/tfmod/internal/errlib"
)

type BaseDir struct {
	raw  string
	abs  string
	rel  string
	work string
}

func NewBaseDir(raw string) *BaseDir {
	return &BaseDir{
		raw:  raw,
		abs:  "",
		rel:  ".",
		work: "",
	}
}

func (d *BaseDir) Abs() string {
	if d.abs == "" {
		d.abs = d.generateAbs()
	}
	return d.abs
}

func (d *BaseDir) generateAbs() string {
	if filepath.IsAbs(d.raw) {
		return d.raw
	}
	return filepath.Clean(filepath.Join(d.Work(), d.raw))
}

func (d *BaseDir) Rel() string {
	return d.rel
}

func (d *BaseDir) RelByWork() string {
	result, err := filepath.Rel(d.Work(), d.Abs())
	if err != nil {
		log.Fatalf("%+v", errlib.Wrapf(err, "cannot resolve rel, work: %s, abs: %s", d.Work(), d.Abs()))
	}
	return filepath.Clean(result)
}

func (d *BaseDir) Work() string {
	if d.work == "" {
		d.work = d.generateWork()
	}
	return d.work
}

func (d *BaseDir) generateWork() string {
	result, err := os.Getwd()
	if err != nil {
		log.Fatalf("%+v", errlib.Wrapf(err, "cannot resolve work dir"))
	}
	return filepath.Clean(result)
}

func (d *BaseDir) CreateDir(raw string) *Dir {
	return NewDir(raw, d)
}

func (d *BaseDir) FilterDirs(paths []string) ([]*Dir, error) {
	items := collection.NewTreeSet()

	for _, path := range paths {
		abs := d.CreateDir(path).Abs()
		stat, err := os.Stat(abs)
		if err != nil {
			return nil, err
		}
		if !stat.IsDir() {
			abs = filepath.Dir(abs)
		}
		items.Add(abs)
	}

	sliceItems := items.Slice()
	result := make([]*Dir, 0, len(sliceItems))
	for _, item := range sliceItems {
		result = append(result, d.CreateDir(item))
	}
	return result, nil
}

func (d *BaseDir) ListSubDirs(ext string, exclude string) ([]*Dir, error) {
	absDirs := collection.NewTreeSet()

	err := filepath.WalkDir(d.Abs(), func(absFilepath string, entry fs.DirEntry, err error) error {
		if err != nil {
			return errlib.Wrapf(err, "invalid base dir: %#v", d)
		}
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ext && !strings.Contains(absFilepath, exclude) {
			absDirs.Add(filepath.Dir(absFilepath))
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	result := make([]*Dir, 0, 64)
	for _, absDir := range absDirs.Slice() {
		result = append(result, d.CreateDir(absDir))
	}
	return result, nil
}

func (d *BaseDir) String() string {
	return d.Abs()
}

func (d *BaseDir) GoString() string {
	return d.String()
}

func (d *BaseDir) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}
