# cool

EBNF grammar, lexer and parser for the
[Cool programming language](https://en.wikipedia.org/wiki/Cool_(programming_language)).

## Grammar

The full language grammar is described in [cool.tm](cool.tm) using the
EBNF-like format specified by [TextMapper](https://github.com/inspirer/textmapper).

## How to generate lexer and parser

Firstly we need to make sure the [TextMapper](https://github.com/inspirer/textmapper) tool is available:

```console
$ go get github.com/inspirer/textmapper/tm-go/cmd/textmapper
```

Then, we can regenerate all the `go` code for both the lexer and parser:

```console
$ git clone https://github.com/nazavode/cool.git
$ cd cool
$ make
```
