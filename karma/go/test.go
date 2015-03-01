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

		jasmine.XIt("works with a single root element", func() {
			// Parse some source html into a tree
			html := "<div></div>"
			tree := setUpDOM(html, sandbox)

			// Use the selector calculated by the virutal element
			// to select the corresponding real element in the DOM
			vEl := tree.Roots[0].(*vdom.Element)
			gotEl := sandbox.QuerySelector(vEl.Selector())
			expectedEl := sandbox.ChildNodes()[0]
			expectExistsInDom(gotEl)
			jasmine.Expect(gotEl).ToEqual(expectedEl)
		})

		jasmine.It("works with a ul and nested lis", func() {
			// Parse some html into a tree
			html := "<ul><li>one</li><li>two</li><li>three</li></ul>"
			tree := setUpDOM(html, sandbox)

			// Use the selector calculated by the virutal element
			// to select the corresponding real element in the DOM
			vEl := tree.Roots[0].(*vdom.Element)
			gotEl := sandbox.QuerySelector(vEl.Selector())
			expectedEl := sandbox.ChildNodes()[0]
			expectExistsInDom(gotEl)
			jasmine.Expect(gotEl).ToEqual(expectedEl)

			// Now do the same thing for each child li element
			for i, vNode := range vEl.Children() {
				vLi := vNode.(*vdom.Element)
				gotLi := sandbox.QuerySelector(vLi.Selector())
				expectedLi := expectedEl.ChildNodes()[i]
				expectExistsInDom(gotLi)
				jasmine.Expect(gotLi).ToEqual(expectedLi)
			}
		})

		jasmine.It("works with a form with autoclosed tags", func() {
			// Parse some html into a tree
			html := `<form method="post"><input type="text" name="firstName"><input type="text" name="lastName"></form>`
			tree := setUpDOM(html, sandbox)

			// Use the selector calculated by the virutal element
			// to select the corresponding real element in the DOM
			vEl := tree.Roots[0].(*vdom.Element)
			gotEl := sandbox.QuerySelector(vEl.Selector())
			expectedEl := sandbox.ChildNodes()[0]
			expectExistsInDom(gotEl)
			jasmine.Expect(gotEl).ToEqual(expectedEl)

			// Now do the same thing for each child input element
			for i, vNode := range vEl.Children() {
				vInput := vNode.(*vdom.Element)
				gotInput := sandbox.QuerySelector(vInput.Selector())
				expectedInput := expectedEl.ChildNodes()[i]
				expectExistsInDom(gotInput)
				jasmine.Expect(gotInput).ToEqual(expectedInput)
			}
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
