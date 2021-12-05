package parser

import (
	"strconv"

	"github.com/masp/hoser/lexer"
)

func (s *parserState) parseIdentifier(token lexer.Token) (*Identifier, error) {
	return &Identifier{Token: token}, nil
}

func (s *parserState) parseString(token lexer.Token) (*String, error) {
	return &String{Token: token}, nil
}

func (s *parserState) parseInteger(token lexer.Token) (*Integer, error) {
	v, err := strconv.ParseInt(token.Value, 0, 64)
	if err != nil {
		return nil, err
	}
	return &Integer{Token: token, Value: v}, nil
}

func (s *parserState) parseFloat(token lexer.Token) (*Float, error) {
	v, err := strconv.ParseFloat(token.Value, 64)
	if err != nil {
		return nil, err
	}
	return &Float{Token: token, Value: v}, nil
}
