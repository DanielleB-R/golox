package interpreter

import "github.com/DanielleB-R/golox/interpreter/token"

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
	default:
		s.errors = append(s.errors, NewSourceError(s.line, "", "Unexpected character."))
	}
}
