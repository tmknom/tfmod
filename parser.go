package tfmod

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Parser struct {
	BaseDir
	Store
}

func NewParser(baseDir BaseDir, store Store) *Parser {
	return &Parser{
		BaseDir: baseDir,
		Store:   store,
	}
}

func (p *Parser) ParseAll(tfDirs *TfDirs) error {
	for _, tfDir := range tfDirs.List() {
		raw, err := os.ReadFile(filepath.Join(tfDir, ModulesPath))
		if err != nil {
			return err
		}

		err = p.Parse(tfDir, raw)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Parser) Parse(tfDir TfDir, raw []byte) error {
	var modulesJson ModulesJson

	err := json.Unmarshal(raw, &modulesJson)
	if err != nil {
		return err
	}

	for _, module := range modulesJson.Modules {
		if module.Dir == "." {
			continue
		}

		absModuleDir, err := filepath.Abs(filepath.Join(tfDir, module.Dir))
		if err != nil {
			return err
		}

		relModuleDir, err := filepath.Rel(p.BaseDir.String(), absModuleDir)
		if err != nil {
			return err
		}

		p.Store.Save(relModuleDir, tfDir)
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
