package lexer

import "bytes"

func (s *lexerState) lex() (Token, error) {
	for {
		s.token = s.cursor
/*!re2c
		re2c:yyfill:enable = 0;
		re2c:flags:nested-ifs = 1;
		re2c:define:YYCTYPE = byte;
		re2c:define:YYPEEK = "s.text[s.cursor]";
		re2c:define:YYSKIP = "s.cursor += 1";
		re2c:define:YYBACKUP = "s.marker = s.cursor";
		re2c:define:YYRESTORE = "s.cursor = s.marker";

		end = [\x00];
		end { return s.createToken(Eof), nil }
		* { return Token{}, ErrBadToken }

		// Whitespace and new lines
		("\r\n" | "\n") { 
			if s.lexEol() == Semicolon {
				s.cursor = s.token // Has the effect of "inserting" the semicolon in the input
				return s.createToken(Semicolon), nil
			} else {
				s.line += 1
				s.lineStart = s.cursor
				continue
			}
		}
		[ \t]+ {
			continue
		}

		// Keywords
		"return" { return s.createToken(Return), nil }

		// Operators and punctuation
		"(" { return s.createToken(LParen), nil }
		")" { return s.createToken(RParen), nil }
		"{" { return s.createToken(LCurlyBrack), nil }
		"}" { return s.createToken(RCurlyBrack), nil }
		"=" { return s.createToken(Equals), nil }
		"," { return s.createToken(Comma), nil }
		":" { return s.createToken(Colon), nil }
		";" { return s.createToken(Semicolon), nil }

		// Integer literals
		dec = [1-9][0-9]*;
		dec { return s.createToken(Integer), nil }

		// Floating point numbers
		// from excellent https://re2c.org/examples/c/real_world/example_cxx98.html
		frc = [0-9]* "." [0-9]+ | [0-9]+ ".";
		exp = 'e' [+-]? [0-9]+;
		flt = (frc exp? | [0-9]+ exp);
		flt { return s.createToken(Float), nil }

		// Strings
		["] { return s.lexString('"') } 


		// Identifiers
		id = [a-zA-Z_][a-zA-Z_0-9]*;
		id { return s.createToken(Ident), nil }
*/		
	}
}

func (s *lexerState) lexString(quote byte) (Token, error) {
	var buf bytes.Buffer
	for {
		var u byte
/*!re2c
		re2c:yyfill:enable = 0;
		re2c:flags:nested-ifs = 1;
		re2c:define:YYCTYPE = byte;
		re2c:define:YYPEEK = "s.text[s.cursor]";
		re2c:define:YYSKIP = "s.cursor += 1";

		*                    { return Token{}, ErrInvalidString }
		[^\n\\]              {
			u = yych
			if (u == quote) {
				tok := s.createToken(String)
				tok.Value = string(buf.Bytes())
				return tok, nil
			}
			buf.WriteByte(u)
			continue
		}
		"\\a"                { buf.WriteByte('\a'); continue }
		"\\b"                { buf.WriteByte('\b'); continue }
		"\\f"                { buf.WriteByte('\f'); continue }
		"\\n"                { buf.WriteByte('\n'); continue }
		"\\r"                { buf.WriteByte('\r'); continue }
		"\\t"                { buf.WriteByte('\t'); continue }
		"\\v"                { buf.WriteByte('\v'); continue }
		"\\\\"               { buf.WriteByte('\\'); continue }
		"\\'"                { buf.WriteByte('\''); continue }
		"\\\""               { buf.WriteByte('"'); continue }
		"\\?"                { buf.WriteByte('?'); continue }
*/		
	}
}