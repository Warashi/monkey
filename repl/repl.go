package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/Warashi/implement-interpreter-with-go/evaluator"
	"github.com/Warashi/implement-interpreter-with-go/lexer"
	"github.com/Warashi/implement-interpreter-with-go/parser"
)

const PROMPT = ">> "

func Start(r io.Reader, w io.Writer) {
	s := bufio.NewScanner(r)
	fmt.Fprint(w, PROMPT)
	for s.Scan() {
		l := lexer.New(s.Text())
		p := parser.New(l)
		program := p.Parse()
		if errs := p.Errors(); len(errs) != 0 {
			for _, err := range errs {
				fmt.Fprintf(w, "\t%s\n", err)
				continue
			}
		}
		if e := evaluator.Eval(program); e != nil {
			fmt.Fprintln(w, e.Inspect())
		}
		fmt.Fprint(w, PROMPT)
	}
}
