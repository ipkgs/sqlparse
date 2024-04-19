package main

import (
	"bufio"
	"fmt"
	"github.com/ipkgs/sqlparse"
	"io"
	"os"
	"strconv"
	"strings"
)

func usage(out io.Writer) {
	fmt.Fprintf(out, "usage: sqlparse [options] <sql>\n")
	fmt.Fprintf(out, "the sql query can be replaced by '-' to read from stdin\n")
	fmt.Fprintf(out, "options:\n")
	fmt.Fprintf(out, "  -h, --help: show this help message\n")
	fmt.Fprintf(out, "  -f, --format: formats the sql query\n")
	fmt.Fprintf(out, "  -r, --reident: reindent the sql query\n")
	fmt.Fprintf(out, "  -c, --from-break-count: number of line breaks after FROM clause (use -c multiple times to increase\n")
	fmt.Fprintf(out, "                          the number of fields, or use the long form with a number parameter)\n")
	fmt.Fprintf(out, "  -C, --remove-comments: remove comments from the sql query\n")
}

type options struct {
	help           bool
	format         bool
	reident        bool
	fromCount      int
	removeComments bool
}

func run(out io.Writer, args ...string) error {
	if len(args) < 1 {
		usage(out)
		return fmt.Errorf("too few arguments")
	}

	startPos := 0
	var o options
	for startPos < len(args) && strings.HasPrefix(args[startPos], "-") {
		currentOption := args[startPos]

		if currentOption == "-" {
			break
		}
		if currentOption == "--" {
			startPos++
			break
		}

		if currentOption[1] != '-' {
			// short options
			for i := 1; i < len(currentOption); i++ {
				switch currentOption[i] {
				case 'h':
					o.help = true
				case 'f':
					o.format = true
				case 'r':
					o.reident = true
				case 'c':
					o.fromCount++
				case 'C':
					o.removeComments = true
				default:
					return fmt.Errorf("unknown option: -%c", currentOption[i])
				}
			}
		} else {
			// long options
			switch currentOption {
			case "--help":
				o.help = true
			case "--format":
				o.format = true
			case "--reident":
				o.reident = true
			case "--from-break-count":
				if startPos+1 >= len(args) {
					return fmt.Errorf("missing parameter for --from-break-count")
				}
				nextParam := args[startPos+1]
				if strings.HasPrefix(nextParam, "-") {
					return fmt.Errorf("missing parameter for --from-break-count")
				}
				count, err := strconv.Atoi(nextParam)
				if err != nil {
					return fmt.Errorf("invalid parameter for --from-break-count: %s", nextParam)
				}
				o.fromCount = count
			case "--remove-comments":
				o.removeComments = true
			default:
				return fmt.Errorf("unknown option: %s", currentOption)
			}
		}

		startPos++
	}

	if o.help {
		usage(out)
		return nil
	}

	if len(args) <= startPos {
		usage(out)
		return fmt.Errorf("missing sql query")
	}

	var query string
	if args[startPos] == "-" {
		reader := bufio.NewReader(os.Stdin)
		queryBytes, err := io.ReadAll(reader)
		if err != nil {
			return fmt.Errorf("error read stdin: %w", err)
		}
		query = string(queryBytes)
	} else {
		query = strings.Join(args[startPos:], " ")
	}

	tokens, err := sqlparse.GetTokens(query)

	if err != nil {
		return fmt.Errorf("error: %w", err)
	}

	if o.format {
		var formatOptions []sqlparse.FormatOption

		if o.reident {
			formatOptions = append(formatOptions, sqlparse.FormatOptionReident(true))
		}
		if o.fromCount > 0 {
			formatOptions = append(formatOptions, sqlparse.FormatOptionFromBreakCount(o.fromCount))
		}
		if o.removeComments {
			formatOptions = append(formatOptions, sqlparse.FormatOptionRemoveComments(true))
		}

		formattedQuery := sqlparse.Format(tokens, formatOptions...)
		fmt.Fprintf(out, formattedQuery)
		fmt.Fprintf(out, "\n")
		return nil
	}

	for _, token := range tokens {
		fmt.Fprintf(out, "%s: %s\n", token.Type, token.Value)
	}

	return nil
}

func main() {
	fmt.Println("sqlparse")

	if err := run(os.Stdout, os.Args[1:]...); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
