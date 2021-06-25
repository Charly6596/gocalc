package main

import (
	"fmt"
	"gocalc/repl"

	"github.com/gowebapi/webapi"
	"github.com/gowebapi/webapi/html"
	"github.com/gowebapi/webapi/html/htmlevent"
)

type myDiv html.HTMLDivElement

func (div myDiv) Write(s []byte) (n int, err error) {
	elem := webapi.GetWindow().Document().CreateElement("xmp", nil)
	elem.SetInnerHTML(fmt.Sprintf("%s", s))
	div.AppendChild(&elem.Node)
	div.SetScrollTop(float64(div.ScrollHeight()))
	n = len(s)
	return
}

func main() {
	var ch chan string = make(chan string)

	element := webapi.GetWindow().Document().GetElementById("stdout")
	output := html.HTMLDivElementFromJS(element)
	out := myDiv(*output)

	element2 := webapi.GetWindow().Document().GetElementById("stdin")
	input := html.HTMLInputElementFromJS(element2)

	input.SetOnKeyPress(func(event *htmlevent.KeyboardEvent, currentTarget *html.HTMLElement) {
		if event.KeyCode() == 13 {
			iput := html.HTMLInputElementFromJS(currentTarget)
			ch <- input.Value()
			iput.SetValue("")
		}
	})

	welcome := "GoCalc. A command line calculator written in Go"
	output.SetTextContent(&welcome)

	repl.Start(ch, out)
}
