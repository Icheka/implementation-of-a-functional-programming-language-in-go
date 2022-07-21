package ast

import (
	"monkey/token"
	"strings"
	"testing"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{
					Type: token.LET,
					Value: strings.ToLower(token.LET),
				},
				Name: &Identifier{
					Token: token.Token{
						Type: token.IDENT,
						Value: "name",
					},
					Value: "name",
				},
				Value: &Identifier{
					Token: token.Token{
						Type: token.IDENT,
						Value: "fullName",
					},
					Value: "fullName",
				},
			},
		},
	}

	test := "let name = fullName;"
	if program.String() != test {
		t.Fatalf("program.String() wrong. Expected %q, got %q", test, program.String())
	}
}