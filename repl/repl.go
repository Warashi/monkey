package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/Warashi/implement-interpreter-with-go/lexer"
	"github.com/Warashi/implement-interpreter-with-go/token"
)

const PROMPT = ">> "

func Start(r io.Reader, w io.Writer) {
	s := bufio.NewScanner(r)
	fmt.Print(PROMPT)
	for s.Scan() {
		l := lexer.New(s.Text())

		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Printf("%+v\n", tok)
		}
	}
}
