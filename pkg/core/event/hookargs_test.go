package event

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"golang.org/x/tools/go/packages"
)

// hookArgsIgnore silences every violation in the statement it trails, or in the statement below
// it when it's on a line of its own. It must be followed by a reason.
const hookArgsIgnore = "//hookargs:ignore"

const modulePath = "github.com/genshinsim/gcsim"

// hookArgsAllowed are the functions a hook may pass args to: they only copy the elements.
var hookArgsAllowed = map[string]bool{
	"slices.Clone": true,
	"(*" + modulePath + "/pkg/core/event.Handler).Emit": true,
}

// TestHookArgs enforces the Hook contract: Emit reuses the args buffer, so a hook must not keep
// args after returning or assign to its elements. It type-checks internal/ and pkg/ and checks
// every function with a Hook's signature (a single ...any parameter and no results), however
// it's registered.
func TestHookArgs(t *testing.T) {
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedSyntax |
			packages.NeedTypes | packages.NeedTypesInfo,
		Dir:   filepath.Join("..", "..", ".."),
		Tests: true,
	}
	pkgs, err := packages.Load(cfg, "./internal/...", "./pkg/...")
	if err != nil {
		t.Fatal(err)
	}
	if packages.PrintErrors(pkgs) > 0 {
		t.Fatal("failed to load packages")
	}
	// with Tests, a package's files also appear in its test variant
	seen := map[string]bool{}
	files := 0
	for _, pkg := range pkgs {
		for _, f := range pkg.Syntax {
			name := pkg.Fset.File(f.Pos()).Name()
			if seen[name] {
				continue
			}
			seen[name] = true
			files++
			src, err := os.ReadFile(name)
			if err != nil {
				t.Fatal(err)
			}
			for _, v := range checkHookArgs(pkg.Fset, f, src, pkg.Types, pkg.TypesInfo) {
				t.Error(v)
			}
		}
	}
	if files < 1000 {
		t.Fatalf("checked only %v files; is the repo root %v?", files, cfg.Dir)
	}
}

// hookArgsPrelude is prepended to the snippets checkSource checks.
const hookArgsPrelude = `package p

import (
	"fmt"
	"slices"
)

var (
	keep  []any
	keepP *any
	hookV func(...any)
	_     = fmt.Sprint
	_     = slices.Clone[[]any]
)

type T struct{ n int }

func other(...any)                    {}
func keepAll(a ...any) error          { keep = a; return nil }
func takesSlice([]any)                {}
func generic[E any](...E)             {}
`

