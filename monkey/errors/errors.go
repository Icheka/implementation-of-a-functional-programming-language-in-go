package errors

import (
	"fmt"
	"monkey/token"
)

func ExpectedNextTokenToBe(expected token.TokenType, got token.TokenType) string {
	return fmt.Sprintf("Expected next token to be %s got %s", expected, string(got))
}

func CouldNotParseInteger(value string) string {
	return fmt.Sprintf("Could not parse %s as integer", value)
}

func NoPrefixParseError(t token.TokenType) string {
	return fmt.Sprintf("No prefix parse function for %s found", t)
}
