package main

import (
	"github.com/JohannWeging/jasmine"
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

	jasmine.Describe("Replace patch", func() {

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
			testPatch(`<div id="old"></div>`, sandbox, newReplaceRootPatcher(`<div id="new"></div>`), func() {
				// Test that the patch was applied
				children := sandbox.ChildNodes()
				jasmine.Expect(len(children)).ToBe(1)
				div := children[0].(dom.Element)
				jasmine.Expect(div.ID()).ToEqual("new")
			})
		})

		jasmine.It("works with a text root", func() {
			testPatch("Old", sandbox, newReplaceRootPatcher("New"), func() {
				// Test that the patch was applied
				children := sandbox.ChildNodes()
				jasmine.Expect(len(children)).ToBe(1)
				textNode := children[0].(*dom.Text)
				jasmine.Expect(textNode.NodeValue()).ToBe("New")
			})
		})

		jasmine.It("works with a comment root", func() {
			testPatch("<!--old-->", sandbox, newReplaceRootPatcher("<!--new-->"), func() {
				// Test that the patch was applied
				children := sandbox.ChildNodes()
				jasmine.Expect(len(children)).ToBe(1)
				commentNode := sandbox.ChildNodes()[0].(*dom.BasicHTMLElement)
				jasmine.Expect(commentNode.NodeValue()).ToEqual("new")
			})
		})

		jasmine.It("works with nested siblings", func() {
			testPatch("<ul><li>one</li><li>two</li><li>three</li></ul>", sandbox, func(tree *vdom.Tree) vdom.Patcher {
				// Create a new tree, which only consists of one of the lis
				// We want to change it from one to uno
				newTree, err := vdom.Parse([]byte("<li>uno</li>"))
				jasmine.Expect(err).ToBe(nil)
				// Create a patch manually
				return &vdom.Replace{
					Old: tree.Roots[0].Children()[0],
					New: newTree.Roots[0],
				}
			}, func() {
				// Test that the patch was applied
				ul := sandbox.ChildNodes()[0].(*dom.HTMLUListElement)
				jasmine.Expect(ul.InnerHTML()).ToBe("<li>uno</li><li>two</li><li>three</li>")
			})
		})
	})

	jasmine.Describe("Remove patch", func() {

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
			testPatch("<div></div>", sandbox, newRemoveRootPatcher(), testRemoveRootPatch(sandbox))
		})

		jasmine.It("works with a text root", func() {
			testPatch("Text", sandbox, newRemoveRootPatcher(), testRemoveRootPatch(sandbox))
		})

		jasmine.It("works with a comment root", func() {
			testPatch("<!--comment-->", sandbox, newRemoveRootPatcher(), testRemoveRootPatch(sandbox))
		})

		jasmine.It("works with nested siblings", func() {
			testPatch("<ul><li>one</li><li>two</li><li>three</li></ul>", sandbox, func(tree *vdom.Tree) vdom.Patcher {
				return &vdom.Remove{
					Node: tree.Roots[0].Children()[1],
				}
			}, func() {
				// Test that the patch was applied by checking the innerHTML
				// property of the ul node.
				ul := sandbox.ChildNodes()[0].(*dom.HTMLUListElement)
				jasmine.Expect(ul.InnerHTML()).ToBe("<li>one</li><li>three</li>")
			})
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

func testSelectors(tree *vdom.Tree, root dom.Element) {
	for i, vRoot := range tree.Roots {
		if vEl, ok := vRoot.(*vdom.Element); ok {
			// If vRoot is an element, test its Selector method
			expectedEl := root.ChildNodes()[i].(dom.Element)
			testSelector(vEl, root, expectedEl)
		}
	}
}

func testPatch(html string, root dom.Element, createPatch func(tree *vdom.Tree) vdom.Patcher, testPatch func()) {
	// Parse some source html into a tree
	tree := setUpDOM(html, root)
	// Create the patch using the provided function
	patch := createPatch(tree)
	// Apply the patch using the provided root
	err := patch.Patch(root)
	jasmine.Expect(err).ToBe(nil)
	// Test the patch using the provided function
	testPatch()
}

func newReplaceRootPatcher(newHtml string) func(tree *vdom.Tree) vdom.Patcher {
	return func(tree *vdom.Tree) vdom.Patcher {
		// Create a new tree with the given html
		newTree, err := vdom.Parse([]byte(newHtml))
		jasmine.Expect(err).ToBe(nil)
		// Return a new patch to the root of the old tree with
		// the root of the new tree
		return &vdom.Replace{
			Old: tree.Roots[0],
			New: newTree.Roots[0],
		}
	}
}

func newRemoveRootPatcher() func(tree *vdom.Tree) vdom.Patcher {
	return func(tree *vdom.Tree) vdom.Patcher {
		// Return a new patch to the root of the old tree with
		// the root of the new tree
		return &vdom.Remove{
			Node: tree.Roots[0],
		}
	}
}

func testRemoveRootPatch(root dom.Element) func() {
	return func() {
		// Test that the patch was applied by testing that the
		// root has no children
		children := root.ChildNodes()
		jasmine.Expect(len(children)).ToBe(0)
	}
}
