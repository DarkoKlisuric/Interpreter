package repl

import (
	"bufio"
	"fmt"
	"github.com/DarkoKlisuric/interpreter/evaluator"
	"github.com/DarkoKlisuric/interpreter/lexer"
	"github.com/DarkoKlisuric/interpreter/object"
	"github.com/DarkoKlisuric/interpreter/parser"
	"io"
)

const Prompt = ">> "

// REPL stands for “Read Eval Print Loop”
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for {
		fmt.Println(Prompt)
		scanned := scanner.Scan()

		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program, env)

		if evaluated != nil {
			_, _ = io.WriteString(out, evaluated.Inspect())
			_, _ = io.WriteString(out, "\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		_, _ = io.WriteString(out, "\t"+msg+"\n")
	}
}
