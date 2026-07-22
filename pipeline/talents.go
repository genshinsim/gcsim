package main

import (
	"bytes"
	"errors"
	"fmt"
	"maps"
	"reflect"
	"slices"
	"strconv"
	"strings"

	"github.com/shizukayuki/excel-hk4e"
)

type AttributeTuple struct {
	*AttributeSpec
	pos int
}

func ParseTalent(attributes []*AttributeSpec, hint string) (*AttributeTuple, error) {
	var find func(attr *AttributeSpec) bool
	t := &AttributeTuple{pos: -1}
	param := -1

	switch {
	case hint == "":
		return nil, errors.New("talent is empty")
	case strings.Contains(hint, "(") && hint[len(hint)-1] == ')':
		t.pos = 0
		s := hint[:len(hint)-1]
		typ, s, _ := strings.Cut(s, "(")
		s, num, ok := strings.Cut(s, "|")
		if ok {
			var err error
			if t.pos, err = strconv.Atoi(num); err != nil {
				return nil, fmt.Errorf("failed to parse talent index: %v", hint)
			}
		}
		partial := strings.HasPrefix(s, "?")
		s = excel.Slug(s)
		find = func(attr *AttributeSpec) bool {
			if typ != "" && attr.Type != typ {
				return false
			}
			desc := excel.Slug(attr.Desc)
			if partial {
				desc += excel.Slug(attr.ParamDesc)
				return strings.Contains(desc, s)
			}
			return desc == s
		}
	case strings.Contains(hint, "#"):
		var err error
		typ, num, _ := strings.Cut(hint, "#")
		if param, err = strconv.Atoi(num); err != nil {
			return nil, fmt.Errorf("failed to parse talent param: %v", hint)
		}
		param = max(param, -1)
		find = func(attr *AttributeSpec) bool { return attr.Type == typ && slices.Contains(attr.Index, param) }
	default:
		return nil, fmt.Errorf("unknown talent format: %v", hint)
	}

	if attrs := excel.Filter(attributes, find); len(attrs) == 1 {
		t.AttributeSpec = attrs[0]
	} else {
		return nil, fmt.Errorf("talent results in attrs=%v but we expect 1: %v", len(attrs), hint)
	}
	if param != -1 {
		t.pos = slices.Index(t.Index, param)
	}

	if t.pos < 0 && t.pos >= len(t.Index) {
		return nil, fmt.Errorf("talent index/param not found in attribute: %v", hint)
	}
	return t, nil
}

func flatten(v any, n int, convert func(string) (*AttributeTuple, error)) (any, int, int, error) {
	switch v := v.(type) {
	case []any, []string:
		var arr []any
		var lv int
		base := n
		rv := reflect.ValueOf(v)
		for i := range rv.Len() {
			v := rv.Index(i).Interface()
			v, levels, nest, err := flatten(v, base+1, convert)
			if err != nil {
				return nil, 0, 0, err
			}
			lv, n = max(lv, levels), max(n, nest)
			arr = append(arr, v)
		}
		if len(arr) == 1 {
			return arr[0], lv, max(0, n-1), nil
		}
		return arr, lv, max(0, n), nil
	case string:
		t, err := convert(v)
		if err != nil {
			return nil, 0, 0, err
		}
		return t, len(t.Values[t.pos]), max(0, n), nil
	default:
		panic("unreachable")
	}
}

func emitTalents(b *bytes.Buffer, talents []map[string]any, attributes []*AttributeSpec) error {
	var emit func(b *bytes.Buffer, v any, n, levels, nest int)
	emit = func(b *bytes.Buffer, v any, n, levels, nest int) {
		switch v := v.(type) {
		case []any:
			if levels > 0 {
				b.WriteString("\n")
			}
			last := len(v) - 1
			for i, v := range v {
				r := nest - n - 1
				for range r {
					b.WriteString("{")
				}
				emit(b, v, n+1, levels, nest)
				for range r {
					b.WriteString("}")
				}
				if i != last || levels > 0 {
					b.WriteString(",")
				}
				if levels > 0 {
					b.WriteString("\n")
				}
			}
		case *AttributeTuple:
			if levels == 0 {
				b.WriteString(strconv.FormatFloat(v.Const[v.pos], 'g', -1, 64))
			} else {
				values := v.Values[v.pos]
				if values == nil {
					values = slices.Repeat([]float64{v.Const[v.pos]}, levels)
				}
				s := dumpGo(values, false)
				s = strings.TrimSpace(s)
				s = strings.TrimPrefix(s, "[]float64")
				if n != 0 && s[0] == '{' {
					desc := fmt.Sprintf("// %s: %s - %s: %s", v.Type, excel.Slug(v.Name), v.Desc, v.ParamDesc)
					s = "{ " + desc + s[1:]
				}
				b.WriteString(s)
			}
		default:
			panic("unreachable")
		}
	}

	var consts []string
	var vars []string
	for _, talent := range talents {
		if len(talent) != 1 {
			return fmt.Errorf("talent has multiple keys: %v", slices.Sorted(maps.Keys(talent)))
		}
		name := slices.Collect(maps.Keys(talent))[0]
		arr, levels, nest, err := flatten(talent[name], 0, func(hint string) (*AttributeTuple, error) {
			return ParseTalent(attributes, hint)
		})
		if err != nil {
			return err
		}

		t, ok := arr.(*AttributeTuple)
		if !ok && nest == 0 {
			panic("unreachable")
		}

		var desc string
		if t != nil {
			desc = fmt.Sprintf("// %s: %s", t.Type, excel.Slug(t.Name))
			if t.ParamDesc != "" {
				desc += fmt.Sprintf(" - %s: %s", t.Desc, t.ParamDesc)
			}
		}

		b := bytes.NewBuffer(nil)
		b.WriteString(name)
		b.WriteString(" = ")
		if nest > 0 || levels > 0 {
			b.WriteString(strings.Repeat("[]", min(levels, 1)+nest))
			b.WriteString("float64")
		}
		if nest > 0 {
			b.WriteString("{")
		}
		emit(b, arr, 0, levels, nest)
		if nest > 0 {
			b.WriteString("}")
		}

		if nest > 0 || levels > 0 {
			v := strings.Join([]string{desc, ""}, "\n") + b.String()
			vars = append(vars, strings.TrimSpace(v))
		} else {
			v := strings.Join([]string{b.String(), desc}, " ")
			consts = append(consts, strings.TrimSpace(v))
		}
	}

	if s := strings.Join(consts, "\n"); s != "" {
		b.WriteString("const (\n")
		b.WriteString(s)
		b.WriteString("\n)\n")
	}
	if s := strings.Join(vars, "\n"); s != "" {
		b.WriteString("var (\n")
		b.WriteString(s)
		b.WriteString("\n)\n")
	}
	return nil
}
