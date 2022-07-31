package main

import (
	"fmt"
	"os"
	"text/template"

	"github.com/genshinsim/gcsim/pkg/core/keys"
)

type charData struct {
	Name string
	Key  string
}

func main() {
	t, err := template.New("mdx").Parse(tmplstr)
	if err != nil {
		panic(err)
	}
	for i := keys.NoChar + 1; i < keys.EndCharKeys; i++ {
		switch i {
		case keys.TravelerDelim:
			continue
		case keys.TravelerDendro:
			continue
		case keys.TravelerHydro:
			continue
		case keys.TravelerPyro:
			continue
		case keys.TravelerMale:
			continue
		case keys.TravelerFemale:
			continue
		}
		d := charData{
			Name: i.Pretty(),
			Key:  i.String(),
		}
		fp := fmt.Sprintf("./out/%v.md", d.Key)
		os.Remove(fp)
		f, err := os.Create(fp)
		if err != nil {
			panic(err)
		}
		err = t.Execute(f, d)
		if err != nil {
			panic(err)
		}
	}
}

const tmplstr = `---
title: {{.Name}}
---

import HitlagTable from "@site/src/components/Hitlag/HitlagTable";
import FieldsTable from "@site/src/components/Fields/FieldsTable";
import ParamsTable from "@site/src/components/Params/ParamsTable";
import FramesTable from "@site/src/components/Frames/FramesTable";
import IssuesTable from "@site/src/components/Issues/IssuesTable";

## Frames

<FramesTable character="{{.Key}}" />

## Hitlag Data

<HitlagTable character="{{.Key}}" />

## Known issues

<IssuesTable character="{{.Key}}" />

## Params

<ParamsTable character="{{.Key}}" />

## Fields

<FieldsTable character="{{.Key}}" />
`
