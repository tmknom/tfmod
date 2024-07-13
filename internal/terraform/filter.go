package terraform

import (
	"github.com/tmknom/tfmod/internal/dir"
)

type Filter struct {
	baseDir *dir.BaseDir
}

func NewFilter(baseDir *dir.BaseDir) *Filter {
	return &Filter{
		baseDir: baseDir,
	}
}

func (t *Filter) SubDirs() ([]*dir.Dir, error) {
	return t.baseDir.ListSubDirs(tfExt, ModulesDir)
}
