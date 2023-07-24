package repl

import (
	"bufio"
	"fmt"
	"io"
	"log"

	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/gcs"
	"github.com/genshinsim/gcsim/pkg/gcs/ast"
)

const Prompt = ">> "

func Eval(s string, log *log.Logger) {

	simActions := make(chan *action.ActionEval)
	done := make(chan bool)
	go handleSimActions(simActions, done)

	p := ast.New(s)
	res, err := p.Parse()

	if err != nil {
		fmt.Println("Error parsing input:")
		fmt.Printf("\t%v\n", err)
		return
	}

	if len(res.Errors) != 0 {
		fmt.Println("Errors encountered in config:")
		for _, v := range res.Errors {
			fmt.Printf("\t%v\n", v)
		}
	}

	fmt.Println("Program parsed:")
	fmt.Println(res.Program.String())

	if len(res.Errors) != 0 {
		//don't run the program if there are errors
		return
	}
	fmt.Println("Running program...:")
	eval := gcs.Eval{
		AST:  res.Program,
		Next: done,
		Work: simActions,
		Log:  log,
	}

	result := eval.Run()

	if eval.Err() != nil {
		fmt.Printf("Program finished with err: %v", eval.Err())
	}
	fmt.Println(result.Inspect())
}

func Start(in io.Reader, out io.Writer, log *log.Logger, showProgram bool) {
	scanner := bufio.NewScanner(in)

	for {
		simActions := make(chan *action.ActionEval)
		next := make(chan bool)
		go handleSimActions(simActions, next)

		fmt.Print(Prompt)
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
			Next: next,
			Work: simActions,
			Log:  log,
		}
		result := eval.Run()

		if eval.Err() != nil {
			fmt.Printf("Program finished with err: %v", eval.Err())
		}
		fmt.Println(result.Inspect())
	}
}

func handleSimActions(in chan *action.ActionEval, next chan bool) {
	for {
		next <- true
		x, ok := <-in
		if !ok {
			return
		}
		fmt.Printf("\tExecuting: %v\n", x.Action.String())
	}
}
