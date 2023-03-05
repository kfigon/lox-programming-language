package lexer

import (
	"fmt"
	"unicode"
)

type TokenType int

const (
	Opening TokenType = iota
	Closing
	Operator
	Number
	Boolean
	Keyword
	Identifier
	StringLiteral
	Semicolon
	Comma
	null
)

func (t TokenType) String() string {
	return [...]string{
		"opening",
		"closing",
		"operator",
		"number",
		"boolean",
		"keyword",
		"identifier",
		"stringLiteral",
		"semicolon",
		"comma",
		"null",
	}[t]
}

type Token struct {
	TokType TokenType
	Lexeme  string
	Line    int
}

func (t Token) String() string {
	return fmt.Sprintf("(%v, %v)", t.TokType, t.Lexeme)
}

func CheckToken(tok Token, tokType TokenType, lexeme string) bool {
	return CheckTokenType(tok, tokType) && tok.Lexeme == lexeme
}

func CheckTokenType(tok Token, tokType TokenType) bool {
	return tok.TokType == tokType
}

func isKeyword(word string) bool {
	return word == "let" || word == "while" || word == "return" || word == "else" || word == "if" || word == "function"
}

func Lex(input string) ([]Token, error) {
	out := []Token{}
	idx := 0
	lineNumer := 1
	peek := func() (rune, bool) {
		if idx+1 >= len(input) {
			return 0, false
		}
		return rune(input[idx+1]), true
	}

	currentChar := func() (rune, bool) {
		if idx >= len(input) {
			return 0, false
		}
		return rune(input[idx]), true
	}

	addTok := func(tokTyp TokenType, lexeme string) {
		out = append(out, Token{TokType: tokTyp, Lexeme: lexeme, Line: lineNumer})
	}

	for current, ok := currentChar(); ok; current, ok = currentChar() {
		if unicode.IsSpace(current) {
			if current == '\n' {
				lineNumer++
			}
		} else if current == ')' || current == '}' {
			addTok(Closing, string(current))
		} else if current == ';' {
			addTok(Semicolon, string(current))
		} else if current == ',' {
			addTok(Comma, string(current))
		} else if current == '(' || current == '{' {
			addTok(Opening, string(current))
		} else if current == '+' || current == '-' || current == '*' || current == '/' || current == '%' {
			addTok(Operator, string(current))
		} else if current == '!' || current == '<' || current == '>' || current == '=' {
			if next, ok := peek(); ok && next == '=' {
				idx++
				addTok(Operator, string(current)+"=")
			} else {
				addTok(Operator, string(current))
			}
		} else if current == '|' || current == '&' {
			if next, ok := peek(); ok && next == current {
				idx++
				addTok(Operator, string(current)+string(next))
			} else {
				return nil, fmt.Errorf("invalid boolean operator on line %d", lineNumer)
			}
		} else if current == '"' {
			idx++
			word := readUntil(input, &idx, func(r rune) bool { return r != '"' })
			if next, ok := peek(); ok && next == '"' {
				idx++
				addTok(StringLiteral, word)
			} else {
				return nil, fmt.Errorf("invalid token at line %d: \"%s\"", lineNumer, word)
			}
		} else if unicode.IsDigit(current) {
			num := readUntil(input, &idx, unicode.IsDigit)
			addTok(Number, num)
		} else {
			word := readUntil(input, &idx, unicode.IsLetter)
			tokType := classifyWord(word)
			addTok(tokType, word)
		}
		idx++
	}
	return out, nil
}

func classifyWord(word string) TokenType {
	if isKeyword(word) {
		return Keyword
	} else if word == "true" || word == "false" {
		return Boolean
	} else if word == "null" {
		return null
	}
	return Identifier
}

func readUntil(input string, idx *int, fn func(rune) bool) string {
	out := ""
	out += string(input[*idx])
	for *idx+1 < len(input) {
		next := input[*idx+1]
		if fn(rune(next)) {
			*idx++
			out += string(next)
		} else {
			break
		}
	}
	return out
}
