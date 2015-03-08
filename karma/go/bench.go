package main

import (
	"fmt"
	"github.com/albrow/vdom"
	"github.com/gopherjs/gopherjs/js"
	"honnef.co/go/js/dom"
	"strconv"
	"strings"
)

var (
	document = dom.GetWindow().Document()
	sandbox  dom.Element
)

func init() {
	sandbox = document.CreateElement("div")
	sandbox.SetID("sandbox")
	document.QuerySelector("body").AppendChild(sandbox)
}

func main() {
	benchmarkRenderList(sandbox, 3, 1)
	benchmarkRenderList(sandbox, 10, 1)
	benchmarkRenderList(sandbox, 10, 5)
	benchmarkRenderList(sandbox, 10, 10)
	benchmarkRenderList(sandbox, 100, 1)
	benchmarkRenderList(sandbox, 100, 50)
	benchmarkRenderList(sandbox, 100, 100)
}

// generateList returns html for an unordered list of n items.
func generateList(n int) []byte {
	result := []byte("<ul>")
	for i := 0; i < n; i++ {
		result = append(result, []byte("<li>")...)
		result = append(result, []byte(strconv.Itoa(i))...)
		result = append(result, []byte("</li>")...)
	}
	result = append(result, []byte("</ul>")...)
	return result
}

// changeList accepts html that contains an unordered list of n
// items and changes the first n items. It returns the new html
// with the changes applied.
func changeList(old []byte, n int) []byte {
	newHTML := string(old)
	for i := 0; i < n && i < len(old); i++ {
		oldItem := fmt.Sprintf("<li>%d</li>", i)
		newItem := fmt.Sprintf("<li>new %d</li>", i)
		newHTML = strings.Replace(newHTML, oldItem, newItem, 1)
	}
	return []byte(newHTML)
}

// benchmarkRenderList generates html for an unordered list of numItems
// items and sets it as the inner html of root. Then it re-renders the list,
// making numChanges changes. It compares re-rendering with the virtual DOM
// vs. re-rendering via setHTML.
func benchmarkRenderList(root dom.Element, numItems, numChanges int) {
	oldHTML := generateList(numItems)
	root.SetInnerHTML(string(oldHTML))
	newHTML := changeList(oldHTML, numChanges)
	oldTree, err := vdom.Parse(oldHTML)
	if err != nil {
		panic(err)
	}
	suiteName := fmt.Sprintf("Re-render list with %d items after %d changes", numItems, numChanges)
	js.Global.Call("suite", suiteName, func() {

		js.Global.Call("benchmark", "with SetInnerHTML", func() {
			root.SetInnerHTML(string(newHTML))
		})

		js.Global.Call("benchmark", "with virtual DOM", func() {
			newTree, err := vdom.Parse(newHTML)
			if err != nil {
				panic(err)
			}
			patches, err := vdom.Diff(oldTree, newTree)
			if err != nil {
				panic(err)
			}
			if err := patches.Patch(root); err != nil {
				panic(err)
			}
		})
	}, js.MakeWrapper(map[string]interface{}{
		"teardown": func() {
			root.SetInnerHTML("")
		},
	}))
}
