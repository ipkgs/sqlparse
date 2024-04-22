package sqlparse

import (
	"strings"
)

type jsonToken struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type Encoder interface {
	Encode(v any) error
}

// Encode writes the tokens to the encoder in JSON format.
//
// One of the many ways to use can be as simple as `sqlparse.Encode(json.NewEncoder(os.Stdout), tokens)`.
func Encode(w Encoder, tokens []Token) error {
	var tokenList []jsonToken
	for _, t := range tokens {
		tokenList = append(
			tokenList,
			jsonToken{
				strings.ToLower(t.Type.String()),
				t.Value,
			},
		)
	}

	return w.Encode(tokenList)
}
