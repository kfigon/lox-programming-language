package parser

import (
	"lox/lexer"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func parseIt(t *testing.T, input string) []Statement {
	toks, err := lexer.Lex(input)
	require.NoError(t, err, "got lexer error")

	got, errs := NewParser(toks).Parse()
	require.Empty(t, errs, "got parser errors")
	if len(errs) != 0 {
		return nil
	}
	return got
}

func TestParseSingleExpressions(t *testing.T) {
	testCases := []struct {
		desc     string
		input    string
		expected Expression
	}{
		{
			desc:     "simple literal expression",
			input:    "3;",
			expected: Literal(lexer.Token{lexer.Number, "3", 1}),
		},
		{
			desc:     "simple math expression",
			input:    "1 + 3;",
			expected: Binary{Op: lexer.Token{lexer.Operator, "+", 1}, Left: Literal(lexer.Token{lexer.Number, "1", 1}), Right: Literal(lexer.Token{lexer.Number, "3", 1})},
		},
		{
			desc:  "more complicated math expression",
			input: "8*1 + 3 * 2;",
			expected: Binary{
				Op: lexer.Token{lexer.Operator, "+", 1},
				Left: Binary{
					Op:    lexer.Token{lexer.Operator, "*", 1},
					Left:  Literal(lexer.Token{lexer.Number, "8", 1}),
					Right: Literal(lexer.Token{lexer.Number, "1", 1}),
				},
				Right: Binary{
					Op:    lexer.Token{lexer.Operator, "*", 1},
					Left:  Literal(lexer.Token{lexer.Number, "3", 1}),
					Right: Literal(lexer.Token{lexer.Number, "2", 1}),
				},
			},
		},
		{
			desc:  "grouped math expression",
			input: "8*1 / (3 + 2);",
			expected: Binary{
				Op: lexer.Token{lexer.Operator, "/", 1},
				Left: Binary{
					Op:    lexer.Token{lexer.Operator, "*", 1},
					Left:  Literal(lexer.Token{lexer.Number, "8", 1}),
					Right: Literal(lexer.Token{lexer.Number, "1", 1}),
				},
				Right: Binary{
					Op:    lexer.Token{lexer.Operator, "+", 1},
					Left:  Literal(lexer.Token{lexer.Number, "3", 1}),
					Right: Literal(lexer.Token{lexer.Number, "2", 1}),
				},
			},
		},
		{
			desc:  "unary math expression",
			input: "-3;",
			expected: Unary{
				Op: lexer.Token{lexer.Operator, "-", 1},
				Ex: Literal(lexer.Token{lexer.Number, "3", 1}),
			},
		},
		{
			desc:  "binary with unary math expression",
			input: "-3 + 4;",
			expected: Binary{
				Op: lexer.Token{lexer.Operator, "+", 1},
				Left: Unary{
					Op: lexer.Token{lexer.Operator, "-", 1},
					Ex: Literal(lexer.Token{lexer.Number, "3", 1}),
				},
				Right: Literal(lexer.Token{lexer.Number, "4", 1}),
			},
		},
		{
			desc:  "primary with literals",
			input: "3 + foo;",
			expected: Binary{
				Op:    lexer.Token{lexer.Operator, "+", 1},
				Left:  Literal(lexer.Token{lexer.Number, "3", 1}),
				Right: Literal(lexer.Token{lexer.Identifier, "foo", 1}),
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got := parseIt(t, tC.input)
			require.Len(t, got, 1, "single expression expected")

			assert.Equal(t, []Statement{StatementExpression{tC.expected}}, got)
		})
	}
}

func TestParserErrors(t *testing.T) {
	testCases := []struct {
		desc     string
		input    string
		expected string
	}{
		{
			desc:  "unmatched paren on grouping",
			input: "(1+3",
		},
		{
			desc:  "eof on grouping",
			input: "(1+",
		},
		{
			desc:  "eof on binary",
			input: "1+",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			toks, err := lexer.Lex(tC.input)
			require.NoError(t, err, "got lexer error")

			p := NewParser(toks)
			p.Parse()
			_, errs := p.Parse()
			require.NotEmpty(t, errs, "expected parser errors")
		})
	}
}

