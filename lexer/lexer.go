package lexer

import (
	"errors"
	"fmt"
)

// Convert re2c files into go files
//go:generate re2go hoser.re -i -W --tags -o hoser.go

// Format "hoser.go" file (re2c does not do a good job)
//go:generate go fmt .

var (
	ErrBadToken      = errors.New("bad token")
	ErrInvalidString = errors.New("unexpected end to string")
)

// Scan will scan the entire text and push all tokens found to tokens. If there is
// a scanning error, it will be pushed immediately to errCh and scanning will stop. If there
// are no more tokens to scan, nil will be sent on errCh and scanning will stop.
func Scan(text string) (tokens chan Token, errCh chan error) {
	if text[len(text)-1] != '\x00' {
		text = text + "\x00"
	}

	tokens = make(chan Token)
	errCh = make(chan error)
	go func() {
		s := newLexerState(text)
		for {
			token, err := s.lex()
			if err != nil {
				errCh <- fmt.Errorf("failed at input '%v' on line %v col %v: %w", s.currToken(), s.line, s.col(), err)
				return
			}

			if token.Kind == Eof {
				errCh <- nil
				return
			}

			select {
			case tokens <- token:
				s.prevToken = token
			case <-errCh:
				return
			}
		}

	}()
	return tokens, errCh
}

// ScanAll will scan the entire text and return all the tokens in a single array. If an error
// happened, a partial array will be returned with all the tokens that were scanned successfully
// with a non-nil error.
func ScanAll(text string) ([]Token, error) {
	tokens, errCh := Scan(text)

	ret := make([]Token, 0)
	for {
		select {
		case token := <-tokens:
			ret = append(ret, token)
		case err := <-errCh:
			return ret, err
		}
	}
}

type lexerState struct {
	text          string // text is the full piece of code we are parsing
	cursor, token int    // text[token:cursor] = currently parsed token
	marker        int    // for re2c YYBACKUP and YYRESTORE
	line          int    // line starts at 1 for first line and increments after each "\n"
	lineStart     int    // lineStart points to the start of the line. Used to calculate column.
	prevToken     Token  // prevToken tracks the previously given Token (Invalid if none)
}

func newLexerState(text string) *lexerState {
	return &lexerState{
		text:      text,
		cursor:    0,
		token:     0,
		line:      1,
		lineStart: 0,
	}
}

// createToken will create a single token from the current cursor position
func (s *lexerState) createToken(kind TokenKind) Token {
	return Token{
		Kind:  kind,
		Value: s.currToken(),
		Line:  s.line,
		Col:   s.col(),
	}
}

// eol is called when a newline token is read (\n or \r\n). It is responsible for resetting the
// line information and injecting semicolons when appropriate.
func (s *lexerState) lexEol() TokenKind {
	// insert semicolon if the previous token seems like it should have a semicolon after it.
	// follows these rules: https://golang.org/doc/effective_go#semicolons
	switch s.prevToken.Kind {
	case Ident, RParen, RCurlyBrack:
		return Semicolon
	default:
		return Invalid
	}
}

func (s *lexerState) col() int {
	return (s.token - s.lineStart) + 1
}

func (s *lexerState) currToken() string {
	return s.text[s.token:s.cursor]
}
