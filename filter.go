package tfmod

import (
	"io/fs"
	"path/filepath"
)

type Filter struct {
	BaseDir
}

func NewFilter(baseDir BaseDir) *Filter {
	return &Filter{
		BaseDir: baseDir,
	}
}

func (d *Filter) FilterTfDirs() (TfDirs, error) {
	result := make(map[string]bool, 100)
	err := filepath.WalkDir(d.BaseDir, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(info.Name()) == ".tf" {
			dirName := filepath.Dir(path)
			result[dirName] = true
		}
		return nil
	})
	return result, err
}
