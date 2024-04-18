package sqlparse

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFormatWithReident(t *testing.T) {
	tests := []struct {
		query    string
		expected string
	}{
		{
			query:    `SELECT * FROM foo`,
			expected: `SELECT * FROM foo`,
		},
		{
			query:    `SELECT * FROM (SELECT bar FROM foo)`,
			expected: "SELECT * FROM (\n  SELECT bar FROM foo\n)",
		},
		{
			query:    `WITH foo AS (SELECT foos, bars FROM foo_list), bar AS (SELECT * FROM foo ORDER BY bars DESC) SELECT * FROM bar`,
			expected: "WITH \nfoo AS (\n  SELECT foos, bars \n  FROM foo_list\n), \nbar AS (\n  SELECT * FROM foo \n  ORDER BY bars DESC\n) \nSELECT * FROM bar",
		},
		{
			query:    `WITH foo AS (SELECT foos, bars FROM foo_list WHERE foos IN (SELECT foo_id FROM ids WHERE active = true)), bar AS (SELECT * FROM foo ORDER BY bars DESC) SELECT * FROM bar`,
			expected: "WITH \nfoo AS (\n  SELECT foos, bars \n  FROM foo_list \n  WHERE foos IN (\n    SELECT foo_id FROM ids \n    WHERE active = true\n  )\n), \nbar AS (\n  SELECT * FROM foo \n  ORDER BY bars DESC\n) \nSELECT * FROM bar",
		},
		{
			query:    `SELECT foo,bar,baz,abc,xyz FROM table_name`,
			expected: "SELECT foo,bar,baz,abc,xyz \nFROM table_name",
		},
		{
			query:    "SELECT (1 + 1)\nAS result",
			expected: "SELECT (1 + 1)\nAS result",
		},
		{
			query:    `SELECT (1+1), (2+2)`,
			expected: `SELECT (1+1), (2+2)`,
		},
		{
			query:    `SELECT (1+1),(2+2)`,
			expected: `SELECT (1+1),(2+2)`,
		},
	}

	for _, test := range tests {
		t.Run(test.query, func(t *testing.T) {
			tokens, err := GetTokens(test.query)
			resultReident := Format(
				tokens,
				FormatOptionReident(true),
				FormatOptionFromBreakCount(3),
			)
			resultNoOptions := Format(tokens)

			require.NoError(t, err, "GetTokens")
			assert.Equal(t, test.query, resultNoOptions)
			assert.Equal(t, test.expected, resultReident)
		})
	}
}
