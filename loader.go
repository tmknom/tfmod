package tfmod

import (
	"log"
	"path/filepath"

	"github.com/tmknom/tfmod/internal/dir"
)

type Loader struct {
	Store
	*dir.BaseDir
	enableTf bool
}

func NewLoader(store Store, baseDir *dir.BaseDir, enableTf bool) *Loader {
	return &Loader{
		Store:    store,
		BaseDir:  baseDir,
		enableTf: enableTf,
	}
}

func (l *Loader) Load() error {
	log.Printf("BaseDir: %v", l.BaseDir)

	subDirs, err := l.BaseDir.FilterSubDirs(".tf", filepath.Dir(TerraformModulesPath))
	if err != nil {
		return err
	}
	sourceDirs := l.toSourceDirs(subDirs)
	log.Printf("Source dirs: %v", sourceDirs)

	terraform := NewTerraform()
	err = terraform.ExecuteGetAll(sourceDirs, l.enableTf)
	if err != nil {
		return err
	}
	log.Printf("Execute terraform get to: %v", sourceDirs)

	parser := NewParser(l.Store)
	err = parser.ParseAll(sourceDirs)
	if err != nil {
		return err
	}
	log.Printf("Parse to: %v", sourceDirs)

	return nil
}

func (l *Loader) toSourceDirs(dirs []string) []*SourceDir {
	sourceDirs := make([]*SourceDir, 0, len(dirs))
	for _, subDir := range dirs {
		sourceDirs = append(sourceDirs, NewSourceDir(subDir, l.BaseDir))
	}
	return sourceDirs
}
