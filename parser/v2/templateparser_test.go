package parser

import (
	"testing"

	"github.com/a-h/lexical/input"
	"github.com/google/go-cmp/cmp"
)

func TestTemplateParser(t *testing.T) {
	var tests = []struct {
		name        string
		input       string
		expected    HTMLTemplate
		expectError bool
	}{
		{
			name: "template: no parameters",
			input: `templ Name() {
}`,
			expected: HTMLTemplate{
				Expression: Expression{
					Value: "Name()",
					Range: Range{
						From: Position{
							Index: 6,
							Line:  0,
							Col:   6,
						},
						To: Position{
							Index: 12,
							Line:  0,
							Col:   12,
						},
					},
				},
				Children: []Node{},
			},
		},
		{
			name: "template: with receiver",
			input: `templ (data Data) Name() {
}`,
			expected: HTMLTemplate{
				Expression: Expression{
					Value: "(data Data) Name()",
					Range: Range{
						From: Position{
							Index: 6,
							Line:  0,
							Col:   6,
						},
						To: Position{
							Index: 24,
							Line:  0,
							Col:   24,
						},
					},
				},
				Children: []Node{},
			},
		},
		{
			name: "template: no spaces",
			input: `templ Name(){
}`,
			expected: HTMLTemplate{
				Expression: Expression{
					Value: "Name()",
					Range: Range{
						From: Position{
							Index: 6,
							Line:  0,
							Col:   6,
						},
						To: Position{
							Index: 12,
							Line:  0,
							Col:   12,
						},
					},
				},
				Children: []Node{},
			},
		},
		{
			name: "template: single parameter",
			input: `templ Name(p Parameter) {
}`,
			expected: HTMLTemplate{
				Expression: Expression{
					Value: "Name(p Parameter)",
					Range: Range{
						From: Position{
							Index: 6,
							Line:  0,
							Col:   6,
						},
						To: Position{
							Index: 23,
							Line:  0,
							Col:   23,
						},
					},
				},
				Children: []Node{},
			},
		},
		{
			name: "template: containing element",
			input: `templ Name(p Parameter) {
<span>{ "span content" }</span>
}`,
			expected: HTMLTemplate{
				Expression: Expression{
					Value: "Name(p Parameter)",
					Range: Range{
						From: Position{
							Index: 6,
							Line:  0,
							Col:   6,
						},
						To: Position{
							Index: 23,
							Line:  0,
							Col:   23,
						},
					},
				},
				Children: []Node{
					Element{
						Name:       "span",
						Attributes: []Attribute{},
						Children: []Node{
							StringExpression{
								Expression: Expression{
									Value: `"span content"`,
									Range: Range{
										From: Position{
											Index: 34,
											Line:  1,
											Col:   8,
										},
										To: Position{
											Index: 48,
											Line:  1,
											Col:   22,
										},
									},
								},
							},
						},
					},
					Whitespace{
						Value: "\n",
					},
				},
			},
		},
		{
			name: "template: containing nested elements",
			input: `templ Name(p Parameter) {
<div>
  { "div content" }
  <span>
	{ "span content" }
  </span>
</div>
}`,
			expected: HTMLTemplate{
				Expression: Expression{
					Value: "Name(p Parameter)",
					Range: Range{
						From: Position{
							Index: 6,
							Line:  0,
							Col:   6,
						},
						To: Position{
							Index: 23,
							Line:  0,
							Col:   23,
						},
					},
				},
				Children: []Node{
					Element{
						Name:       "div",
						Attributes: []Attribute{},
						Children: []Node{
							Whitespace{Value: "\n  "},
							StringExpression{
								Expression: Expression{
									Value: `"div content"`,
									Range: Range{
										From: Position{
											Index: 36,
											Line:  2,
											Col:   4,
										},
										To: Position{
											Index: 49,
											Line:  2,
											Col:   17,
										},
									},
								},
							},
							Whitespace{Value: "\n  "},
							Element{
								Name:       "span",
								Attributes: []Attribute{},
								Children: []Node{
									Whitespace{Value: "\n\t"},
									StringExpression{
										Expression: Expression{
											Value: `"span content"`,
											Range: Range{
												From: Position{
													Index: 64,
													Line:  4,
													Col:   3,
												},
												To: Position{
													Index: 78,
													Line:  4,
													Col:   17,
												},
											},
										},
									},
									Whitespace{Value: "\n  "},
								},
							},
							Whitespace{Value: "\n"},
						},
					},
					Whitespace{Value: "\n"},
				},
			},
		},
		{
			name: "template: containing if element",
			input: `templ Name(p Parameter) {
	if p.Test {
		<span>
			{ "span content" }
		</span>
	}
}`,
			expected: HTMLTemplate{
				Expression: Expression{
					Value: "Name(p Parameter)",
					Range: Range{
						From: Position{
							Index: 6,
							Line:  0,
							Col:   6,
						},
						To: Position{
							Index: 23,
							Line:  0,
							Col:   23,
						},
					},
				},
				Children: []Node{
					Whitespace{Value: "\t"},
					IfExpression{
						Expression: Expression{
							Value: `p.Test`,
							Range: Range{
								From: Position{
									Index: 30,
									Line:  1,
									Col:   4,
								},
								To: Position{
									Index: 36,
									Line:  1,
									Col:   10,
								},
							},
						},
						Then: []Node{
							Whitespace{Value: "\t\t"},
							Element{
								Name:       "span",
								Attributes: []Attribute{},
								Children: []Node{
									Whitespace{"\n\t\t\t"},
									StringExpression{
										Expression: Expression{
											Value: `"span content"`,
											Range: Range{
												From: Position{
													Index: 53,
													Line:  3,
													Col:   5,
												},
												To: Position{
													Index: 67,
													Line:  3,
													Col:   19,
												},
											},
										},
									},
									Whitespace{"\n\t\t"},
								},
							},
							Whitespace{
								Value: "\n\t",
							},
						},
						Else: []Node{},
					},
					Whitespace{
						Value: "\n",
					},
				},
			},
		},
		{
			name: "template: inputs",
			input: `templ Name(p Parameter) {
	<input type="text" value="a" />
	<input type="text" value="b" />
}`,
			expected: HTMLTemplate{
				Expression: Expression{
					Value: "Name(p Parameter)",
					Range: Range{
						From: Position{
							Index: 6,
							Line:  0,
							Col:   6,
						},
						To: Position{
							Index: 23,
							Line:  0,
							Col:   23,
						},
					},
				},
				Children: []Node{
					Whitespace{Value: "\t"},
					Element{
						Name: "input",
						Attributes: []Attribute{
							ConstantAttribute{Name: "type", Value: "text"},
							ConstantAttribute{Name: "value", Value: "a"},
						},
					},
					Whitespace{Value: "\n\t"},
					Element{
						Name: "input",
						Attributes: []Attribute{
							ConstantAttribute{Name: "type", Value: "text"},
							ConstantAttribute{Name: "value", Value: "b"},
						},
					},
					Whitespace{Value: "\n"},
				},
			},
		},
		{
			name: "template: doctype",
			input: `templ Name() {
<!DOCTYPE html>
}`,
			expected: HTMLTemplate{
				Expression: Expression{
					Value: "Name()",
					Range: Range{
						From: Position{
							Index: 6,
							Line:  0,
							Col:   6,
						},
						To: Position{
							Index: 12,
							Line:  0,
							Col:   12,
						},
					},
				},
				Children: []Node{
					DocType{
						Value: "html",
					},
					Whitespace{Value: "\n"},
				},
			},
		},
		{
			name: "template: incomplete open tag",
			input: `templ Name() {
				        <div
						{"some string"}
					</div>
}`,
			expected:    HTMLTemplate{},
			expectError: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			input := input.NewFromString(tt.input)
			result := template.Parse(input)
			diff := cmp.Diff(tt.expected, result.Item)
			switch {
			case tt.expectError && result.Error == nil:
				t.Error("expected an error got nil")
			case !tt.expectError && result.Error != nil:
				t.Errorf("parser error: %v", result.Error)
			case tt.expectError == result.Success:
				t.Errorf("Success=%v want=%v", result.Success, !tt.expectError)
			case !tt.expectError && diff != "":
				t.Errorf(diff)
			}
		})
	}
}

func TestTemplateParserErrors(t *testing.T) {
	var tests = []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "template: containing element",
			input: `templ Name(p Parameter) {
<span
}`,
			expected: "closing brace not found",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			input := input.NewFromString(tt.input)
			result := template.Parse(input)
			if result.Error == nil {
				t.Fatalf("expected error %q, got nil", tt.expected)
			}
			if diff := cmp.Diff(tt.expected, result.Error.Error()); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
