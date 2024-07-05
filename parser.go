package tfmod

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"

	"github.com/tmknom/tfmod/internal/errlib"
)

type Parser struct {
	Store
}

func NewParser(store Store) *Parser {
	return &Parser{
		Store: store,
	}
}

func (p *Parser) ParseAll(sourceDirs []*SourceDir) error {
	for _, sourceDir := range sourceDirs {
		raw, err := os.ReadFile(filepath.Join(sourceDir.Abs(), ModulesPath))
		if err != nil {
			return errlib.Wrapf(err, "not readfile")
		}

		moduleDirs, err := p.Parse(sourceDir, raw)
		if err != nil {
			return err
		}

		for _, moduleDir := range moduleDirs {
			p.Store.Save(moduleDir, sourceDir.ToTfDir())
		}
	}
	return nil
}

func (p *Parser) Parse(sourceDir *SourceDir, raw []byte) ([]*ModuleDir, error) {
	var modulesJson ModulesJson

	err := json.Unmarshal(raw, &modulesJson)
	if err != nil {
		return nil, errlib.Wrapf(err, "invalid json: %s", string(raw))
	}

	relModuleDirs := make([]*ModuleDir, 0, len(modulesJson.Modules))
	for _, module := range modulesJson.Modules {
		if module.Dir == "." {
			continue
		}

		absModuleDir, err := filepath.Abs(filepath.Join(sourceDir.Abs(), module.Dir))
		if err != nil {
			return nil, errlib.Wrapf(err, "invalid json at Modules.Dir: %s", module.Dir)
		}

		relModuleDir, err := filepath.Rel(sourceDir.AbsBaseDir(), absModuleDir)
		if err != nil {
			return nil, errlib.Wrapf(err, "invalid absolute module dir: %s", absModuleDir)
		}
		relModuleDirs = append(relModuleDirs, NewModuleDir(relModuleDir, sourceDir.baseDir))
	}

	sort.Slice(relModuleDirs, func(i, j int) bool {
		return relModuleDirs[i].Rel() < relModuleDirs[j].Rel()
	})
	return relModuleDirs, nil
}

type Module struct {
	Key    string `json:"Key"`
	Source string `json:"Source"`
	Dir    string `json:"Dir"`
}

type ModulesJson struct {
	Modules []Module `json:"Modules"`
}
