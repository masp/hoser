package token

import (
	"bytes"
)

type Token int

const (
	Invalid Token = iota

	// Keywords
	Return
	IntType
	StringType

	// Literals
	Ident
	String
	Integer
	Float

	// Operators
	Equals

	// Other
	Comma
	Colon
	Semicolon
	LCurlyBrack
	RCurlyBrack
	LParen
	RParen

	Eof Token = 999 // should always be at end
)

func (tk Token) String() string {
	switch tk {
	case Invalid:
		return "Invalid"
	case Return:
		return "Return"
	case IntType:
		return "IntType"
	case StringType:
		return "StringType"
	case Ident:
		return "Ident"
	case String:
		return "String"
	case Integer:
		return "Integer"
	case Float:
		return "Float"
	case LParen:
		return "("
	case RParen:
		return ")"
	case LCurlyBrack:
		return "{"
	case RCurlyBrack:
		return "}"
	case Equals:
		return "="
	case Semicolon:
		return ";"
	case Colon:
		return ":"
	case Comma:
		return ","
	default:
		return "Invalid"
	}
}

func (tk *Token) FromString(s string) error {
	switch s {
	case "Invalid":
		*tk = Invalid
	case "Return":
		*tk = Return
	case "IntType":
		*tk = IntType
	case "StringType":
		*tk = StringType
	case "Ident":
		*tk = Ident
	case "String":
		*tk = String
	case "Integer":
		*tk = Integer
	case "Float":
		*tk = Float
	case "(":
		*tk = LParen
	case ")":
		*tk = RParen
	case "{":
		*tk = LCurlyBrack
	case "}":
		*tk = RCurlyBrack
	case "=":
		*tk = Equals
	case ";":
		*tk = Semicolon
	case ":":
		*tk = Colon
	case ",":
		*tk = Comma
	default:
		*tk = Invalid
	}
	return nil
}

func (tk Token) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(tk.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}
