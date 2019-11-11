package cool

import (
	"github.com/google/go-cmp/cmp"
	"io/ioutil"
	"log"
	"path"
	"testing"
)

var testFiles = []struct {
	file   string
	tokens []Token
}{
	// {"testdata/arith.cool", []Token{}},
	// {"testdata/atoi.cool", []Token{}},
	{"testdata/backslash2.cool", []Token{INVALID_TOKEN, INVALID_TOKEN, INVALID_TOKEN, INVALID_TOKEN, STRINGLITERAL, STRINGLITERAL, STRINGLITERAL}},
	{"testdata/backslash.cool", []Token{STRINGLITERAL, STRINGLITERAL, STRINGLITERAL}},
	{"testdata/badidentifiers.cool", []Token{OBJECTID, DOT, OBJECTID, DOT, OBJECTID, DOT, OBJECTID, INTEGERLITERAL, OBJECTID, OBJECTID, MINUS, OBJECTID, TYPEID, MINUS, TYPEID, INVALID_TOKEN, OBJECTID, TYPEID, MINUS, OBJECTID}},
	{"testdata/badkeywords.cool", []Token{OBJECTID, OBJECTID, OBJECTID, OBJECTID, SEMICOLON, OBJECTID, OBJECTID, TYPEID, OBJECTID, OBJECTID, OBJECTID, OBJECTID, OBJECTID, TYPEID, TYPEID, TYPEID, OBJECTID, TYPEID, THEN, OBJECTID, OBJECTID, TYPEID, OBJECTID}},
	// {"testdata/book_list.cl.cool", []Token{}},
	{"testdata/bothcomments.cool", []Token{OBJECTID, OBJECTID, OBJECTID, OBJECTID, OBJECTID, OBJECTID, INVALID_TOKEN, DOT, IF, INTEGERLITERAL, THEN}},
	{"testdata/comment_in_string.cl.cool", []Token{STRINGLITERAL, STRINGLITERAL, STRINGLITERAL}},
	{"testdata/endcomment.cool", []Token{TYPEID, OBJECTID, OBJECTID, OBJECTID, OBJECTID, INVALID_TOKEN}},
	{"testdata/eofstring.cool", []Token{INVALID_TOKEN}},
	{"testdata/escaped_chars_in_comment.cl.cool", []Token{OBJECTID}},
	{"testdata/escapedeof.cool", []Token{INVALID_TOKEN}},
	// FAIL {"testdata/escapednull.cool", []Token{INVALID_TOKEN}},
	{"testdata/escapedquote.cool", []Token{INVALID_TOKEN}},
	{"testdata/escapedunprintables.cool", []Token{STRINGLITERAL, STRINGLITERAL, STRINGLITERAL, STRINGLITERAL, STRINGLITERAL}},
	// {"testdata/hairyscary.cool", []Token{}},
	{"testdata/integers2.cool", []Token{INTEGERLITERAL, INTEGERLITERAL, INTEGERLITERAL, INTEGERLITERAL, INTEGERLITERAL, INTEGERLITERAL, OBJECTID, INTEGERLITERAL, INTEGERLITERAL, MINUS, INTEGERLITERAL, INTEGERLITERAL, OBJECTID}},
	{"testdata/invalidcharacters.cool", []Token{INVALID_TOKEN, INVALID_TOKEN, INVALID_TOKEN, INVALID_TOKEN, INVALID_TOKEN, INVALID_TOKEN, INVALID_TOKEN, INVALID_TOKEN, INVALID_TOKEN, INVALID_TOKEN, INVALID_TOKEN, INVALID_TOKEN, INVALID_TOKEN, INVALID_TOKEN}},
	{"testdata/invalidinvisible.cool", []Token{OBJECTID, INVALID_TOKEN, INVALID_TOKEN, INVALID_TOKEN, INVALID_TOKEN, OBJECTID}},
	// {"testdata/io.cool", []Token{}},
	{"testdata/keywords.cool", []Token{CASE, CLASS, ELSE, ESAC, BOOLLITERAL, BOOLLITERAL, FI, IF, IN, INHERITS, ISVOID, LET, LOOP, NEW, NOT, OF, POOL, THEN, WHILE}},
	// {"testdata/life.cool", []Token{}},
	{"testdata/lineno2.cool", []Token{INVALID_TOKEN, OBJECTID, STRINGLITERAL}},
	{"testdata/lineno3.cool", []Token{PLUS, OBJECTID, TYPEID, STRINGLITERAL}},
	{"testdata/longcomment.cool", []Token{TYPEID, OBJECTID, TYPEID, OBJECTID, TYPEID, OBJECTID, TYPEID, OBJECTID, TYPEID, TYPEID, OBJECTID, OBJECTID, OBJECTID, SEMICOLON, OBJECTID}},
	{"testdata/longstring_escapedbackslashes.cool", []Token{STRINGLITERAL, STRINGLITERAL}},
	{"testdata/multilinecomment.cool", []Token{NOT, IN, OBJECTID, OBJECTID, LPAREN, MULT, NOT, OBJECTID, OBJECTID, MULT, RPAREN}},
	{"testdata/nestedcomment.cool", nil},
	// {"testdata/new_complex.cool", []Token{}},
	{"testdata/null_in_code.cl.cool", []Token{OBJECTID, OBJECTID, OBJECTID, OBJECTID, ASSIGNGT, INVALID_TOKEN, LTMINUS, RPAREN}},
	// FAIL {"testdata/null_in_string.cl.cool", []Token{INVALID_TOKEN}},
	// FAIL {"testdata/null_in_string_followed_by_tokens.cl.cool", []Token{INVALID_TOKEN, OBJECTID, PLUS}},
	{"testdata/null_in_string_unescaped_newline.cl.cool", []Token{INVALID_TOKEN, OBJECTID, PLUS}},
	// {"testdata/objectid.test.cool", []Token{}},
	{"testdata/opencomment.cool", []Token{INVALID_TOKEN}},
	{"testdata/operators.cool", []Token{SEMICOLON, LBRACE, RBRACE, LPAREN, COMMA, RPAREN, COLON, ATSIGN, DOT, PLUS, MINUS, MULT, DIV, TILDE, LT, ASSIGN, LTMINUS, ASSIGNGT, LTASSIGN, LT, LTASSIGN, LTASSIGN, ASSIGN, LTASSIGN, ASSIGNGT, LT, LTMINUS}},
	// {"testdata/palindrome.cool", []Token{}},
	{"testdata/pathologicalstrings.cool", []Token{STRINGLITERAL, STRINGLITERAL, STRINGLITERAL, STRINGLITERAL, STRINGLITERAL}},
	{"testdata/s03.test.cool", nil},
	{"testdata/s04.test.cool", []Token{INTEGERLITERAL, INTEGERLITERAL, INTEGERLITERAL, INTEGERLITERAL, INTEGERLITERAL, INTEGERLITERAL}},
	{"testdata/s05.test.cool", []Token{STRINGLITERAL}},
	{"testdata/s14.test.cool", []Token{OBJECTID, OBJECTID}},
	{"testdata/s16.test.cool", []Token{INTEGERLITERAL, INTEGERLITERAL, INTEGERLITERAL, INTEGERLITERAL, INTEGERLITERAL, INTEGERLITERAL, INTEGERLITERAL}},
	{"testdata/s19.test.cool", []Token{INVALID_TOKEN, OBJECTID, INVALID_TOKEN}},
	{"testdata/s25.test.cool", []Token{OBJECTID, OBJECTID, OBJECTID, OBJECTID, TYPEID, TYPEID, TYPEID, TYPEID}},
	{"testdata/s26.test.cool", []Token{INTEGERLITERAL, OBJECTID}},
	{"testdata/s31.test.cool", []Token{INVALID_TOKEN}},
	{"testdata/s32.test.cool", []Token{INVALID_TOKEN}},
	{"testdata/s33.test.cool", []Token{OBJECTID, INVALID_TOKEN}},
	{"testdata/s34.test.cool", []Token{OBJECTID, INVALID_TOKEN}},
	{"testdata/simplestrings.cool", []Token{STRINGLITERAL, STRINGLITERAL, STRINGLITERAL, STRINGLITERAL, STRINGLITERAL}},
	// {"testdata/sort_list.cl.cool", []Token{}},
	{"testdata/stringcomment.cool", []Token{STRINGLITERAL, STRINGLITERAL}},
	{"testdata/stringwithescapes.cool", []Token{STRINGLITERAL, INTEGERLITERAL}},
	{"testdata/twice_512_nested_comments.cl.cool", []Token{OBJECTID}},
	{"testdata/validcharacters.cool", []Token{PLUS, DIV, MINUS, MULT, ASSIGN, LT, DOT, TILDE, COMMA, SEMICOLON, COLON, LPAREN, RPAREN, ATSIGN, LBRACE, RBRACE}},
	{"testdata/weirdcharcomment.cool", []Token{OBJECTID}},
	{"testdata/wq0607-c1.cool", []Token{STRINGLITERAL}},
	{"testdata/wq0607-c2.cool", []Token{STRINGLITERAL, INVALID_TOKEN, STRINGLITERAL}},
	{"testdata/wq0607-c3.cool", []Token{STRINGLITERAL, STRINGLITERAL, INVALID_TOKEN, INVALID_TOKEN, STRINGLITERAL, INVALID_TOKEN, STRINGLITERAL, INVALID_TOKEN, STRINGLITERAL, STRINGLITERAL, INVALID_TOKEN, STRINGLITERAL, INVALID_TOKEN, INVALID_TOKEN}},
	{"testdata/wq0607-c4.cool", []Token{STRINGLITERAL, MINUS, STRINGLITERAL, STRINGLITERAL, STRINGLITERAL, INVALID_TOKEN, MINUS, INVALID_TOKEN, STRINGLITERAL, INVALID_TOKEN, STRINGLITERAL, INVALID_TOKEN, STRINGLITERAL, STRINGLITERAL, STRINGLITERAL, INVALID_TOKEN, INVALID_TOKEN, STRINGLITERAL, INVALID_TOKEN, STRINGLITERAL, INVALID_TOKEN, STRINGLITERAL, STRINGLITERAL, INVALID_TOKEN, STRINGLITERAL, INVALID_TOKEN, INVALID_TOKEN, STRINGLITERAL, INVALID_TOKEN}},
	{"testdata/all_else_true.cl.cool", []Token{ELSE, ELSE, ELSE, ELSE, ELSE, ELSE, ELSE, ELSE, ELSE, ELSE, ELSE, ELSE, ELSE, ELSE, ELSE, ELSE, BOOLLITERAL, BOOLLITERAL, BOOLLITERAL, BOOLLITERAL, BOOLLITERAL, BOOLLITERAL, BOOLLITERAL, BOOLLITERAL, TYPEID, TYPEID, TYPEID, TYPEID, TYPEID, TYPEID, TYPEID, TYPEID}},
}

