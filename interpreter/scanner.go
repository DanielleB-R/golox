package interpreter

import (
	"strconv"

	"github.com/DanielleB-R/golox/interpreter/token"
)

type SourceScanner struct {
	source string
	tokens []*token.Token

	start   int
	current int
	line    int
	errors  SourceErrors
}

func NewSourceScanner(source string) SourceScanner {
	return SourceScanner{source: source, tokens: []*token.Token{}, start: 0, current: 0, line: 1, errors: SourceErrors{}}
}

func (s *SourceScanner) ScanTokens() ([]*token.Token, error) {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	s.tokens = append(s.tokens, token.NewToken(token.EOF, "", nil, s.line))
	if len(s.errors) > 0 {
		return nil, s.errors
	}
	return s.tokens, nil
}

func (s *SourceScanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *SourceScanner) advance() byte {
	r := s.source[s.current]
	s.current += 1
	return r
}

func (s *SourceScanner) match(expected byte) bool {
	if s.isAtEnd() {
		return false
	}
	if s.source[s.current] != expected {
		return false
	}

	s.current += 1
	return true
}

func (s *SourceScanner) peek() byte {
	if s.isAtEnd() {
		return 0
	}
	return s.source[s.current]
}

func (s *SourceScanner) peekNext() byte {
	if (s.current + 1) >= len(s.source) {
		return 0
	}
	return s.source[s.current+1]
}

func (s *SourceScanner) addToken(tokenType int, literal interface{}) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, token.NewToken(tokenType, text, literal, s.line))
}

func (s *SourceScanner) scanToken() {
	c := s.advance()

	switch c {
	case '(':
		s.addToken(token.LEFT_PAREN, nil)
	case ')':
		s.addToken(token.RIGHT_PAREN, nil)
	case '{':
		s.addToken(token.LEFT_BRACE, nil)
	case '}':
		s.addToken(token.RIGHT_BRACE, nil)
	case ',':
		s.addToken(token.COMMA, nil)
	case '.':
		s.addToken(token.DOT, nil)
	case '-':
		s.addToken(token.MINUS, nil)
	case '+':
		s.addToken(token.PLUS, nil)
	case ';':
		s.addToken(token.SEMICOLON, nil)
	case '*':
		s.addToken(token.STAR, nil)
	case '!':
		if s.match('=') {
			s.addToken(token.BANG_EQUAL, nil)
		} else {
			s.addToken(token.BANG, nil)
		}
	case '=':
		if s.match('=') {
			s.addToken(token.EQUAL_EQUAL, nil)

		} else {
			s.addToken(token.EQUAL, nil)
		}
	case '<':
		if s.match('=') {
			s.addToken(token.LESS_EQUAL, nil)
		} else {
			s.addToken(token.LESS, nil)
		}
	case '>':
		if s.match('=') {
			s.addToken(token.GREATER_EQUAL, nil)

		} else {
			s.addToken(token.GREATER, nil)
		}
	case '/':
		if s.match('/') {
			// Skip comment, which goes to the end of the line
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(token.SLASH, nil)
		}
	case ' ', '\r', '\t':
		break
	case '\n':
		s.line += 1
	case '"':
		s.string()
	default:
		if isDigit(c) {
			s.number()
		} else if isAlpha(c) {
			s.identifier()
		} else {
			s.errors = append(s.errors, NewSourceError(s.line, "", "Unexpected character."))
		}
	}
}

func (s *SourceScanner) string() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line += 1
		}
		s.advance()
	}

	if s.isAtEnd() {
		s.errors = append(s.errors, NewSourceError(s.line, "", "Unterminated string."))
		return
	}

	// Consume the closing quote
	s.advance()

	value := s.source[s.start+1 : s.current-1]
	s.addToken(token.STRING, value)
}

func (s *SourceScanner) number() {
	for isDigit(s.peek()) {
		s.advance()
	}

	if s.peek() == '.' && isDigit(s.peekNext()) {
		s.advance()

		for isDigit(s.peek()) {
			s.advance()
		}
	}

	value, err := strconv.ParseFloat(s.source[s.start:s.current], 64)
	if err != nil {
		panic("Unexpected numeric conversion error")
	}
	s.addToken(token.NUMBER, value)
}

func (s *SourceScanner) identifier() {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}

	text := s.source[s.start:s.current]
	tokenType, ok := token.Keywords[text]
	if !ok {
		tokenType = token.IDENTIFIER
	}

	s.addToken(tokenType, nil)
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func isAlphaNumeric(c byte) bool {
	return isDigit(c) || isAlpha(c)
}
