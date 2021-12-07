package parser

import (
	"errors"

	"github.com/masp/hoser/ast"
	"github.com/masp/hoser/lexer"
)

var (
	ErrInvalidMapValue   = errors.New("invalid value in map, must be 'key: value' pair")
	ErrInvalidEntryKey   = errors.New("key of an entry must be an identifier")
	ErrInvalidEntryValue = errors.New("value of an entry must be an identifier")
	ErrExpectedEndOfMap  = errors.New("expected ] to end map")
)

func (s *parserState) parseEntry(left ast.Expression, token lexer.Token) (*ast.Entry, error) {
	right, err := s.parseExpression(token.Kind)
	if err != nil {
		return nil, err
	}

	if left, ok := left.(*ast.Identifier); ok {
		return &ast.Entry{Key: *left, Val: right}, nil
	}
	return nil, ErrInvalidEntryKey
}

// parseEntries is called when the first , is encountered in a map. parseEntry is called for the
// first entry and then for every entry after that.
//
// key: value -> Entry
// key: value, key2: value2 -> EntryList
func (s *parserState) parseEntries(left ast.Expression, token lexer.Token) (*ast.Map, error) {
	if first, ok := left.(*ast.Entry); ok {
		entries := []ast.Entry{*first}
		for {
			next, err := s.parseExpression(lexer.Comma)
			if err != nil {
				return nil, err
			}

			if entry, ok := next.(*ast.Entry); ok {
				entries = append(entries, *entry)
			} else {
				return nil, ErrInvalidEntryValue
			}

			comma, err := s.peek()
			if err != nil {
				return nil, err
			}

			if comma.Kind != lexer.Comma {
				break
			}
		}
		key, _ := first.Span()
		return &ast.Map{
			StartToken: key,
			Entries:    entries,
		}, nil
	} else {
		return nil, ErrInvalidMapValue
	}
}
