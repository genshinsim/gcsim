package docs

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/genshinsim/gcsim/pipeline/pkg/translation"
	"github.com/genshinsim/gcsim/pkg/model"
)

type Generator struct {
	GeneratorConfig
	Names translation.OutData
}

type GeneratorConfig struct {
	Excels     string
	Characters []*model.AvatarData
	Weapons    []*model.WeaponData
	Artifacts  []*model.ArtifactData
	Enemies    []*model.MonsterData
}

type docsData struct {
	Name string
	Key  string
	Type string
}

func NewGenerator(cfg GeneratorConfig) (*Generator, error) {
	gen := &Generator{
		GeneratorConfig: cfg,
	}

	transCfg := translation.GeneratorConfig{
		Characters: cfg.Characters,
		Weapons:    cfg.Weapons,
		Artifacts:  cfg.Artifacts,
		Enemies:    cfg.Enemies,
		Languages: map[string]string{
			"English": filepath.Join(cfg.Excels, "..", "TextMap", "TextMapEN.json"),
		},
	}
	ts, err := translation.NewGenerator(transCfg)
	if err != nil {
		return nil, err
	}
	names, err := ts.GetNames("English")
	if err != nil {
		return nil, err
	}
	gen.Names = names

	return gen, nil
}

func (g *Generator) GenerateDocsPages(path string) {
	charsPath := filepath.Join(path, "/characters/")
	weaponsPath := filepath.Join(path, "/weapons/")
	artifactsPath := filepath.Join(path, "/artifacts/")
	enemiesPath := filepath.Join(path, "/enemies/")

	for _, v := range g.Characters {
		if v.SubId > 0 {
			// skip traveler
			continue
		}

		g.generateDocPage(charsPath, docsData{
			Name: g.Names.CharacterNames[v.Key],
			Key:  v.Key,
		}, charDocsPageTmpl)
	}
	for _, v := range g.Weapons {
		g.generateDocPage(weaponsPath, docsData{
			Name: g.Names.WeaponNames[v.Key],
			Key:  v.Key,
			Type: "weapon",
		}, docPageTmpl)
	}
	for _, v := range g.Artifacts {
		g.generateDocPage(artifactsPath, docsData{
			Name: g.Names.ArtifactNames[v.Key],
			Key:  v.Key,
			Type: "artifact",
		}, docPageTmpl)
	}
	for _, v := range g.Enemies {
		g.generateDocPage(enemiesPath, docsData{
			Name: g.Names.EnemyNames[v.Key],
			Key:  v.Key,
			Type: "enemy",
		}, docPageTmpl)
	}
}

func (g *Generator) generateDocPage(path string, data docsData, tmpl string) error {
	t, err := template.New("docs_page").Parse(tmpl)
	if err != nil {
		return fmt.Errorf("error compiling data template: %w", err)
	}
	buff := new(bytes.Buffer)
	err = t.Execute(buff, data)
	if err != nil {
		return fmt.Errorf("failed to execute template for %v: %w", data.Key, err)
	}
	src := buff.Bytes()
	os.WriteFile(fmt.Sprintf("%v/%v.md", path, data.Key), src, 0o644)
	return nil
}

const docPageTmpl = `---
title: {{ .Name }}
---

import AoETable from "@site/src/components/AoE/AoETable";
import IssuesTable from "@site/src/components/Issues/IssuesTable";
import NamesList from "@site/src/components/Names/NamesList";
import ParamsTable from "@site/src/components/Params/ParamsTable";
import FieldsTable from "@site/src/components/Fields/FieldsTable";

## AoE Data

<AoETable item_key="{{ .Key }}" data_src="{{ .Type }}" />

## Known issues

<IssuesTable item_key="{{ .Key }}" data_src="{{ .Type }}" />

## Names

<NamesList item_key="{{ .Key }}" data_src="{{ .Type }}" />

## Params

<ParamsTable item_key="{{ .Key }}" data_src="{{ .Type }}" />

## Fields

<FieldsTable item_key="{{ .Key }}" data_src="{{ .Type }}" />
`

const charDocsPageTmpl = `---
title: {{ .Name }}
---

import HitlagTable from "@site/src/components/Hitlag/HitlagTable";
import FieldsTable from "@site/src/components/Fields/FieldsTable";
import ParamsTable from "@site/src/components/Params/ParamsTable";
import FramesTable from "@site/src/components/Frames/FramesTable";
import IssuesTable from "@site/src/components/Issues/IssuesTable";
import AoETable from "@site/src/components/AoE/AoETable";
import NamesList from "@site/src/components/Names/NamesList";
import ActionsTable from "@site/src/components/Actions/ActionsTable";

## Frames

<FramesTable item_key="{{ .Key }}" />

## Hitlag Data

<HitlagTable item_key="{{ .Key }}" />

## AoE Data

<AoETable item_key="{{ .Key }}" />

## Known issues

<IssuesTable item_key="{{ .Key }}" />

## Names

<NamesList item_key="{{ .Key }}" />

## Legal Actions

<ActionsTable item_key="{{ .Key }}" />

## Params

<ParamsTable item_key="{{ .Key }}" />

## Fields

<FieldsTable item_key="{{ .Key }}" />
`
