package vdom

import (
	"fmt"
	"reflect"
)

func Diff(t, other *Tree) (PatchSet, error) {
	patches := []Patcher{}
	numOtherChildren := len(other.Children)
	numChildren := len(t.Children)
	minNumChildren := numOtherChildren
	if numOtherChildren > numChildren {
		// other has more first-level children than t.
		// We should append the additional children.
		for _, otherChild := range other.Children[numChildren:] {
			patches = append(patches, &Append{
				Child: otherChild,
			})
		}
		minNumChildren = numChildren
	} else if numChildren > numOtherChildren {
		// t has more first-level children than other.
		// We should remove the additional children.
		for _, child := range t.Children[numOtherChildren:] {
			patches = append(patches, &Remove{
				Node: child,
			})
		}
		minNumChildren = numOtherChildren
	}
	for i := 0; i < minNumChildren; i++ {
		otherChild := other.Children[i]
		child := t.Children[i]
		if reflect.TypeOf(child) != reflect.TypeOf(otherChild) {
			// The types don't match. We should replace child
			// with other child
			patches = append(patches, &Replace{
				Old: child,
				New: otherChild,
			})
		}
		// If we've reached here, the types do match. We should compare
		// based on the type
		switch otherChild.(type) {
		case *Element:
			otherEl := otherChild.(*Element)
			el := child.(*Element)
			if otherEl.Name != el.Name {
				// The elements have different tag names. We should replace
				// el with otherEl
				patches = append(patches, &Replace{
					Old: el,
					New: otherEl,
				})
				continue
			}
			// If we've reached here, the elements have the same tag name
			// Next, we should compare the attributes
			// TODO: use the attr names as keys to make this comparison
			// more effecient.
			numOtherAttrs := len(otherEl.Attrs)
			numAttrs := len(el.Attrs)
			minNumAttrs := numOtherAttrs
			if numOtherAttrs > numAttrs {
				// otherEl has more attributes than el
				// We should add the additional attributes.
				for _, otherAttr := range otherEl.Attrs[numAttrs:] {
					patches = append(patches, &SetAttr{
						Node: el,
						Attr: &otherAttr,
					})
				}
				minNumAttrs = numAttrs
			} else if numAttrs > numOtherAttrs {
				// el has more attributes than otherEl
				// We should remove the additional attributes.
				for _, attr := range el.Attrs[numOtherAttrs:] {
					patches = append(patches, &RemoveAttr{
						Node:     el,
						AttrName: attr.Name,
					})
				}
				minNumAttrs = numOtherAttrs
			}
			for i := 0; i < minNumAttrs; i++ {
				// Compare each individual shared attribute
				otherAttr := otherEl.Attrs[i]
				attr := el.Attrs[i]
				if otherAttr.Name != attr.Name || otherAttr.Value != attr.Value {
					// The attributes don't match. We should replace attr
					// with other attr
					patches = append(patches, &RemoveAttr{
						Node:     el,
						AttrName: attr.Name,
					})
					patches = append(patches, &SetAttr{
						Node: el,
						Attr: &otherAttr,
					})
				}
			}
		case *Text:
			otherText := otherChild.(*Text)
			text := child.(*Text)
			if string(otherText.Value) != string(text.Value) {
				// The text nodes don't match. We should replace
				// text with otherText
				patches = append(patches, &Replace{
					Old: text,
					New: otherText,
				})
			}
		case *Comment:
			otherComment := otherChild.(*Comment)
			comment := child.(*Comment)
			if string(otherComment.Value) != string(comment.Value) {
				// The comment nodes don't match. We should replace
				// comment with otherComment
				patches = append(patches, &Replace{
					Old: comment,
					New: otherComment,
				})
			}
		default:
			return nil, fmt.Errorf("Don't know how to compare node of type %T", otherChild)
		}
	}
	return patches, nil
}
