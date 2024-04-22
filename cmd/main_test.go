package main

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestVersion(t *testing.T) {
	var buf bytes.Buffer
	err := run(&buf, "-v")
	require.NoError(t, err)
	require.Contains(t, buf.String(), "sqlparse dev\n")
}

func TestNoParams(t *testing.T) {
	var buf bytes.Buffer
	err := run(&buf)
	require.Error(t, err)
}

func TestSingleSelectQueryFormatted(t *testing.T) {
	var buf bytes.Buffer
	const query = "SELECT * FROM foo"
	err := run(&buf, "-f", query)
	require.NoError(t, err)
	require.Equal(t, query+"\n", buf.String())
}

func TestSingleSelectQueryFormattedLong(t *testing.T) {
	var buf bytes.Buffer
	const query = "SELECT * FROM foo"
	err := run(&buf, "--format", query)
	require.NoError(t, err)
	require.Equal(t, query+"\n", buf.String())
}

func TestSingleSelectQueryJSON(t *testing.T) {
	var buf bytes.Buffer
	const query = "SELECT * FROM foo"
	err := run(&buf, "-j", query)
	require.NoError(t, err)
	require.Equal(t, `[{"type":"keyword","value":"SELECT"},{"type":"whitespace","value":" "},{"type":"wildcard","value":"*"},{"type":"whitespace","value":" "},{"type":"keyword","value":"FROM"},{"type":"whitespace","value":" "},{"type":"name","value":"foo"}]`+"\n", buf.String())
}

func TestQueryReident(t *testing.T) {
	var buf bytes.Buffer
	const query = "SELECT bar, baz, baj, xyz FROM foo"
	err := run(&buf, "-fr", query)
	require.NoError(t, err)
	require.Equal(t, "SELECT bar, baz, baj, xyz\nFROM foo\n", buf.String())
}

func TestQueryReidentLong(t *testing.T) {
	var buf bytes.Buffer
	const query = "SELECT bar, baz, baj, xyz FROM foo"
	err := run(&buf, "--format", "--reident", query)
	require.NoError(t, err)
	require.Equal(t, "SELECT bar, baz, baj, xyz\nFROM foo\n", buf.String())
}

func TestRemoveComments(t *testing.T) {
	var buf bytes.Buffer
	const query = "SELECT bar, baz, baj, xyz FROM foo -- comment"
	err := run(&buf, "-fC", query)
	require.NoError(t, err)
	require.Equal(t, "SELECT bar, baz, baj, xyz FROM foo \n", buf.String())
}

func TestUppercaseKeywords(t *testing.T) {
	var buf bytes.Buffer
	const query = "select bar, baz, baj, xyz from foo"
	err := run(&buf, "-fU", query)
	require.NoError(t, err)
	require.Equal(t, "SELECT bar, baz, baj, xyz FROM foo\n", buf.String())
}

func TestRemoveCommentsJSON(t *testing.T) {
	var buf bytes.Buffer
	const query = "SELECT bar, baz, baj, xyz FROM foo -- comment"
	err := run(&buf, "-jC", query)
	require.NoError(t, err)
	require.Equal(t, `[{"type":"keyword","value":"SELECT"},{"type":"whitespace","value":" "},{"type":"name","value":"bar"},{"type":"punctuation","value":","},{"type":"whitespace","value":" "},{"type":"name","value":"baz"},{"type":"punctuation","value":","},{"type":"whitespace","value":" "},{"type":"name","value":"baj"},{"type":"punctuation","value":","},{"type":"whitespace","value":" "},{"type":"name","value":"xyz"},{"type":"whitespace","value":" "},{"type":"keyword","value":"FROM"},{"type":"whitespace","value":" "},{"type":"name","value":"foo"},{"type":"whitespace","value":" "}]`+"\n", buf.String())
}
