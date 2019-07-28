package lexer

import (
	"bufio"
	"io"
	"strings"
	"unicode"

	"github.com/ljanyst/monkey/pkg/token"
)

type Lexer struct {
	reader  *bufio.Reader
	line    uint32
	column  uint32
	curRune rune
}

func NewLexerFromString(input string) *Lexer {
	return NewLexerFromReader(strings.NewReader(input))
}

func NewLexerFromReader(input io.Reader) *Lexer {
	l := new(Lexer)
	l.reader = bufio.NewReader(input)
	l.line = 1
	return l
}

func (l *Lexer) mkToken(t token.TokenType) token.Token {
	return token.Token{t, string(l.curRune), l.line, l.column}
}

func (l *Lexer) maybeConsume(r rune) bool {
	return l.maybeConsumePred(func(readRune rune) bool { return readRune == r })
}

func (l *Lexer) maybeConsumePred(pred func(rune) bool) bool {
	readRune, _, err := l.reader.ReadRune()
	if err != nil {
		return false
	}
	if pred(readRune) {
		l.column++
		l.curRune = readRune
		return true
	}
	l.reader.UnreadRune()
	return false
}

func (l *Lexer) gather(pred func(rune) bool) string {
	group := []rune{l.curRune}

	for {
		if !l.maybeConsumePred(pred) {
			break
		}
		group = append(group, l.curRune)
	}

	return string(group)
}

func (l *Lexer) readIdentifier() token.Token {
	startCol := l.column
	pred := func(r rune) bool {
		return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'
	}
	ident := l.gather(pred)

	return token.Token{token.IDENT, ident, l.line, startCol}
}

func (l *Lexer) readNumber() token.Token {
	startCol := l.column
	number := l.gather(unicode.IsDigit)

	return token.Token{token.INT, number, l.line, startCol}
}

func (l *Lexer) readString() token.Token {
	startCol := l.column
	str := []rune{}

	pred := func(r rune) bool {
		return r != '\n'
	}

	for {
		if !l.maybeConsumePred(pred) {
			return token.Token{token.INVALID, string(append([]rune{'"'}, str...)), l.line, startCol}
		}

		if l.curRune == '"' {
			break
		}
		str = append(str, l.curRune)
	}
	return token.Token{token.STRING, string(str), l.line, startCol}
}

func (l *Lexer) NextToken() token.Token {
	for {
		var err error
		l.curRune, _, err = l.reader.ReadRune()
		if err != nil {
			return token.Token{token.EOF, "", l.line, l.column}
		}

		if l.curRune == '\n' {
			l.column = 0
			l.line++
			continue
		}
		l.column++

		if unicode.IsSpace(l.curRune) {
			continue
		}

		switch l.curRune {
		case '=':
			if l.maybeConsume('=') {
				return token.Token{token.EQ, "==", l.line, l.column - 1}
			}
			return l.mkToken(token.ASSIGN)
		case ';':
			return l.mkToken(token.SEMICOLON)
		case '(':
			return l.mkToken(token.LPAREN)
		case ',':
			return l.mkToken(token.COMMA)
		case ')':
			return l.mkToken(token.RPAREN)
		case '{':
			return l.mkToken(token.LBRACE)
		case '+':
			return l.mkToken(token.PLUS)
		case '}':
			return l.mkToken(token.RBRACE)
		case '!':
			if l.maybeConsume('=') {
				return token.Token{token.NOT_EQ, "!=", l.line, l.column - 1}
			}
			return l.mkToken(token.BANG)
		case '-':
			return l.mkToken(token.MINUS)
		case '/':
			return l.mkToken(token.SLASH)
		case '*':
			return l.mkToken(token.ASTERISK)
		case '<':
			if l.maybeConsume('=') {
				return token.Token{token.LE, "<=", l.line, l.column - 1}
			}
			return l.mkToken(token.LT)
		case '>':
			if l.maybeConsume('=') {
				return token.Token{token.GE, ">=", l.line, l.column - 1}
			}
			return l.mkToken(token.GT)
		default:
			if unicode.IsLetter(l.curRune) {
				ident := l.readIdentifier()
				tokenType := token.LookupKeyword(ident.Literal)
				ident.Type = tokenType
				return ident
			} else if unicode.IsDigit(l.curRune) {
				return l.readNumber()
			} else if l.curRune == '"' {
				return l.readString()
			}
			return l.mkToken(token.INVALID)
		}
	}
}
