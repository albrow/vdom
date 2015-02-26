package vdom

import (
	"bytes"
	"encoding/xml"
	"testing"
)

func TestParse(t *testing.T) {
	buf := bytes.NewBuffer([]byte("<ul><li>one</li><li>two</li><li>three</li></ul>"))
	gotTree, err := Parse(buf)
	if err != nil {
		t.Errorf("Unexpected error in Parse: %s", err.Error())
	}
	expectedTree := &Tree{
		Root: &Element{
			Name:  "ul",
			Attrs: []xml.Attr{},
			children: []Node{
				&Element{
					Name: "li",
					children: []Node{
						&Text{
							Value: []byte("one"),
						},
					},
				},
				&Element{
					Name: "li",
					children: []Node{
						&Text{
							Value: []byte("two"),
						},
					},
				},
				&Element{
					Name: "li",
					children: []Node{
						&Text{
							Value: []byte("three"),
						},
					},
				},
			},
		},
	}
	if match, msg := expectedTree.Compare(gotTree); !match {
		t.Errorf("HTML was not parsed correctly.\n%s", msg)
	}
}
