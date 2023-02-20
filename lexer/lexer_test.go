package lexer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLex(t *testing.T) {
	testCases := []struct {
		desc     string
		input    string
		expected []Token
	}{
		{
			desc:  "identifiers",
			input: ` somevalue ife`,
			expected: []Token{
				{TokType: Identifier, Lexeme: "somevalue"},
				{TokType: Identifier, Lexeme: "ife"},
			},
		},
		{
			desc:  "number",
			input: ` 1234;`,
			expected: []Token{
				{TokType: Number, Lexeme: "1234"},
				{TokType: Semicolon, Lexeme: ";"},
			},
		},
		{
			desc:  "whitespaces",
			input: " \t \n 123\t",
			expected: []Token{
				{TokType: Number, Lexeme: "123"},
			},
		},
		{
			desc: "whitespaces and string",
			input: ` 	 
" fo
o	"	`,
			expected: []Token{
				{TokType: StringLiteral, Lexeme: " fo\no\t"},
			},
		},
		{
			desc:  "string literal",
			input: "\" hello world if 123\" 123",
			expected: []Token{
				{TokType: StringLiteral, Lexeme: " hello world if 123"},
				{TokType: Number, Lexeme: "123"},
			},
		},
		{
			desc:  "operators",
			input: `= == < <= > >= ! !! != || &&`,
			expected: []Token{
				{TokType: Operator, Lexeme: "="},
				{TokType: Operator, Lexeme: "=="},
				{TokType: Operator, Lexeme: "<"},
				{TokType: Operator, Lexeme: "<="},
				{TokType: Operator, Lexeme: ">"},
				{TokType: Operator, Lexeme: ">="},
				{TokType: Operator, Lexeme: "!"},
				{TokType: Operator, Lexeme: "!"},
				{TokType: Operator, Lexeme: "!"},
				{TokType: Operator, Lexeme: "!="},
				{TokType: Operator, Lexeme: "||"},
				{TokType: Operator, Lexeme: "&&"},
			},
		},
		{
			desc:  "operators without spaces",
			input: `==<<=>>=||&&!!!!=`,
			expected: []Token{
				{TokType: Operator, Lexeme: "=="},
				{TokType: Operator, Lexeme: "<"},
				{TokType: Operator, Lexeme: "<="},
				{TokType: Operator, Lexeme: ">"},
				{TokType: Operator, Lexeme: ">="},
				{TokType: Operator, Lexeme: "||"},
				{TokType: Operator, Lexeme: "&&"},
				{TokType: Operator, Lexeme: "!"},
				{TokType: Operator, Lexeme: "!"},
				{TokType: Operator, Lexeme: "!"},
				{TokType: Operator, Lexeme: "!="},
			},
		},
	}
	for _, tC := range testCases {
		stringify := func(xs []Token) []string {
			var out []string
			for _, x := range xs {
				out = append(out, x.String())
			}
			return out
		}
		t.Run(tC.desc, func(t *testing.T) {
			got, err := Lex(tC.input)
			require.NoError(t, err)

			assert.Equal(t, stringify(tC.expected), stringify(got))
		})
	}
}

func TestInvalidInput(t *testing.T) {
	t.Run("invalid string", func(t *testing.T) {
		input := `" hello world `

		_, err := Lex(input)
		assert.Error(t, err)
		assert.Equal(t, "Invalid token at line 1: \" hello world \"", err.Error())
	})
}
