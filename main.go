// main.go
package main

import (
	"fmt"
	"gocalc/repl"
	"os"
)

func main() {
	fmt.Printf("GoCalc. A command line calculator written in Go\n")
	repl.Start(os.Stdin, os.Stdout)
}
