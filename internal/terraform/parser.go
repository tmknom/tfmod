package terraform

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/tmknom/tfmod/internal/dir"
	"github.com/tmknom/tfmod/internal/errlib"
)

type Parser struct {
	ParserStore
}

func NewParser(store ParserStore) *Parser {
	return &Parser{
		ParserStore: store,
	}
}

type ParserStore interface {
	Save(moduleDir *ModuleDir, stateDir *StateDir)
}

func (p *Parser) ParseAll(sourceDirs []*dir.Dir) error {
	for _, sourceDir := range sourceDirs {
		raw, err := os.ReadFile(filepath.Join(sourceDir.Abs(), ModulesJsonPath))
		if err != nil {
			continue
		}

		moduleDirs, err := p.Parse(sourceDir, raw)
		if err != nil {
			return err
		}

		for _, moduleDir := range moduleDirs {
			p.ParserStore.Save(moduleDir, NewStateDir(sourceDir.Rel(), sourceDir.BaseDir()))
		}
	}
	return nil
}

func (p *Parser) Parse(sourceDir *dir.Dir, raw []byte) ([]*ModuleDir, error) {
	var terraformModulesJson TerraformModulesJson
	err := json.Unmarshal(raw, &terraformModulesJson)
	if err != nil {
		return nil, errlib.Wrapf(err, "invalid json: %s", string(raw))
	}

	relModuleDirs := make([]*ModuleDir, 0, len(terraformModulesJson.Modules))
	for _, module := range terraformModulesJson.Modules {
		rawDir := module.Dir
		if rawDir == RootModuleDir || strings.Contains(rawDir, ModulesDir) {
			continue
		}
		moduleDir := NewModuleDir(filepath.Join(sourceDir.Abs(), rawDir), sourceDir.BaseDir())
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
