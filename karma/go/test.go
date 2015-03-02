package main

import (
	"github.com/JohannWeging/jasmine"
	"github.com/albrow/vdom"
	"honnef.co/go/js/dom"
)

var (
	document = dom.GetWindow().Document()
)

func main() {

	var body dom.Element

	jasmine.BeforeEach(func() {
		if body == nil {
			body = document.QuerySelector("body")
		}
	})

	jasmine.AfterEach(func() {
		body.SetInnerHTML("")
	})

	jasmine.Describe("Tests", func() {
		jasmine.It("can be loaded", func() {
			jasmine.Expect(true).ToBe(true)
		})
	})

	jasmine.Describe("Selector", func() {

		jasmine.It("works with a single root element", func() {
			// Parse some source html into a tree
			html := "<div></div>"
			tree := setUpDOM(html, body)
			testSelectors(tree, body)
		})

		jasmine.It("works with a ul and nested lis", func() {
			// Parse some html into a tree
			html := "<ul><li>one</li><li>two</li><li>three</li></ul>"
			tree := setUpDOM(html, body)
			testSelectors(tree, body)
		})

		jasmine.It("works with a form with autoclosed tags", func() {
			// Parse some html into a tree
			html := `<form method="post"><input type="text" name="firstName"><input type="text" name="lastName"></form>`
			tree := setUpDOM(html, body)
			testSelectors(tree, body)
		})
	})

	jasmine.Describe("Append", func() {

		jasmine.It("works with a single root element", func() {
			testAppendRootPatcher(body, "<div></div>")
		})

		jasmine.It("works with a text root", func() {
			testAppendRootPatcher(body, "Text")
		})

		jasmine.It("works with a comment root", func() {
			testAppendRootPatcher(body, "<!--comment-->")
		})

		jasmine.It("works with nested siblings", func() {
			createAndApplyPatcher(body, "<ul><li>one</li><li>two</li></ul>", func(tree *vdom.Tree) vdom.Patcher {
				// Create a new tree, which only consists of a new li element
				// which we want to append
				newTree, err := vdom.Parse([]byte("<li>three</li>"))
				jasmine.Expect(err).ToBe(nil)
				// Create a patch manually
				return &vdom.Append{
					Child:  newTree.Roots[0],
					Parent: tree.Roots[0],
				}
			})
			// Test that the patch was applied
			ul := body.ChildNodes()[0].(*dom.HTMLUListElement)
			jasmine.Expect(ul.InnerHTML()).ToBe("<li>one</li><li>two</li><li>three</li>")
		})
	})

	jasmine.Describe("Replace", func() {

		jasmine.It("works with a single root element", func() {
			testReplaceRootPatcher(body, `<div id="old"></div>`, `<div id="new"></div>`)
		})

		jasmine.It("works with a text root", func() {
			testReplaceRootPatcher(body, "Old", "New")
		})

		jasmine.It("works with a comment root", func() {
			testReplaceRootPatcher(body, "<!--old-->", "<!--new-->")
		})

		jasmine.It("works with nested siblings", func() {
			createAndApplyPatcher(body, "<ul><li>one</li><li>two</li><li>three</li></ul>", func(tree *vdom.Tree) vdom.Patcher {
				// Create a new tree, which only consists of one of the lis
				// We want to change it from one to uno
				newTree, err := vdom.Parse([]byte("<li>uno</li>"))
				jasmine.Expect(err).ToBe(nil)
				// Create a patch manually
				return &vdom.Replace{
					Old: tree.Roots[0].Children()[0],
					New: newTree.Roots[0],
				}
			})
			// Test that the patch was applied
			ul := body.ChildNodes()[0].(*dom.HTMLUListElement)
			jasmine.Expect(ul.InnerHTML()).ToBe("<li>uno</li><li>two</li><li>three</li>")
		})
	})

	jasmine.Describe("Remove", func() {

		jasmine.It("works with a single root element", func() {
			testRemoveRootPatcher(body, "<div></div>")
		})

		jasmine.It("works with a text root", func() {
			testRemoveRootPatcher(body, "Text")
		})

		jasmine.It("works with a comment root", func() {
			testRemoveRootPatcher(body, "<!--comment-->")
		})

		jasmine.It("works with nested siblings", func() {
			createAndApplyPatcher(body, "<ul><li>one</li><li>two</li><li>three</li></ul>", func(tree *vdom.Tree) vdom.Patcher {
				return &vdom.Remove{
					Node: tree.Roots[0].Children()[1],
				}
			})
			// Test that the patch was applied by checking the innerHTML
			// property of the ul node.
			ul := body.ChildNodes()[0].(*dom.HTMLUListElement)
			jasmine.Expect(ul.InnerHTML()).ToBe("<li>one</li><li>three</li>")
		})
	})

	jasmine.Describe("SetAttr", func() {

		jasmine.It("works on a root element", func() {
			createAndApplyPatcher(body, "<div></div>", func(tree *vdom.Tree) vdom.Patcher {
				return &vdom.SetAttr{
					Node: tree.Roots[0],
					Attr: &vdom.Attr{
						Name:  "id",
						Value: "foo",
					},
				}
			})
			// Test that the patch was applied
			jasmine.Expect(body.InnerHTML()).ToBe(`<div id="foo"></div>`)
		})

		jasmine.It("works on a nested element", func() {
			createAndApplyPatcher(body, "<ul><li>one</li><li>two</li><li>three</li></ul>", func(tree *vdom.Tree) vdom.Patcher {
				return &vdom.SetAttr{
					Node: tree.Roots[0].Children()[1],
					Attr: &vdom.Attr{
						Name:  "data-value",
						Value: "two",
					},
				}
			})
			// Test that the patch was applied
			ul := body.ChildNodes()[0].(*dom.HTMLUListElement)
			jasmine.Expect(ul.InnerHTML()).ToBe(`<li>one</li><li data-value="two">two</li><li>three</li>`)
		})
	})

	jasmine.Describe("RemoveAttr", func() {

		jasmine.It("works on a root element", func() {
			createAndApplyPatcher(body, `<div id="foo"></div>`, func(tree *vdom.Tree) vdom.Patcher {
				return &vdom.RemoveAttr{
					Node:     tree.Roots[0],
					AttrName: "id",
				}
			})
			// Test that the patch was applied
			jasmine.Expect(body.InnerHTML()).ToBe("<div></div>")
		})

		jasmine.It("works on a nested element", func() {
			createAndApplyPatcher(body, `<ul><li>one</li><li data-value="two">two</li><li>three</li></ul>`, func(tree *vdom.Tree) vdom.Patcher {
				return &vdom.RemoveAttr{
					Node:     tree.Roots[0].Children()[1],
					AttrName: "data-value",
				}
			})
			// Test that the patch was applied
			ul := body.ChildNodes()[0].(*dom.HTMLUListElement)
			jasmine.Expect(ul.InnerHTML()).ToBe("<li>one</li><li>two</li><li>three</li>")
		})

	})
}

