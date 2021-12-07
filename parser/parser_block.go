package parser

import (
	"errors"
	"fmt"

	"github.com/masp/hoser/ast"
	"github.com/masp/hoser/lexer"
)

var (
	ErrExpectedFnBody    = errors.New("expected start of function body with '{'")
	ErrInputOutputNotMap = errors.New("input & output of function must be maps")
	ErrInvalidBlockName  = errors.New("invalid block name, expected identifier")
	// An empty function body is different than no function body
	// `main() {}` -> empty
	// vs.
	// `main()` -> nil
	EmptyFnBody = []ast.Expression{}
)

func (s *parserState) parseBlock() (*ast.Block, error) {
	ident, err := s.eatOnly(lexer.Ident)
	if err == ErrUnexpectedEnd {
		// It is fine if no more tokens in this case (no more functions)
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	name, err := s.parseIdentifier(ident)
	if err != nil {
		return nil, err
	}

	inputs, err := s.parseArgs()
	if err != nil {
		return nil, err
	}

	// output is optional
	// e.g. `main() () {}`` is equivalent to `main() {}`
	next, err := s.peek()
	if err != nil {
		return nil, err
	}

	outputs := &ast.Map{StartToken: next}
	if next.Kind == lexer.LParen {
		// parse output definition
		if outputs, err = s.parseArgs(); err != nil {
			return nil, err
		}
	} else if next.Kind != lexer.LCurlyBrack {
		return nil, ErrExpectedFnBody
	}

	body, err := s.parseFnBody()
	if err != nil {
		return nil, err
	}

	return &ast.Block{
		Name:    name,
		Inputs:  inputs,
		Outputs: outputs,
		Body:    body,
	}, nil
}

// parseArgs takes either the input or output arguments specification and converts it to a Map
// example:
// ([name: string, value: int]) -> Map{{Key: name, Val: string}, {Key: value, Val: int}}
func (s *parserState) parseArgs() (*ast.Map, error) {
	if _, err := s.eatOnly(lexer.LParen); err != nil {
		return nil, err
	}

	next, err := s.peek()
	if err != nil {
		return nil, err
	}

	if next.Kind != lexer.RParen {
		inputs, err := s.parseExpression(lexer.Invalid)
		if err != nil {
			return nil, err
		}

		if _, err := s.eatOnly(lexer.RParen); err != nil {
			return nil, err
		}

		switch ent := inputs.(type) {
		case *ast.Map:
			return ent, nil
		case *ast.Entry:
			key, _ := ent.Span()
			return &ast.Map{StartToken: key, Entries: []ast.Entry{*ent}}, nil
		default:
			return nil, ErrInputOutputNotMap
		}
	} else {
		s.eat()
	}
	return &ast.Map{StartToken: next}, nil
}

func (s *parserState) parseFnBody() ([]ast.Expression, error) {
	if next, err := s.peek(); next.Kind == lexer.Semicolon || err == ErrUnexpectedEnd {
		// handle empty block declarations like:
		// integer(v: int) (t: int)
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	if _, err := s.eatOnly(lexer.LCurlyBrack); err != nil {
		return nil, err
	}

	exprs := make([]ast.Expression, 0)
	for {
		token, err := s.peek()
		if err != nil {
			return nil, err
		}

		if token.Kind == lexer.RCurlyBrack {
			// no more expressions in body
			break
		}

		expr, err := s.parseExpression(lexer.Invalid)
		if err != nil {
			return nil, err
		}
		exprs = append(exprs, expr)

		// expressions always end in either ; or }. If there's one or more ;, we should eat them.
		s.eatAll(lexer.Semicolon)
	}

	if _, err := s.eatOnly(lexer.RCurlyBrack); err != nil {
		return nil, err
	}
	return exprs, nil
}

func (s *parserState) parseBlockCall(left ast.Expression, token lexer.Token) (*ast.BlockCall, error) {
	if name, ok := left.(*ast.Identifier); ok {
		next, err := s.peek()
		if err != nil {
			return nil, err
		}

		result := &ast.BlockCall{Name: name, Args: &ast.Map{StartToken: token}}
		if next.Kind != lexer.RParen {
			// Let's parse the args (non-empty), and we continue parsing until we match our end paren and reset precedence
			args, err := s.parseExpression(lexer.Invalid)
			if err != nil {
				return nil, err
			}

			switch ent := args.(type) {
			case *ast.Map:
				result.Args = ent
			case *ast.Entry:
				key, _ := ent.Span()
				result.Args = &ast.Map{StartToken: key, Entries: []ast.Entry{*ent}}
			default:
				return nil, fmt.Errorf("invalid argument syntax to call a block %T, must be a Map", args)
			}
		}
		if _, err := s.eatOnly(lexer.RParen); err != nil {
			return nil, err
		}
		return result, nil
	}
	return nil, ErrInvalidBlockName
}

func (s *parserState) parseReturn(token lexer.Token) (*ast.Return, error) {
	val, err := s.parseExpression(lexer.Invalid)
	if err != nil {
		return nil, err
	}

	return &ast.Return{
		Token: token,
		Value: val,
	}, nil
}
