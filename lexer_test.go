package sqlparse

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestProcess(t *testing.T) {
	tests := []struct {
		piece         string
		expectedMatch string
		expectedType  TokenType
	}{
		{"=", "=", TokenOperator},
		{"= ", "=", TokenOperator},
		{"foo", "foo", TokenName},
		{"foo ", "foo", TokenName},
		{"foo =", "foo", TokenName},
		{"foo=", "foo", TokenName},
		{"foo,", "foo", TokenName},
		{"foo ,", "foo", TokenName},
		{"foo !=\"*", "foo", TokenName},
		{"ORDER BY", "ORDER BY", TokenKeyword},
		{"ORDER BY ", "ORDER BY", TokenKeyword},
		{"ORDER", "ORDER", TokenKeyword},
		{"ORDER ", "ORDER", TokenKeyword},
		{"1", "1", TokenNumberInteger},
		{"1,", "1", TokenNumberInteger},
		{"1E99", "1E99", TokenNumberFloat},
		{"1.99", "1.99", TokenNumberFloat},
		{"-1.99", "-1.99", TokenNumberFloat},
	}

	lexer := defaultLexer()

	for _, test := range tests {
		t.Run(test.piece, func(t *testing.T) {
			token := lexer.process(test.piece)
			assert.Equal(t, test.expectedMatch, token.Value, "token.Value")
			assert.Equalf(t, test.expectedType, token.Type, "token.Type: expected %s is not %s", test.expectedType, token.Type)
		})
	}
}

func TestDefaultLexer(t *testing.T) {
	tests := []struct {
		query         string
		expectedCount int
	}{
		{
			query:         `SELECT * FROM bar`,
			expectedCount: 7,
		},
		{
			query:         `SELECT foo FROM bar`,
			expectedCount: 7,
		},
		{
			query:         `SELECT foo, baz FROM bar`,
			expectedCount: 10,
		},
		{
			query:         `SELECT foo, baz FROM bar WHERE foo = 1`,
			expectedCount: 18,
		},
		{
			query:         `SELECT foo, baz FROM bar WHERE foo =1`,
			expectedCount: 17,
		},
		{
			query:         `SELECT foo, baz FROM bar WHERE foo= 1`,
			expectedCount: 17,
		},
		{
			query:         `SELECT foo, baz FROM bar WHERE foo=1`,
			expectedCount: 16,
		},
		{
			query:         `SELECT distance, baz FROM bar WHERE distance >= 3.1415`,
			expectedCount: 18,
		},
		{
			query:         `SELECT distance, baz FROM bar WHERE distance >= 314.15E-2`,
			expectedCount: 18,
		},
		{
			query:         `SELECT distance, baz FROM bar WHERE distance >= 31415E-4`,
			expectedCount: 18,
		},
		{
			query:         `SELECT foo, baz FROM bar WHERE foo = 99`,
			expectedCount: 18,
		},
		{
			query:         `SELECT foo, baz FROM bar WHERE foo = 99 AND baz = 'hello'`,
			expectedCount: 26,
		},
		{
			query:         `SELECT foo, baz FROM bar WHERE foo = 99 AND baz = 'hello world'`,
			expectedCount: 26,
		},
		{
			query:         "SELECT foo, baz\nFROM bar\nWHERE foo = 99 AND baz = 'hello world'",
			expectedCount: 26,
		},
		{
			query:         `WITH cte AS (SELECT * FROM xyz WHERE k = 0) SELECT * FROM cte ORDER BY k`,
			expectedCount: 35,
		},
		{
			query:         `SELECT (1 + 1) AS result`,
			expectedCount: 13,
		},
		{
			query:         `SELECT (1+1) AS result`,
			expectedCount: 11,
		},
		{
			query:         `SELECT (1+1), (2+2)`,
			expectedCount: 14,
		},
		{
			query:         `SELECT (1+1),(2+2)`,
			expectedCount: 13,
		},
		{
			query:         "SELECT *\n-- testing comment\nFROM bar",
			expectedCount: 8,
		},
		{
			query:         "SELECT *, xyz, abc FROM `scope.group.table_name`",
			expectedCount: 13,
		},
	}

	for _, test := range tests {
		t.Run(test.query, func(t *testing.T) {
			resp, err := GetTokens(test.query)
			require.NoError(t, err, "sqlparse.GetTokens")

			assert.Len(t, resp, test.expectedCount, "resp")

			var sb strings.Builder
			for _, token := range resp {
				sb.WriteString(token.Value)
			}
			assert.Equal(t, test.query, sb.String(), "query")
		})
	}
}
