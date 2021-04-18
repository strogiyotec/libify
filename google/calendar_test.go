package google

import (
	"fmt"
	"testing"
)

func TestReadToken(t *testing.T) {
	token := readToken("test_token.json")
	fmt.Println("HERE")
	if token.TokenType != "Bearer" {
		t.Errorf("Token type has to be Bearer")
	}
}
