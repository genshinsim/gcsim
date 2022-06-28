package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/davecgh/go-spew/spew"
	"golang.org/x/tools/go/ast/astutil"
)

func main() {
	// read every file in directory

	//check args if we want to print ast out
	dump := false
	if len(os.Args) > 1 {
		dump = true
	}

	files, err := ioutil.ReadDir("./")
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		err = fix(file.Name(), dump)
		if err != nil {
			panic(err)
		}
	}
}

func fix(path string, dump bool) error {
	//do nothing
	if filepath.Ext(path) != ".go" {
		return nil
	}
	fs := token.NewFileSet()
	f, err := parser.ParseFile(fs, path, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	if dump {
		spew.Dump(f)
	}

	// fix any core package names
	astutil.Apply(f, func(cr *astutil.Cursor) bool {
		found, next := fixStatMod(cr.Node())
		if !found {
			return true
		}
		cr.Replace(next)
		return false
	}, nil)

	// Print result
	out, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer out.Close()
	printer.Fprint(out, fs, f)

	// printer.Fprint(os.Stdout, fs, f)
	return nil
}

func fixAttackMod(n ast.Node) (bool, ast.Node) {
	if expr, ok := n.(*ast.ExprStmt); ok {
		block, ok := expr.X.(*ast.CallExpr)
		if !ok {
			return false, nil
		}

		//FUN should be a SelectorExpr
		fun, ok := block.Fun.(*ast.SelectorExpr)
		if !ok {
			return false, nil
		}

		//Sel should be AddPreDamageMod
		if fun.Sel.Name != "AddPreDamageMod" {
			return false, nil
		}

		fmt.Println("found pre damage block")

		//work through the args and find amount, expiry, and key
		//args should be len 1
		if len(block.Args) != 1 {
			fmt.Println("unexpected args length > 1")
			return false, nil
		}

		//check to make sure it's a CompositeLit
		lit, ok := block.Args[0].(*ast.CompositeLit)
		if !ok {
			fmt.Println("unexpected arg type, not a composite lit")
			return false, nil
		}

		//loop through Elts to find amount, expiry and key
		var amount, expiry, key ast.Expr
		for _, v := range lit.Elts {
			t, ok := v.(*ast.KeyValueExpr)
			if !ok {
				continue
			}
			switch t.Key.(*ast.Ident).Name {
			case "Amount":
				amount = t.Value
			case "Expiry":
				expiry = t.Value
			case "Key":
				key = t.Value
			}
		}

		caller, ok := fun.X.(*ast.Ident)
		if !ok {
			fmt.Println("unexpected fun.X type, not an ident")
		}

		next := &ast.ExprStmt{
			X: &ast.CallExpr{
				Fun:  ast.NewIdent(fmt.Sprintf("%v.AddAttackMod", caller.Name)),
				Args: []ast.Expr{key, expiry, amount},
			},
		}

		return true, next
	}
	return false, nil
}

func fixStatMod(n ast.Node) (bool, ast.Node) {
	if expr, ok := n.(*ast.ExprStmt); ok {
		block, ok := expr.X.(*ast.CallExpr)
		if !ok {
			return false, nil
		}

		//FUN should be a SelectorExpr
		fun, ok := block.Fun.(*ast.SelectorExpr)
		if !ok {
			return false, nil
		}

		if fun.Sel.Name != "AddStatMod" {
			return false, nil
		}

		//work through the args and find amount, expiry, and key
		//args should be len 1
		if len(block.Args) != 4 {
			fmt.Println("unexpected args length != 4")
			return false, nil
		}

		//the args are BasicLit, UnaryExpr, SelectorExpr, FuncLit

		//base stuff -> key + dur
		base := &ast.KeyValueExpr{
			Key: ast.NewIdent("Base"),
			// Value: &ast.CompositeLit{
			// 	Type: &ast.SelectorExpr{
			// 		X:   ast.NewIdent("modifier"),
			// 		Sel: ast.NewIdent("Base"),
			// 	},
			// 	Elts: []ast.Expr{
			// 		&ast.KeyValueExpr{
			// 			Key:   ast.NewIdent("ModKey"),
			// 			Value: block.Args[0],
			// 		},
			// 		&ast.KeyValueExpr{
			// 			Key:   ast.NewIdent("Dur"),
			// 			Value: block.Args[1],
			// 		},
			// 	},
			// },
			Value: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   ast.NewIdent("modifier"),
					Sel: ast.NewIdent("NewBase"),
				},
				Args: []ast.Expr{
					block.Args[0],
					block.Args[1],
				},
			},
		}
		affstat := &ast.KeyValueExpr{
			Key:   ast.NewIdent("AffectedStat"),
			Value: block.Args[2],
		}
		amtfun := &ast.KeyValueExpr{
			Key:   ast.NewIdent("Amount"),
			Value: block.Args[3],
		}

		//should be compositelit with ident
		next := &ast.ExprStmt{
			X: &ast.CallExpr{
				Fun: block.Fun,
				Args: []ast.Expr{
					&ast.CompositeLit{
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent("character"),
							Sel: ast.NewIdent("StatMod"),
						},
						Elts: []ast.Expr{
							base, affstat, amtfun},
					},
				},
			},
		}

		return true, next
	}
	return false, nil
}