func TestStatements(t *testing.T) {
	testCases := []struct {
		desc     string
		input    string
		expected []Statement
	}{
		{
			desc:  "let statement",
			input: "let foo = -123;",
			expected: []Statement{
				LetStatement{AssignmentStatement{
					"foo",
					Unary{
						lexer.Token{lexer.Operator, "-", 1},
						Literal(lexer.Token{lexer.Number, "123", 1}),
					},
				},
				}},
		},
		{
			desc:  "assignment statement",
			input: "foo = -123;",
			expected: []Statement{
				AssignmentStatement{
					"foo",
					Unary{
						lexer.Token{lexer.Operator, "-", 1},
						Literal(lexer.Token{lexer.Number, "123", 1}),
					},
				},
			},
		},
		{
			desc:  "block statement",
			input: `let foo = 123;
			{
				foo = 4;
				foo = 18;
			}
			x = true;`,
			expected: []Statement{
				LetStatement{
					AssignmentStatement{ "foo", Literal(lexer.Token{lexer.Number, "123", 1})},
				},
				BlockStatement{
					[]Statement{
						AssignmentStatement{ "foo", Literal(lexer.Token{lexer.Number, "4", 3})},
						AssignmentStatement{ "foo", Literal(lexer.Token{lexer.Number, "18", 4})},
					},
				},
				AssignmentStatement{"x", Literal(lexer.Token{lexer.Boolean, "true", 6})},
			},
		},
		{
			desc:  "single if",
			input: `if(foo == 123) {
				foo = 18;
			}`,
			expected: []Statement{
				IfStatement{
					Ifs: []IfBlock{
						{
							Predicate: Binary{
								Op: lexer.Token{lexer.Operator, "==", 1},
								Left: Literal(lexer.Token{lexer.Identifier, "foo", 1}),
								Right: Literal(lexer.Token{lexer.Number, "123", 1}),
							}, 
							Body: BlockStatement{
								[]Statement{
									AssignmentStatement{"foo", Literal(lexer.Token{lexer.Number, "18", 2})},
								},
							},
						},
					},
				},
			},
		},
		{
			desc:  "else",
			input: `if(foo == 123) {
				foo = 18;
			} else {
				foo = 2;
			}`,
			expected: []Statement{
				IfStatement{
					Ifs: []IfBlock{
						{
							Predicate: Binary{
								Op: lexer.Token{lexer.Operator, "==", 1},
								Left: Literal(lexer.Token{lexer.Identifier, "foo", 1}),
								Right: Literal(lexer.Token{lexer.Number, "123", 1}),
							}, 
							Body: BlockStatement{
								[]Statement{
									AssignmentStatement{"foo", Literal(lexer.Token{lexer.Number, "18", 2})},
								},
							},
						},
						{
							Predicate: Literal(lexer.Token{lexer.Boolean, "true", 3}),
							Body: BlockStatement{
								[]Statement{
									AssignmentStatement{"foo", Literal(lexer.Token{lexer.Number, "2", 4})},
								},
							},
						},
					},
				},
			},
		},
		{
			desc:  "else if",
			input: `if(foo == 123) {
				foo = 18;
			} else if (foo < 3) {
				foo = 2;
			}`,
			expected: []Statement{
				IfStatement{
					Ifs: []IfBlock{
						{
							Predicate: Binary{
								Op: lexer.Token{lexer.Operator, "==", 1},
								Left: Literal(lexer.Token{lexer.Identifier, "foo", 1}),
								Right: Literal(lexer.Token{lexer.Number, "123", 1}),
							}, 
							Body: BlockStatement{
								[]Statement{
									AssignmentStatement{"foo", Literal(lexer.Token{lexer.Number, "18", 2})},
								},
							},
						},
						{
							Predicate: Binary{
								Op: lexer.Token{lexer.Operator, "<", 3},
								Left: Literal(lexer.Token{lexer.Identifier, "foo", 3}),
								Right: Literal(lexer.Token{lexer.Number, "3", 3}),
							}, 
							Body: BlockStatement{
								[]Statement{
									AssignmentStatement{"foo", Literal(lexer.Token{lexer.Number, "2", 4})},
								},
							},
						},
					},
				},
			},
		},
		{
			desc:  "else ifs",
			input: `if(foo == 123) {
				foo = 18;
			} else if (foo < 3) {
				foo = 2;
			} else {
				foo = 1;
			}`,
			expected: []Statement{
				IfStatement{
					Ifs: []IfBlock{
						{
							Predicate: Binary{
								Op: lexer.Token{lexer.Operator, "==", 1},
								Left: Literal(lexer.Token{lexer.Identifier, "foo", 1}),
								Right: Literal(lexer.Token{lexer.Number, "123", 1}),
							}, 
							Body: BlockStatement{
								[]Statement{
									AssignmentStatement{"foo", Literal(lexer.Token{lexer.Number, "18", 2})},
								},
							},
						},
						{
							Predicate: Binary{
								Op: lexer.Token{lexer.Operator, "<", 3},
								Left: Literal(lexer.Token{lexer.Identifier, "foo", 3}),
								Right: Literal(lexer.Token{lexer.Number, "3", 3}),
							}, 
							Body: BlockStatement{
								[]Statement{
									AssignmentStatement{"foo", Literal(lexer.Token{lexer.Number, "2", 4})},
								},
							},
						},
						{
							Predicate: Literal(lexer.Token{lexer.Boolean, "true", 5}),
							Body: BlockStatement{
								[]Statement{
									AssignmentStatement{"foo", Literal(lexer.Token{lexer.Number, "1", 6})},
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got := parseIt(t, tC.input)
			assert.Equal(t, tC.expected, got)
		})
	}
}
