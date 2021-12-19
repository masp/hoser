package lexer

import (
	"errors"

	"github.com/masp/hoser/token"
)

// Convert re2c files into go files
//go:generate re2go hoser.re -i -W --tags -o hoser.go

// Format "hoser.go" file (re2c does not do a good job)
//go:generate go fmt .

var (
	ErrBadToken      = errors.New("bad token")
	ErrInvalidString = errors.New("unexpected end to string")
)

// ScanAll will scan the entire text and return all the tokens in a single array. If an error
// happened, a partial array will be returned with all the tokens that were scanned successfully
// with a non-nil error.
func ScanAll(file *token.File, text []byte) (ret []token.Token, err error) {
	scanner := NewScanner(file, text)
	for {
		_, tok, _ := scanner.Next()
		if scanner.Errors.Len() > 0 {
			err = scanner.Errors[0]
			return
		}

		if tok == token.Eof {
			return
		}
		ret = append(ret, tok)
	}
}

type Scanner struct {
	file          *token.File
	text          []byte      // text is the full piece of code we are parsing
	cursor, token int         // text[token:cursor] = currently parsed token
	marker        int         // for re2c YYBACKUP and YYRESTORE
	prevToken     token.Token // prevToken tracks the previously given Token (Invalid if none)
	Errors        token.ErrorList
}

func NewScanner(file *token.File, text []byte) *Scanner {
	// Insert null character at end if none exists for lexer to know when to terminate
	if text[len(text)-1] != '\x00' {
		text = append(text, '\x00')
	}
	return &Scanner{
		file:   file,
		text:   text,
		cursor: 0,
		token:  0,
	}
}

func (s *Scanner) Next() (pos token.Pos, tok token.Token, lit string) {
	pos, tok, lit, err := s.lex()
	if err != nil {
		s.Errors.Add(s.file.Position(pos), err)
		return s.Next()
	}

	if tok != token.Comment {
		// ignore comments for the sake of semicolon insertion
		s.prevToken = tok
	}
	return pos, tok, lit
}

// eol is called when a newline token is read (\n or \r\n). It is responsible for resetting the
// line information and injecting semicolons when appropriate.
func (s *Scanner) insertSemi() bool {
	// insert semicolon if the previous token seems like it should have a semicolon after it.
	// follows these rules: https://golang.org/doc/effective_go#semicolons
	switch s.prevToken {
	case token.Ident, token.RParen, token.RCurlyBrack:
		return true
	default:
		return false
	}
}

func (s *Scanner) literal() string          { return string(s.text[s.token:s.cursor]) }
func (s *Scanner) pos() token.Pos           { return s.file.Pos(s.cursor) }
func (s *Scanner) position() token.Position { return s.file.Position(s.pos()) }
