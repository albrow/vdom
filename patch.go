package vdom

import (
	"fmt"
	"honnef.co/go/js/dom"
)

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
		el := root.QuerySelector(vEl.Selector())
		el.SetInnerHTML(string(p.inner))
	default:
		return fmt.Errorf("Don't know how to apply SetInnerHTML patch with vnode of type %T", p.vnode)
	}
	return nil
}
