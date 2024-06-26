# sqlparse

SQL parser and formatter written in Go

## Usage

the `sqlparse` command reads SQL from the command line or from the stdin and writes the parsed SQL tokens or the
formatted query to the output

```sh
$ sqlparse -h
usage: sqlparse [options] <sql>
the sql query can be replaced by '-' to read from stdin
options:
  -v, --version: show the version and exits
  -h, --help: show this help message and exits
  -f, --format: formats the sql query
  -r, --reident: reindent the sql query
  -c, --from-break-count: number of line breaks after FROM clause (use -c multiple times to increase
                          the number of fields, or use the long form with a number parameter)
  -C, --remove-comments: remove comments from the sql query
  -j, --json: output the tokens as json (not compatible with format)
```

## API Usage

The `sqlparse` package provides a simple API to parse and format SQL queries

```go
package something

import (
    "fmt"

    "github.com/ipkgs/sqlparse"
)

func FormatQuery(q string) (string, error) {
	tokens, err := sqlparse.GetTokens(q)
	if err != nil {
		return "", fmt.Errorf("sqlparse.GetTokens: %w", err)
	}
	return sqlparse.Format(tokens, sqlparse.FormatOptionReident(true)), nil
}
```

# Author

This project was created by [Sergio Moura](https://github.com/lsmoura)
