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

// BUG: doesn't work if old has sibling elements
func (p *Replace) Patch(root dom.Element) error {
	var parent dom.Element
	if p.Old.Parent() != nil {
		parent = root.QuerySelector(p.Old.Parent().Selector())
	} else {
		parent = root
	}
	// TODO: To fix the bug, we need to find both nodes in the DOM
	// and use ReplaceChild instead of SetInnerHTML
	parent.SetInnerHTML(string(p.New.HTML()))
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
