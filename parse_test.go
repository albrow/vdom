package vdom

import (
	"bytes"
	"io"
	"testing"
)

func TestParse(t *testing.T) {
	// We'll use table-driven testing here.
	// Each test case consists of a human-readable name,
	// a reader which holds html data, and the tree structure
	// that we expect after calling Parse.
	testCases := []struct {
		name         string
		reader       io.Reader
		expectedTree *Tree
	}{
		{
			name:   "Element root",
			reader: bytes.NewBuffer([]byte("<div></div>")),
			expectedTree: &Tree{
				Root: &Element{
					Name: "div",
				},
			},
		},
		{
			name:   "Text root",
			reader: bytes.NewBuffer([]byte("Hello")),
			expectedTree: &Tree{
				Root: &Text{
					Value: []byte("Hello"),
				},
			},
		},
		{
			name:   "Comment root",
			reader: bytes.NewBuffer([]byte("<!--comment-->")),
			expectedTree: &Tree{
				Root: &Comment{
					Value: []byte("comment"),
				},
			},
		},
		{
			name:   "ProcInst root",
			reader: bytes.NewBuffer([]byte("<?target inst?>")),
			expectedTree: &Tree{
				Root: &ProcInst{
					Target: "target",
					Inst:   []byte("inst"),
				},
			},
		},
		{
			name:   "Directive root",
			reader: bytes.NewBuffer([]byte("<!doctype html>")),
			expectedTree: &Tree{
				Root: &Directive{
					Value: []byte("doctype html"),
				},
			},
		},
		{
			name:   "ul with nested li's",
			reader: bytes.NewBuffer([]byte("<ul><li>one</li><li>two</li><li>three</li></ul>")),
			expectedTree: &Tree{
				Root: &Element{
					Name: "ul",
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
		{
			name:   "Element with attrs",
			reader: bytes.NewBuffer([]byte(`<div class="container" id="main" data-custom-attr="foo"></div>`)),
			expectedTree: &Tree{
				Root: &Element{
					Name: "div",
					Attrs: []Attr{
						{Name: "class", Value: "container"},
						{Name: "id", Value: "main"},
						{Name: "data-custom-attr", Value: "foo"},
					},
				},
			},
		},
		{
			name:   "Script tag with escaped characters",
			reader: bytes.NewBuffer([]byte(`<script type="text/javascript">function((){console.log("&lt;Hello brackets&gt;")})()</script>`)),
			expectedTree: &Tree{
				Root: &Element{
					Name: "script",
					Attrs: []Attr{
						{Name: "type", Value: "text/javascript"},
					},
					children: []Node{
						&Text{
							Value: []byte(`function((){console.log("<Hello brackets>")})()`),
						},
					},
				},
			},
		},
	}
	// Iterate through each test case
	for i, tc := range testCases {
		// Parse the input from tc.reader
		gotTree, err := Parse(tc.reader)
		if err != nil {
			t.Errorf("Unexpected error in Parse: %s", err.Error())
		}
		// Check that the resulting tree matches what we expect
		if match, msg := tc.expectedTree.Compare(gotTree); !match {
			t.Errorf("Error in test case %d (%s): HTML was not parsed correctly.\n%s", i, tc.name, msg)
		}
	}
}
