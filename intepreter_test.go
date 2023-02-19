package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInterpretExpression(t *testing.T) {
	testCases := []struct {
		desc     string
		input    string
		expected loxObject
	}{
		{
			desc:     "simple numeric literal",
			input:    "15;",
			expected: toLoxObj(15),
		},
		{
			desc:     "string literal",
			input:    "\"hello\";",
			expected: toLoxObj("hello"),
		},
		{
			desc:     "bool literal",
			input:    "true;",
			expected: toLoxObj(true),
		},
		{
			desc:     "simple expression1",
			input:    "15+5;",
			expected: toLoxObj(20),
		},
		{
			desc:     "simple expression2",
			input:    "15-5;",
			expected: toLoxObj(10),
		},
		{
			desc:     "simple expression3",
			input:    "5-17;",
			expected: toLoxObj(-12),
		},
		{
			desc:     "simple expression4",
			input:    "3*5;",
			expected: toLoxObj(15),
		},
		{
			desc:     "simple expression5",
			input:    "15/5;",
			expected: toLoxObj(3),
		},
		{
			desc:     "unary expr1",
			input:    "-15;",
			expected: toLoxObj(-15),
		},
		{
			desc:     "unary expr2",
			input:    "!false;",
			expected: toLoxObj(true),
		},
		{
			desc:     "complicated expr",
			input:    "5*3+1;",
			expected: toLoxObj(16),
		},
		{
			desc:     "complicated expr2",
			input:    "1+5*3;",
			expected: toLoxObj(16),
		},
		{
			desc:     "complicated expr3",
			input:    "(1+5)*3 + 2;",
			expected: toLoxObj(20),
		},
		{
			desc:     "complicated expr4",
			input:    "1+5*3 + 2;",
			expected: toLoxObj(18),
		},
		{
			desc:     "boolean expr",
			input:    "5*3+1 == 16;",
			expected: toLoxObj(true),
		},
		{
			desc:     "boolean expr2",
			input:    "true == !true;",
			expected: toLoxObj(false),
		},
		{
			desc:     "boolean expr3",
			input:    "!false != !true;",
			expected: toLoxObj(true),
		},
		{
			desc:     "string concantenation and comparison",
			input:    `"foo" + "bar" == "foobar";`,
			expected: toLoxObj(true),
		},
		{
			desc:     "string concantenation",
			input:    `"foo" + "bar";`,
			expected: toLoxObj("foobar"),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			stmts := parseIt(t, tC.input)
			require.Len(t, stmts, 1, "expect single statement")

			expStmt, ok := stmts[0].(expression)
			require.True(t, ok, "single expression statement")

			got, err := expStmt.acceptExpr(newInterpreter())
			require.NoError(t, err, "got intepreter error")
			assert.Equal(t, tC.expected, got)
		})
	}
}

func TestInvalidExpressions(t *testing.T) {
	perform := func(t *testing.T, input string) error {
		toks, err := lex(input)
		require.NoError(t, err, "got lexer error")

		p := NewParser(toks)
		p.Parse()
		exps, errs := p.Parse()
		require.Empty(t, errs, "got parser errors")

		return interpret(exps)
	}
	t.Run("mismatched types", func(t *testing.T) {
		input := `2 * (3 / -"muffin");`
		interpreterErrs := perform(t, input)

		assert.Error(t, interpreterErrs)
	})

	t.Run("not declared variable used", func(t *testing.T) {
		input := `2 + foob;`
		interpreterErrs := perform(t, input)

		assert.Error(t, interpreterErrs)
	})
}

func TestInterpreterWithVariables(t *testing.T) {
	t.Run("let and eval", func(t *testing.T) {
		input := `let x = 4; 
		5 + x;`
		stmts := parseIt(t, input)
		assert.Equal(t, []statement{
			letStatement{assignmentStatement{"x", literal(token{number, "4", 1})}},
			statementExpression{binary{
				op:    token{operator, "+", 2},
				left:  literal(token{number, "5", 2}),
				right: literal(token{identifier, "x", 2}),
			},
			},
		}, stmts)
	})
}
