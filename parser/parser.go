package parser

import (
	"errors"
	"fmt"

	"github.com/masp/hoser/lexer"
)

type parserState struct {
	tokens    <-chan lexer.Token
	errCh     <-chan error
	nextToken lexer.Token
}

var (
	ErrUnexpectedEnd      = errors.New("unexpected end of text")
	ErrExpectedExpression = errors.New("expected expression")
)

func Scan(tokens <-chan lexer.Token, errCh chan error) (*Module, error) {
	s := parserState{
		tokens: tokens,
		errCh:  errCh,
	}

	return s.parseModule()
}

func (s *parserState) fetchToken() (lexer.Token, error) {
	select {
	case token := <-s.tokens:
		return token, nil
	case err := <-s.errCh:
		if err == nil {
			// The lexer will send nil on the errChan once input is finished
			return lexer.Token{}, ErrUnexpectedEnd
		}
		return lexer.Token{}, err
	}
}

func (s *parserState) peek() (lexer.Token, error) {
	if s.nextToken.Kind != lexer.Invalid {
		return s.nextToken, nil
	}

	next, err := s.fetchToken()
	if err != nil {
		return lexer.Token{}, err
	}
	s.nextToken = next // cache it so when we eat we eat the cached and not from the token stream
	return next, nil
}

func (s *parserState) eat() (lexer.Token, error) {
	next, err := s.peek()
	if err != nil {
		return lexer.Token{}, err
	}
	s.nextToken = lexer.Token{}
	return next, nil
}

func (s *parserState) eatOnly(tk lexer.TokenKind) (lexer.Token, error) {
	token, err := s.eat()
	if err != nil {
		return lexer.Token{}, err
	}
	if token.Kind != tk {
		return lexer.Token{}, fmt.Errorf("expected token %v, got %v", tk, token)
	}
	return token, nil
}

func (s *parserState) eatAll(tk lexer.TokenKind) error {
	for {
		next, err := s.peek()
		if err != nil {
			return err
		}

		if next.Kind != tk {
			return nil
		}
		s.eat()
	}
}

func precedence(kind lexer.TokenKind) int {
	switch kind {
	case lexer.Semicolon:
		return -1 // End of expression always
	case lexer.Invalid:
		return 0 // lexer.Invalid is a special token meaning no parent
	case lexer.Equals:
		return 1
	case lexer.Comma:
		return 3
	case lexer.Colon:
		return 4
	case lexer.LParen:
		return 5
	default:
		// Every other token is lower precedence than these and signal an end to an expression
		return -1
	}
}

func (s *parserState) parseExpression(parent lexer.TokenKind) (Expression, error) {
	left, err := s.parsePrefix()
	if err != nil {
		return nil, err
	}

	for {
		next, err := s.peek()
		if err != nil {
			return nil, err
		}

		if precedence(parent) >= precedence(next.Kind) {
			return left, nil
		}

		left, err = s.parseInfixExpr(left)
		if err != nil {
			return nil, err
		}
	}
}

func (s *parserState) parsePrefix() (Expression, error) {
	token, err := s.eat()
	if err != nil {
		return nil, err
	}

	switch token.Kind {
	case lexer.LParen:
		return s.parseLParen(token)
	case lexer.Ident:
		return s.parseIdentifier(token)
	case lexer.String:
		return s.parseString(token)
	case lexer.Integer:
		return s.parseInteger(token)
	case lexer.Float:
		return s.parseFloat(token)
	case lexer.Return:
		return s.parseReturn(token)
	default:
		return nil, ErrExpectedExpression
	}
}

func (s *parserState) parseInfixExpr(left Expression) (Expression, error) {
	token, err := s.eat()
	if err != nil {
		return nil, err
	}

	switch token.Kind {
	case lexer.Equals:
		return s.parseEquals(left, token)
	case lexer.Colon:
		return s.parseEntry(left, token)
	case lexer.Comma:
		return s.parseEntries(left, token)
	case lexer.LParen:
		return s.parseBlockCall(left, token)
	default:
		return nil, ErrExpectedExpression
	}
}
