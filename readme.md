# Lox programming language

based on https://craftinginterpreters.com/contents.html

## run
```
make run
```
### more params
* run without params to run interpreter
* `go run . ` with a file name to run the interpreter on the file itself


## test
```
make test
```


## Grammar

Recursive descend parser - top-down approach - go through the grammar from the top. Precedence - top - lowest, bottom - highest

symbols:
* `|` - or
* `+` - at least once ( > 1)
* `?` - at most once ( <= 1)
* `*` - 0 or more ( >= 0)

```
program        → statement* ;

statement      → letDecl
               | assignment
               | block
               | exprStmt 
               | ifStmt;

block          → "{" statement* "}" ;
letDecl        → "let" assignment
assignment     → IDENTIFIER "=" exprStmt

ifStmt         → "if" "(" expression ")" block
                 ( "else" ifStmt )* 
                 ( "else" block )?;

exprStmt       → expression ";" ;

expression     → equality ;
equality       → comparison ( ( "!=" | "==" | "||" | "&&" ) comparison )* ;
comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
term           → factor ( ( "-" | "+" ) factor )* ;
factor         → unary ( ( "/" | "*" | "%" ) unary )* ;
unary          → ( "!" | "-" ) unary
               | primary ;
primary        → NUMBER | STRING | "true" | "false" | "nil"
               | "(" expression ")" 
               | IDENTIFIER ;
```

some notes:
* in C languages assignments are expessions, not statements, so we can do
`newPoint(x + 2, 0).y = 3;`, but here it's a statement