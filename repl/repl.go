package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/Warashi/monkey/compiler"
	"github.com/Warashi/monkey/lexer"
	"github.com/Warashi/monkey/parser"
	"github.com/Warashi/monkey/vm"
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
		compiler := compiler.New()
		if err := compiler.Compile(program); err != nil {
			fmt.Fprintf(w, "failed compile: %+v", err)
		}

		machine := vm.New(compiler.Bytecode())
		if err := machine.Run(); err != nil {
			fmt.Fprintf(w, "failed run: %+v", err)
		}

		fmt.Fprintln(w, machine.LastPopedStackElem().Inspect())
		fmt.Fprint(w, PROMPT)
	}
}
