package main

import (
	"github.com/albrow/jasmine"
	"github.com/albrow/vdom"
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/jquery"
)

var jq = jquery.NewJQuery

func main() {
	jasmine.Describe("Tests", func() {
		jasmine.It("can be loaded", func() {
			jasmine.Expect(true).ToBe(true)
		})

		jasmine.It("can create a sandbox", func() {
			// Try creating a sandbox
			sandbox := js.Global.Call("sandbox")
			js.Global.Call("setFixtures", sandbox)
			// Make sure the sandbox was created correctly
			js.Global.Call("expect", sandbox).Call("toExist")
			js.Global.Call("expect", sandbox).Call("toBeInDOM")
		})
	})

	jasmine.Describe("PartialSelector", func() {

		// sandbox is a div with id = sandbox. It will be
		// automatically created and cleaned up for each test.
		var sandbox *js.Object

		jasmine.BeforeEach(func() {
			sandbox = js.Global.Call("sandbox")
			js.Global.Call("setFixtures", sandbox)
		})

		jasmine.It("works with a single root element", func() {
			// Parse some source html into a tree
			html := "<div></div>"
			tree, el := setUpDOM(html, sandbox)

			// Use the selector calculated by the virutal element
			// to select the corresponding real element in the DOM
			vEl := tree.Roots[0].(*vdom.Element)
			// Recall that we need to append the partial selector to some
			// parent selector, in this case the sandbox div.
			gotEl := jq("#sandbox" + vEl.PartialSelector())
			expectExistsInDom(jq(el))
			jasmine.Expect(el).ToEqual(gotEl)
		})

		jasmine.It("works with a ul and nested lis", func() {
			// Parse some html into a tree
			html := "<ul><li>one</li><li>two</li><li>three</li></ul>"
			tree, el := setUpDOM(html, sandbox)

			// Use the selector calculated by the virutal element
			// to select the corresponding real element in the DOM
			vEl := tree.Roots[0].(*vdom.Element)
			// Recall that we need to append the partial selector to some
			// parent selector, in this case the sandbox div.
			gotEl := jq("#sandbox" + vEl.PartialSelector())
			expectExistsInDom(gotEl)
			jasmine.Expect(el).ToEqual(gotEl)

			// Now do the same thing for each child li element
			for _, vNode := range vEl.Children() {
				vLi := vNode.(*vdom.Element)
				gotLi := jq("#sandbox" + vLi.PartialSelector())
				expectExistsInDom(gotLi)
			}
		})

		jasmine.It("works with a form with autoclosed tags", func() {
			// Parse some html into a tree
			html := `<form method="post"><input type="text" name="firstName"><input type="text" name="lastName"></form>`
			tree, el := setUpDOM(html, sandbox)

			// Use the selector calculated by the virutal element
			// to select the corresponding real element in the DOM
			vEl := tree.Roots[0].(*vdom.Element)
			// Recall that we need to append the partial selector to some
			// parent selector, in this case the sandbox div.
			gotEl := jq("#sandbox" + vEl.PartialSelector())
			expectExistsInDom(gotEl)
			jasmine.Expect(el).ToEqual(gotEl)

			// Now do the same thing for each child li element
			for _, vNode := range vEl.Children() {
				vInput := vNode.(*vdom.Element)
				gotInput := jq("#sandbox" + vInput.PartialSelector())
				expectExistsInDom(gotInput)
			}
		})
	})
}

// setUpDOM parses html into a virtual tree, then adds it to the
// actual dom by appending to sandbox. It returns both the virtual
// tree and the root element in the actual DOM.
func setUpDOM(html string, sandbox *js.Object) (tree *vdom.Tree, root []interface{}) {
	// Parse the html into a virtual tree
	tree, err := vdom.Parse([]byte(html))
	jasmine.Expect(err).ToBe(nil)

	// Add html to the actual DOM
	el := jquery.ParseHTML(html)
	jq(sandbox).Append(el)
	expectExistsInDom(jq(el))

	return tree, el
}

func expectExistsInDom(el jquery.JQuery) {
	js.Global.Call("expect", el).Call("toExist")
	js.Global.Call("expect", el).Call("toBeInDOM")
}
