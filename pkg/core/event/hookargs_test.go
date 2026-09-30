package event

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strings"
	"testing"
)

// hookArgsIgnore on the offending line, or the line above it, silences a violation. It must
// be followed by a reason.
const hookArgsIgnore = "//hookargs:ignore"

// TestHookArgs enforces the Hook contract: Emit reuses the args buffer, so a hook must not
// keep args after returning or assign to its elements. It scans every function shaped like a
// Hook (a single ...any parameter and no results) in internal/ and pkg/, however it's registered.
func TestHookArgs(t *testing.T) {
	root := filepath.Join("..", "..", "..")
	fset := token.NewFileSet()
	files := 0
	for _, dir := range []string{"internal", "pkg"} {
		err := filepath.WalkDir(filepath.Join(root, dir), func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				if d.Name() == "testdata" {
					return filepath.SkipDir
				}
				return nil
			}
			if !strings.HasSuffix(path, ".go") {
				return nil
			}
			f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
			if err != nil {
				return err
			}
			files++
			for _, v := range checkHookArgs(fset, f) {
				t.Error(v)
			}
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	}
	if files < 1000 {
		t.Fatalf("scanned only %v files; is the repo root %v?", files, root)
	}
}

func TestCheckHookArgs(t *testing.T) {
	cases := []struct {
		name string
		body string
		want []string // substrings of the expected violations, in order
	}{
		{"read element", `x := args[0].(int); _ = x`, nil},
		{"len", `if len(args) > 1 { return }`, nil},
		{"range", `for _, a := range args { _ = a }`, nil},
		{"modify pointee", `args[0].(*T).n = 1`, nil},
		{"assign element", `args[0] = 1`, []string{"assigns to an element"}},
		{"compound assign element", `args[0] += 1`, []string{"assigns to an element"}},
		{"incdec element", `args[0]++`, []string{"assigns to an element"}},
		{"address of element", `p := &args[0]; _ = p`, []string{"address of an element"}},
		{"closure", `go func() { _ = args[0] }()`, []string{"inside a closure"}},
		{"forward to another hook", `other(args...)`, nil},
		{"forward with other args", `other(1, args...)`, []string{"uses args whole"}},
		{"pass as a slice", `other(args)`, []string{"uses args whole"}},
		{"store", `keep = args`, []string{"uses args whole"}},
		{"slice", `x := args[1:]; _ = x`, []string{"uses args whole"}},
		{"ignored", "keep = args " + hookArgsIgnore + " copied before the next emit", nil},
		{"ignored line above", hookArgsIgnore + " copied before the next emit\nkeep = args", nil},
		{"ignore without reason", "keep = args " + hookArgsIgnore, []string{"needs a reason", "uses args whole"}},
		{"shadowed", `{ args := []any{1}; keep = args }`, nil},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			src := fmt.Sprintf("package p\n\nfunc hook(args ...any) {\n%v\n}\n", c.body)
			got := checkSource(t, src)
			if len(got) != len(c.want) {
				t.Fatalf("got %q, want %d violations matching %q", got, len(c.want), c.want)
			}
			for i, w := range c.want {
				if !strings.Contains(got[i], w) {
					t.Errorf("violation %d = %q, want it to contain %q", i, got[i], w)
				}
			}
		})
	}
}

func TestCheckHookArgsShapes(t *testing.T) {
	cases := []struct {
		name string
		src  string
		want int
	}{
		{"func literal", `var _ = func(args ...any) { keep = args }`, 1},
		{"method", `func (c *C) hook(args ...any) { keep = args }`, 1},
		{"returned by a factory", `func f() func(...any) { return func(args ...any) { keep = args } }`, 1},
		{"interface{} param", `func hook(args ...interface{}) { keep = args }`, 1},
		{"has results", `func g(args ...any) []any { return args }`, 0},
		{"has another param", `func g(s string, args ...any) { keep = args }`, 0},
		{"not variadic", `func g(args []any) { keep = args }`, 0},
		{"unnamed param", `var _ = func(...any) {}`, 0},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := checkSource(t, "package p\n\n"+c.src+"\n")
			if len(got) != c.want {
				t.Errorf("got %q, want %d violations", got, c.want)
			}
		})
	}
}

func checkSource(t *testing.T, src string) []string {
	t.Helper()
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "p.go", src, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}
	return checkHookArgs(fset, f)
}

