package vdom

import (
	"fmt"
	"honnef.co/go/js/dom"
)

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

type SetInnerHTML struct {
	Node  Node
	Inner []byte
}

func (p *SetInnerHTML) Patch(root dom.Element) error {
	switch p.Node.(type) {
	case (*Element):
		vEl := p.Node.(*Element)
		el := root.QuerySelector(vEl.Selector())
		el.SetInnerHTML(string(p.Inner))
	default:
		return fmt.Errorf("Don't know how to apply SetInnerHTML patch with Node of type %T", p.Node)
	}
	return nil
}
