package rsql

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLexer_Parse_Functions(t *testing.T) {
	lexer := NewLexer("1aA A <>=!")
	assert.Equal(t, lexer.isBlank(3), true)
	assert.Equal(t, lexer.isBlank(4), false)
	assert.Equal(t, lexer.isBlankBefore(0), true)
	assert.Equal(t, lexer.isBlankBefore(1), false)
	assert.Equal(t, lexer.isBlankBefore(4), true)
	assert.Equal(t, lexer.isBlankAfter(0), false)
	assert.Equal(t, lexer.isBlankAfter(2), true)
	assert.Equal(t, lexer.isBlankAfter(lexer.buflen-1), true)
	assert.Equal(t, lexer.isDigit(0), true, "digit 0")
	assert.Equal(t, lexer.isDigit(1), false, "digit 1")
	assert.Equal(t, lexer.isAlpha(1), true, "alpha 1")
	assert.Equal(t, lexer.isAlpha(0), false, "alpha 0")
	assert.Equal(t, lexer.isString(0, "1aA"), true, "string 0 -> 1")
	assert.Equal(t, lexer.isString(0, "1aA "), true, "string 0 -> 2")
	assert.Equal(t, lexer.isString(lexer.buflen-2, ">=!"), false, "string -2")
	assert.Equal(t, lexer.isString(lexer.buflen-3, ">=!"), true, "string -3")
}

func TestLexer_Parse_Operator(t *testing.T) {
	lex := NewLexer("and or")
	assert.Equal(t, lex.nextToken().Type, AndToken)
	assert.Equal(t, lex.nextToken().Type, OrToken)
	assert.Equal(t, lex.nextToken().Type, EOFToken)
}

func TestLexer_Parse_Boolean(t *testing.T) {
	lex := NewLexer("true false")
	assert.Equal(t, lex.nextToken().Type, BooleanToken)
	assert.Equal(t, lex.nextToken().Type, BooleanToken)
	assert.Equal(t, lex.nextToken().Type, EOFToken)
}

func TestLexer_Parse_Numeric(t *testing.T) {
	lex := NewLexer("123 123. 123.2")
	assert.Equal(t, lex.nextToken().Type, IntegerToken)
	assert.Equal(t, lex.nextToken().Type, DoubleToken)
	assert.Equal(t, lex.nextToken().Type, DoubleToken)
	assert.Equal(t, lex.nextToken().Type, EOFToken)
}

func TestLexer_Parse_Date(t *testing.T) {
	lex := NewLexer("1985-02-25 1985-02-25T22:35:00 1985-02-25T22:35:00+0100 1985-02-25T22:35:00+01:00 1985-02-25T22:35:00Z")
	assert.Equal(t, lex.nextToken().Type, DateToken)
	assert.Equal(t, lex.nextToken().Type, DateTimeToken)
	assert.Equal(t, lex.nextToken().Type, DateTimeToken)
	assert.Equal(t, lex.nextToken().Type, DateTimeToken)
	assert.Equal(t, lex.nextToken().Type, DateTimeToken)
	assert.Equal(t, lex.nextToken().Type, EOFToken)
}

func TestLexer_Parse_Identifier(t *testing.T) {
	lex := NewLexer("a1 a1.a2 a1.a2.a3")
	assert.Equal(t, lex.nextToken().Type, IdentifierToken)
	assert.Equal(t, lex.nextToken().Type, IdentifierToken)
	assert.Equal(t, lex.nextToken().Type, IdentifierToken)
	assert.Equal(t, lex.nextToken().Type, EOFToken)
}

func TestLexer_Parse_String(t *testing.T) {
	lex := NewLexer(`"123" "123'12" "123\"12" '123' '123"12' '123\'12'`)
	token := lex.nextToken()
	assert.Equal(t, token.Type, StringToken)
	assert.Equal(t, token.Value, "123")
	token = lex.nextToken()
	assert.Equal(t, token.Type, StringToken)
	assert.Equal(t, token.Value, "123'12")
	token = lex.nextToken()
	assert.Equal(t, token.Type, StringToken)
	assert.Equal(t, token.Value, "123\"12")
	token = lex.nextToken()
	assert.Equal(t, token.Type, StringToken)
	assert.Equal(t, token.Value, "123")
	token = lex.nextToken()
	assert.Equal(t, token.Type, StringToken)
	assert.Equal(t, token.Value, "123\"12")
	token = lex.nextToken()
	assert.Equal(t, token.Type, StringToken)
	assert.Equal(t, token.Value, "123'12")
	assert.Equal(t, lex.nextToken().Type, EOFToken)
}

func TestLexer_Parse_Reserved(t *testing.T) {
	lex := NewLexer(`== != !=~ ==~ =in= =out= > >= < <= , ( )`)
	assert.Equal(t, lex.nextToken().Type, EqualsToken)
	assert.Equal(t, lex.nextToken().Type, NotEqualsToken)
	assert.Equal(t, lex.nextToken().Type, NotLikeToken)
	assert.Equal(t, lex.nextToken().Type, LikeToken)
	assert.Equal(t, lex.nextToken().Type, InToken)
	assert.Equal(t, lex.nextToken().Type, NotInToken)
	assert.Equal(t, lex.nextToken().Type, GreaterToken)
	assert.Equal(t, lex.nextToken().Type, GreaterOrEqualsToken)
	assert.Equal(t, lex.nextToken().Type, LessToken)
	assert.Equal(t, lex.nextToken().Type, LessOrEqualsToken)
	assert.Equal(t, lex.nextToken().Type, CommaToken)
	assert.Equal(t, lex.nextToken().Type, LeftParenToken)
	assert.Equal(t, lex.nextToken().Type, RightParenToken)

	assert.Equal(t, lex.nextToken().Type, EOFToken)
}