func TestCheckHookArgs(t *testing.T) {
	cases := []struct {
		name string
		body string
		want []string // substrings of the expected violations, in order
	}{
		{"read element", `x := args[0].(int); _ = x`, nil},
		{"len and cap", `_, _ = len(args), cap(args)`, nil},
		{"range", `for _, a := range args { _ = a }`, nil},
		{"modify pointee", `args[0].(*T).n = 1`, nil},
		{"copy an element out", `keep = []any{args[0]}`, nil},
		{"assign element", `args[0] = 1`, []string{"assigns to an element"}},
		{"assign parenthesized element", `(args[0]) = 1`, []string{"assigns to an element"}},
		{"range key into element", `for args[0] = range []int{1} {}`, []string{"assigns to an element"}},
		{"range value into element", `for _, args[0] = range []int{1} {}`, []string{"assigns to an element"}},
		{"address of element", `keepP = &args[0]`, []string{"address of an element"}},
		{"address of parenthesized element", `keepP = &(args[0])`, []string{"address of an element"}},
		{"closure", `go func() { _ = args[0] }()`, []string{"inside a closure"}},
		{"forward to a hook", `other(args...)`, nil},
		{"forward to a hook variable", `hookV(args...)`, nil},
		{"forward to a generic hook", `generic(args...)`, nil},
		{"forward in a goroutine", `go other(args...)`, []string{"goroutine"}},
		{"forward to a non-hook", `_ = keepAll(args...)`, []string{"not a hook"}},
		{"forward to fmt", `_ = fmt.Sprint(args...); _ = fmt.Sprintf("%v", args...)`, nil},
		{"append args after", `keep = append(keep[:0:0], args...)`, nil},
		{"append onto args", `keep = append(args, 1)`, []string{"uses args whole"}},
		{"clone", `keep = slices.Clone(args)`, nil},
		{"pass as a slice", `takesSlice(args)`, []string{"uses args whole"}},
		{"store", `keep = args`, []string{"uses args whole"}},
		{"store parenthesized", `keep = (args)`, []string{"uses args whole"}},
		{"convert", `keep = []any(args)`, []string{"uses args whole"}},
		{"slice", `keep = args[1:]`, []string{"uses args whole"}},
		{"ignored", "keep = args " + hookArgsIgnore + " copied before the next emit", nil},
		{"ignored line above", hookArgsIgnore + " copied before the next emit\nkeep = args", nil},
		{"ignored multiline statement", hookArgsIgnore + " copied before the next emit\ntakesSlice(\nargs)", nil},
		{"trailing ignore doesn't reach the next line", "keep = args " + hookArgsIgnore + " ok\nkeep = args", []string{"uses args whole"}},
		{"ignore needs the whole tag", "keep = args " + hookArgsIgnore + "me", []string{"uses args whole"}},
		{"ignore without reason", "keep = args " + hookArgsIgnore, []string{"needs a reason", "uses args whole"}},
		{"shadowed", `{ args := []any{1}; keep = args }`, nil},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := checkSource(t, fmt.Sprintf("%v\nfunc hook(args ...any) {\n%v\n}\n", hookArgsPrelude, c.body))
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
		{"method", `func (c *T) hook(args ...any) { keep = args }`, 1},
		{"returned by a factory", `func f() func(...any) { return func(args ...any) { keep = args } }`, 1},
		{"interface{} param", `func hook(args ...interface{}) { keep = args }`, 1},
		{"alias of any", "type A = any\nfunc hook(args ...A) { keep = args }", 1},
		{"generic", `func hook[E any](args ...E) { _ = any(args) }`, 1},
		{"forward to a hook that keeps args", `func hook(args ...any) { keeper(args...) }; func keeper(a ...any) { keep = a }`, 1},
		{"has results", `func g(args ...any) []any { return args }`, 0},
		{"has another param", `func g(s string, args ...any) { keep = args }`, 0},
		{"not variadic", `func g(args []any) { keep = args }`, 0},
		{"not any", `func g(args ...int) { _ = args }`, 0},
		{"unnamed param", `var _ = func(...any) {}`, 0},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := checkSource(t, hookArgsPrelude+"\n"+c.src+"\n")
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
	info := &types.Info{
		Types:      map[ast.Expr]types.TypeAndValue{},
		Defs:       map[*ast.Ident]types.Object{},
		Uses:       map[*ast.Ident]types.Object{},
		Selections: map[*ast.SelectorExpr]*types.Selection{},
		Instances:  map[*ast.Ident]types.Instance{},
	}
	conf := types.Config{Importer: importer.Default()}
	pkg, err := conf.Check("p", fset, []*ast.File{f}, info)
	if err != nil {
		t.Fatal(err)
	}
	return checkHookArgs(fset, f, []byte(src), pkg, info)
}

// checkHookArgs reports the hooks in f that keep, alias or assign to args.
func checkHookArgs(fset *token.FileSet, f *ast.File, src []byte, pkg *types.Package, info *types.Info) []string {
	lines := strings.Split(string(src), "\n")
	type ignore struct {
		reason string
		alone  bool // on a line of its own
	}
	ignores := map[int]ignore{}
	var out []string
	for _, g := range f.Comments {
		for _, c := range g.List {
			rest, ok := strings.CutPrefix(c.Text, hookArgsIgnore)
			if !ok || (rest != "" && rest[0] != ' ') {
				continue
			}
			p := fset.Position(c.Pos())
			ig := ignore{
				reason: strings.TrimSpace(rest),
				alone:  strings.HasPrefix(strings.TrimSpace(lines[p.Line-1]), hookArgsIgnore),
			}
			ignores[p.Line] = ig
			if ig.reason == "" {
				out = append(out, fmt.Sprintf("%v: %v needs a reason", p, hookArgsIgnore))
			}
		}
	}
	report := func(pos token.Pos, stmt ast.Stmt, msg string) {
		p := fset.Position(pos)
		start, end := fset.Position(stmt.Pos()).Line, fset.Position(stmt.End()).Line
		for l := start; l <= end; l++ {
			if ig, ok := ignores[l]; ok && ig.reason != "" {
				return
			}
		}
		if ig, ok := ignores[start-1]; ok && ig.alone && ig.reason != "" {
			return
		}
		out = append(out, fmt.Sprintf("%v: hook %v", p, msg))
	}

	ast.Inspect(f, func(n ast.Node) bool {
		var sig *types.Signature
		var ft *ast.FuncType
		var body *ast.BlockStmt
		switch fn := n.(type) {
		case *ast.FuncLit:
			sig, _ = info.TypeOf(fn).(*types.Signature)
			ft, body = fn.Type, fn.Body
		case *ast.FuncDecl:
			if obj, ok := info.Defs[fn.Name].(*types.Func); ok {
				sig, _ = obj.Type().(*types.Signature)
			}
			ft, body = fn.Type, fn.Body
		default:
			return true
		}
		if body == nil || sig == nil || !isHookSig(sig) || len(ft.Params.List[0].Names) == 0 {
			return true
		}
		param, ok := info.Defs[ft.Params.List[0].Names[0]].(*types.Var)
		if ok {
			checkHookBody(param, body, pkg, info, report)
		}
		return true
	})
	return out
}

// isHookSig reports whether sig has the shape of a Hook: func(...any), or func(...E) for a
// type parameter E, which can be instantiated as a Hook.
func isHookSig(sig *types.Signature) bool {
	if sig.Results().Len() != 0 || sig.Params().Len() != 1 || !sig.Variadic() {
		return false
	}
	s, ok := sig.Params().At(0).Type().(*types.Slice)
	if !ok {
		return false
	}
	if _, ok := s.Elem().(*types.TypeParam); ok {
		return true
	}
	return types.Identical(s.Elem(), types.Universe.Lookup("any").Type())
}

// checkHookBody allows args only as an element read, len(args), cap(args), range args, an
// argument that is copied (append(s, args...), fmt, hookArgsAllowed), or forwarded as
// f(args...) to another hook, which is checked in turn.
func checkHookBody(param *types.Var, body *ast.BlockStmt, pkg *types.Package, info *types.Info, report func(token.Pos, ast.Stmt, string)) {
	var stack []ast.Node
	ast.Inspect(body, func(n ast.Node) bool {
		if n == nil {
			stack = stack[:len(stack)-1]
			return true
		}
		stack = append(stack, n)
		id, ok := n.(*ast.Ident)
		if !ok || info.Uses[id] != param {
			return true
		}
		var stmt ast.Stmt
		for _, s := range slices.Backward(stack) {
			if st, ok := s.(ast.Stmt); ok {
				stmt = st
				break
			}
		}
		rep := func(msg string) { report(id.Pos(), stmt, msg) }
		for _, s := range stack[:len(stack)-1] {
			if _, ok := s.(*ast.FuncLit); ok {
				rep("uses args inside a closure, which may run after the hook returns")
				return true
			}
		}
		// parent returns the nearest ancestor of stack[i] that isn't a paren, and the child of it
		// that stack[i] is or is wrapped in
		parent := func(i int) (ast.Node, ast.Node, int) {
			child := stack[i]
			for i--; i >= 0; i-- {
				if _, ok := stack[i].(*ast.ParenExpr); !ok {
					return stack[i], child, i
				}
				child = stack[i]
			}
			return nil, nil, -1
		}

		p, child, pi := parent(len(stack) - 1)
		switch p := p.(type) {
		case *ast.IndexExpr:
			if p.X != child {
				break
			}
			gp, pchild, _ := parent(pi)
			switch gp := gp.(type) {
			case *ast.AssignStmt:
				if slices.Contains(gp.Lhs, pchild.(ast.Expr)) {
					rep("assigns to an element of args")
				}
			case *ast.RangeStmt:
				if gp.Key == pchild || gp.Value == pchild {
					rep("assigns to an element of args")
				}
			case *ast.UnaryExpr:
				if gp.Op == token.AND {
					rep("takes the address of an element of args")
				}
			}
			return true
		case *ast.RangeStmt:
			if p.X == child {
				return true
			}
		case *ast.CallExpr:
			if msg := checkCallArg(p, child, stack[:pi], pkg, info); msg != "" {
				rep(msg)
			}
			return true
		}
		rep("uses args whole (passed on, sliced or stored)")
		return true
	})
}

// checkCallArg checks args passed as child to call, whose ancestors are stack. It returns the
// violation, if any.
func checkCallArg(call *ast.CallExpr, child ast.Node, stack []ast.Node, pkg *types.Package, info *types.Info) string {
	const whole = "uses args whole (passed on, sliced or stored)"
	if call.Fun == child {
		return whole
	}
	spread := call.Ellipsis.IsValid() && call.Args[len(call.Args)-1] == child
	fun := ast.Unparen(call.Fun)
	if ix, ok := fun.(*ast.IndexExpr); ok { // explicit instantiation
		fun = ix.X
	}
	var obj types.Object
	switch fun := fun.(type) {
	case *ast.Ident:
		obj = info.Uses[fun]
	case *ast.SelectorExpr:
		obj = info.Uses[fun.Sel]
	}
	switch obj := obj.(type) {
	case *types.Builtin:
		switch obj.Name() {
		case "len", "cap":
			return ""
		case "append":
			if spread {
				return ""
			}
		}
		return whole
	case *types.Func:
		if obj.Pkg() != nil && obj.Pkg().Path() == "fmt" || hookArgsAllowed[obj.Origin().FullName()] {
			return ""
		}
	}
	if !spread || len(call.Args) != 1 {
		return whole
	}
	if len(stack) > 0 {
		if _, ok := stack[len(stack)-1].(*ast.GoStmt); ok {
			return "forwards args to a goroutine, which may run after the hook returns"
		}
	}
	sig, ok := info.TypeOf(call.Fun).(*types.Signature)
	if !ok || !isHookSig(sig) {
		return "forwards args to a function that's not a hook"
	}
	// a hook declared outside the repo isn't checked
	if fn, ok := obj.(*types.Func); ok && fn.Pkg() != pkg &&
		!strings.HasPrefix(fn.Pkg().Path(), modulePath+"/internal/") &&
		!strings.HasPrefix(fn.Pkg().Path(), modulePath+"/pkg/") {
		return "forwards args to a function that's not a hook"
	}
	return ""
}
