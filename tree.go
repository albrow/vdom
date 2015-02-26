package vdom

import (
	"fmt"
	"reflect"
)

// A Tree is an in-memory representation of a DOM tree
type Tree struct {
	Root Node
}

// A Node is an element inside a tree.
type Node interface {
	// Parent returns the parent node or nil if there is none
	Parent() Node
	// Children returns a slice of child nodes or nil if there
	// are none
	Children() []Node
}

// Attr is an xml/html attribute
type Attr struct {
	Name  string
	Value string
}

// Element is an xml/html element, e.g., <div></div>. Name does not include the
// <, >, or / symbols.
type Element struct {
	Name     string
	Attrs    []Attr
	parent   Node
	children []Node
}

func (e *Element) Parent() Node {
	return e.parent
}

func (e *Element) Children() []Node {
	return e.children
}

// String satisfies fmt.Stringer
func (e *Element) String() string {
	// TODO: improve this by recursively casting children
	// Currently they just show up as pointer values, which
	// are not very informative.
	return fmt.Sprintf("%+v", *e)
}

// Compare non-recursively compares e to other. It does not check
// the children or parent fields since they can be a Node with
// any underlying type. If you want to compare the parent and children
// fields, use CompareNodes.
func (e *Element) Compare(other *Element) (bool, string) {
	if e.Name != other.Name {
		return false, fmt.Sprintf("e.Name was %s but other.Name was %s", e.Name, other.Name)
	}
	attrs := e.Attrs
	otherAttrs := other.Attrs
	if len(attrs) != len(otherAttrs) {
		return false, fmt.Sprintf("n has %d attrs but other has %d attrs.", len(attrs), len(otherAttrs))
	}
	for i, attr := range attrs {
		otherAttr := otherAttrs[i]
		if attr != otherAttr {
			return false, fmt.Sprintf("e.Attrs[%d] was %s but other.Attrs[%d] was %s", i, attr, i, otherAttr)
		}
	}
	return true, ""
}

// Text is a text node inside an xml/html document, i.e. anything
// not surrounded by tags.
type Text struct {
	Value  []byte
	parent Node
}

func (t *Text) Parent() Node {
	return t.parent
}

func (t *Text) Children() []Node {
	// A text node can't have any children
	return nil
}

// String satisfies fmt.Stringer
func (t *Text) String() string {
	// TODO: improve this by recursively casting children
	// Currently they just show up as pointer values, which
	// are not very informative.
	return fmt.Sprintf("%+v", *t)
}

// Compare non-recursively compares t to other. It does not check
// the parent fields since they can be a Node with any underlying type.
// If you want to compare the parent fields, use CompareNodes.
func (t *Text) Compare(other *Text) (bool, string) {
	if string(t.Value) != string(other.Value) {
		return false, fmt.Sprintf("t.Value was %s but other.Value was %s", string(t.Value), string(other.Value))
	}
	return true, ""
}

// Comment is an xml/html comment of the form <!-- value -->.
// Value does not include the <!-- and --> markers.
type Comment struct {
	Value  []byte
	parent Node
}

func (c *Comment) Parent() Node {
	return c.parent
}

func (c *Comment) Children() []Node {
	// A commet node can't have any children
	return nil
}

// String satisfies fmt.Stringer
func (c *Comment) String() string {
	// TODO: improve this by recursively casting children
	// Currently they just show up as pointer values, which
	// are not very informative.
	return fmt.Sprintf("%+v", *c)
}

// Compare non-recursively compares c to other. It does not check
// the parent fields since they can be a Node with any underlying type.
// If you want to compare the parent fields, use CompareNodes.
func (c *Comment) Compare(other *Comment) (bool, string) {
	if string(c.Value) != string(other.Value) {
		return false, fmt.Sprintf("c.Value was %s but other.Value was %s", string(c.Value), string(other.Value))
	}
	return true, ""
}

// ProcInst is an xml/html processing instruction of the form <?target inst?>
type ProcInst struct {
	Target string
	Inst   []byte
	parent Node
}

func (p *ProcInst) Parent() Node {
	return p.parent
}

func (p *ProcInst) Children() []Node {
	// A processing instruction node can't have any children
	return nil
}

