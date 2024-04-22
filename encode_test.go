package sqlparse

import (
	"encoding/json"
	"fmt"
	"github.com/pmezard/go-difflib/difflib"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestEncode(t *testing.T) {
	tests := []struct {
		query    string
		expected string
	}{
		{
			`SELECT * FROM foo`,
			`[{"type":"keyword","value":"SELECT"},{"type":"whitespace","value":" "},{"type":"wildcard","value":"*"},{"type":"whitespace","value":" "},{"type":"keyword","value":"FROM"},{"type":"whitespace","value":" "},{"type":"name","value":"foo"}]`,
		},
		{
			query:    "WHERE id = 99999 AND exists",
			expected: `[{"type":"keyword","value":"WHERE"},{"type":"whitespace","value":" "},{"type":"name","value":"id"},{"type":"whitespace","value":" "},{"type":"operator","value":"="},{"type":"whitespace","value":" "},{"type":"numberinteger","value":"99999"},{"type":"whitespace","value":" "},{"type":"keyword","value":"AND"},{"type":"whitespace","value":" "},{"type":"name","value":"exists"}]`,
		},
		{
			query:    "WHERE id = 99999\nAND exists",
			expected: `[{"type":"keyword","value":"WHERE"},{"type":"whitespace","value":" "},{"type":"name","value":"id"},{"type":"whitespace","value":" "},{"type":"operator","value":"="},{"type":"whitespace","value":" "},{"type":"numberinteger","value":"99999"},{"type":"newline","value":"\n"},{"type":"keyword","value":"AND"},{"type":"whitespace","value":" "},{"type":"name","value":"exists"}]`,
		},
		{
			query:    "WHERE id = 99999\tAND exists",
			expected: `[{"type":"keyword","value":"WHERE"},{"type":"whitespace","value":" "},{"type":"name","value":"id"},{"type":"whitespace","value":" "},{"type":"operator","value":"="},{"type":"whitespace","value":" "},{"type":"numberinteger","value":"99999"},{"type":"whitespace","value":"\t"},{"type":"keyword","value":"AND"},{"type":"whitespace","value":" "},{"type":"name","value":"exists"}]`,
		},
		{
			query:    "IF(x IN (1,1,2,3,5,8,13,21), 'A', 'B')",
			expected: `[{"type":"name","value":"IF"},{"type":"punctuation","value":"("},{"type":"name","value":"x"},{"type":"whitespace","value":" "},{"type":"keyword","value":"IN"},{"type":"whitespace","value":" "},{"type":"punctuation","value":"("},{"type":"numberinteger","value":"1"},{"type":"punctuation","value":","},{"type":"numberinteger","value":"1"},{"type":"punctuation","value":","},{"type":"numberinteger","value":"2"},{"type":"punctuation","value":","},{"type":"numberinteger","value":"3"},{"type":"punctuation","value":","},{"type":"numberinteger","value":"5"},{"type":"punctuation","value":","},{"type":"numberinteger","value":"8"},{"type":"punctuation","value":","},{"type":"numberinteger","value":"13"},{"type":"punctuation","value":","},{"type":"numberinteger","value":"21"},{"type":"punctuation","value":")"},{"type":"punctuation","value":","},{"type":"whitespace","value":" "},{"type":"string","value":"'A'"},{"type":"punctuation","value":","},{"type":"whitespace","value":" "},{"type":"string","value":"'B'"},{"type":"punctuation","value":")"}]`,
		},
	}

	for _, test := range tests {
		testName := strings.Replace(test.query, "\n", " ", -1)
		if len(testName) > 15 {
			testName = testName[:12] + "..."
		}

		t.Run(testName, func(t *testing.T) {
			var sb strings.Builder
			encoder := json.NewEncoder(&sb)
			tokens, err := GetTokens(test.query)

			require.NoError(t, err, "GetTokens")
			require.NoError(t, Encode(encoder, tokens), "EncodeJSON")

			result := sb.String()
			// Remove any generated leading or trailing newlines and line breaks
			result = strings.TrimSpace(result)

			//assert.Equal(t, test.expected, result)

			if test.expected != result {
				var expectedTokens []jsonToken
				var resultTokens []jsonToken
				for _, token := range tokens {
					resultTokens = append(resultTokens, jsonToken{
						Type:  strings.ToLower(token.Type.String()),
						Value: token.Value,
					})
				}

				err := json.Unmarshal([]byte(test.expected), &expectedTokens)
				require.NoError(t, err, "Unmarshal error for expected string")

				formattedExpectedBytes, err := json.MarshalIndent(expectedTokens, "", "  ")
				require.NoError(t, err, "MarshalIndent error for expected string")

				formattedResultBytes, err := json.MarshalIndent(resultTokens, "", "  ")
				require.NoError(t, err, "MarshalIndent error for result string")

				diff := difflib.UnifiedDiff{
					A:        difflib.SplitLines(string(formattedExpectedBytes)),
					B:        difflib.SplitLines(string(formattedResultBytes)),
					FromFile: "Expected",
					ToFile:   "Actual",
					Context:  3,
				}
				text, err := difflib.GetUnifiedDiffString(diff)
				require.NoError(t, err, "GetUnifiedDiffString error")

				t.Errorf("Not equal")
				fmt.Println(result)
				fmt.Printf(text)
			}
		})
	}
}