// checkHookArgs reports the hook-shaped functions in f that keep, alias or assign to args.
// It works on syntax alone: f(args...) is allowed on the assumption that f is another hook,
// which is checked in turn when it's declared in the repo.
func checkHookArgs(fset *token.FileSet, f *ast.File) []string {
	ignores := map[int]string{} // line -> reason
	var out []string
	for _, g := range f.Comments {
		for _, c := range g.List {
			reason, ok := strings.CutPrefix(c.Text, hookArgsIgnore)
			if !ok {
				continue
			}
			reason = strings.TrimSpace(reason)
			ignores[fset.Position(c.Pos()).Line] = reason
			if reason == "" {
				out = append(out, fmt.Sprintf("%v: %v needs a reason", fset.Position(c.Pos()), hookArgsIgnore))
			}
		}
	}
	report := func(pos token.Pos, msg string) {
		p := fset.Position(pos)
		if r, ok := ignores[p.Line]; ok && r != "" {
			return
		}
		if r, ok := ignores[p.Line-1]; ok && r != "" {
			return
		}
		out = append(out, fmt.Sprintf("%v: hook %v", p, msg))
	}

	ast.Inspect(f, func(n ast.Node) bool {
		var ft *ast.FuncType
		var body *ast.BlockStmt
		switch fn := n.(type) {
		case *ast.FuncLit:
			ft, body = fn.Type, fn.Body
		case *ast.FuncDecl:
			ft, body = fn.Type, fn.Body
		default:
			return true
		}
		if body == nil {
			return true
		}
		if p := hookParam(ft); p != nil {
			checkHookBody(p, body, report)
		}
		return true
	})
	return out
}

// hookParam returns the args parameter if ft has the shape of a Hook.
func hookParam(ft *ast.FuncType) *ast.Ident {
	if ft.Results != nil || len(ft.Params.List) != 1 {
		return nil
	}
	p := ft.Params.List[0]
	if len(p.Names) != 1 || p.Names[0].Name == "_" {
		return nil
	}
	e, ok := p.Type.(*ast.Ellipsis)
	if !ok {
		return nil
	}
	switch elt := e.Elt.(type) {
	case *ast.Ident:
		if elt.Name != "any" {
			return nil
		}
	case *ast.InterfaceType:
		if len(elt.Methods.List) != 0 {
			return nil
		}
	default:
		return nil
	}
	return p.Names[0]
}

// checkHookBody allows args only as len(args), range args, an element read, or forwarded
// whole to another hook as f(args...).
func checkHookBody(param *ast.Ident, body *ast.BlockStmt, report func(token.Pos, string)) {
	var stack []ast.Node
	ast.Inspect(body, func(n ast.Node) bool {
		if n == nil {
			stack = stack[:len(stack)-1]
			return true
		}
		stack = append(stack, n)
		id, ok := n.(*ast.Ident)
		if !ok || id.Obj == nil || id.Obj != param.Obj {
			return true
		}
		for _, s := range stack[:len(stack)-1] {
			if _, ok := s.(*ast.FuncLit); ok {
				report(id.Pos(), "uses args inside a closure, which may run after the hook returns")
				return true
			}
		}
		switch p := stack[len(stack)-2].(type) {
		case *ast.IndexExpr:
			if p.X != id {
				break
			}
			switch gp := stack[len(stack)-3].(type) {
			case *ast.AssignStmt:
				for _, l := range gp.Lhs {
					if l == p {
						report(id.Pos(), "assigns to an element of args")
						return true
					}
				}
			case *ast.IncDecStmt:
				report(id.Pos(), "assigns to an element of args")
				return true
			case *ast.UnaryExpr:
				if gp.Op == token.AND {
					report(id.Pos(), "takes the address of an element of args")
					return true
				}
			}
			return true
		case *ast.CallExpr:
			if fn, ok := p.Fun.(*ast.Ident); ok && fn.Name == "len" && !p.Ellipsis.IsValid() {
				return true
			}
			if p.Ellipsis.IsValid() && len(p.Args) == 1 {
				return true
			}
		case *ast.RangeStmt:
			if p.X == id {
				return true
			}
		}
		report(id.Pos(), "uses args whole (passed on, sliced or stored)")
		return true
	})
}
