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
	raw string
	abs string
}

func NewBaseDir(raw string) *BaseDir {
	return &BaseDir{
		raw: raw,
		abs: "",
	}
}

func (d *BaseDir) CreateDir(raw string) *Dir {
	return NewDir(raw, d)
}

func (d *BaseDir) ConvertDirs(dirs []string) []*Dir {
	result := make([]*Dir, 0, len(dirs))
	for _, dir := range dirs {
		result = append(result, d.CreateDir(dir))
	}
	return result
}

func (d *BaseDir) Abs() string {
	if d.abs != "" {
		return d.abs
	}
	return d.generateAbs()
}

func (d *BaseDir) generateAbs() string {
	dir := d.raw
	if len(dir) > 0 && dir[0] != os.PathSeparator {
		currentDir, err := os.Getwd()
		if err != nil {
			log.Fatalf("%+v", errlib.Wrapf(err, "invalid current dir: %s", dir))
		}
		dir = filepath.Join(currentDir, dir)
	}

	result, err := filepath.Abs(dir)
	if err != nil {
		log.Fatalf("%+v", errlib.Wrapf(err, "invalid base dir: %s", dir))
	}

	d.abs = result
	return d.abs
}

func (d *BaseDir) FilterSubDirs(ext string, exclude string) ([]*Dir, error) {
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
