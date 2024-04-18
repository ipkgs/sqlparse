package sqlparse

import (
	"strings"
)

type FormatOption func(*formatOptionList)

func FormatOptionReident(value bool) FormatOption {
	return func(f *formatOptionList) {
		f.reident = value
	}
}
func FormatOptionFromBreakCount(value int) FormatOption {
	return func(f *formatOptionList) {
		f.fromBreakCount = value
	}
}

type formatOptionList struct {
	reident        bool
	fromBreakCount int
}

func Format(tokens []Token, optionList ...FormatOption) string {
	var options formatOptionList

	for _, option := range optionList {
		option(&options)
	}

	var sb strings.Builder
	var parenthesisIdented []bool

	var writtenInThisLine bool
	write := func(s string) {
		if options.reident {
			if s == "\n" {
				writtenInThisLine = false
			} else if !writtenInThisLine {
				var identLevel int
				for _, b := range parenthesisIdented {
					if b {
						identLevel++
					}
				}

				sb.WriteString(strings.Repeat("  ", identLevel))
				writtenInThisLine = true
			}
		}
		sb.WriteString(s)
	}
	writeLinebreak := func() {
		if options.reident {
			write("\n")
		}
	}
	for pos, token := range tokens {
		if token.Value == "(" {
			parenthesisIdented = append(parenthesisIdented, false)
		} else if token.Value == ")" {
			var isParenthesisIdented bool
			if len(parenthesisIdented) > 0 {
				isParenthesisIdented = parenthesisIdented[len(parenthesisIdented)-1]
				parenthesisIdented = parenthesisIdented[:len(parenthesisIdented)-1]
			}
			if isParenthesisIdented {
				writeLinebreak()
			}
		}
		if token.Type == TokenName {
			// if it's a CTE name, the next non-whitespace tokens are always "AS", "(" and "SELECT"
			var nextTokens []Token
			for i := pos + 1; i < len(tokens); i++ {
				if tokens[i].Type != TokenWhitespace && tokens[i].Type != TokenNewline {
					nextTokens = append(nextTokens, tokens[i])
				}
				if len(nextTokens) >= 3 {
					break
				}
			}
			if len(nextTokens) == 3 {
				if nextTokens[0].Value == "AS" && nextTokens[1].Value == "(" && nextTokens[2].Value == "SELECT" {
					writeLinebreak()
				}
			}
		}
		if token.Value == "SELECT" {
			if len(parenthesisIdented) > 0 {
				parenthesisIdented[len(parenthesisIdented)-1] = true
				writeLinebreak()
			} else {
				var lastNonWhitespaceToken *Token
				for i := pos - 1; i >= 0; i-- {
					if tokens[i].Type != TokenWhitespace && tokens[i].Type != TokenNewline {
						lastNonWhitespaceToken = &tokens[i]
						break
					}
				}
				if lastNonWhitespaceToken != nil && lastNonWhitespaceToken.Value == ")" {
					writeLinebreak()
				}
			}
		}
		if token.Value == "FROM" {
			// if we have more than X tokens since the select, break the line
			var tokensSinceSelect int
			for i := pos - 1; i >= 0; i-- {
				if tokens[i].Value == "SELECT" {
					break
				}
				if tokens[i].Type != TokenWhitespace && tokens[i].Type != TokenNewline {
					tokensSinceSelect++
				}
			}

			if tokensSinceSelect >= options.fromBreakCount {
				writeLinebreak()
			}
		}
		if token.Value == "WHERE" || token.Value == "ORDER BY" || token.Value == "GROUP BY" {
			writeLinebreak()
		}

		write(token.Value)
	}

	return sb.String()
}
