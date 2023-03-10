package interpreter

import (
	"fmt"
	"lox/lexer"
	"lox/parser"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInterpretExpression(t *testing.T) {
	testCases := []struct {
		desc     string
		input    string
		expected LoxObject
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
		{
			desc:     "modulo",
			input:    `15 % 3;`,
			expected: toLoxObj(15 % 3),
		},
		{
			desc:     "modulo 2",
			input:    `15 % 2;`,
			expected: toLoxObj(15 % 2),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			stmts := parseIt(t, tC.input)
			require.Len(t, stmts, 1, "expect single statement")

			expStmt, ok := stmts[0].(parser.Expression)
			require.True(t, ok, "single expression statement")

			got, err := expStmt.AcceptExpr(NewInterpreter())
			require.NoError(t, err, "got intepreter error")
			assert.Equal(t, tC.expected, got)
		})
	}
}

func TestInvalidExpressions(t *testing.T) {
	perform := func(t *testing.T, input string) error {
		toks, err := lexer.Lex(input)
		require.NoError(t, err, "got lexer error")

		p := parser.NewParser(toks)
		p.Parse()
		exps, errs := p.Parse()
		require.Empty(t, errs, "got parser errors")

		return Interpret(exps)
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

	t.Run("not declared variable assigned", func(t *testing.T) {
		input := `foo = 4;`
		interpreterErrs := perform(t, input)

		assert.Error(t, interpreterErrs)
	})
}

func TestInterpreterWithVariables(t *testing.T) {
	t.Run("let and eval", func(t *testing.T) {
		input := `let x = 4; 
		5 + x;`
		stmts := parseIt(t, input)
		assert.Equal(t, []parser.Statement{
			parser.LetStatement{parser.AssignmentStatement{"x", parser.Literal(lexer.Token{lexer.Number, "4", 1})}},
			parser.StatementExpression{parser.Binary{
				Op:    lexer.Token{lexer.Operator, "+", 2},
				Left:  parser.Literal(lexer.Token{lexer.Number, "5", 2}),
				Right: parser.Literal(lexer.Token{lexer.Identifier, "x", 2}),
			},
			},
		}, stmts)
	})
}

func TestInterpreter(t *testing.T) {
	testCases := []struct {
		desc     string
		input    string
		expected LoxObject
	}{
		{
			desc:     "simple number declaration",
			input:    `let result = 14;`,
			expected: toLoxObj(14),
		},
		{
			desc:     "simple string declaration",
			input:    `let result = "foobar";`,
			expected: toLoxObj("foobar"),
		},
		{
			desc:     "arithmetic and boolean",
			input:    `let x = 123;
			let y = 5;
			let result = (x + y) == 128;`,
			expected: toLoxObj(true),
		},
		{
			desc:     "arithmetic and boolean 2",
			input:    `let x = 123;
			let y = 6;
			let result = (x + y) == 128;`,
			expected: toLoxObj(false),
		},
		{
			desc:     "arithmetic and boolean 3",
			input:    `let x = 5;
			let y = 6;
			let result = (x + y)*3-1 == 32;`,
			expected: toLoxObj(true),
		},
		{
			desc:     "arithmetic and boolean 4",
			input:    `let x = 5;
			let y = 6;
			let result = (x + y)*3-1 + 12;`,
			expected: toLoxObj(44),
		},
		{
			desc:     "assignment",
			input:    `let result = 3;
			result = result + 4;`,
			expected: toLoxObj(7),
		},
		{
			desc:     "if",
			input:    `let result = 0;
			let x = 1;
			if (x == 1) {result = 1; }`,
			expected: toLoxObj(1),
		},
		{
			desc:     "if not executed",
			input:    `let result = -1;
			let x = 0;
			if (x == 1) {result = 1; }`,
			expected: toLoxObj(-1),
		},
		{
			desc:     "else",
			input:    `let result = 0;
			let x = 0;
			if (x == 1) {result = 1; }
			else {result=2;}`,
			expected: toLoxObj(2),
		},
		{
			desc:     "if else",
			input:    `let result = 0;
			let x = 2;
			if (x == 1) {result = 1; }
			else if (x == 2) {result = 2;}
			else { result=3; }`,
			expected: toLoxObj(2),
		},
		{
			desc:     "if with logic operators",
			input:    `let result = 0;
			let x = 2;
			if (x == 1 || true) {result = 1; }`,
			expected: toLoxObj(1),
		},
		{
			desc:     "if with logic operators",
			input:    `let result = 0;
			let x = 2;
			if (x == 2 && !false) {result = 1; }`,
			expected: toLoxObj(1),
		},
		{
			desc:     "if with logic operators",
			input:    `let result = -2;
			let x = 2;
			if (x == 2 && false) {result = 1; }`,
			expected: toLoxObj(-2),
		},
		{
			desc:     "if with logic more operators",
			input:    `let result = 0;
			if (2 == 1+1 && true && (13 > 2 && 3 >= 3)) {result = 1; }`,
			expected: toLoxObj(1),
		},
		{
			desc:     "while",
			input:    `let result = 0;
			let i = 1;
			while (i < 5) {
				result = result + i;
				i = i + 1;
			}`,
			expected: toLoxObj(10),
		},
		{
			desc:     "function",
			input:    `let result = 0;
			function foo() {
				result = 1;
			}
			foo();`,
			expected: toLoxObj(1),
		},
		{
			desc:     "function with arg",
			input:    `let result = 0;
			function foo(a) {
				result = a;
			}
			foo(12);`,
			expected: toLoxObj(12),
		},
		{
			desc:     "function with args",
			input:    `let result = 0;
			function foo(a,b) {
				result = a + b;
			}
			foo(12, 3+2);`,
			expected: toLoxObj(17),
		},
		{
			desc:     "function double called",
			input:    `function foo(a) {
				result = result + a;
			}
			let result = 0;
			foo(2);
			foo(3);`,
			expected: toLoxObj(5),
		},
		{
			desc:     "std lib call",
			input:    `let result = 0;
			print("foo");`,
			expected: toLoxObj(0),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			statements := parseIt(t, tC.input)
			in := NewInterpreter()

			execute(t, in, statements)

			assertVariable(t, tC.expected, "result", in)
		})
	}
}

func TestScoping(t *testing.T) {
	t.Run("inner variables are not visible", func(t *testing.T) {
		input := `let foo = 4;
		{
			let bar = 123;
			foo = bar;
		}
		foo = foo + 1;`
		statements := parseIt(t, input)
		in := NewInterpreter()

		execute(t, in, statements)

		assertVariable(t, toLoxObj(124), "foo", in)
		assertVariableNotFound(t, "bar", in)
	})
}

func assertVariable[T any](t *testing.T, exp T, name string, i *Interpreter) {
	v, ok := i.env.get(name)
	require.True(t, ok, fmt.Sprintf("%v variable not found", name))
	assert.Equal(t, exp, v)
}

func assertVariableNotFound(t *testing.T, name string, i *Interpreter) {
	v, ok := i.env.get(name)
	assert.False(t, ok, fmt.Sprintf("%v variable found, but not expected, v: %v", name, v))
}

func parseIt(t *testing.T, input string) []parser.Statement {
	toks, err := lexer.Lex(input)
	require.NoError(t, err, "got lexer error")

	got, errs := parser.NewParser(toks).Parse()
	require.Empty(t, errs, "got parser errors")

	return got
}

func execute(t *testing.T, in *Interpreter, stmts []parser.Statement) {
	for _, st := range stmts {
		err := st.AcceptStatement(in)
		require.NoError(t, err)
	}
}
