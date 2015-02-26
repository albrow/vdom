package vdom

import (
	"encoding/xml"
	"fmt"
	"io"
)

// Parse reads escaped html from io.Reader and returns a tree structure representing
// it. It returns an error if there was a problem parsing the html. The html read from
// r must only contain one root node. If it contains more than one root node, only the
// first node and its children will exist in the tree structure.
func Parse(r io.Reader) (*Tree, error) {
	// Create a xml.Decoder to read from r
	dec := xml.NewDecoder(r)
	dec.Entity = xml.HTMLEntity
	dec.Strict = false
	dec.AutoClose = xml.HTMLAutoClose

	// Iterate through each token and construct the tree
	tree := &Tree{}
	var currentParent *Element = nil
	for token, err := dec.Token(); ; token, err = dec.Token() {
		if err != nil {
			if err == io.EOF {
				// We reached the end of the document we were parsing
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

// parseToken parses a single token and adds the appropriate node(s) to the tree. When calling
// parseToken iteratively, you should always capture the nextParent return and use it as the
// currentParent argument in the next iteration.
func parseToken(tree *Tree, token xml.Token, currentParent *Element) (nextParent *Element, err error) {
	switch token.(type) {
	case xml.StartElement:
		// Parse the name and attrs directly from the xml.StartElement
		startEl := token.(xml.StartElement)
		el := &Element{
			Name: parseName(startEl.Name),
		}
		for _, attr := range startEl.Attr {
			el.Attrs = append(el.Attrs, Attr{
				Name:  parseName(attr.Name),
				Value: attr.Value,
			})
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
		// Parse the value from the xml.CharData
		text := &Text{
			Value: []byte(charData.Copy()),
		}
		if currentParent != nil {
			// Set this element's parent
			text.parent = currentParent
			// Add this element to the currentParent's children
			currentParent.children = append(currentParent.children, text)
		}
		if tree.Root == nil {
			// If this is the first element we've come accross, it is
			// the root of the tree
			tree.Root = text
		}
	case xml.Comment:
		xmlComment := token.(xml.Comment)
		// Parse the value from the xml.Comment
		comment := &Comment{
			Value: []byte(xmlComment.Copy()),
		}
		if currentParent != nil {
			// Set this element's parent
			comment.parent = currentParent
			// Add this element to the currentParent's children
			currentParent.children = append(currentParent.children, comment)
		}
		if tree.Root == nil {
			// If this is the first element we've come accross, it is
			// the root of the tree
			tree.Root = comment
		}
	case xml.ProcInst:
		xmlProcInst := token.(xml.ProcInst)
		// Parse the value from the xml.ProcInst
		proc := &ProcInst{
			Target: xmlProcInst.Target,
			Inst:   xmlProcInst.Inst,
		}
		if currentParent != nil {
			// Set this element's parent
			proc.parent = currentParent
			// Add this element to the currentParent's children
			currentParent.children = append(currentParent.children, proc)
		}
		if tree.Root == nil {
			// If this is the first element we've come accross, it is
			// the root of the tree
			tree.Root = proc
		}
	case xml.Directive:
		xmlDir := token.(xml.Directive)
		// Parse the value from the xml.Directive
		dir := &Directive{
			Value: []byte(xmlDir.Copy()),
		}
		if currentParent != nil {
			// Set this element's parent
			dir.parent = currentParent
			// Add this element to the currentParent's children
			currentParent.children = append(currentParent.children, dir)
		}
		if tree.Root == nil {
			// If this is the first element we've come accross, it is
			// the root of the tree
			tree.Root = dir
		}
	}
	return currentParent, nil
}

// parseName converts an xml.Name to a single string name. For our
// purposes we are not interested in the different namespaces, and
// just need to treat the name as a single string.
func parseName(name xml.Name) string {
	if name.Space != "" {
		return fmt.Sprintf("%s:%s", name.Space, name.Local)
	}
	return name.Local
}
