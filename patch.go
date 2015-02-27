// +build js

package vdom

import (
	"github.com/gopherjs/gopherjs/js"
	"honnef.co/go/js/dom"
)

var (
	document dom.Document
)

func init() {
	// Initialize document iff we are running
	// inside some environment that has a window, which
	// means it probably has a DOM.
	if js.Global.Get("window") != js.Undefined {
		document = dom.GetWindow().Document()
	}
}

type Patch interface {
	Apply(parentSelector string)
}

type PatchSet []Patch

func (ps PatchSet) Apply(parentSelector string) {
	for _, patch := range ps {
		patch.Apply(parentSelector)
	}
}

type AppendChild struct {
	parent dom.Node
	child  dom.Node
}

func (ac *AppendChild) Apply() {
	ac.parent.AppendChild(ac.child)
}

func NewAppendChildPatch(child Node) *AppendChild {
	// var domChild dom.Node
	// switch child.(type) {
	// case *Element:
	// 	el := child.(*Element)
	// 	domEl := document.CreateElement(el.Name)
	// 	for _, attr := range el.Attrs {
	// 		domEl.SetAttribute(attr.Name, attr.Value)
	// 	}
	// }
	return &AppendChild{}
}