// String satisfies fmt.Stringer
func (p *ProcInst) String() string {
	// TODO: improve this by recursively casting children
	// Currently they just show up as pointer values, which
	// are not very informative.
	return fmt.Sprintf("%+v", *p)
}

// Compare non-recursively compares p to other. It does not check
// the parent fields since they can be a Node with any underlying type.
// If you want to compare the parent fields, use CompareNodes.
func (p *ProcInst) Compare(other *ProcInst) (bool, string) {
	if p.Target != other.Target {
		return false, fmt.Sprintf("p.Target was %s but other.Target was %s", p.Target, other.Target)
	}
	if string(p.Inst) != string(other.Inst) {
		return false, fmt.Sprintf("p.Inst was %s but other.Inst was %s", string(p.Inst), string(other.Inst))
	}
	return true, ""
}

// Directive is an xml/html directive of the form  <!value>. Value
// does not include the <! and > markers.
type Directive struct {
	Value  []byte
	parent Node
}

func (d *Directive) Parent() Node {
	return d.parent
}

func (d *Directive) Children() []Node {
	// A directive node can't have any children
	return nil
}

// String satisfies fmt.Stringer
func (d *Directive) String() string {
	// TODO: improve this by recursively casting children
	// Currently they just show up as pointer values, which
	// are not very informative.
	return fmt.Sprintf("%+v", *d)
}

// Compare non-recursively compares d to other. It does not check
// the parent fields since they can be a Node with any underlying type.
// If you want to compare the parent fields, use CompareNodes.
func (d *Directive) Compare(other *Directive) (bool, string) {
	if string(d.Value) != string(other.Value) {
		return false, fmt.Sprintf("d.Value was %s but other.Value was %s", string(d.Value), string(other.Value))
	}
	return true, ""
}

// Compare recursively compares t to other. It returns false and a detailed
// message if n does not equal other. Otherwise, it returns true and an empty
// string. NOTE: Comare never checks the parent properties of t's
// children. This is so you can construct a comparable tree inside a literal.
// (You can't set the parent field inside a literal).
func (t *Tree) Compare(other *Tree) (bool, string) {
	if t.Root == nil {
		if other.Root == nil {
			return true, ""
		} else {
			return false, fmt.Sprintf("t.Root was nil, but other.Root was %v", other.Root)
		}
	}
	return CompareNodes(t.Root, other.Root)
}

// CompareNodes recursively compares n to other. It returns false and a detailed
// message if n does not equal other. Otherwise, it returns true and an empty
// string. NOTE: CompareNodes never checks the parent properties of n or n's
// children. This is so you can construct a comparable tree inside a literal.
// (You can't set the parent field inside a literal).
func CompareNodes(n Node, other Node) (bool, string) {
	if reflect.TypeOf(n) != reflect.TypeOf(other) {
		return false, fmt.Sprintf("n has underlying type %T but the other node has underlying type %T", n, other)
	}
	switch n.(type) {
	case *Element:
		el := n.(*Element)
		otherEl := other.(*Element)
		if match, msg := el.Compare(otherEl); !match {
			return false, msg
		}
	case *Text:
		text := n.(*Text)
		otherText := other.(*Text)
		if match, msg := text.Compare(otherText); !match {
			return false, msg
		}
	case *Comment:
		comment := n.(*Comment)
		otherComment := other.(*Comment)
		if match, msg := comment.Compare(otherComment); !match {
			return false, msg
		}
	case *ProcInst:
		proc := n.(*ProcInst)
		otherProc := other.(*ProcInst)
		if match, msg := proc.Compare(otherProc); !match {
			return false, msg
		}
	case *Directive:
		dir := n.(*Directive)
		otherDir := other.(*Directive)
		if match, msg := dir.Compare(otherDir); !match {
			return false, msg
		}
	default:
		return false, fmt.Sprintf("Don't know how to compare n of underlying type %T", n)
	}
	children := n.Children()
	otherChildren := other.Children()
	if len(children) != len(otherChildren) {
		return false, fmt.Sprintf("n has %d children but other has %d children.", len(children), len(otherChildren))
	}
	for i, child := range children {
		otherChild := otherChildren[i]
		if match, msg := CompareNodes(child, otherChild); !match {
			return false, msg
		}
	}
	return true, ""
}
