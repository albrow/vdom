package vdom

import (
	"encoding/xml"
	"fmt"
	"io"
)

func Parse(r io.Reader) (*Tree, error) {
	// Create a xml.Decoder to read from r
	dec := xml.NewDecoder(r)
	dec.Entity = xml.HTMLEntity
	dec.Strict = false
	dec.AutoClose = xml.HTMLAutoClose

	// TODO: Iterate through each token and construct the tree
	tree := &Tree{}
	var currentParent *Element = nil
	for token, err := dec.Token(); ; token, err = dec.Token() {
		if err != nil {
			if err == io.EOF {
				// We reahed the end of the document we were parsing
				break
			} else {
				// There was some unexpected error
				return nil, err
			}
		}
		if nextParent, err := parseToken(tree, token, currentParent); err != nil {
			return nil, err
		} else {
			currentParent = nextParent
		}
	}
	return tree, nil
}

func parseToken(tree *Tree, token xml.Token, currentParent *Element) (nextParent *Element, err error) {
	switch token.(type) {
	case xml.StartElement:
		// Parse the name and attrs directly from the xml.StartElement
		startEl := token.(xml.StartElement)
		el := &Element{
			Name:  parseName(startEl.Name),
			Attrs: startEl.Attr,
		}
		if currentParent != nil {
			// Set this element's parent
			el.parent = currentParent
			// Add this element to the currentParent's children
			currentParent.children = append(currentParent.children, el)
		}
		if tree.Root == nil {
			// If this is the first element we've come accross, it is
			// the root of the tree
			tree.Root = el
		}
		// Set this element to the nextParent. The next node(s) we find
		// are children of this element until we reach xml.EndElement
		return el, nil
	case xml.EndElement:
		// Assuming the xml is well-formed, this marks the end of the current
		// parent
		endEl := token.(xml.EndElement)
		if currentParent == nil {
			// If we reach a closing tag without a corresponding start tag, the
			// xml is malformed
			return nil, fmt.Errorf("XML was malformed: Found closing tag %s before a corresponding opening tag.", parseName(endEl.Name))
		} else if currentParent.Name != parseName(endEl.Name) {
			// Make sure the name of the closing tag matches what we expect
			return nil, fmt.Errorf("XML was malformed: Found closing tag %s before the closing tag for %s", parseName(endEl.Name), currentParent.Name)
		}
		// The currentParent has been closed, so it has no more children.
		// The next node(s) we find must be children of currentParent.parent.
		var parentEl *Element
		if currentParent.parent != nil {
			var ok bool
			parentEl, ok = currentParent.parent.(*Element)
			if !ok {
				return nil, fmt.Errorf("Expected parent to be type *Element, but got type %T", currentParent.parent)
			}
		}
		currentParent = parentEl
	case xml.CharData:
		charData := token.(xml.CharData)
		// Parse the value of the text from the xml.CharData
		text := &Text{
			Value: []byte(charData.Copy()),
		}
		if currentParent != nil {
			// Set this element's parent
			text.parent = currentParent
			// Add this element to the currentParent's children
			currentParent.children = append(currentParent.children, text)
		}
	case xml.Comment:

	case xml.ProcInst:

	case xml.Directive:
	}
	return currentParent, nil
}

func parseName(name xml.Name) string {
	if name.Space != "" {
		return fmt.Sprintf("%s:%s", name.Space, name.Local)
	}
	return name.Local
}
