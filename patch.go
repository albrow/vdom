package vdom

import (
	"fmt"
	"github.com/gopherjs/gopherjs/js"
	"honnef.co/go/js/dom"
)

var document dom.Document

func init() {
	// We only want to initialize document if we are running in the browser.
	// We can detect this by checking if the document is defined.
	if js.Global != nil && js.Global.Get("document") != js.Undefined {
		document = dom.GetWindow().Document()
	}
}

type Patcher interface {
	Patch(root dom.Element) error
}

type PatchSet []Patcher

func (ps PatchSet) Patch(root dom.Element) error {
	for _, patch := range ps {
		if err := patch.Patch(root); err != nil {
			return err
		}
	}
	return nil
}

type Replace struct {
	Old Node
	New Node
}

func (p *Replace) Patch(root dom.Element) error {
	var parent dom.Node
	if p.Old.Parent() != nil {
		parent = findInDOM(p.Old.Parent(), root)
	} else {
		parent = root
	}
	oldChild := findInDOM(p.Old, root)
	newChild := createForDOM(p.New)
	parent.ReplaceChild(newChild, oldChild)
	return nil
}

type Remove struct {
	Node Node
}

func (p *Remove) Patch(root dom.Element) error {
	var parent dom.Node
	if p.Node.Parent() != nil {
		parent = findInDOM(p.Node.Parent(), root)
	} else {
		parent = root
	}
	self := findInDOM(p.Node, root)
	parent.RemoveChild(self)
	return nil
}

// findInDOM finds the node in the actual DOM corresponding
// to the given virtual node, using the given root as a relative
// starting point.
func findInDOM(node Node, root dom.Element) dom.Node {
	el := root.ChildNodes()[node.Index()[0]]
	for _, i := range node.Index()[1:] {
		el = el.ChildNodes()[i]
	}
	return el
}

// createForDOM creates a real node corresponding to the given
// virtual node. It does not insert it into the actual DOM.
func createForDOM(node Node) dom.Node {
	switch node.(type) {
	case *Element:
		vEl := node.(*Element)
		el := document.CreateElement(vEl.Name)
		for _, attr := range vEl.Attrs {
			el.SetAttribute(attr.Name, attr.Value)
		}
		el.SetInnerHTML(string(vEl.InnerHTML()))
		return el
	case *Text:
		vText := node.(*Text)
		textNode := document.CreateTextNode(string(vText.Value))
		return textNode
	case *Comment:
		vComment := node.(*Comment)
		commentNode := document.Underlying().Call("createComment", string(vComment.Value))
		return dom.WrapNode(commentNode)
	default:
		msg := fmt.Sprintf("Don't know how to create node for type %T", node)
		panic(msg)
	}
}
