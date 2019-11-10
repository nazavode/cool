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
error:

# Accept end-of-input in all states.
# <*> eoi: /{eoi}/

whitespace: /[\n\r\t\f\v ]+/ (space)

# Multiline, nested comment blocks
<initial, inComment>
EnterBlockComment:  /\(\*/ (space)
	{ l.enterBlockComment() }

<initial>
invalid_token: /\*\)/

<inComment> {

# TODO report eoi as invalid token in comment
invalid_token: /{eoi}/
	{ l.State = StateInitial }

ExitBlockComment:  /\*\)/ (space)
	{ l.exitBlockComment() }

# TODO
# Still have to figure out how to match \*+ without breaking comments
# like "(***)" (see testdata/longcomment.cool for instance). Just go
# for the slowest solution, at least for now: the rhs is going to
# make the lexer change state at each "*" or "(" or ")" character
# found in a block comment.
# Note: this solution needs no backtracking, lexer tables turned out
# really compact. Still have to do actual measurements to figure out
# whether this one is suboptimal in the real world or not.
BlockComment: /[^\(\)\*]+|[\*\(\)]/ (space)
}

LineComment: /\-\-.*/ (space)

# Identifiers
# All identifier rules conflict with keywords. We cannot use the (class)
# qualifier here since keywords in Cool are not constant words (case
# insensitiveness makes them non-terminal). Just use explicit priority.
ObjectId: /[a-z][_\w]*/ -1
TypeId  : /[A-Z][_\w]*/ -1

# Literals
IntegerLiteral: /\d+/
BoolLiteral   : /t{R}{U}{E}|f{A}{L}{S}{E}/
StringLiteral : /"([^"\n\\]|\\[\\nbtf\n])*"/

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
'>=': />=/
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
A = /(a|A)/
# B = /(b|B)/
C = /(c|C)/
D = /(d|D)/
E = /(e|E)/
F = /(f|F)/
# G = /(g|G)/
H = /(h|H)/
I = /(i|I)/
# J = /(j|J)/
# K = /(k|K)/
L = /(l|L)/
# M = /(m|M)/
N = /(n|N)/
O = /(o|O)/
P = /(p|P)/
# Q = /(q|Q)/
R = /(r|R)/
S = /(s|S)/
T = /(t|T)/
U = /(u|U)/
V = /(v|V)/
W = /(w|W)/
# X = /(x|X)/
# Y = /(y|Y)/
# Z = /(z|Z)/

%%

${template go_lexer.stateVars-}
	commentLevel int // number of open nested block comments
${end}

${template go_lexer.initStateVars-}
	l.commentLevel = 0
${end}

${template newTemplates-}
{{define "onAfterLexer"}}

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
