package lexer

import (
	"monkey/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	// tokens to test against
	input := `let five = 5;
	let ten = 10;
	
	let add = function(x, y) {
		x + y;
	};
	
	let result = add(five, ten);
	!-/*5;
	5 < 10 > 5;
	true false if else return;
	10 == 10;
	10 != 9;
	"foorbar"
	"icheka ozuru"
	[1, 2]`

	tests := []struct {
		expectedType  token.TokenType
		expectedValue string
	}{
		{token.LET, "let"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "function"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.ASTERISK, "*"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.INT, "5"},
		{token.L_THAN, "<"},
		{token.INT, "10"},
		{token.G_THAN, ">"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.TRUE, "true"},
		{token.FALSE, "false"},
		{token.IF, "if"},
		{token.ELSE, "else"},
		{token.RETURN, "return"},
		{token.SEMICOLON, ";"},
		{token.INT, "10"},
		{token.EQUAL, "=="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.INT, "10"},
		{token.NOT_EQUAL, "!="},
		{token.INT, "9"},
		{token.SEMICOLON, ";"},
		{token.STRING, "foobar"},
		{token.STRING, "icheka ozuru"},
		{token.LBRACKET, "["},
		{token.INT, "1"},
		{token.COMMA, ","},
		{token.INT, "2"},
		{token.RBRACKET, "]"},
	}

	// a new lexer instance
	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. Expected=%q, %q. Got=%q, %q", i, tt.expectedType, tt.expectedValue, tok.Type, tok.Value)
		}
	}
}
