package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseSingleExpressions(t *testing.T) {
	testCases := []struct {
		desc	string
		input 	string
		expected expression
	}{
		{
			desc: "simple math expression",
			input: "1 + 3",
			expected: binary{op: token{operator, "+", 1}, left: literal(token{number, "1", 1}), right: literal(token{number, "3", 1})},
		},
		{
			desc: "more complicated math expression",
			input: "8*1 + 3 * 2",
			expected: binary{
				op: token{operator, "+", 1}, 
				left: binary{
					op: token{operator, "*", 1},
					left: literal(token{number, "8", 1}),
					right: literal(token{number, "1", 1}),
				}, 
				right: binary{
					op: token{operator, "*", 1},
					left: literal(token{number, "3", 1}),
					right: literal(token{number, "2", 1}),
				},
			},
		},
		{
			desc: "grouped math expression",
			input: "8*1 / (3 + 2)",
			expected: binary{
				op: token{operator, "/", 1}, 
				left: binary{
					op: token{operator, "*", 1},
					left: literal(token{number, "8", 1}),
					right: literal(token{number, "1", 1}),
				}, 
				right: binary{
					op: token{operator, "+", 1},
					left: literal(token{number, "3", 1}),
					right: literal(token{number, "2", 1}),
				},
			},
		},
		{
			desc: "unary math expression",
			input: "-3",
			expected: unary{
				op: token{operator, "-", 1}, 
				ex: literal(token{number, "3", 1}),
			},
		},
		{
			desc: "binary with unary math expression",
			input: "-3 + 4",
			expected: binary{
				op: token{operator, "+", 1},
				left: unary{
					op: token{operator, "-", 1}, 
					ex: literal(token{number, "3", 1}),
				},
				right: literal(token{number, "4", 1}),
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			toks,err := lex(tC.input)
			require.NoError(t, err, "got lexer error")

			got, errs := NewParser(toks).Parse()
			require.Empty(t, errs, "got parser errors")
			require.Len(t, got, 1, "single expression expected")

			assert.Equal(t, []statement{{tC.expected}}, got)
		})
	}
}

func TestParserErrors(t *testing.T) {
	testCases := []struct {
		desc	string
		input 	string
		expected string
	}{
		{
			desc: "unmatched paren on grouping",
			input: "(1+3",
		},
		{
			desc: "eof on grouping",
			input: "(1+",
		},
		{
			desc: "eof on binary",
			input: "1+",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			toks,err := lex(tC.input)
			require.NoError(t, err, "got lexer error")

			p := NewParser(toks)
			p.Parse()
			_, errs := p.Parse()
			require.NotEmpty(t, errs, "expected parser errors")
		})
	}
}