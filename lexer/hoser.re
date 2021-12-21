package lexer

import "bytes"
import "github.com/masp/hoser/token"

func (s *Scanner) lex() (pos token.Pos, tok token.Token, lit string, err error) {
	for {
		lit = ""
		pos = s.pos()
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
		end { tok = token.Eof; return }
		* { err = ErrBadToken; return }

		// Whitespace and new lines
		eol = ("\r\n" | "\n" | end);
		eol {
			if s.insertSemi() {
				s.cursor = s.token // Has the effect of "inserting" the semicolon in the input
				tok = token.Semicolon
				lit = "\n"
				return
			} else {
				s.file.AddLine(s.token)
				continue
			}
		}
		[ \t]+ {
			continue
		}

		// Comments
		"#" [^\r\n\x00]* { tok = token.Comment; lit = s.literal(); return }

		// Keywords
		"return" { tok = token.Return; lit = "return"; return }
		"module" { tok = token.Module; lit = "module"; return }
		"pipe" { tok = token.Pipe; lit = "pipe"; return }
		"stub" { tok = token.Stub; lit = "stub"; return }

		// Operators and punctuation
		"(" { tok = token.LParen; lit = "("; return }
		")" { tok = token.RParen; lit = ")"; return }
		"{" { tok = token.LCurlyBrack; lit = "{"; return }
		"}" { tok = token.RCurlyBrack; lit = "}"; return }
		"=" { tok = token.Equals; lit = "="; return }
		"." { tok = token.Period; lit = "."; return }
		"," { tok = token.Comma; lit = ","; return }
		":" { tok = token.Colon; lit = ":"; return }
		";" { tok = token.Semicolon; lit = ";"; return }

		// Integer literals
		dec = [1-9][0-9]*;
		dec { tok = token.Integer; lit = s.literal(); return }

		// Floating point numbers
		// from excellent https://re2c.org/examples/c/real_world/example_cxx98.html
		frc = [0-9]* "." [0-9]+ | [0-9]+ ".";
		exp = 'e' [+-]? [0-9]+;
		flt = (frc exp? | [0-9]+ exp);
		flt { tok = token.Float; lit = s.literal(); return }

		// Strings
		["] { return s.lexString('"') } 


		// Identifiers
		id = [a-zA-Z_][a-zA-Z_0-9]*;
		id { tok = token.Ident; lit = s.literal(); return }
*/		
	}
}

func (s *Scanner) lexString(quote byte) (pos token.Pos, tok token.Token, lit string, err error) {
	var buf bytes.Buffer
	for {
		var u byte
/*!re2c
		re2c:yyfill:enable = 0;
		re2c:flags:nested-ifs = 1;
		re2c:define:YYCTYPE = byte;
		re2c:define:YYPEEK = "s.text[s.cursor]";
		re2c:define:YYSKIP = "s.cursor += 1";

		*                    { err = ErrInvalidString; return }
		[^\n\\]              {
			u = yych
			if (u == quote) {
				tok = token.String
				pos = token.Pos(s.token)
				lit = string(buf.Bytes())
				return
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