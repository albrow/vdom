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
			src := "<div></div>"
			tree, err := vdom.Parse([]byte(src))
			jasmine.Expect(err).ToBe(nil)

			// Add html to the actual DOM
			el := jquery.ParseHTML(src)
			jq(sandbox).Append(el)
			js.Global.Call("expect", el).Call("toExist")
			js.Global.Call("expect", el).Call("toBeInDOM")

			// Use the selector calculated by the virutal element
			// to select the corresponding real element in the DOM
			vEl := tree.Roots[0].(*vdom.Element)
			// Recall that we need to append the partial selector to some
			// parent selector, in this case the sandbox div.
			gotEl := jq("#sandbox" + vEl.PartialSelector())
			js.Global.Call("expect", gotEl).Call("toExist")
			js.Global.Call("expect", gotEl).Call("toBeInDOM")
			jasmine.Expect(el).ToEqual(gotEl)
		})

		jasmine.It("works with a ul and nested lis", func() {
			// Parse some html into a tree
			src := "<ul><li>one</li><li>two</li><li>three</li></ul>"
			tree, err := vdom.Parse([]byte(src))
			jasmine.Expect(err).ToBe(nil)

			// Add html to the actual DOM
			el := jquery.ParseHTML(src)
			jq(sandbox).Append(el)
			js.Global.Call("expect", el).Call("toExist")
			js.Global.Call("expect", el).Call("toBeInDOM")

			// Use the selector calculated by the virutal element
			// to select the corresponding real element in the DOM
			vEl := tree.Roots[0].(*vdom.Element)
			// Recall that we need to append the partial selector to some
			// parent selector, in this case the sandbox div.
			gotEl := jq("#sandbox" + vEl.PartialSelector())
			js.Global.Call("expect", gotEl).Call("toExist")
			js.Global.Call("expect", gotEl).Call("toBeInDOM")
			jasmine.Expect(el).ToEqual(gotEl)

			// Now do the same thing for each child li element
			for _, vNode := range vEl.Children() {
				vLi := vNode.(*vdom.Element)
				gotLi := jq("#sandbox" + vLi.PartialSelector())
				js.Global.Call("expect", gotLi).Call("toExist")
				js.Global.Call("expect", gotLi).Call("toBeInDOM")
			}
		})

		jasmine.It("works with a form with autoclosed tags", func() {
			// Parse some html into a tree
			src := `<form method="post"><input type="text" name="firstName"><input type="text" name="lastName"></form>`
			tree, err := vdom.Parse([]byte(src))
			jasmine.Expect(err).ToBe(nil)

			// Add html to the actual DOM
			el := jquery.ParseHTML(src)
			jq(sandbox).Append(el)
			js.Global.Call("expect", el).Call("toExist")
			js.Global.Call("expect", el).Call("toBeInDOM")

			// Use the selector calculated by the virutal element
			// to select the corresponding real element in the DOM
			vEl := tree.Roots[0].(*vdom.Element)
			// Recall that we need to append the partial selector to some
			// parent selector, in this case the sandbox div.
			gotEl := jq("#sandbox" + vEl.PartialSelector())
			js.Global.Call("expect", gotEl).Call("toExist")
			js.Global.Call("expect", gotEl).Call("toBeInDOM")
			jasmine.Expect(el).ToEqual(gotEl)

			// Now do the same thing for each child li element
			for _, vNode := range vEl.Children() {
				vInput := vNode.(*vdom.Element)
				gotInput := jq("#sandbox" + vInput.PartialSelector())
				js.Global.Call("expect", gotInput).Call("toExist")
				js.Global.Call("expect", gotInput).Call("toBeInDOM")
			}
		})
	})
}
