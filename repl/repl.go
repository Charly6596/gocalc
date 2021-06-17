package repl

import (
	"bufio"
	"fmt"
	"gocalc/evaluator"
	"io"
)

const PROMPT = ">>> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	ev := evaluator.New()

	for {
		fmt.Printf(PROMPT)

		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()

		res := ev.Eval(line)

		if res != nil {
			io.WriteString(out, res.String())
			io.WriteString(out, "\n")
		}
	}
}
