package tfmod

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Parser struct {
	Mapping map[string][]string
	BaseDir
}

func NewParser(baseDir BaseDir) *Parser {
	return &Parser{
		Mapping: make(map[string][]string, 64),
		BaseDir: baseDir,
	}
}

func (p *Parser) ParseAll(tfDirs TfDirs) error {
	for tfDir := range tfDirs {
		raw, err := os.ReadFile(filepath.Join(tfDir, ModulesPath))
		if err != nil {
			return err
		}

		err = p.parse(tfDir, raw)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Parser) parse(tfDir string, raw []byte) error {
	var modulesJson ModulesJson

	err := json.Unmarshal(raw, &modulesJson)
	if err != nil {
		return err
	}

	for _, module := range modulesJson.Modules {
		if module.Dir == "." {
			continue
		}

		absModulePath, err := filepath.Abs(filepath.Join(tfDir, module.Dir))
		if err != nil {
			return err
		}

		relModulePath, err := filepath.Rel(p.BaseDir, absModulePath)
		if err != nil {
			return err
		}

		relTfPath, err := filepath.Rel(p.BaseDir, tfDir)
		if err != nil {
			return err
		}

		p.Mapping[relModulePath] = append(p.Mapping[relModulePath], relTfPath)
	}

	return nil
}

type Module struct {
	Key    string `json:"Key"`
	Source string `json:"Source"`
	Dir    string `json:"Dir"`
}

type ModulesJson struct {
	Modules []Module `json:"Modules"`
}
