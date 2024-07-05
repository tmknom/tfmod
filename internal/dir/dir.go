package dir

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/tmknom/tfmod/internal/errlib"
)

type Dir struct {
	raw     string
	abs     string
	rel     string
	baseDir *BaseDir
}

func NewDir(raw string, baseDir *BaseDir) *Dir {
	return &Dir{
		raw:     raw,
		abs:     "",
		rel:     "",
		baseDir: baseDir,
	}
}

func (d *Dir) BaseDir() *BaseDir {
	return d.baseDir
}

func (d *Dir) Abs() string {
	if d.abs != "" {
		return d.abs
	}
	return d.generateAbs()
}

func (d *Dir) generateAbs() string {
	dir := d.raw
	if len(dir) > 0 && dir[0] != os.PathSeparator {
		dir = filepath.Join(d.baseDir.Abs(), dir)
	}

	result, err := filepath.Abs(dir)
	if err != nil {
		log.Fatalf("%+v", errlib.Wrapf(err, "invalid dir: %s", dir))
	}

	d.abs = result
	return d.abs
}

func (d *Dir) Rel() string {
	if d.rel != "" {
		return d.rel
	}
	return d.generateRel()
}

func (d *Dir) generateRel() string {
	rel, err := filepath.Rel(d.baseDir.Abs(), d.Abs())
	if err != nil {
		log.Fatalf("%+v", errlib.Wrapf(err, "invalid path"))
	}

	d.rel = rel
	return d.rel
}

func (d *Dir) String() string {
	return d.Rel()
}

func (d *Dir) GoString() string {
	return d.String()
}

func (d *Dir) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}
