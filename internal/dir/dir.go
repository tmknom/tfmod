package dir

import (
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/tmknom/tfmod/internal/errlib"
)

type Dir struct {
	raw  string
	abs  string
	rel  string
	base *BaseDir
}

func NewDir(raw string, base *BaseDir) *Dir {
	return &Dir{
		raw:  raw,
		abs:  "",
		rel:  "",
		base: base,
	}
}

func (d *Dir) BaseDir() *BaseDir {
	return d.base
}

func (d *Dir) Abs() string {
	if d.abs == "" {
		d.abs = d.generateAbs()
	}
	return d.abs
}

func (d *Dir) generateAbs() string {
	clean := filepath.Clean(d.raw)
	if filepath.IsAbs(clean) {
		return clean
	}

	if strings.HasPrefix(clean, d.base.RelByWork()) {
		return filepath.Clean(filepath.Join(d.base.Work(), clean))
	}
	return filepath.Clean(filepath.Join(d.base.Abs(), clean))
}

func (d *Dir) Rel() string {
	if d.rel == "" {
		d.rel = d.generateRel()
	}
	return d.rel
}

func (d *Dir) generateRel() string {
	result, err := filepath.Rel(d.base.Abs(), d.Abs())
	if err != nil {
		log.Fatalf("%+v", errlib.Wrapf(err, "invalid path: %#v", d))
	}
	return filepath.Clean(result)
}

func (d *Dir) String() string {
	return d.Rel()
}

func (d *Dir) GoString() string {
	return fmt.Sprintf("&dir.Dir{raw: %s, abs: %s, rel: %s, base: %s}", d.raw, d.Abs(), d.Rel(), d.base.Abs())
}

func (d *Dir) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}
