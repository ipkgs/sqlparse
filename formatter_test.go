package sqlparse

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestFormat(t *testing.T) {
	tests := []struct {
		query    string
		expected string
		options  []FormatOption
	}{
		{
			query:    `SELECT * FROM foo`,
			expected: `SELECT * FROM foo`,
			options:  []FormatOption{FormatOptionReident(true), FormatOptionFromBreakCount(3)},
		},
		{
			query:    `SELECT * FROM (SELECT bar FROM foo)`,
			expected: "SELECT * FROM (\n  SELECT bar FROM foo\n)",
			options:  []FormatOption{FormatOptionReident(true), FormatOptionFromBreakCount(3)},
		},
		{
			query:    `WITH foo AS (SELECT foos, bars FROM foo_list), bar AS (SELECT * FROM foo ORDER BY bars DESC) SELECT * FROM bar`,
			expected: "WITH\nfoo AS (\n  SELECT foos, bars\n  FROM foo_list\n),\nbar AS (\n  SELECT * FROM foo\n  ORDER BY bars DESC\n)\nSELECT * FROM bar",
			options:  []FormatOption{FormatOptionReident(true), FormatOptionFromBreakCount(3)},
		},
		{
			query:    `WITH foo AS (SELECT foos, bars FROM foo_list WHERE foos IN (SELECT foo_id FROM ids WHERE active = true)), bar AS (SELECT * FROM foo ORDER BY bars DESC) SELECT * FROM bar`,
			expected: "WITH\nfoo AS (\n  SELECT foos, bars\n  FROM foo_list\n  WHERE foos IN (\n    SELECT foo_id FROM ids\n    WHERE active = true\n  )\n),\nbar AS (\n  SELECT * FROM foo\n  ORDER BY bars DESC\n)\nSELECT * FROM bar",
			options:  []FormatOption{FormatOptionReident(true), FormatOptionFromBreakCount(3)},
		},
		{
			query:    `SELECT foo,bar,baz,abc,xyz FROM table_name`,
			expected: "SELECT foo,bar,baz,abc,xyz\nFROM table_name",
			options:  []FormatOption{FormatOptionReident(true), FormatOptionFromBreakCount(3)},
		},
		{
			query:    "SELECT (1 + 1)\nAS result",
			expected: "SELECT (1 + 1) AS result",
			options:  []FormatOption{FormatOptionReident(true), FormatOptionFromBreakCount(3)},
		},
		{
			query:    `SELECT (1+1), (2+2)`,
			expected: `SELECT (1+1), (2+2)`,
			options:  []FormatOption{FormatOptionReident(true), FormatOptionFromBreakCount(3)},
		},
		{
			query:    `SELECT (1+1),(2+2)`,
			expected: `SELECT (1+1),(2+2)`,
			options:  []FormatOption{FormatOptionReident(true), FormatOptionFromBreakCount(3)},
		},
		{
			query:    "SELECT *\n-- testing comment\nFROM bar",
			expected: "SELECT *\n-- testing comment\nFROM bar",
			options:  []FormatOption{FormatOptionReident(true)},
		},
		{
			query:    "SELECT * -- testing comment\nFROM bar",
			expected: "SELECT *\n-- testing comment\nFROM bar",
			options:  []FormatOption{FormatOptionReident(true)},
		},
		{
			query:    "SELECT * -- testing comment\nFROM bar",
			expected: "SELECT * FROM bar",
			options:  []FormatOption{FormatOptionRemoveComments(true)},
		},
		{
			query:    "SELECT *\n-- testing comment\nFROM bar",
			expected: "SELECT *\nFROM bar",
			options:  []FormatOption{FormatOptionRemoveComments(true)},
		},
		{
			query:    `WITH complicated AS (SELECT some_id AS id, IF(x IN (1,1,2,3,5,8,13,21), 'A', 'B') AS something FROM data_source WHERE i = 9999999999 AND created_at >= DATE('2024-01-01') GROUP BY some_id, IF(x IN (1,1,2,3,5,8,13,21), 'A', 'B'), created_at) SELECT * FROM complicated`,
			expected: "WITH\ncomplicated AS (\n  SELECT some_id AS id, IF(x IN (1,1,2,3,5,8,13,21), 'A', 'B') AS something\n  FROM data_source\n  WHERE i = 9999999999 AND created_at >= DATE('2024-01-01')\n  GROUP BY some_id, IF(x IN (1,1,2,3,5,8,13,21), 'A', 'B'), created_at\n)\nSELECT *\nFROM complicated",
			options:  []FormatOption{FormatOptionReident(true)},
		},
		{
			query:    "SELECT \n * FROM foo",
			expected: `SELECT * FROM foo`,
			options:  []FormatOption{FormatOptionReident(true), FormatOptionFromBreakCount(3)},
		},
		{
			query:    "select * from foo",
			expected: `select * from foo`,
			options:  []FormatOption{FormatOptionUppercaseKeywords(false)},
		},
		{
			query:    "select * from foo",
			expected: `SELECT * FROM foo`,
			options:  []FormatOption{FormatOptionUppercaseKeywords(true)},
		},
	}

	for _, test := range tests {
		testName := test.query
		if len(testName) > 15 {
			testName = testName[:12] + "..."
		}
		t.Run(testName, func(t *testing.T) {
			tokens, err := GetTokens(test.query)
			resultFormatted := Format(tokens, test.options...)
			resultNoOptions := Format(tokens)

			require.NoError(t, err, "GetTokens")
			assert.Equal(t, test.query, resultNoOptions)
			assert.Equal(t, test.expected, resultFormatted)

			if t.Failed() {
				var sb strings.Builder
				encoder := json.NewEncoder(&sb)
				encoder.SetIndent("", "  ")
				if err := Encode(encoder, tokens); err != nil {
					t.Fatal(err)
				}

				t.Log(sb.String())
			}
		})
	}
}
