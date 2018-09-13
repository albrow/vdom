package main

import (
	"github.com/JohannWeging/jasmine"
	"github.com/gowasm/vdom"
	dom "github.com/yml/go-js-dom"
)

var (
	document = dom.GetWindow().Document()
)

func main() {

	// The body element will be used throughout all tests, often as the
	// root or starting point of the virtual tree in the actual DOM.
	var body dom.Element

	// Before each test, instantiate the body variable if it is not
	// already instantiated.
	jasmine.BeforeEach(func() {
		if body == nil {
			body = document.QuerySelector("body")
		}
	})

	// After each test, remove everything inside the body element in order
	// to prepare for the next test.
	jasmine.AfterEach(func() {
		body.SetInnerHTML("")
	})

	// This test is just checks that the code cross-compiled correctly and can
	// be executed by the karma test runner.
	jasmine.Describe("Tests", func() {
		jasmine.It("can be loaded", func() {
			jasmine.Expect(true).ToBe(true)
		})
	})

	// Test the Element.Selector method in the actual DOM with various different
	// html structures.
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

	// Test the Append Patcher in the actual DOM with various different html
	// structures.
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
					Child:  newTree.Children[0],
					Parent: tree.Children[0].(*vdom.Element),
				}
			})
			// Test that the patch was applied
			ul := body.ChildNodes()[0].(*dom.HTMLUListElement)
			jasmine.Expect(ul.InnerHTML()).ToBe("<li>one</li><li>two</li><li>three</li>")
		})
	})

	// Test the Replace Patcher in the actual DOM with various different html
	// structures.
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
					Old: tree.Children[0].Children()[0],
					New: newTree.Children[0],
				}
			})
			// Test that the patch was applied
			ul := body.ChildNodes()[0].(*dom.HTMLUListElement)
			jasmine.Expect(ul.InnerHTML()).ToBe("<li>uno</li><li>two</li><li>three</li>")
		})
	})

	// Test the Remove Patcher in the actual DOM with various different html
	// structures.
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
					Node: tree.Children[0].Children()[1],
				}
			})
			// Test that the patch was applied by checking the innerHTML
			// property of the ul node.
			ul := body.ChildNodes()[0].(*dom.HTMLUListElement)
			jasmine.Expect(ul.InnerHTML()).ToBe("<li>one</li><li>three</li>")
		})
	})

	// Test the SetAttr Patcher in the actual DOM with various different html
	// structures.
	jasmine.Describe("SetAttr", func() {

		jasmine.It("works on a root element", func() {
			createAndApplyPatcher(body, "<div></div>", func(tree *vdom.Tree) vdom.Patcher {
				return &vdom.SetAttr{
					Node: tree.Children[0],
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
					Node: tree.Children[0].Children()[1],
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

	// Test the RemoveAttr Patcher in the actual DOM with various different html
	// structures.
	jasmine.Describe("RemoveAttr", func() {

		jasmine.It("works on a root element", func() {
			createAndApplyPatcher(body, `<div id="foo"></div>`, func(tree *vdom.Tree) vdom.Patcher {
				return &vdom.RemoveAttr{
					Node:     tree.Children[0],
					AttrName: "id",
				}
			})
			// Test that the patch was applied
			jasmine.Expect(body.InnerHTML()).ToBe("<div></div>")
		})

		jasmine.It("works on a nested element", func() {
			createAndApplyPatcher(body, `<ul><li>one</li><li data-value="two">two</li><li>three</li></ul>`, func(tree *vdom.Tree) vdom.Patcher {
				return &vdom.RemoveAttr{
					Node:     tree.Children[0].Children()[1],
					AttrName: "data-value",
				}
			})
			// Test that the patch was applied
			ul := body.ChildNodes()[0].(*dom.HTMLUListElement)
			jasmine.Expect(ul.InnerHTML()).ToBe("<li>one</li><li>two</li><li>three</li>")
		})

	})

	// Test the Diff function in the actual DOM with various different html
	// structures.
	jasmine.Describe("Diff", func() {

		jasmine.It("creates a root element", func() {
			testDiff(body, "", "<div></div>")
		})

		jasmine.It("removes a root element", func() {
			testDiff(body, "<div></div>", "")
		})

		jasmine.It("replaces a root element", func() {
			testDiff(body, "<div></div>", "<span></span>")
		})

		jasmine.It("creates a root text node", func() {
			testDiff(body, "", "Text")
		})

		jasmine.It("removes a root text node", func() {
			testDiff(body, "Text", "")
		})

		jasmine.It("replaces a root text node", func() {
			testDiff(body, "OldText", "NewText")
		})

		jasmine.It("creates a root comment node", func() {
			testDiff(body, "", "<!--comment-->")
		})

		jasmine.It("removes a root comment node", func() {
			testDiff(body, "<!--comment-->", "")
		})

		jasmine.It("replaces a root comment node", func() {
			testDiff(body, "<!--old-->", "<!--new-->")
		})

		jasmine.It("adds a root element attribute", func() {
			testDiff(body, "<div></div>", `<div id="foo"></div>`)
		})

		jasmine.It("removes a root element attribute", func() {
			testDiff(body, `<div id="foo"></div>`, "<div></div>")
		})

		jasmine.It("replaces a root element attribute", func() {
			testDiff(body, `<div id="old"></div>`, `<div id="new"></div>`)
		})

		jasmine.It("creates a nested element", func() {
			testDiff(body, "<div></div>", "<div><div></div></div>")
		})

		jasmine.It("removes a nested element", func() {
			testDiff(body, "<div><div></div></div>", "<div></div>")
		})

		jasmine.It("replaces a nested element", func() {
			testDiff(body, "<div><div></div></div>", "<div><span></span></div>")
		})

		jasmine.It("creates a nested text node", func() {
			testDiff(body, "<div></div>", "<div>Text</div>")
		})

		jasmine.It("removes a nested text node", func() {
			testDiff(body, "<div>Text</div>", "<div></div>")
		})

		jasmine.It("replaces a nested text node", func() {
			testDiff(body, "<div>OldText</div>", "<div>NewText</div>")
		})

		jasmine.It("creates a nested comment node", func() {
			testDiff(body, "<div></div>", "<div><!--comment--></div>")
		})

		jasmine.It("removes a nested comment node", func() {
			testDiff(body, "<div><!--comment--></div>", "<div></div>")
		})

		jasmine.It("replaces a nested comment node", func() {
			testDiff(body, "<div><!--old--></div>", "<div><!--new--></div>")
		})

		jasmine.It("adds a nested element attribute", func() {
			testDiff(body, "<div><div></div></div>", `<div><div id="foo"></div></div>`)
		})

		jasmine.It("removes a nested element attribute", func() {
			testDiff(body, `<div><div id="foo"></div></div>`, "<div><div></div></div>")
		})

		jasmine.It("replaces a nested element attribute", func() {
			testDiff(body, `<div><div id="old"></div></div>`, `<div><div id="new"></div></div>`)
		})

		jasmine.It("creates a nested element with siblings", func() {
			testDiff(body, "<ul><li>one</li><li>three</li></ul>", "<ul><li>one</li><li>two</li><li>three</li></ul>")
		})

		jasmine.It("removes a nested element with siblings", func() {
			testDiff(body, "<ul><li>one</li><li>two</li><li>three</li></ul>", "<ul><li>one</li><li>three</li></ul>")
		})

		jasmine.It("replaces a nested element siblings", func() {
			testDiff(body, "<ul><li>one</li><li>two</li><li>three</li></ul>", "<ul><li>one</li><li>dos</li><li>three</li></ul>")
		})

		jasmine.It("adds/replaces multiple attributes", func() {
			// Since the order of attributes can change, we'll have to do this test
			// manually
			oldHTML := `<div class="foo" id="bar" data-target="self" name="biz"></div>`
			newHTML := `<div class="bar" id="foo" name="biz" onClick="doStuff()"></div>`
			// Parse some source oldHTML into a tree and add it
			// to the actual DOM
			tree := setUpDOM(oldHTML, body)
			// Create a virtual tree with the newHTML
			newTree, err := vdom.Parse([]byte(newHTML))
			jasmine.Expect(err).ToBe(nil)
			// Use the diff function to calculate the difference between
			// the trees and return a patch set
			patches, err := vdom.Diff(tree, newTree)
			jasmine.Expect(err).ToBe(nil)
			// Apply the patches to the body in the actual DOM
			err = patches.Patch(body)
			jasmine.Expect(err).ToBe(nil)
			// Check that the body now has innerHTML equal to newHTML,
			// which would indecate the diff and patch set worked as
			// expected
			jasmine.Expect(len(body.ChildNodes())).ToBe(1)
			div := body.ChildNodes()[0].(dom.Element)
			expectedAttributes := map[string]string{
				"class":   "bar",
				"id":      "foo",
				"name":    "biz",
				"onClick": "doStuff()",
			}
			jasmine.Expect(div.Underlying().Get("attributes").Get("length")).ToBe(len(expectedAttributes))
			for name, value := range expectedAttributes {
				jasmine.Expect(div.HasAttribute(name)).ToBe(true)
				jasmine.Expect(div.GetAttribute(name)).ToBe(value)
			}
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

// expectExistsInDOM invokes jasmine and the dom bindings to check that
// el exists in the DOM. If it does not, jasmine will report an error.
func expectExistsInDOM(root dom.Element, el dom.Element) {
	jasmine.Expect(root.Contains(el)).ToBe(true)
}

// testSelector tests the Selector method for vEl and then recursively
// iterates through its children and tests the Selector method for them
// as well.
func testSelector(vEl *vdom.Element, root, expectedEl dom.Element) {
	gotEl := root.QuerySelector(vEl.Selector())
	expectExistsInDOM(root, gotEl)
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

// testSelectors recursively iterates through the virtual tree and the
// corresponding nodes in the actual DOM and tests the Selector method
// for every element.
func testSelectors(tree *vdom.Tree, root dom.Element) {
	for i, vRoot := range tree.Children {
		if vEl, ok := vRoot.(*vdom.Element); ok {
			// If vRoot is an element, test its Selector method
			expectedEl := root.ChildNodes()[i].(dom.Element)
			testSelector(vEl, root, expectedEl)
		}
	}
}

// createAndApplyPatcher adds the given html to the actual DOM starting at
// the given root. Then it invokes the createPatch function to create a Patcher
// that will act on the new nodes that were created.
func createAndApplyPatcher(root dom.Element, html string, createPatch func(tree *vdom.Tree) vdom.Patcher) {
	// Parse some source html into a tree
	tree := setUpDOM(html, root)
	// Create the patch using the provided function
	patch := createPatch(tree)
	// Apply the patch using the provided root
	err := patch.Patch(root)
	jasmine.Expect(err).ToBe(nil)
}

// newAppendRootPatcher returns an Append Patcher which will simply append a
// new node with the given newHTML to the DOM at the root of the tree.
func newAppendRootPatcher(newHTML string) func(tree *vdom.Tree) vdom.Patcher {
	return func(tree *vdom.Tree) vdom.Patcher {
		// Create a new tree with the given html
		newTree, err := vdom.Parse([]byte(newHTML))
		jasmine.Expect(err).ToBe(nil)
		// Return a new patch to append to the root
		return &vdom.Append{
			Child: newTree.Children[0],
		}
	}
}

// testAppendRootPatcher will set up the DOM with an empty root, then create
// and apply a patch that should append a new element directly to the root, and
// finally it tests that the patch was applied correctly.
func testAppendRootPatcher(root dom.Element, newHTML string) {
	createAndApplyPatcher(root, "", newAppendRootPatcher(newHTML))
	jasmine.Expect(root.InnerHTML()).ToBe(newHTML)
}

// newReplaceRootPatcher returns an Replace Patcher which will simply replace
// the contents of the root with a new element created with newHTML.
func newReplaceRootPatcher(newHTML string) func(tree *vdom.Tree) vdom.Patcher {
	return func(tree *vdom.Tree) vdom.Patcher {
		// Create a new tree with the given html
		newTree, err := vdom.Parse([]byte(newHTML))
		jasmine.Expect(err).ToBe(nil)
		// Return a new patch to replace the root of the old tree with
		// the root of the new tree
		return &vdom.Replace{
			Old: tree.Children[0],
			New: newTree.Children[0],
		}
	}
}

// testReplaceRootPatcher will set up the DOM with a root element created from the
// given oldHTML, then it will create and apply a patch that should replace the content
// of the root with a new element created from newHTML, and finally it tests that the
// patch was applied correctly.
func testReplaceRootPatcher(root dom.Element, oldHTML string, newHTML string) {
	createAndApplyPatcher(root, oldHTML, newReplaceRootPatcher(newHTML))
	// Test that the patch was applied
	children := root.ChildNodes()
	jasmine.Expect(len(children)).ToBe(1)
	jasmine.Expect(root.InnerHTML()).ToBe(newHTML)
}

// newRemoveRootPatcher returns an Remove Patcher which will simply remove
// the first child of the root.
func newRemoveRootPatcher() func(tree *vdom.Tree) vdom.Patcher {
	return func(tree *vdom.Tree) vdom.Patcher {
		// Return a new patch to remove the root from the tree
		return &vdom.Remove{
			Node: tree.Children[0],
		}
	}
}

// testReplaceRootPatcher will set up the DOM with a root element created from the
// given html, then it will create and apply a patch that should remove the first child
// of the root, and finally it tests that the patch was applied correctly.
func testRemoveRootPatcher(root dom.Element, html string) {
	createAndApplyPatcher(root, html, newRemoveRootPatcher())
	// Test that the patch was applied by testing that the
	// root has no children
	children := root.ChildNodes()
	jasmine.Expect(len(children)).ToBe(0)
}

func testDiff(root dom.Element, oldHTML string, newHTML string) {
	// Parse some source oldHTML into a tree and add it
	// to the actual DOM
	tree := setUpDOM(oldHTML, root)
	// Create a virtual tree with the newHTML
	newTree, err := vdom.Parse([]byte(newHTML))
	jasmine.Expect(err).ToBe(nil)
	// Use the diff function to calculate the difference between
	// the trees and return a patch set
	patches, err := vdom.Diff(tree, newTree)
	jasmine.Expect(err).ToBe(nil)
	// Apply the patches to the root in the actual DOM
	err = patches.Patch(root)
	jasmine.Expect(err).ToBe(nil)
	// Check that the root now has innerHTML equal to newHTML,
	// which would indecate the diff and patch set worked as
	// expected
	jasmine.Expect(root.InnerHTML()).ToBe(newHTML)
}
