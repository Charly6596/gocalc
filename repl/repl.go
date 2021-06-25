package repl

import (
	"fmt"
	"gocalc/evaluator"
	"io"
)

const PROMPT = ">>> "

func Start(in chan string, out io.Writer) {
	ev := evaluator.New()

	for {
		fmt.Printf(PROMPT)

		line := <-in

		res := ev.Eval(line)

		if res != nil {
			io.WriteString(out, res.String())
		}
	}
}
