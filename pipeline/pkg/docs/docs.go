package docs

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
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
		Languages: map[string][]string{
			"English": {filepath.Join(cfg.Excels, "..", "TextMap", "TextMap_MediumEN.json")},
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
			Name: strconv.Quote(g.Names.CharacterNames[v.Key]),
			Key:  v.Key,
		}, charDocsPageTmpl)
	}
	for _, v := range g.Weapons {
		g.generateDocPage(weaponsPath, docsData{
			Name: strconv.Quote(g.Names.WeaponNames[v.Key]),
			Key:  v.Key,
			Type: "weapon",
		}, docPageTmpl)
	}
	for _, v := range g.Artifacts {
		g.generateDocPage(artifactsPath, docsData{
			Name: strconv.Quote(g.Names.ArtifactNames[v.Key]),
			Key:  v.Key,
			Type: "artifact",
		}, docPageTmpl)
	}
	for _, v := range g.Enemies {
		g.generateDocPage(enemiesPath, docsData{
			Name: strconv.Quote(g.Names.MonsterNames[v.Key]),
			Key:  v.Key,
			Type: "monster",
		}, enemyPageTmpl)
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

const enemyPageTmpl = `---
title: {{ .Name }}
---

import HPTable from "@site/src/components/HP/HPTable";
import NamesList from "@site/src/components/Names/NamesList";
import ParticleTable from "@site/src/components/Particle/ParticleTable";
import ResistTable from "@site/src/components/Resist/ResistTable";

## Names

<NamesList item_key="{{ .Key }}" data_src="{{ .Type }}" />

## Resist Data

<ResistTable item_key="{{ .Key }}" data_src="{{ .Type }}" />

## Particle Data

<ParticleTable item_key="{{ .Key }}" data_src="{{ .Type }}" />

## HP Data

<HPTable item_key="{{ .Key }}" data_src="{{ .Type }}" />
`

const docPageTmpl = `---
title: {{ .Name }}
---

import AoETable from "@site/src/components/AoE/AoETable";
import FieldsTable from "@site/src/components/Fields/FieldsTable";
import IssuesTable from "@site/src/components/Issues/IssuesTable";
import NamesList from "@site/src/components/Names/NamesList";
import ParamsTable from "@site/src/components/Params/ParamsTable";

## Known issues

<IssuesTable item_key="{{ .Key }}" data_src="{{ .Type }}" />

## Names

<NamesList item_key="{{ .Key }}" data_src="{{ .Type }}" />

## Params

<ParamsTable item_key="{{ .Key }}" data_src="{{ .Type }}" />

## Fields

<FieldsTable item_key="{{ .Key }}" data_src="{{ .Type }}" />

## AoE Data

<AoETable item_key="{{ .Key }}" data_src="{{ .Type }}" />
`

const charDocsPageTmpl = `---
title: {{ .Name }}
---

import ActionsTable from "@site/src/components/Actions/ActionsTable";
import AoETable from "@site/src/components/AoE/AoETable";
import FieldsTable from "@site/src/components/Fields/FieldsTable";
import FramesTable from "@site/src/components/Frames/FramesTable";
import HitlagTable from "@site/src/components/Hitlag/HitlagTable";
import IssuesTable from "@site/src/components/Issues/IssuesTable";
import NamesList from "@site/src/components/Names/NamesList";
import ParamsTable from "@site/src/components/Params/ParamsTable";

## Known issues

<IssuesTable item_key="{{ .Key }}" />

## Names

<NamesList item_key="{{ .Key }}" />

## Frames

<FramesTable item_key="{{ .Key }}" />

## Hitlag Data

<HitlagTable item_key="{{ .Key }}" />

## AoE Data

<AoETable item_key="{{ .Key }}" />

## Legal Actions

<ActionsTable item_key="{{ .Key }}" />

## Params

<ParamsTable item_key="{{ .Key }}" />

## Fields

<FieldsTable item_key="{{ .Key }}" />
`
