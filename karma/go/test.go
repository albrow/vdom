package main

import (
	"github.com/albrow/jasmine"
	"github.com/albrow/vdom"
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/jquery"
	"honnef.co/go/js/dom"
)

var (
	document = dom.GetWindow().Document()
	jq       = jquery.NewJQuery
)

func main() {
	jasmine.Describe("Tests", func() {
		jasmine.It("can be loaded", func() {
			jasmine.Expect(true).ToBe(true)
		})
	})

	jasmine.Describe("Selector", func() {

		// sandbox is a div with id = sandbox. It will be
		// created and cleaned up for each test.
		var sandbox dom.Element

		jasmine.BeforeEach(func() {
			if sandbox == nil {
				sandbox = document.CreateElement("div")
				sandbox.SetAttribute("id", "sandbox")
			}
			document.QuerySelector("body").AppendChild(sandbox)
		})

		jasmine.AfterEach(func() {
			document.QuerySelector("body").RemoveChild(sandbox)
		})

		jasmine.It("works with a single root element", func() {
			// Parse some source html into a tree
			html := "<div></div>"
			tree := setUpDOM(html, sandbox)
			testSelectors(tree, sandbox)
		})

		jasmine.It("works with a ul and nested lis", func() {
			// Parse some html into a tree
			html := "<ul><li>one</li><li>two</li><li>three</li></ul>"
			tree := setUpDOM(html, sandbox)
			testSelectors(tree, sandbox)
		})

		jasmine.It("works with a form with autoclosed tags", func() {
			// Parse some html into a tree
			html := `<form method="post"><input type="text" name="firstName"><input type="text" name="lastName"></form>`
			tree := setUpDOM(html, sandbox)
			testSelectors(tree, sandbox)
		})
	})
}

// setUpDOM parses html into a virtual tree, then adds it to the
// actual dom by appending to sandbox. It returns both the virtual
// tree.
func setUpDOM(html string, sandbox dom.Element) *vdom.Tree {
	// Parse the html into a virtual tree
	vtree, err := vdom.Parse([]byte(html))
	jasmine.Expect(err).ToBe(nil)

	// Add html to the actual DOM
	sandbox.SetInnerHTML(html)
	return vtree
}

func expectExistsInDom(el dom.Element) {
	jqEl := jq(el)
	js.Global.Call("expect", jqEl).Call("toExist")
	js.Global.Call("expect", jqEl).Call("toBeInDOM")
}

func testSelectors(tree *vdom.Tree, root dom.Element) {
	for i, vRoot := range tree.Roots {
		if vEl, ok := vRoot.(*vdom.Element); ok {
			// If vRoot is an element, test its Selector method
			expectedEl := root.ChildNodes()[i].(dom.Element)
			testSelector(vEl, root, expectedEl)
		}
	}
}

func testSelector(vEl *vdom.Element, root, expectedEl dom.Element) {
	gotEl := root.QuerySelector(vEl.Selector())
	expectExistsInDom(gotEl)
	jasmine.Expect(gotEl).ToEqual(expectedEl)
	// Test vEl's children recursively
	for i, vChild := range vEl.Children() {
		if vChildEl, ok := vChild.(*vdom.Element); ok {
			// If vRoot is an element, test its Selector method
			expectedChildEl := expectedEl.ChildNodes()[i].(dom.Element)
			testSelector(vChildEl, root, expectedChildEl)
		}
	}
}
