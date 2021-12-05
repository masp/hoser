package lexer

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type TokenKind int

type Token struct {
	Kind      TokenKind
	Value     string
	Line, Col int
}

func (t Token) String() string {
	return fmt.Sprintf("[%v %v %v:%v]", t.Kind, t.Value, t.Line, t.Col)
}

const (
	Invalid TokenKind = iota

	// Keywords
	Return

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

	Eof TokenKind = 999 // should always be at end
)

func (tk TokenKind) String() string {
	switch tk {
	case Invalid:
		return "Invalid"
	case Return:
		return "Return"
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

func (tk *TokenKind) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	switch s {
	case "Invalid":
		*tk = Invalid
	case "Return":
		*tk = Return
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

func (tk TokenKind) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(tk.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}
