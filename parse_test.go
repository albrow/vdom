package vdom

import (
	"bytes"
	"encoding/xml"
	"io"
	"testing"
)

func TestParse(t *testing.T) {
	testCases := []struct {
		reader       io.Reader
		expectedTree *Tree
	}{
		{
			// A single text node
			reader: bytes.NewBuffer([]byte("Hello")),
			expectedTree: &Tree{
				Root: &Text{
					Value: []byte("Hello"),
				},
			},
		},
		{
			// HTML with nested elements
			reader: bytes.NewBuffer([]byte("<ul><li>one</li><li>two</li><li>three</li></ul>")),
			expectedTree: &Tree{
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
			},
		},
	}
	for _, tc := range testCases {
		// Parse the input from tc.reader
		gotTree, err := Parse(tc.reader)
		if err != nil {
			t.Errorf("Unexpected error in Parse: %s", err.Error())
		}
		// Check that the resulting tree matches what we expect
		if match, msg := tc.expectedTree.Compare(gotTree); !match {
			t.Errorf("HTML was not parsed correctly.\n%s", msg)
		}
	}
}
