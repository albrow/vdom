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

type Replace struct {
	Old Node
	New Node
}

func (p *Replace) Patch(root dom.Element) error {
	switch p.Old.(type) {
	case (*Element):
		oldEl := p.Old.(*Element)
		parent := root.QuerySelector(oldEl.Selector()).ParentElement()
		parent.SetInnerHTML(string(p.New.HTML()))
	default:
		return fmt.Errorf("Don't know how to apply Replace patch with Node of type %T", p.Old)
	}
	return nil
}

type Remove struct {
	Node Node
}

func (p *Remove) Patch(root dom.Element) error {
	switch p.Node.(type) {
	case (*Element):
		vEl := p.Node.(*Element)
		el := root.QuerySelector(vEl.Selector())
		el.ParentNode().RemoveChild(el)
	default:
		return fmt.Errorf("Don't know how to apply Remove patch with Node of type %T", p.Node)
	}
	return nil
}
