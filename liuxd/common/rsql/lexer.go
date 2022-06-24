package rsql

import (
	"fmt"
	"strings"
)

type TokenType string

const (
	UnknownToken         = TokenType("UnknownToken")
	IntegerToken         = TokenType("IntegerToken")
	DoubleToken          = TokenType("DoubleToken")
	DateToken            = TokenType("DateToken")
	DateTimeToken        = TokenType("DateTimeToken")
	BooleanToken         = TokenType("BooleanToken")
	IdentifierToken      = TokenType("IdentifierToken")
	StringToken          = TokenType("StringToken")
	LeftParenToken       = TokenType("LeftParenToken")
	RightParenToken      = TokenType("RightParenToken")
	OrToken              = TokenType("OrToken")
	AndToken             = TokenType("AndToken")
	EqualsToken          = TokenType("EqualsToken")
	NotEqualsToken       = TokenType("NotEqualsToken")
	LikeToken            = TokenType("LikeToken")
	NotLikeToken         = TokenType("NotLikeToken")
	GreaterToken         = TokenType("GreaterToken")
	GreaterOrEqualsToken = TokenType("GreaterOrEqualsToken")
	LessToken            = TokenType("LessToken")
	LessOrEqualsToken    = TokenType("LessOrEqualsToken")
	InToken              = TokenType("InToken")
	NotInToken           = TokenType("NotInToken")
	CommaToken           = TokenType("CommaToken")
	EOFToken             = TokenType("EOFToken")
)

type Token struct {
	Type  TokenType
	Value string
	Pos   int
	info  string
}

type Lexer struct {
	pos    int
	buf    string
	buflen int
}

func unknownToken(messages ...string) Token {
	return Token{
		Type:  UnknownToken,
		Value: "",
		Pos:   0,
		info:  strings.Join(messages, " "),
	}
}

func NewLexer(input string) *Lexer {
	return &Lexer{
		buf:    input,
		pos:    0,
		buflen: len(input),
	}
}

func (t *Lexer) Parse() ([]Token, error) {
	var res []Token
	var token Token
	var err error
	for {
		token = t.nextToken()
		if err != nil {
			return nil, fmt.Errorf("not a valid token at %d", t.pos)
		}
		if token.Type == UnknownToken {
			return nil, fmt.Errorf("not a valid token at %d (%s)", t.pos, token.info)
		}
		res = append(res, token)
		if token.Type == EOFToken {
			break
		}
	}

	return res, nil
}

func (t *Lexer) nextToken() (token Token) {
	defer func() {
		if r := recover(); r != nil {
			token = unknownToken(fmt.Sprintf("%+v", r))
		}
	}()
	t.skipBlank()
	if t.pos >= t.buflen {
		return Token{
			Type:  EOFToken,
			Value: "",
			Pos:   t.buflen,
		}
	}

	if token = t.processBool(); token.Type != UnknownToken {
		return token
	} else if token = t.processOperator(); token.Type != UnknownToken {
		return token
	} else if token = t.processDate(); token.Type != UnknownToken {
		return token
	} else if token = t.processNumber(); token.Type != UnknownToken {
		return token
	} else if token = t.processString(); token.Type != UnknownToken {
		return token
	} else if token = t.processIdentifier(); token.Type != UnknownToken {
		return token
	} else if token = t.processReserved(); token.Type != UnknownToken {
		return token
	} else {
		return unknownToken()
	}
}

func (t *Lexer) processBool() Token {
	if t.isString(t.pos, "true") {
		return t.generateToken(BooleanToken, t.pos+4)
	} else if t.isString(t.pos, "false") {
		return t.generateToken(BooleanToken, t.pos+5)
	}
	return unknownToken()
}

func (t *Lexer) processOperator() Token {
	if t.isBlankBefore(t.pos) && (t.isString(t.pos, "and") || t.isString(t.pos, "AND")) && t.isBlankAfter(t.pos+2) {
		return t.generateToken(AndToken, t.pos+3)
	} else if t.isBlankBefore(t.pos) && (t.isString(t.pos, "or") || t.isString(t.pos, "OR")) && t.isBlankAfter(t.pos+1) {
		return t.generateToken(OrToken, t.pos+2)
	}
	return unknownToken()

}

func (t *Lexer) processNumber() Token {
	idx := t.pos
	if t.isDigit(idx) {
		typ := IntegerToken
		idx++
		process := func() {
			for t.isDigit(idx) {
				idx++
			}
		}
		process()
		if t.charAt(idx) == '.' {
			idx++
			typ = DoubleToken
			process()
		}
		return t.generateToken(typ, idx)
	}
	return unknownToken()
}

func (t *Lexer) processIdentifier() Token {
	idx := t.pos
	if t.isAlpha(idx) {
		idx++
		process := func() {
			for t.isDigit(idx) || t.isAlpha(idx) {
				idx++
			}
		}
		process()
		for t.charAt(idx) == '.' && t.isAlpha(idx+1) {
			idx += 2
			process()
		}
		return t.generateToken(IdentifierToken, idx)
	}
	return unknownToken()
}

