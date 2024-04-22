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
func FormatOptionRemoveComments(value bool) FormatOption {
	return func(f *formatOptionList) {
		f.removeComments = value
	}
}

func FormatOptionUppercaseKeywords(value bool) FormatOption {
	return func(f *formatOptionList) {
		f.uppercaseKeywords = value
	}
}

type formatOptionList struct {
	reident           bool
	fromBreakCount    int
	removeComments    bool
	uppercaseKeywords bool

	parenthesisIdented []bool
	writtenInThisLine  bool
	lastWrittenToken   *Token
	spaceQueued        string
	buf                strings.Builder
}

func (f *formatOptionList) write(s string) {
	if f.reident {
		if trimmedSpace := strings.Trim(s, " \t"); trimmedSpace == "" {
			f.spaceQueued += s
			return
		}
		if s == "\n" {
			f.spaceQueued = ""
			f.writtenInThisLine = false
		} else if !f.writtenInThisLine {
			var identLevel int
			for _, b := range f.parenthesisIdented {
				if b {
					identLevel++
				}
			}

			f.buf.WriteString(strings.Repeat("  ", identLevel))
			f.writtenInThisLine = true
		}
	}

	if f.spaceQueued != "" {
		f.buf.WriteString(f.spaceQueued)
		f.spaceQueued = ""
	}

	f.buf.WriteString(s)
}

func (f *formatOptionList) writeLinebreak() {
	if f.reident {
		f.write("\n")
	}
}

func (f *formatOptionList) writeToken(tokens []Token, pos int) {
	tokenType := tokens[pos].Type
	tokenValue := tokens[pos].Value
	if tokenType == TokenComment && f.removeComments {
		return
	}

	if f.reident {
		if tokenValue == "(" {
			var shouldIdent bool

			// if next keyword is SELECT, then we should ident. If we find any punctuation before the keyword, bail
			var nextKeywordToken *Token
			for i := pos + 1; i < len(tokens); i++ {
				if tokens[i].Type == TokenPunctuation {
					break
				}
				if tokens[i].Type == TokenKeyword {
					nextKeywordToken = &tokens[i]
					break
				}
			}
			if nextKeywordToken != nil && nextKeywordToken.Value == "SELECT" {
				shouldIdent = true
			}

			f.parenthesisIdented = append(f.parenthesisIdented, shouldIdent)
		} else if tokenValue == ")" {
			var isParenthesisIdented bool
			if len(f.parenthesisIdented) > 0 {
				isParenthesisIdented = f.parenthesisIdented[len(f.parenthesisIdented)-1]
				f.parenthesisIdented = f.parenthesisIdented[:len(f.parenthesisIdented)-1]
			}
			if isParenthesisIdented {
				f.writeLinebreak()
			}
		}

		if tokenType == TokenWhitespace || tokenType == TokenNewline {
			if !f.writtenInThisLine {
				return
			}
			var nextToken *Token
			if pos+1 < len(tokens) {
				nextToken = &tokens[pos+1]
			}
			if nextToken != nil {
				if nextToken.Type == TokenNewline || nextToken.Type == TokenComment {
					return
				}
			}
			if f.lastWrittenToken != nil && f.lastWrittenToken.Type != TokenWhitespace && f.lastWrittenToken.Type != TokenNewline {
				f.lastWrittenToken = &tokens[pos]
				f.write(" ")
			}
			return
		}
		if tokenType == TokenComment {
			f.writeLinebreak()
		}
		if tokenValue == "SELECT" {
			if len(f.parenthesisIdented) > 0 {
				f.writeLinebreak()
			} else {
				// get previous non-whitespace, non-newline token
				var prevToken *Token
				for i := pos - 1; i >= 0; i-- {
					if tokens[i].Type != TokenWhitespace && tokens[i].Type != TokenNewline {
						prevToken = &tokens[i]
						break
					}
				}
				if prevToken != nil && prevToken.Value == ")" {
					f.writeLinebreak()
				}
			}
		}

		if tokenType == TokenName {
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
					f.writeLinebreak()
				}
			}
		}

		if tokenValue == "FROM" {
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

			if tokensSinceSelect >= f.fromBreakCount {
				f.writeLinebreak()
			}
		}

		if tokenValue == "WHERE" || tokenValue == "ORDER BY" || tokenValue == "GROUP BY" || tokenValue == "UNION ALL" || tokenValue == "LEFT OUTER JOIN" {
			f.writeLinebreak()
		}

		tokenValue = strings.TrimSpace(tokenValue)
	}

	if f.uppercaseKeywords && tokenType == TokenKeyword {
		tokenValue = strings.ToUpper(tokenValue)
	}

	f.lastWrittenToken = &tokens[pos]
	f.write(tokenValue)
}

func (f *formatOptionList) formattedQuery(tokens []Token) string {
	for i := 0; i < len(tokens); i++ {
		f.writeToken(tokens, i)
	}

	return f.buf.String()

}

func Format(tokens []Token, optionList ...FormatOption) string {
	var options formatOptionList

	for _, option := range optionList {
		option(&options)
	}

	return options.formattedQuery(tokens)
}
