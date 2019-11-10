// Package cool provides both a lexer and parser for the Cool
// programming language.
// Here is a simple example, reading a stream of tokens from
// a Cool source string:
//
//	var scan cool.Lexer
//	scan.Init("class MyType inherits IO")
//	for tok := scan.Next(); tok != cool.EOI; tok = scan.Next() {
//		fmt.Printf("token: %s\n", tok)
//	}
//
// Note: All the code in this package (except for the EBNF grammar
// specification and tests) is generated via TextMapper tool. For
// further details please have a look to https://textmapper.org/
package cool
