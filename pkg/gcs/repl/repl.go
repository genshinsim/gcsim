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

	p := ast.New(s)
	res, gcsl, err := p.Parse()

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
	fmt.Println(gcsl.String())

	if len(res.Errors) != 0 {
		//don't run the program if there are errors
		return
	}
	fmt.Println("Running program...:")
	eval, _ := gcs.NewEvaluator(gcsl, nil)
	eval.Log = log
	resultChan := make(chan gcs.Obj)
	errChan := make(chan error)
	go func() {
		res, err := eval.Run()
		// fmt.Printf("done with result: %v, err: %v\n", res, err)
		resultChan <- res
		errChan <- err
	}()

	for {
		a, err := eval.NextAction()
		if a == nil {
			break
		}
		if err != nil {
			fmt.Printf("Program finished with err: %v", err)
			return
		}
	}

	result := <-resultChan
	err = <-errChan
	if err != nil {
		fmt.Printf("Program finished with err: %v", err)
	}
	fmt.Println(result.Inspect())
}

func Start(in io.Reader, out io.Writer, log *log.Logger, showProgram bool) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Print(Prompt)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		Eval(scanner.Text(), log)
	}
}
