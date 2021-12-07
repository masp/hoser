package parser

import (
	"strconv"

	"github.com/masp/hoser/ast"
	"github.com/masp/hoser/lexer"
)

func (s *parserState) parseIdentifier(token lexer.Token) (*ast.Identifier, error) {
	return &ast.Identifier{Token: token}, nil
}

func (s *parserState) parseString(token lexer.Token) (*ast.String, error) {
	return &ast.String{Token: token}, nil
}

func (s *parserState) parseInteger(token lexer.Token) (*ast.Integer, error) {
	v, err := strconv.ParseInt(token.Value, 0, 64)
	if err != nil {
		return nil, err
	}
	return &ast.Integer{Token: token, Value: v}, nil
}

func (s *parserState) parseFloat(token lexer.Token) (*ast.Float, error) {
	v, err := strconv.ParseFloat(token.Value, 64)
	if err != nil {
		return nil, err
	}
	return &ast.Float{Token: token, Value: v}, nil
}
