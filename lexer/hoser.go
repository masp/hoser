// Code generated by re2c 2.2 on Mon Dec 27 21:38:22 2021, DO NOT EDIT.
package lexer

import "bytes"
import "github.com/masp/hoser/token"

func (s *Scanner) lex() (pos token.Pos, tok token.Token, lit string, err error) {
	for {
		lit = ""
		pos = s.pos()
		s.token = s.cursor

		{
			var yych byte
			yyaccept := 0
			yych = s.text[s.cursor]
			switch yych {
			case 0x00:
				goto yy2
			case '\t':
				fallthrough
			case ' ':
				goto yy6
			case '\n':
				goto yy9
			case '\r':
				goto yy11
			case '"':
				goto yy12
			case '#':
				goto yy14
			case '(':
				goto yy17
			case ')':
				goto yy19
			case ',':
				goto yy21
			case '.':
				goto yy23
			case '0':
				goto yy25
			case '1':
				fallthrough
			case '2':
				fallthrough
			case '3':
				fallthrough
			case '4':
				fallthrough
			case '5':
				fallthrough
			case '6':
				fallthrough
			case '7':
				fallthrough
			case '8':
				fallthrough
			case '9':
				goto yy26
			case ':':
				goto yy29
			case ';':
				goto yy31
			case '=':
				goto yy33
			case 'A':
				fallthrough
			case 'B':
				fallthrough
			case 'C':
				fallthrough
			case 'D':
				fallthrough
			case 'E':
				fallthrough
			case 'F':
				fallthrough
			case 'G':
				fallthrough
			case 'H':
				fallthrough
			case 'I':
				fallthrough
			case 'J':
				fallthrough
			case 'K':
				fallthrough
			case 'L':
				fallthrough
			case 'M':
				fallthrough
			case 'N':
				fallthrough
			case 'O':
				fallthrough
			case 'P':
				fallthrough
			case 'Q':
				fallthrough
			case 'R':
				fallthrough
			case 'S':
				fallthrough
			case 'T':
				fallthrough
			case 'U':
				fallthrough
			case 'V':
				fallthrough
			case 'W':
				fallthrough
			case 'X':
				fallthrough
			case 'Y':
				fallthrough
			case 'Z':
				fallthrough
			case '_':
				fallthrough
			case 'a':
				fallthrough
			case 'b':
				fallthrough
			case 'c':
				fallthrough
			case 'd':
				fallthrough
			case 'e':
				fallthrough
			case 'f':
				fallthrough
			case 'g':
				fallthrough
			case 'h':
				fallthrough
			case 'j':
				fallthrough
			case 'k':
				fallthrough
			case 'l':
				fallthrough
			case 'n':
				fallthrough
			case 'o':
				fallthrough
			case 'q':
				fallthrough
			case 't':
				fallthrough
			case 'u':
				fallthrough
			case 'v':
				fallthrough
			case 'w':
				fallthrough
			case 'x':
				fallthrough
			case 'y':
				fallthrough
			case 'z':
				goto yy35
			case 'i':
				goto yy38
			case 'm':
				goto yy39
			case 'p':
				goto yy40
			case 'r':
				goto yy41
			case 's':
				goto yy42
			case '{':
				goto yy43
			case '}':
				goto yy45
			default:
				goto yy4
			}
		yy2:
			s.cursor += 1
			{
				tok = token.Eof
				return
			}
		yy4:
			s.cursor += 1
		yy5:
			{
				err = ErrBadToken
				return
			}
		yy6:
			s.cursor += 1
			yych = s.text[s.cursor]
			if yych == '\t' {
				goto yy6
			}
			if yych == ' ' {
				goto yy6
			}
			{
				continue
			}
		yy9:
			s.cursor += 1
			{
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
		yy11:
			s.cursor += 1
			yych = s.text[s.cursor]
			if yych == '\n' {
				goto yy9
			}
			goto yy5
		yy12:
			s.cursor += 1
			{
				return s.lexString('"')
			}
		yy14:
			s.cursor += 1
			yych = s.text[s.cursor]
			if yych <= '\n' {
				if yych <= 0x00 {
					goto yy16
				}
				if yych <= '\t' {
					goto yy14
				}
			} else {
				if yych != '\r' {
					goto yy14
				}
			}
		yy16:
			{
				tok = token.Comment
				lit = s.literal()
				return
			}
		yy17:
			s.cursor += 1
			{
				tok = token.LParen
				lit = "("
				return
			}
		yy19:
			s.cursor += 1
			{
				tok = token.RParen
				lit = ")"
				return
			}
		yy21:
			s.cursor += 1
			{
				tok = token.Comma
				lit = ","
				return
			}
		yy23:
			s.cursor += 1
			yych = s.text[s.cursor]
			if yych <= '/' {
				goto yy24
			}
			if yych <= '9' {
				goto yy47
			}
		yy24:
			{
				tok = token.Period
				lit = "."
				return
			}
		yy25:
			yyaccept = 0
			s.cursor += 1
			s.marker = s.cursor
			yych = s.text[s.cursor]
			if yych <= '9' {
				if yych == '.' {
					goto yy47
				}
				if yych <= '/' {
					goto yy5
				}
				goto yy50
			} else {
				if yych <= 'E' {
					if yych <= 'D' {
						goto yy5
					}
					goto yy53
				} else {
					if yych == 'e' {
						goto yy53
					}
					goto yy5
				}
			}
		yy26:
			yyaccept = 1
			s.cursor += 1
			s.marker = s.cursor
			yych = s.text[s.cursor]
			if yych <= '9' {
				if yych == '.' {
					goto yy47
				}
				if yych >= '0' {
					goto yy26
				}
			} else {
				if yych <= 'E' {
					if yych >= 'E' {
						goto yy53
					}
				} else {
					if yych == 'e' {
						goto yy53
					}
				}
			}
		yy28:
			{
				tok = token.Integer
				lit = s.literal()
				return
			}
		yy29:
			s.cursor += 1
			{
				tok = token.Colon
				lit = ":"
				return
			}
		yy31:
			s.cursor += 1
			{
				tok = token.Semicolon
				lit = ";"
				return
			}
		yy33:
			s.cursor += 1
			{
				tok = token.Equals
				lit = "="
				return
			}
		yy35:
			s.cursor += 1
			yych = s.text[s.cursor]
		yy36:
			if yych <= 'Z' {
				if yych <= '/' {
					goto yy37
				}
				if yych <= '9' {
					goto yy35
				}
				if yych >= 'A' {
					goto yy35
				}
			} else {
				if yych <= '_' {
					if yych >= '_' {
						goto yy35
					}
				} else {
					if yych <= '`' {
						goto yy37
					}
					if yych <= 'z' {
						goto yy35
					}
				}
			}
		yy37:
			{
				tok = token.Ident
				lit = s.literal()
				return
			}
		yy38:
			s.cursor += 1
			yych = s.text[s.cursor]
			if yych == 'm' {
				goto yy54
			}
			goto yy36
		yy39:
			s.cursor += 1
			yych = s.text[s.cursor]
			if yych == 'o' {
				goto yy55
			}
			goto yy36
		yy40:
			s.cursor += 1
			yych = s.text[s.cursor]
			if yych == 'i' {
				goto yy56
			}
			goto yy36
		yy41:
			s.cursor += 1
			yych = s.text[s.cursor]
			if yych == 'e' {
				goto yy57
			}
			goto yy36
		yy42:
			s.cursor += 1
			yych = s.text[s.cursor]
			if yych == 't' {
				goto yy58
			}
			goto yy36
		yy43:
			s.cursor += 1
			{
				tok = token.LCurlyBrack
				lit = "{"
				return
			}
		yy45:
			s.cursor += 1
			{
				tok = token.RCurlyBrack
				lit = "}"
				return
			}
		yy47:
			yyaccept = 2
			s.cursor += 1
			s.marker = s.cursor
			yych = s.text[s.cursor]
			if yych <= 'D' {
				if yych <= '/' {
					goto yy49
				}
				if yych <= '9' {
					goto yy47
				}
			} else {
				if yych <= 'E' {
					goto yy53
				}
				if yych == 'e' {
					goto yy53
				}
			}
		yy49:
			{
				tok = token.Float
				lit = s.literal()
				return
			}
		yy50:
			s.cursor += 1
			yych = s.text[s.cursor]
			if yych <= '9' {
				if yych == '.' {
					goto yy47
				}
				if yych >= '0' {
					goto yy50
				}
			} else {
				if yych <= 'E' {
					if yych >= 'E' {
						goto yy53
					}
				} else {
					if yych == 'e' {
						goto yy53
					}
				}
			}
		yy52:
			s.cursor = s.marker
			if yyaccept <= 1 {
				if yyaccept == 0 {
					goto yy5
				} else {
					goto yy28
				}
			} else {
				goto yy49
			}
		yy53:
			s.cursor += 1
			yych = s.text[s.cursor]
			if yych <= ',' {
				if yych == '+' {
					goto yy59
				}
				goto yy52
			} else {
				if yych <= '-' {
					goto yy59
				}
				if yych <= '/' {
					goto yy52
				}
				if yych <= '9' {
					goto yy60
				}
				goto yy52
			}
		yy54:
			s.cursor += 1
			yych = s.text[s.cursor]
			if yych == 'p' {
				goto yy62
			}
			goto yy36
		yy55:
			s.cursor += 1
			yych = s.text[s.cursor]
			if yych == 'd' {
				goto yy63
			}
			goto yy36
		yy56:
			s.cursor += 1
			yych = s.text[s.cursor]
			if yych == 'p' {
				goto yy64
			}
			goto yy36
		yy57:
			s.cursor += 1
			yych = s.text[s.cursor]
			if yych == 't' {
				goto yy65
			}
			goto yy36
		yy58:
			s.cursor += 1
			yych = s.text[s.cursor]
			if yych == 'u' {
				goto yy66
			}
			goto yy36
		yy59:
			s.cursor += 1
			yych = s.text[s.cursor]
			if yych <= '/' {
				goto yy52
			}
			if yych >= ':' {
				goto yy52
			}
		yy60:
			s.cursor += 1
			yych = s.text[s.cursor]
			if yych <= '/' {
				goto yy49
			}
			if yych <= '9' {
				goto yy60
			}
			goto yy49
		yy62:
			s.cursor += 1
			yych = s.text[s.cursor]
			if yych == 'o' {
				goto yy67
			}
			goto yy36
		yy63:
			s.cursor += 1
			yych = s.text[s.cursor]
			if yych == 'u' {
				goto yy68
			}
			goto yy36
		yy64:
			s.cursor += 1
			yych = s.text[s.cursor]
			if yych == 'e' {
				goto yy69
			}
			goto yy36
		yy65:
			s.cursor += 1
			yych = s.text[s.cursor]
			if yych == 'u' {
				goto yy71
			}
			goto yy36
		yy66:
			s.cursor += 1
			yych = s.text[s.cursor]
			if yych == 'b' {
				goto yy72
			}
			goto yy36
		yy67:
			s.cursor += 1
			yych = s.text[s.cursor]
			if yych == 'r' {
				goto yy74
			}
			goto yy36
		yy68:
			s.cursor += 1
			yych = s.text[s.cursor]
			if yych == 'l' {
				goto yy75
			}
			goto yy36
		yy69:
			s.cursor += 1
			yych = s.text[s.cursor]
			if yych <= 'Z' {
				if yych <= '/' {
					goto yy70
				}
				if yych <= '9' {
					goto yy35
				}
				if yych >= 'A' {
					goto yy35
				}
			} else {
				if yych <= '_' {
					if yych >= '_' {
						goto yy35
					}
				} else {
					if yych <= '`' {
						goto yy70
					}
					if yych <= 'z' {
						goto yy35
					}
				}
			}
		yy70:
			{
				tok = token.Pipe
				lit = "pipe"
				return
			}
		yy71:
			s.cursor += 1
			yych = s.text[s.cursor]
			if yych == 'r' {
				goto yy76
			}
			goto yy36
		yy72:
			s.cursor += 1
			yych = s.text[s.cursor]
			if yych <= 'Z' {
				if yych <= '/' {
					goto yy73
				}
				if yych <= '9' {
					goto yy35
				}
				if yych >= 'A' {
					goto yy35
				}
			} else {
				if yych <= '_' {
					if yych >= '_' {
						goto yy35
					}
				} else {
					if yych <= '`' {
						goto yy73
					}
					if yych <= 'z' {
						goto yy35
					}
				}
			}
		yy73:
			{
				tok = token.Stub
				lit = "stub"
				return
			}
		yy74:
			s.cursor += 1
			yych = s.text[s.cursor]
			if yych == 't' {
				goto yy77
			}
			goto yy36
		yy75:
			s.cursor += 1
			yych = s.text[s.cursor]
			if yych == 'e' {
				goto yy79
			}
			goto yy36
		yy76:
			s.cursor += 1
			yych = s.text[s.cursor]
			if yych == 'n' {
				goto yy81
			}
			goto yy36
		yy77:
			s.cursor += 1
			yych = s.text[s.cursor]
			if yych <= 'Z' {
				if yych <= '/' {
					goto yy78
				}
				if yych <= '9' {
					goto yy35
				}
				if yych >= 'A' {
					goto yy35
				}
			} else {
				if yych <= '_' {
					if yych >= '_' {
						goto yy35
					}
				} else {
					if yych <= '`' {
						goto yy78
					}
					if yych <= 'z' {
						goto yy35
					}
				}
			}
		yy78:
			{
				tok = token.Import
				lit = "import"
				return
			}
		yy79:
			s.cursor += 1
			yych = s.text[s.cursor]
			if yych <= 'Z' {
				if yych <= '/' {
					goto yy80
				}
				if yych <= '9' {
					goto yy35
				}
				if yych >= 'A' {
					goto yy35
				}
			} else {
				if yych <= '_' {
					if yych >= '_' {
						goto yy35
					}
				} else {
					if yych <= '`' {
						goto yy80
					}
					if yych <= 'z' {
						goto yy35
					}
				}
			}
		yy80:
			{
				tok = token.Module
				lit = "module"
				return
			}
		yy81:
			s.cursor += 1
			yych = s.text[s.cursor]
			if yych <= 'Z' {
				if yych <= '/' {
					goto yy82
				}
				if yych <= '9' {
					goto yy35
				}
				if yych >= 'A' {
					goto yy35
				}
			} else {
				if yych <= '_' {
					if yych >= '_' {
						goto yy35
					}
				} else {
					if yych <= '`' {
						goto yy82
					}
					if yych <= 'z' {
						goto yy35
					}
				}
			}
		yy82:
			{
				tok = token.Return
				lit = "return"
				return
			}
		}

	}
}

func (s *Scanner) lexString(quote byte) (pos token.Pos, tok token.Token, lit string, err error) {
	var buf bytes.Buffer
	for {
		var u byte

		{
			var yych byte
			yych = s.text[s.cursor]
			if yych == '\n' {
				goto yy87
			}
			if yych == '\\' {
				goto yy89
			}
			s.cursor += 1
			{
				u = yych
				if u == quote {
					tok = token.String
					pos = token.Pos(s.token)
					lit = string(buf.Bytes())
					return
				}
				buf.WriteByte(u)
				continue
			}
		yy87:
			s.cursor += 1
		yy88:
			{
				err = ErrInvalidString
				return
			}
		yy89:
			s.cursor += 1
			yych = s.text[s.cursor]
			if yych <= 'b' {
				if yych <= '>' {
					if yych <= '"' {
						if yych <= '!' {
							goto yy88
						}
					} else {
						if yych == '\'' {
							goto yy92
						}
						goto yy88
					}
				} else {
					if yych <= '\\' {
						if yych <= '?' {
							goto yy94
						}
						if yych <= '[' {
							goto yy88
						}
						goto yy96
					} else {
						if yych <= '`' {
							goto yy88
						}
						if yych <= 'a' {
							goto yy98
						}
						goto yy100
					}
				}
			} else {
				if yych <= 'q' {
					if yych <= 'f' {
						if yych <= 'e' {
							goto yy88
						}
						goto yy102
					} else {
						if yych == 'n' {
							goto yy104
						}
						goto yy88
					}
				} else {
					if yych <= 't' {
						if yych <= 'r' {
							goto yy106
						}
						if yych <= 's' {
							goto yy88
						}
						goto yy108
					} else {
						if yych == 'v' {
							goto yy110
						}
						goto yy88
					}
				}
			}
			s.cursor += 1
			{
				buf.WriteByte('"')
				continue
			}
		yy92:
			s.cursor += 1
			{
				buf.WriteByte('\'')
				continue
			}
		yy94:
			s.cursor += 1
			{
				buf.WriteByte('?')
				continue
			}
		yy96:
			s.cursor += 1
			{
				buf.WriteByte('\\')
				continue
			}
		yy98:
			s.cursor += 1
			{
				buf.WriteByte('\a')
				continue
			}
		yy100:
			s.cursor += 1
			{
				buf.WriteByte('\b')
				continue
			}
		yy102:
			s.cursor += 1
			{
				buf.WriteByte('\f')
				continue
			}
		yy104:
			s.cursor += 1
			{
				buf.WriteByte('\n')
				continue
			}
		yy106:
			s.cursor += 1
			{
				buf.WriteByte('\r')
				continue
			}
		yy108:
			s.cursor += 1
			{
				buf.WriteByte('\t')
				continue
			}
		yy110:
			s.cursor += 1
			{
				buf.WriteByte('\v')
				continue
			}
		}

	}
}
