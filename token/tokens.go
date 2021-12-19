package token

import "strconv"

type Token int

const (
	Invalid Token = iota

	Comment

	// Keywords
	Return
	Module

	// Literals
	literal_begin
	Ident
	String
	Integer
	Float
	literal_end

	// Operators
	Equals

	// Other
	Period
	Comma
	Colon
	Semicolon
	LCurlyBrack
	RCurlyBrack
	LParen
	RParen

	Eof Token = 999 // should always be at end
)

var tokens = [...]string{
	Invalid: "INVALID",

	Comment: "COMMENT",

	// Keywords
	Return: "return",
	Module: "module",

	// Literals
	Ident:   "IDENT",
	String:  "STRING",
	Integer: "INT",
	Float:   "FLOAT",

	// Operators
	Equals: "=",

	// Other
	Period:      ".",
	Comma:       ",",
	Colon:       ":",
	Semicolon:   ";",
	LCurlyBrack: "{",
	RCurlyBrack: "}",
	LParen:      "(",
	RParen:      ")",

	Eof: "EOF",
}

func (tok Token) String() string {
	s := ""
	if 0 <= tok && tok < Token(len(tokens)) {
		s = tokens[tok]
	}
	if s == "" {
		s = "token(" + strconv.Itoa(int(tok)) + ")"
	}
	return s
}

func (tok Token) IsLiteral() bool {
	return literal_begin < tok && tok < literal_end
}