var testSnippets = []struct {
	name   string
	source string
	tokens []Token
}{
	{"LineComment", "class -- this is a class\n  class", []Token{CLASS, CLASS}},
	{"BlockComment0", "(**)", nil},
	{"BlockComment1", "(*(**)*)", nil},
	{"BlockComment2", "(*(*(**)*)*)", nil},
	{"BlockComment3", "(*(*(**)*)(*(**)*)*)", nil},
	{"BlockComment4", "class (*else(*then(*if*)*)class(*(*loop*)pool*)case*) class", []Token{CLASS, CLASS}},
	{"BlockComment5", "class(************)class", []Token{CLASS, CLASS}},
	{"BlockComment6", "class(***)class", []Token{CLASS, CLASS}},
	{"Empty", "", nil},
	{"Identifier", "object Type oBJECT", []Token{OBJECTID, TYPEID, OBJECTID}},
	{"IntegerLiteral", "0 000 0000 01234567890", []Token{INTEGERLITERAL, INTEGERLITERAL, INTEGERLITERAL, INTEGERLITERAL}},
	{"StringLiteral", "\"\" \" \" \" foo \"", []Token{STRINGLITERAL, STRINGLITERAL, STRINGLITERAL}},
	{"Whitespace", "    \t\t \f \v \r\r\r\n\n      ", nil},
	{"BoolLiteral", "true false tRUE fALSE True False", []Token{BOOLLITERAL, BOOLLITERAL, BOOLLITERAL, BOOLLITERAL, TYPEID, TYPEID}},
	{"KeywordClass", "class CLASS Class cLASS", []Token{CLASS, CLASS, CLASS, CLASS}},
	{"KeywordElse", "else ELSE Else eLSE", []Token{ELSE, ELSE, ELSE, ELSE}},
	{"KeywordIf", "if IF If iF", []Token{IF, IF, IF, IF}},
	{"KeywordFi", "fi FI Fi fI", []Token{FI, FI, FI, FI}},
	{"KeywordIn", "in IN In iN", []Token{IN, IN, IN, IN}},
	{"KeywordInherits", "inherits INHERITS Inherits iNHERITS", []Token{INHERITS, INHERITS, INHERITS, INHERITS}},
	{"KeywordIsvoid", "isvoid ISVOID Isvoid iSVOID", []Token{ISVOID, ISVOID, ISVOID, ISVOID}},
	{"KeywordLet", "let LET Let lET", []Token{LET, LET, LET, LET}},
	{"KeywordLoop", "loop LOOP Loop lOOP", []Token{LOOP, LOOP, LOOP, LOOP}},
	{"KeywordPool", "pool POOL Pool pOOL", []Token{POOL, POOL, POOL, POOL}},
	{"KeywordThen", "then THEN Then tHEN", []Token{THEN, THEN, THEN, THEN}},
	{"KeywordWhile", "while WHILE While wHILE", []Token{WHILE, WHILE, WHILE, WHILE}},
	{"KeywordCase", "case CASE Case cASE", []Token{CASE, CASE, CASE, CASE}},
	{"KeywordEsac", "esac ESAC Esac eSAC", []Token{ESAC, ESAC, ESAC, ESAC}},
	{"KeywordNew", "new NEW New nEW", []Token{NEW, NEW, NEW, NEW}},
	{"KeywordOf", "of OF Of oF", []Token{OF, OF, OF, OF}},
	{"KeywordNot", "not NOT Not nOT", []Token{NOT, NOT, NOT, NOT}},
	{"NullInCode", "class\000if", []Token{CLASS, INVALID_TOKEN, IF}},
	{"TerminalPlus", "+", []Token{PLUS}},
	{"TerminalDiv", "/", []Token{DIV}},
	{"TerminalMinus", "-", []Token{MINUS}},
	{"TerminalMult", "*", []Token{MULT}},
	{"TerminalAssign", "=", []Token{ASSIGN}},
	{"TerminalLt", "<", []Token{LT}},
	{"TerminalDot", ".", []Token{DOT}},
	{"TerminalTilde", "~", []Token{TILDE}},
	{"TerminalComma", ",", []Token{COMMA}},
	{"TerminalSemicolon", ";", []Token{SEMICOLON}},
	{"TerminalColon", ":", []Token{COLON}},
	{"TerminalLParen", "(", []Token{LPAREN}},
	{"TerminalRParen", ")", []Token{RPAREN}},
	{"TerminalAt", "@", []Token{ATSIGN}},
	{"TerminalLBrace", "{", []Token{LBRACE}},
	{"TerminalRBrace", "}", []Token{RBRACE}},
	{"InvalidBang", "!", []Token{INVALID_TOKEN}},
	{"InvalidHash", "#", []Token{INVALID_TOKEN}},
	{"InvalidDollar", "$", []Token{INVALID_TOKEN}},
	{"InvalidPercent", "%", []Token{INVALID_TOKEN}},
	{"InvalidHat", "^", []Token{INVALID_TOKEN}},
	{"InvalidPound", "&", []Token{INVALID_TOKEN}},
	{"InvalidUnderscore", "_", []Token{INVALID_TOKEN}},
	{"InvalidGT", ">", []Token{INVALID_TOKEN}},
	{"InvalidQuestion", "?", []Token{INVALID_TOKEN}},
	{"InvalidBacktick", "`", []Token{INVALID_TOKEN}},
	{"InvalidLSub", "[", []Token{INVALID_TOKEN}},
	{"InvalidRSub", "]", []Token{INVALID_TOKEN}},
	{"InvalidBackslash", "\\", []Token{INVALID_TOKEN}},
	{"InvalidPipe", "|", []Token{INVALID_TOKEN}},
	{"TokenAndInvalidBang", "a ! a", []Token{OBJECTID, INVALID_TOKEN, OBJECTID}},
	{"TokenAndInvalidHash", "a # a", []Token{OBJECTID, INVALID_TOKEN, OBJECTID}},
	{"TokenAndInvalidDollar", "a $ a", []Token{OBJECTID, INVALID_TOKEN, OBJECTID}},
	{"TokenAndInvalidPercent", "a % a", []Token{OBJECTID, INVALID_TOKEN, OBJECTID}},
	{"TokenAndInvalidHat", "a ^ a", []Token{OBJECTID, INVALID_TOKEN, OBJECTID}},
	{"TokenAndInvalidPound", "a & a", []Token{OBJECTID, INVALID_TOKEN, OBJECTID}},
	{"TokenAndInvalidUnderscore", "a _ a", []Token{OBJECTID, INVALID_TOKEN, OBJECTID}},
	{"TokenAndInvalidGT", "a > a", []Token{OBJECTID, INVALID_TOKEN, OBJECTID}},
	{"TokenAndInvalidQuestion", "a ? a", []Token{OBJECTID, INVALID_TOKEN, OBJECTID}},
	{"TokenAndInvalidBacktick", "a ` a", []Token{OBJECTID, INVALID_TOKEN, OBJECTID}},
	{"TokenAndInvalidLSub", "a [ a", []Token{OBJECTID, INVALID_TOKEN, OBJECTID}},
	{"TokenAndInvalidRSub", "a ] a", []Token{OBJECTID, INVALID_TOKEN, OBJECTID}},
	{"TokenAndInvalidBackslash", "a \\ a", []Token{OBJECTID, INVALID_TOKEN, OBJECTID}},
	{"TokenAndInvalidPipe", "a | a", []Token{OBJECTID, INVALID_TOKEN, OBJECTID}},
}

func scan(source string) []Token {
	var lex Lexer
	lex.Init(source)
	var tokens []Token
	for cur := lex.Next(); cur != EOI; cur = lex.Next() {
		tokens = append(tokens, cur)
	}
	return tokens
}

func TestLexerSnippets(t *testing.T) {
	for _, tt := range testSnippets {
		t.Run(tt.name, func(t *testing.T) {
			got := scan(tt.source)
			if diff := cmp.Diff(tt.tokens, got); diff != "" {
				t.Errorf("lex mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestLexerFiles(t *testing.T) {
	for _, tt := range testFiles {
		name := path.Base(tt.file)
		data, err := ioutil.ReadFile(tt.file)
		if err != nil {
			log.Fatalln(err)
		}
		source := string(data)
		t.Run(name, func(t *testing.T) {
			got := scan(source)
			if diff := cmp.Diff(tt.tokens, got); diff != "" {
				t.Errorf("lex mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