// setUpDOM parses html into a virtual tree, then adds it to the
// actual dom by appending to body. It returns both the virtual
// tree.
func setUpDOM(html string, body dom.Element) *vdom.Tree {
	// Parse the html into a virtual tree
	vtree, err := vdom.Parse([]byte(html))
	jasmine.Expect(err).ToBe(nil)
	// Add html to the actual DOM
	body.SetInnerHTML(html)
	return vtree
}

func expectExistsInDOM(el dom.Element) {
	jasmine.Expect(document.Contains(el)).ToBe(true)
}

func testSelector(vEl *vdom.Element, root, expectedEl dom.Element) {
	gotEl := root.QuerySelector(vEl.Selector())
	expectExistsInDOM(gotEl)
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

func testSelectors(tree *vdom.Tree, root dom.Element) {
	for i, vRoot := range tree.Roots {
		if vEl, ok := vRoot.(*vdom.Element); ok {
			// If vRoot is an element, test its Selector method
			expectedEl := root.ChildNodes()[i].(dom.Element)
			testSelector(vEl, root, expectedEl)
		}
	}
}

func createAndApplyPatcher(root dom.Element, html string, createPatch func(tree *vdom.Tree) vdom.Patcher) {
	// Parse some source html into a tree
	tree := setUpDOM(html, root)
	// Create the patch using the provided function
	patch := createPatch(tree)
	// Apply the patch using the provided root
	err := patch.Patch(root)
	jasmine.Expect(err).ToBe(nil)
}

func newAppendRootPatcher(newHTML string) func(tree *vdom.Tree) vdom.Patcher {
	return func(tree *vdom.Tree) vdom.Patcher {
		// Create a new tree with the given html
		newTree, err := vdom.Parse([]byte(newHTML))
		jasmine.Expect(err).ToBe(nil)
		// Return a new patch to append to the root
		return &vdom.Append{
			Child: newTree.Roots[0],
		}
	}
}

func testAppendRootPatcher(root dom.Element, newHTML string) {
	createAndApplyPatcher(root, "", newAppendRootPatcher(newHTML))
	jasmine.Expect(root.InnerHTML()).ToBe(newHTML)
}

func newReplaceRootPatcher(newHTML string) func(tree *vdom.Tree) vdom.Patcher {
	return func(tree *vdom.Tree) vdom.Patcher {
		// Create a new tree with the given html
		newTree, err := vdom.Parse([]byte(newHTML))
		jasmine.Expect(err).ToBe(nil)
		// Return a new patch to replace the root of the old tree with
		// the root of the new tree
		return &vdom.Replace{
			Old: tree.Roots[0],
			New: newTree.Roots[0],
		}
	}
}

func testReplaceRootPatcher(root dom.Element, oldHTML string, newHTML string) {
	createAndApplyPatcher(root, oldHTML, newReplaceRootPatcher(newHTML))
	// Test that the patch was applied
	children := root.ChildNodes()
	jasmine.Expect(len(children)).ToBe(1)
	jasmine.Expect(root.InnerHTML()).ToBe(newHTML)
}

func newRemoveRootPatcher() func(tree *vdom.Tree) vdom.Patcher {
	return func(tree *vdom.Tree) vdom.Patcher {
		// Return a new patch to remove the root from the tree
		return &vdom.Remove{
			Node: tree.Roots[0],
		}
	}
}

func testRemoveRootPatcher(root dom.Element, html string) {
	createAndApplyPatcher(root, html, newRemoveRootPatcher())
	// Test that the patch was applied by testing that the
	// root has no children
	children := root.ChildNodes()
	jasmine.Expect(len(children)).ToBe(0)
}
