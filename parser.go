package tfmod

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/tmknom/tfmod/internal/dir"
	"github.com/tmknom/tfmod/internal/errlib"
	"github.com/tmknom/tfmod/internal/terraform"
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
	Save(moduleDir *terraform.ModuleDir, tfDir *terraform.TfDir)
}

func (p *Parser) ParseAll(sourceDirs []*dir.Dir) error {
	for _, sourceDir := range sourceDirs {
		raw, err := os.ReadFile(filepath.Join(sourceDir.Abs(), terraform.ModulesJsonPath))
		if err != nil {
			return errlib.Wrapf(err, "not readfile")
		}

		moduleDirs, err := p.Parse(sourceDir, raw)
		if err != nil {
			return err
		}

		for _, moduleDir := range moduleDirs {
			p.ParserStore.Save(moduleDir, terraform.NewTfDir(sourceDir.Rel(), sourceDir.BaseDir()))
		}
	}
	return nil
}

func (p *Parser) Parse(sourceDir *dir.Dir, raw []byte) ([]*terraform.ModuleDir, error) {
	var terraformModulesJson TerraformModulesJson
	err := json.Unmarshal(raw, &terraformModulesJson)
	if err != nil {
		return nil, errlib.Wrapf(err, "invalid json: %s", string(raw))
	}

	relModuleDirs := make([]*terraform.ModuleDir, 0, len(terraformModulesJson.Modules))
	for _, module := range terraformModulesJson.Modules {
		rawDir := module.Dir
		if rawDir == terraform.RootModuleDir || strings.Contains(rawDir, terraform.ModulesDir) {
			continue
		}
		moduleDir := terraform.NewModuleDir(filepath.Join(sourceDir.Abs(), rawDir), sourceDir.BaseDir())
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
