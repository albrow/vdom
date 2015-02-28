// +build js

package vdom

import (
	"fmt"
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

type Patcher interface {
	Patch(root dom.Element)
}

type PatchSet []Patcher

func (ps PatchSet) Patch(root dom.Element) {
	for _, patch := range ps {
		patch.Patch(root)
	}
}

type SetInnerHTML struct {
	vnode Node
	inner []byte
}

func (p *SetInnerHTML) Patch(root dom.Element) error {
	switch p.vnode.(type) {
	case (*Element):
		vEl := p.vnode.(*Element)
		el := root.QuerySelector(vEl.PartialSelector())
		el.SetInnerHTML(string(p.inner))
	default:
		return fmt.Errorf("Don't know how to apply SetInnerHTML patch with vnode of type %T", p.vnode)
	}
	return nil
}
