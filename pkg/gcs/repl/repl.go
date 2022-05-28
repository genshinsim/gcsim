package repl

import (
	"bufio"
	"fmt"
	"io"
	"log"

	"github.com/genshinsim/gcsim/pkg/gcs"
	"github.com/genshinsim/gcsim/pkg/gcs/ast"
)

const Prompt = ">> "

func Eval(s string, log *log.Logger) {

	simActions := make(chan ast.ActionStmt)
	done := make(chan bool)
	go handleSimActions(simActions, done)

	p := ast.New(s)
	res, err := p.Parse()

	if err != nil {
		fmt.Println("Error parsing input:")
		fmt.Printf("\t%v\n", err)
		return
	}

	fmt.Println("Program parsed:")
	fmt.Println(res.Program.String())

	eval := gcs.Eval{
		AST:  res.Program,
		Next: done,
		Work: simActions,
		Log:  log,
	}

	result := eval.Run()

	fmt.Println("Program results:")
	fmt.Println(result.Inspect())
}

func Start(in io.Reader, out io.Writer, log *log.Logger, showProgram bool) {
	scanner := bufio.NewScanner(in)

	simActions := make(chan ast.ActionStmt)
	done := make(chan bool)
	go handleSimActions(simActions, done)

	for {
		fmt.Printf(Prompt)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		p := ast.New(line)
		res, err := p.Parse()

		if err != nil {
			fmt.Println("Error parsing input:")
			fmt.Printf("\t%v\n", err)
			continue
		}

		if showProgram {
			fmt.Println("Program parsed:")
			fmt.Println(res.Program.String())
		}

		eval := gcs.Eval{
			AST:  res.Program,
			Next: done,
			Work: simActions,
			Log:  log,
		}

		result := eval.Run()

		fmt.Println(result.Inspect())
	}
}

func handleSimActions(in chan ast.ActionStmt, done chan bool) {
	for {
		x := <-in
		fmt.Printf("\tExecuting: %v\n", x.String())
		done <- true
	}
}
