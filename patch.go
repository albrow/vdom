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
	VNode Node
	Inner []byte
}

func (p *SetInnerHTML) Patch(root dom.Element) error {
	fmt.Printf("%+v\n", root)
	switch p.VNode.(type) {
	case (*Element):
		vEl := p.VNode.(*Element)
		el := root.QuerySelector(vEl.Selector())
		el.SetInnerHTML(string(p.Inner))
	default:
		return fmt.Errorf("Don't know how to apply SetInnerHTML patch with VNode of type %T", p.VNode)
	}
	return nil
}
