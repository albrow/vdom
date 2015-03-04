package main

import (
	"bytes"
	"encoding/xml"
	"github.com/albrow/vdom"
	"github.com/gopherjs/gopherjs/js"
	"honnef.co/go/js/dom"
)

var (
	document = dom.GetWindow().Document()
)

func main() {
	js.Global.Call("suite", "Parsing html string", func() {
		js.Global.Call("benchmark", "xml Decode", func() {
			buf := bytes.NewBuffer([]byte("<ul><li>one</li><li>two</li><li>three</li></ul>"))
			dec := xml.NewDecoder(buf)
			for _, err := dec.Token(); err == nil; _, err = dec.Token() {
			}
		})

		js.Global.Call("benchmark", "Parse method", func() {
			vdom.Parse([]byte("<ul><li>one</li><li>two</li><li>three</li></ul>"))
		})

		js.Global.Call("benchmark", "document.createElement", func() {
			ul := document.CreateElement("ul")
			ul.SetInnerHTML("<li>one</li><li>two</li><li>three</li>")
		})
	})

	js.Global.Call("suite", "Diff algorithm", func() {

		oldTree, _ := vdom.Parse([]byte("<ul><li>one</li><li>two</li><li>three</li></ul>"))
		newTree, _ := vdom.Parse([]byte("<ul><li>uno</li><li>two</li><li>three</li></ul>"))

		js.Global.Call("benchmark", "Diff function", func() {
			vdom.Diff(oldTree, newTree)
		})
	})
}
