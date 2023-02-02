package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLex(t *testing.T) {
	testCases := []struct {
		desc	string
		input	string
		expected []token
	}{
		{
			desc: "identifiers",
			input: ` somevalue ife`,
			expected: []token{
				{ tokType: identifier, lexeme: "somevalue"},
				{ tokType: identifier, lexeme: "ife"},
			},
		},
		{
			desc: "number",
			input: ` 1234;`,
			expected: []token{
				{ tokType: number, lexeme: "1234"},
				{ tokType: semicolon, lexeme: ";"},
			},
		},
		{
			desc: "whitespaces",
			input: " \t \n 123\t",
			expected: []token{
				{tokType: number, lexeme: "123"},
			},
		},
		{
			desc: "whitespaces and string",
			input: ` 	 
" fo
o	"	`,
			expected: []token{
				{tokType: stringLiteral, lexeme: " fo\no\t"},
			},
		},
		{
			desc: "string literal",
			input: "\" hello world if 123\" 123",
			expected: []token{
				{tokType: stringLiteral, lexeme:  " hello world if 123"},
				{tokType: number, lexeme:  "123"},
			},
		},
		{
			desc: "operators",
			input: `= == < <= > >= ! !! != || &&`,
			expected: []token{
				{ tokType: operator, lexeme: "="},
				{ tokType: operator, lexeme: "=="},
				{ tokType: operator, lexeme: "<"},
				{ tokType: operator, lexeme: "<="},
				{ tokType: operator, lexeme: ">"},
				{ tokType: operator, lexeme: ">="},
				{ tokType: operator, lexeme: "!"},
				{ tokType: operator, lexeme: "!"},
				{ tokType: operator, lexeme: "!"},
				{ tokType: operator, lexeme: "!="},
				{ tokType: operator, lexeme: "||"},
				{ tokType: operator, lexeme: "&&"},
			},
		},
		{
			desc: "operators without spaces",
			input: `==<<=>>=||&&!!!!=`,
			expected: []token{
				{ tokType: operator, lexeme: "=="},
				{ tokType: operator, lexeme: "<"},
				{ tokType: operator, lexeme: "<="},
				{ tokType: operator, lexeme: ">"},
				{ tokType: operator, lexeme: ">="},
				{ tokType: operator, lexeme: "||"},
				{ tokType: operator, lexeme: "&&"},
				{ tokType: operator, lexeme: "!"},
				{ tokType: operator, lexeme: "!"},
				{ tokType: operator, lexeme: "!"},
				{ tokType: operator, lexeme: "!="},
			},
		},
	}
	for _, tC := range testCases {
		stringify := func(xs []token) []string {
			var out []string
			for _, x := range xs {
				out = append(out, x.String())
			}
			return out
		}
		t.Run(tC.desc, func(t *testing.T) {
			got,err := lex(tC.input)
			require.NoError(t, err)
			
			assert.Equal(t, stringify(tC.expected), stringify(got))
		})
	}
}

func TestInvalidInput(t *testing.T) {
	t.Run("invalid string", func(t *testing.T) {
		input := `" hello world `

		_, err := lex(input)
		assert.Error(t, err)
		assert.Equal(t, "Invalid token at line 1: \" hello world \"", err.Error())
	})
}