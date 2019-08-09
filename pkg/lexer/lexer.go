package lexer

import (
	"bufio"
	"io"
	"strings"
	"unicode"
)

type Lexer struct {
	reader     *bufio.Reader
	line       uint32
	column     uint32
	curRune    rune
	curToken   Token
	retCurrent bool
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

func (l *Lexer) mkToken(t TokenType) Token {
	return Token{t, string(l.curRune), l.line, l.column}
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

func (l *Lexer) readIdentifier() Token {
	startCol := l.column
	pred := func(r rune) bool {
		return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'
	}
	ident := l.gather(pred)

	return Token{IDENT, ident, l.line, startCol}
}

func (l *Lexer) readNumber() Token {
	startCol := l.column
	number := l.gather(unicode.IsDigit)

	return Token{INT, number, l.line, startCol}
}

func (l *Lexer) readString() Token {
	startCol := l.column
	str := []rune{}

	pred := func(r rune) bool {
		return r != '\n'
	}

	for {
		if !l.maybeConsumePred(pred) {
			return Token{INVALID, string(append([]rune{'"'}, str...)), l.line, startCol}
		}

		if l.curRune == '"' {
			break
		}
		str = append(str, l.curRune)
	}
	return Token{STRING, string(str), l.line, startCol}
}

func (l *Lexer) nextToken() Token {
	for {
		var err error
		l.curRune, _, err = l.reader.ReadRune()
		if err != nil {
			return Token{EOF, "", l.line, l.column}
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
				return Token{EQ, "==", l.line, l.column - 1}
			}
			return l.mkToken(ASSIGN)
		case ';':
			return l.mkToken(SEMICOLON)
		case '(':
			return l.mkToken(LPAREN)
		case ',':
			return l.mkToken(COMMA)
		case ')':
			return l.mkToken(RPAREN)
		case '{':
			return l.mkToken(LBRACE)
		case '+':
			return l.mkToken(PLUS)
		case '}':
			return l.mkToken(RBRACE)
		case '!':
			if l.maybeConsume('=') {
				return Token{NOT_EQ, "!=", l.line, l.column - 1}
			}
			return l.mkToken(BANG)
		case '-':
			return l.mkToken(MINUS)
		case '/':
			return l.mkToken(SLASH)
		case '*':
			return l.mkToken(ASTERISK)
		case '<':
			if l.maybeConsume('=') {
				return Token{LE, "<=", l.line, l.column - 1}
			}
			return l.mkToken(LT)
		case '>':
			if l.maybeConsume('=') {
				return Token{GE, ">=", l.line, l.column - 1}
			}
			return l.mkToken(GT)
		default:
			if unicode.IsLetter(l.curRune) {
				ident := l.readIdentifier()
				tokenType := LookupKeyword(ident.Literal)
				ident.Type = tokenType
				return ident
			} else if unicode.IsDigit(l.curRune) {
				return l.readNumber()
			} else if l.curRune == '"' {
				return l.readString()
			}
			return l.mkToken(INVALID)
		}
	}
}

func (l *Lexer) ReadToken() Token {
	if l.retCurrent {
		l.retCurrent = false
	} else {
		l.curToken = l.nextToken()
	}
	return l.curToken
}

func (l *Lexer) UnreadToken() {
	l.retCurrent = true
}