func (t *Lexer) processDate() Token {
	idx := t.pos
	if t.isDigit(idx) && t.isDigit(idx+1) && t.isDigit(idx+2) && t.isDigit(idx+3) &&
		t.isString(idx+4, "-") &&
		t.isDigit(idx+5) && t.isDigit(idx+6) &&
		t.isString(idx+7, "-") &&
		t.isDigit(idx+8) && t.isDigit(idx+9) {
		// TODO Validate month and day ?
		typ := DateToken
		idx += 10
		if t.charAt(idx) == 'T' &&
			t.isDigit(idx+1) && t.isDigit(idx+2) && t.isString(idx+3, ":") &&
			t.isDigit(idx+4) && t.isDigit(idx+5) && t.isString(idx+6, ":") &&
			t.isDigit(idx+7) && t.isDigit(idx+8) {
			idx += 9
			typ = DateTimeToken
			if t.charAt(idx) == 'Z' {
				idx++
			} else if t.isString(idx, "+") || t.isString(idx, "-") {
				if t.isDigit(idx+1) && t.isDigit(idx+2) {
					if t.isString(idx+3, ":") && t.isDigit(idx+4) && t.isDigit(idx+5) {
						idx += 6
					} else if t.isDigit(idx+3) && t.isDigit(idx+4) {
						idx += 5
					}
				}
			}
		}
		return t.generateToken(typ, idx)
	}
	return unknownToken()
}

func (t *Lexer) processString() Token {
	if t.charAt(t.pos) == '\'' || t.charAt(t.pos) == '"' {
		quote := t.charAt(t.pos)
		idx := strings.IndexByte(t.buf[t.pos+1:], quote) + t.pos + 1
		for idx != -1 && t.charAt(idx-1) == '\\' {
			idx = strings.IndexByte(t.buf[idx+1:], quote) + idx + 1
		}
		if idx == -1 {
			panic(fmt.Errorf("unterminated quote %d, %d", t.pos, idx))
			// t.error('Unterminated quote', t.pos, idx)
		}
		token := Token{
			Type:  StringToken,
			Value: strings.Join(strings.Split(t.buf[t.pos+1:idx], "\\"+string(quote)), string(quote)),
			Pos:   t.pos,
		}
		t.pos = idx + 1
		return token
	}
	return unknownToken()
}

func (t *Lexer) processReserved() Token {
	idx := t.pos
	if t.isString(idx, "(") {
		return t.generateToken(LeftParenToken, idx+1)
	} else if t.isString(idx, ")") {
		return t.generateToken(RightParenToken, idx+1)
	} else if t.isString(idx, ",") {
		return t.generateToken(CommaToken, idx+1)
	} else if t.isString(idx, "!=~") {
		return t.generateToken(NotLikeToken, idx+3)
	} else if t.isString(idx, "!=") {
		return t.generateToken(NotEqualsToken, idx+2)
	} else if t.isString(idx, "==~") {
		return t.generateToken(LikeToken, idx+3)
	} else if t.isString(idx, "==") {
		return t.generateToken(EqualsToken, idx+2)
	} else if t.isString(idx, ">=") {
		return t.generateToken(GreaterOrEqualsToken, idx+2)
	} else if t.isString(idx, ">") {
		return t.generateToken(GreaterToken, idx+1)
	} else if t.isString(idx, "<=") {
		return t.generateToken(LessOrEqualsToken, idx+2)
	} else if t.isString(idx, "<") {
		return t.generateToken(LessToken, idx+1)
	} else if t.isString(idx, "=in=") {
		return t.generateToken(InToken, idx+4)
	} else if t.isString(idx, "=out=") {
		return t.generateToken(NotInToken, idx+5)
	}
	return unknownToken()
}

func (t *Lexer) skipBlank() {
	for t.pos < t.buflen && t.isBlank(t.pos) {
		t.pos++
	}
}

func (t *Lexer) charAt(pos int) uint8 {
	if pos >= t.buflen {
		return 0
	}
	return t.buf[pos]
}

func (t *Lexer) generateToken(ty TokenType, newPos int) Token {
	res := Token{
		Type:  ty,
		Value: t.buf[t.pos:newPos],
		Pos:   t.pos,
	}
	t.pos = newPos
	return res
}

func (t *Lexer) isBlank(pos int) bool {
	c := t.charAt(pos)
	return c == ' ' || c == '\n' || c == '\t' || c == '\r'
}

func (t *Lexer) isBlankBefore(pos int) bool {
	return pos == 0 || t.isBlank(pos-1)
}

func (t *Lexer) isBlankAfter(pos int) bool {
	return pos == t.buflen-1 || pos < t.buflen-1 && t.isBlank(pos+1)
}

func (t *Lexer) isString(idx int, value string) bool {
	endIdx := idx + len(value)
	if endIdx > t.buflen {
		return false
	}
	return t.buf[idx:endIdx] == value
}

func (t *Lexer) isAlpha(idx int) bool {
	if idx >= t.buflen {
		return false
	}
	c := t.charAt(idx)
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c == '_' || c == '$'
}

func (t *Lexer) isDigit(idx int) bool {
	if idx >= t.buflen {
		return false
	}
	c := t.charAt(idx)
	return c >= '0' && c <= '9'
}
