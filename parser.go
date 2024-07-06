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
		raw, err := os.ReadFile(filepath.Join(sourceDir.Abs(), TerraformModulesPath))
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
	var terraformModulesJson TerraformModulesJson
	err := json.Unmarshal(raw, &terraformModulesJson)
	if err != nil {
		return nil, errlib.Wrapf(err, "invalid json: %s", string(raw))
	}

	relModuleDirs := make([]*ModuleDir, 0, len(terraformModulesJson.Modules))
	for _, module := range terraformModulesJson.Modules {
		if module.Dir == "." {
			continue
		}
		moduleDir := NewModuleDir(filepath.Join(sourceDir.Abs(), module.Dir), sourceDir.BaseDir())
		relModuleDirs = append(relModuleDirs, moduleDir)
	}

	sort.Slice(relModuleDirs, func(i, j int) bool {
		return relModuleDirs[i].Rel() < relModuleDirs[j].Rel()
	})
	return relModuleDirs, nil
}

type TerraformModule struct {
	Key    string `json:"Key"`
	Source string `json:"Source"`
	Dir    string `json:"Dir"`
}

type TerraformModulesJson struct {
	Modules []TerraformModule `json:"Modules"`
}

const (
	TerraformModulesPath = ".terraform/modules/modules.json"
)
