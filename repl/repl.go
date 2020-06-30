package repl

import (
	"bufio"
	"fmt"
	"github.com/DarkoKlisuric/interpreter/lexer"
	"github.com/DarkoKlisuric/interpreter/token"
	"io"
)

const PROMPT = ">> "

// REPL stands for “Read Eval Print Loop”
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Println(PROMPT)
		scanned := scanner.Scan()

		if !scanned {
			return
		}

		line := scanner.Text()

		l := lexer.New(line)

		// Looping through inpunts (tokens)
		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Printf("%+v\n", tok)
		}
	}
}
