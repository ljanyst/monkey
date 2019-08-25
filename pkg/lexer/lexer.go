package lexer

import (
	"bufio"
	"io"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Lexer struct {
	reader     *bufio.Reader
	line       uint32
	column     uint32
	curRune    rune
	curToken   Token
	retCurrent bool
	fileName   string
}

func NewLexerFromString(input, name string) *Lexer {
	return NewLexerFromReader(strings.NewReader(input), name)
}

func NewLexerFromReader(input io.Reader, name string) *Lexer {
	l := new(Lexer)
	l.reader = bufio.NewReader(input)
	l.line = 1
	l.fileName = name
	return l
}

func (l *Lexer) mkToken(t TokenType) Token {
	return Token{t, string(l.curRune), l.line, l.column, &l.fileName}
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

	return Token{IDENT, ident, l.line, startCol, &l.fileName}
}

func (l *Lexer) readNumber() Token {
	startCol := l.column
	number := l.gather(unicode.IsDigit)

	return Token{INT, number, l.line, startCol, &l.fileName}
}

func (l *Lexer) readString(delimiter rune) Token {
	startCol := l.column
	str := []rune{}

	pred := func(r rune) bool {
		return r != '\n'
	}

	for {
		if !l.maybeConsumePred(pred) {
			return Token{INVALID, string(append([]rune{'"'}, str...)), l.line, startCol, &l.fileName}
		}

		if l.curRune == delimiter {
			break
		}
		str = append(str, l.curRune)
	}
	return Token{STRING, string(str), l.line, startCol, &l.fileName}
}

func (l *Lexer) readRune() Token {
	tok := l.readString('\'')
	if tok.Type == INVALID || utf8.RuneCountInString(tok.Literal) != 1 {
		tok.Type = INVALID
		return tok
	}
	tok.Type = RUNE
	return tok
}

func (l *Lexer) nextToken() Token {
	for {
		var err error
		l.curRune, _, err = l.reader.ReadRune()
		if err != nil {
			return Token{EOF, "", l.line, l.column, &l.fileName}
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
				return Token{EQ, "==", l.line, l.column - 1, &l.fileName}
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
				return Token{NOT_EQ, "!=", l.line, l.column - 1, &l.fileName}
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
				return Token{LE, "<=", l.line, l.column - 1, &l.fileName}
			}
			return l.mkToken(LT)
		case '>':
			if l.maybeConsume('=') {
				return Token{GE, ">=", l.line, l.column - 1, &l.fileName}
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
				return l.readString('"')
			} else if l.curRune == '\'' {
				return l.readRune()
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
