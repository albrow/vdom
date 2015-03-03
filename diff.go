package vdom

import (
	"fmt"
	"reflect"
)

func Diff(t, other *Tree) (PatchSet, error) {
	patches := []Patcher{}
	if err := recursiveDiff(&patches, t.Children, other.Children); err != nil {
		return nil, err
	}
	return patches, nil
}

func recursiveDiff(patches *[]Patcher, nodes, otherNodes []Node) error {

	numOtherNodes := len(otherNodes)
	numNodes := len(nodes)
	minNumNodes := numOtherNodes
	if numOtherNodes > numNodes {
		// other has more first-level children than t.
		// We should append the additional children.
		for _, otherNode := range otherNodes[numNodes:] {
			*patches = append(*patches, &Append{
				Parent: otherNode.Parent(),
				Child:  otherNode,
			})
		}
		minNumNodes = numNodes
	} else if numNodes > numOtherNodes {
		// t has more first-level children than other.
		// We should remove the additional children.
		for _, node := range nodes[numOtherNodes:] {
			*patches = append(*patches, &Remove{
				Node: node,
			})
		}
		minNumNodes = numOtherNodes
	}
	for i := 0; i < minNumNodes; i++ {
		otherNode := otherNodes[i]
		node := nodes[i]
		if reflect.TypeOf(node) != reflect.TypeOf(otherNode) {
			// The types don't match. We should replace node
			// with other node
			*patches = append(*patches, &Replace{
				Old: node,
				New: otherNode,
			})
		}
		// If we've reached here, the types do match. We should compare
		// based on the type
		switch otherNode.(type) {
		case *Element:
			otherEl := otherNode.(*Element)
			el := node.(*Element)
			if otherEl.Name != el.Name {
				// The elements have different tag names. We should replace
				// el with otherEl
				*patches = append(*patches, &Replace{
					Old: el,
					New: otherEl,
				})
				continue
			}
			// If we've reached here, the elements have the same tag name
			// Next, we should compare the attributes...
			// Any attributes in el that are not in otherEl should be removed
			otherAttrs := otherEl.AttrMap()
			attrs := el.AttrMap()
			for attrName := range attrs {
				if _, found := otherAttrs[attrName]; !found {
					*patches = append(*patches, &RemoveAttr{
						Node:     el,
						AttrName: attrName,
					})
				}
			}
			// Now iterate through the attributes in otherEl
			for name, otherValue := range otherAttrs {
				value, found := attrs[name]
				if !found {
					// The attribute exists in otherEl but not in el,
					// we should add it.
					*patches = append(*patches, &SetAttr{
						Node: el,
						Attr: &Attr{
							Name:  name,
							Value: otherValue,
						},
					})
				} else if value != otherValue {
					// The attribute exists in el but has a differnt value
					// than it does in otherEl. We should set it to the value
					// in otherEl.
					*patches = append(*patches, &SetAttr{
						Node: el,
						Attr: &Attr{
							Name:  name,
							Value: otherValue,
						},
					})
				}
			}
			// Recursively apply diff algorithm to each element's children
			recursiveDiff(patches, el.Children(), otherEl.Children())
		case *Text:
			otherText := otherNode.(*Text)
			text := node.(*Text)
			if string(otherText.Value) != string(text.Value) {
				// The text nodes don't match. We should replace
				// text with otherText
				*patches = append(*patches, &Replace{
					Old: text,
					New: otherText,
				})
			}
		case *Comment:
			otherComment := otherNode.(*Comment)
			comment := node.(*Comment)
			if string(otherComment.Value) != string(comment.Value) {
				// The comment nodes don't match. We should replace
				// comment with otherComment
				*patches = append(*patches, &Replace{
					Old: comment,
					New: otherComment,
				})
			}
		default:
			return fmt.Errorf("Don't know how to compare node of type %T", otherNode)
		}
	}
	return nil
}