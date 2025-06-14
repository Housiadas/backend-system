package userapp

import "encoding/json"

// Token represents the user token when requested.
type Token struct {
	Token string `json:"token"`
}

// Encode implements the encoder interface.
func (t Token) Encode() ([]byte, string, error) {
	data, err := json.Marshal(t)
	return data, "application/json", err
}

func toToken(v string) Token {
	return Token{
		Token: v,
	}
}
