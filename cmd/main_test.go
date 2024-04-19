package main

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"testing"
)

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
	require.Equal(t, "SELECT * FROM foo\n", buf.String())
}

func TestSingleSelectQueryFormattedLong(t *testing.T) {
	var buf bytes.Buffer
	const query = "SELECT * FROM foo"
	err := run(&buf, "--format", query)
	require.NoError(t, err)
	require.Equal(t, "SELECT * FROM foo\n", buf.String())
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
