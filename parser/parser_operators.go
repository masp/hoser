package parser

import (
	"github.com/masp/hoser/ast"
	"github.com/masp/hoser/lexer"
)

func (s *parserState) parseEquals(left ast.Expression, token lexer.Token) (ast.Expression, error) {
	right, err := s.parseExpression(token.Kind)
	if err != nil {
		return nil, err
	}
	return &ast.AssignmentExpr{Left: left, Right: right}, nil
}

func (s *parserState) parseLParen(token lexer.Token) (ast.Expression, error) {
	expr, err := s.parseExpression(lexer.Invalid)
	if err != nil {
		return nil, err
	}

	if _, err := s.eatOnly(lexer.RParen); err != nil {
		return nil, err
	}
	return expr, nil
}
