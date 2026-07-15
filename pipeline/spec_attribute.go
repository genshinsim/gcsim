package main

import (
	"bytes"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

type AttributeSpec struct {
	Type      string      `yaml:"type,omitempty"`
	Name      string      `yaml:"name,omitempty"`
	Desc      string      `yaml:"desc,omitempty"`
	ParamDesc string      `yaml:"param_desc,omitempty"`
	Config    string      `yaml:"config,omitempty"`
	Index     []int       `yaml:"index,flow,omitempty"`
	Const     []float64   `yaml:"const,flow,omitempty"`
	Values    [][]float64 `yaml:"values,flow,omitempty"`
}

func NewTalentIndex(params []float64) []int {
	i := slices.Index(params, 0)
	if i == -1 {
		i = len(params)
	}
	index := make([]int, len(params[:i]))
	for i := range index {
		index[i] = i
	}
	return index
}

func IndexFromParams(s string) ([]int, string) {
	const prefix = "{param"
	var params []int
	b := bytes.NewBuffer(nil)
	for {
		var ok bool
		var cut string
		if cut, s, ok = strings.Cut(s, prefix); !ok {
			b.WriteString(cut)
			break
		}
		var num string
		if num, s, ok = strings.Cut(s, ":"); !ok {
			continue
		}
		if _, s, ok = strings.Cut(s, "}"); !ok {
			continue
		}
		if param, err := strconv.Atoi(num); err == nil {
			param--
			params = append(params, param)
			b.WriteString(cut)
			fmt.Fprintf(b, "{%d}", param)
		}
	}
	return slices.Compact(params), b.String()
}

func (s *AttributeSpec) SetValues(count int, level func(i int) []float64) {
	s.Desc = cleanText(s.Desc)
	s.Const = make([]float64, len(s.Index))
	s.Values = make([][]float64, len(s.Index))
	for ind, param := range s.Index {
		values := make([]float64, count)
		for i := range values {
			values[i] = level(i)[param]
		}
		if cpy := slices.Compact(append([]float64(nil), values...)); len(cpy) == 1 {
			s.Const[ind] = cpy[0]
		} else if len(values) > 0 {
			s.Values[ind] = values
		}
	}
}

func (s *AttributeSpec) EmitDesc(prefix string) string {
	b := bytes.NewBuffer(nil)
	fmt.Fprintf(b, "%s: %s", s.Type, s.Name)
	if s.Config != "" {
		fmt.Fprintf(b, " # %s", s.Config)
	}
	b.WriteString("\n\n")

	desc := s.Desc
	desc = strings.ReplaceAll(desc, ". ", ".\n")
	desc = strings.ReplaceAll(desc, ": ", ":\n")
	for line := range strings.Lines(desc) {
		line = strings.TrimSpace(line)
		fmt.Fprintf(b, " %s\n", line)
	}
	lines := b.String()

	b.Reset()
	for line := range strings.Lines(lines) {
		line = strings.TrimSuffix(line, "\n")
		fmt.Fprintf(b, "%s%s\n", prefix, line)
	}
	return b.String()
}
