// generated by Textmapper; DO NOT EDIT

package cool

import (
	"fmt"
)

// Token is an enum of all terminal symbols of the cool language.
type Token int

// Token values.
const (
	UNAVAILABLE Token = iota - 1
	EOI
	INVALID_TOKEN
	ERROR
	WHITESPACE
	ENTERBLOCKCOMMENT // (*
	EXITBLOCKCOMMENT  // *)
	BLOCKCOMMENT
	LINECOMMENT
	OBJECTID
	TYPEID
	INTEGERLITERAL
	BOOLLITERAL
	STRINGLITERAL
	CLASS
	ELSE
	IF
	FI
	IN
	INHERITS
	ISVOID
	LET
	LOOP
	POOL
	THEN
	WHILE
	CASE
	ESAC
	NEW
	OF
	NOT
	LBRACE    // {
	RBRACE    // }
	LPAREN    // (
	RPAREN    // )
	DOT       // .
	SEMICOLON // ;
	COMMA     // ,
	LT        // <
	LTASSIGN  // <=
	ASSIGNGT  // =>
	ATSIGN    // @
	PLUS      // +
	MINUS     // -
	MULT      // *
	DIV       // /
	LTMINUS   // <-
	TILDE     // ~
	COLON     // :
	ASSIGN    // =

	NumTokens
)

var tokenStr = [...]string{
	"EOI",
	"INVALID_TOKEN",
	"ERROR",
	"WHITESPACE",
	"(*",
	"*)",
	"BLOCKCOMMENT",
	"LINECOMMENT",
	"OBJECTID",
	"TYPEID",
	"INTEGERLITERAL",
	"BOOLLITERAL",
	"STRINGLITERAL",
	"CLASS",
	"ELSE",
	"IF",
	"FI",
	"IN",
	"INHERITS",
	"ISVOID",
	"LET",
	"LOOP",
	"POOL",
	"THEN",
	"WHILE",
	"CASE",
	"ESAC",
	"NEW",
	"OF",
	"NOT",
	"{",
	"}",
	"(",
	")",
	".",
	";",
	",",
	"<",
	"<=",
	"=>",
	"@",
	"+",
	"-",
	"*",
	"/",
	"<-",
	"~",
	":",
	"=",
}

func (tok Token) String() string {
	if tok >= 0 && int(tok) < len(tokenStr) {
		return tokenStr[tok]
	}
	return fmt.Sprintf("token(%d)", tok)
}
