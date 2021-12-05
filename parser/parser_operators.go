package parser

import "github.com/masp/hoser/lexer"

func (s *parserState) parseEquals(left Expression, token lexer.Token) (Expression, error) {
	right, err := s.parseExpression(token.Kind)
	if err != nil {
		return nil, err
	}
	return &AssignmentExpr{Left: left, Right: right}, nil
}

func (s *parserState) parseLParen(token lexer.Token) (Expression, error) {
	expr, err := s.parseExpression(lexer.Invalid)
	if err != nil {
		return nil, err
	}

	if _, err := s.eatOnly(lexer.RParen); err != nil {
		return nil, err
	}
	return expr, nil
}
