# Lox programming language

based on https://craftinginterpreters.com/contents.html

## run
```
go run .
```
### more params
* run without params to run interpreter
* run with a file name to run the interpreter on the file itself


## test
```
go test ./...
```


## Grammar

Recursive descend parser - top-down approach - go through the grammar from the top. Precedence - top - lowest, bottom - highest

symbols:
* `|` - or
* `+` - at least once ( > 1)
* `?` - at most once ( <= 1)
* `*` - 0 or more ( >= 0)

```
expression     → equality ;
equality       → comparison ( ( "!=" | "==" ) comparison )* ;
comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
term           → factor ( ( "-" | "+" ) factor )* ;
factor         → unary ( ( "/" | "*" ) unary )* ;
unary          → ( "!" | "-" ) unary
               | primary ;
primary        → NUMBER | STRING | "true" | "false" | "nil"
               | "(" expression ")" ;
```