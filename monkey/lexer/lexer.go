package lexer

import "monkey/token"

type Lexer struct {
	input        string
	currentIndex int
	nextIndex    int
	character    byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readCharacter()
	return l
}

func (l *Lexer) readCharacter() {
	if l.nextIndex >= len(l.input) {
		l.character = 0
	} else {
		l.character = l.input[l.nextIndex]
	}
	l.currentIndex = l.nextIndex
	l.nextIndex++
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	l.eatWhitespace()

	switch l.character {
	case '=':
		if l.peakNextCharacter() == '=' {
			ch := l.character
			l.readCharacter()
			tok.Type = token.EQUAL
			tok.Value = string(ch) + string(l.character)
		} else {
			tok = newToken(token.ASSIGN, l.character)
		}
	case '+':
		tok = newToken(token.PLUS, l.character)
	case '(':
		tok = newToken(token.LPAREN, l.character)
	case ')':
		tok = newToken(token.RPAREN, l.character)
	case '{':
		tok = newToken(token.LBRACE, l.character)
	case '}':
		tok = newToken(token.RBRACE, l.character)
	case ';':
		tok = newToken(token.SEMICOLON, l.character)
	case ',':
		tok = newToken(token.COMMA, l.character)
	case '!':
		if l.peakNextCharacter() == '=' {
			tok.Type = token.NOT_EQUAL
			ch := l.character
			l.readCharacter()
			tok.Value = string(ch) + string(l.character)
		} else {
			tok = newToken(token.BANG, l.character)
		}
	case '-':
		tok = newToken(token.MINUS, l.character)
	case '/':
		tok = newToken(token.SLASH, l.character)
	case '*':
		tok = newToken(token.ASTERISK, l.character)
	case '<':
		tok = newToken(token.L_THAN, l.character)
	case '>':
		tok = newToken(token.G_THAN, l.character)
	case 0:
		tok = newToken(token.EOF, l.character)
	case '"':
		tok.Type = token.STRING
		tok.Value = l.readString()
	case '[':
		tok = newToken(token.LBRACKET, l.character)
	case ']':
		tok = newToken(token.RBRACKET, l.character)
	default:
		if isLetter(l.character) {
			tok.Value = l.readIdentifier()
			tok.Type = token.LookupIdentifier(tok.Value)
			return tok
		} else if isDigit(l.character) {
			tok.Type = token.INT
			tok.Value = l.readNumber()
			return tok
		} else {
			print(l.character)
			tok = newToken(token.ILLEGAL, l.character)
		}
	}

	l.readCharacter()
	return tok
}

func (l *Lexer) readString() string {
	pos := l.nextIndex
	for {
		l.readCharacter()
		if l.character == '"' || l.character == 0 {
			break
		}
	}
	return l.input[pos:l.currentIndex]
}

func (l *Lexer) peakNextCharacter() byte {
	if l.nextIndex < len(l.input) {
		return l.input[l.nextIndex]
	}
	return 0
}

func (l *Lexer) readNumber() string {
	index := l.currentIndex
	for isDigit(l.character) {
		l.readCharacter()
	}
	return l.input[index:l.currentIndex]
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) eatWhitespace() {
	for l.character == ' ' || l.character == '\t' || l.character == '\n' || l.character == '\r' {
		l.readCharacter()
	}
}

func (l *Lexer) readIdentifier() string {
	index := l.currentIndex
	for isLetter(l.character) {
		l.readCharacter()
	}
	return l.input[index:l.currentIndex]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func newToken(tokenType token.TokenType, character byte) token.Token {
	return token.Token{Type: tokenType, Value: string(character)}
}
