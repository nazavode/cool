# Grammar of the Cool programming language
# as specified here:
# http://theory.stanford.edu/~aiken/software/cool/cool.html

language cool(go);

lang = "cool"
package = "github.com/nazavode/cool"

:: lexer

%s initial;
%x inComment;

invalid_token:

invalid_token: /\x00/
	{ l.invalidTokenClass = InvalidTokenNullCharInCode }

whitespace: /[\n\r\t\f\v ]+/ (space)

# Multiline, nested comment blocks
<initial, inComment>
EnterBlockComment:  /\(\*/ (space)
	{ l.enterBlockComment() }

<initial>
invalid_token: /\*\)/
{ l.invalidTokenClass = InvalidTokenUnmatchedBlockComment }

# Sublexer dealing with the stack of open blocks
<inComment> {

	invalid_token: /{eoi}/ { 
		l.State = StateInitial
		l.invalidTokenClass = InvalidTokenEoiInComment }

	ExitBlockComment:  /\*\)/ (space)
		{ l.exitBlockComment() }

	# TODO
	# Still have to figure out how to match \*+ without breaking comments
	# like "(***)" (see testdata/longcomment.cool for instance). Just go
	# for the slowest solution, at least for now: the rhs is going to
	# match each "*" or "(" or ")" character found in a block comment.
	# Note: this solution needs no backtracking, lexer tables turned out
	# really compact compared to a more greedy attempt. Still have to carry
	# out actual measurements to figure out whether this one is really
	# suboptimal with actual source files or not.
	BlockComment: /[^\(\)\*]+|[\*\(\)]/ (space)

} # <inComment>

LineComment: /\-\-.*/ (space)

# Identifiers
# All identifier rules conflict with keywords. We cannot use the (class)
# qualifier here since keywords in Cool are not constant words (case
# insensitiveness makes them kinda non-terminals for TextMapper). Just
# use explicit priority.
ObjectId: /[a-z]\w*/ -1
TypeId  : /[A-Z]\w*/ -1

# Literals
IntegerLiteral: /\d+/
BoolLiteral   : /t{R}{U}{E}|f{A}{L}{S}{E}/
## String literal:
strEscape    = /\\[^\x00]/
strChar      = /[^"\n\\\x00]/
strRune      = /{strChar}|{strEscape}/
## Make sure to report an ill formed string literal as a single
## invalid_token to make the lexer restart scanning right after
## the closing \". A string literal is ill-formed when:
##   1. contains at least one '\0' (both escaped and raw):
invalid_token: /"({strRune}*\x00{strRune}*)+"/
	{ l.invalidTokenClass = InvalidTokenNullCharInString }
invalid_token: /"({strRune}*\\\x00{strRune}*)+"/
	{ l.invalidTokenClass = InvalidTokenEscapedNullCharInString }
##   2. contains end-of-input:
invalid_token: /"{strRune}*{eoi}/
	{ l.invalidTokenClass = InvalidTokenEoiInString }
##   3. contains at least one raw (non-escaped) '\n':
#    Note: It's unclear from the language spec whether multiple unescaped '\n'
#          should produce a single invalid token or not. No golden files with
#          this case are available but 's19.test.cool' shows that a single '\n'
#          splits the invalid literal in two lexable halves. Leaving the rule
#          commented out while looking for clarifications.
# invalid_token: /"({strRune}*([^\\]?\n){strRune}*)+"/  # <- This needs backtracking!
StringLiteral: /"{strRune}*"/

# Keywords (case insensitive)
class   : /{C}{L}{A}{S}{S}/
else    : /{E}{L}{S}{E}/
if      : /{I}{F}/
fi      : /{F}{I}/
in      : /{I}{N}/
inherits: /{I}{N}{H}{E}{R}{I}{T}{S}/
isvoid  : /{I}{S}{V}{O}{I}{D}/
let     : /{L}{E}{T}/
loop    : /{L}{O}{O}{P}/
pool    : /{P}{O}{O}{L}/
then    : /{T}{H}{E}{N}/
while   : /{W}{H}{I}{L}{E}/
case    : /{C}{A}{S}{E}/
esac    : /{E}{S}{A}{C}/
new     : /{N}{E}{W}/
of      : /{O}{F}/
not     : /{N}{O}{T}/

# Punctuation
'{': /\{/
'}': /\}/
'(': /\(/
')': /\)/
'.': /\./
';': /;/
',': /,/
'<': /</
'<=': /<=/
'=>': /=>/
'@': /@/
'+': /\+/
'-': /-/
'*': /\*/
'/': /\//
'<-': /<-/
'~': /~/
':': /:/
'=': /=/

# Case insensitive chars
A = /a|A/
# B = /b|B/
C = /c|C/
D = /d|D/
E = /e|E/
F = /f|F/
# G = /g|G/
H = /h|H/
I = /i|I/
# J = /j|J/
# K = /k|K/
L = /l|L/
# M = /m|M/
N = /n|N/
O = /o|O/
P = /p|P/
# Q = /q|Q/
R = /r|R/
S = /s|S/
T = /t|T/
U = /u|U/
V = /v|V/
W = /w|W/
# X = /x|X/
# Y = /y|Y/
# Z = /z|Z/

%%

${template go_lexer.stateVars-}
	commentLevel int // number of open nested block comments
	invalidTokenClass InvalidTokenClass // reason for the last invalid token found
${end}

${template go_lexer.initStateVars-}
	l.commentLevel = 0
	l.invalidTokenClass = InvalidTokenUnknown
${end}

${template newTemplates-}
{{define "onAfterLexer"}}

type InvalidTokenClass int

const (
	InvalidTokenUnknown = iota - 1
	InvalidTokenEoiInComment
        InvalidTokenEoiInString
	InvalidTokenUnterminatedStringLiteral
	InvalidTokenNullCharInString
        InvalidTokenEscapedNullCharInString
	InvalidTokenNullCharInCode
	InvalidTokenUnmatchedBlockComment
)

// InvalidTokenReason returns the error class that led to the
// last invalid token found during lexing.
func (l *Lexer) InvalidTokenReason() InvalidTokenClass {
	return l.invalidTokenClass
}

// enterBlockComment marks the beginning of a comment block
// and makes the lexer to transition to "inComment" state.
func (l *Lexer) enterBlockComment() {
	l.commentLevel++
	l.State = StateInComment
}

// exitBlockComment marks the end of a comment block
// and makes the lexer to transition to "initial" state
// if no other blocks are still open.
func (l *Lexer) exitBlockComment() {
	l.commentLevel--
	if l.commentLevel <= 0 {
		l.State = StateInitial
	}
}
{{end}}
${end}
