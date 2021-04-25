package google

import (
	"testing"
)

func TestReadToken(t *testing.T) {
	token := readToken("test_token.json")
	if token.TokenType != "Bearer" {
		t.Errorf("Token type has to be Bearer, but was %s", token.TokenType)
	}
	if token.AccessToken != "token" {
		t.Errorf("Access token has to be token , but was %s", token.AccessToken)
	}
}
