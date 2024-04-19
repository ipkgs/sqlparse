package sqlparse

import (
	"regexp"
	"strings"
)

//go:generate stringer -type=TokenType -trimprefix=Token
type TokenType int

const (
	TokenUnknown TokenType = iota
	TokenWhitespace
	TokenNewline
	TokenKeyword
	TokenKeywordCTE
	TokenOperator
	TokenUseAsKeyword
	TokenPunctuation
	TokenName
	TokenWildcard
	TokenNumberInteger
	TokenNumberFloat
	TokenString
	TokenComment
)

type matchInstruction[T any] struct {
	value           T
	instructionType TokenType
}

var defaultRegexChecks = []matchInstruction[*regexp.Regexp]{
	{regexp.MustCompile(`[\r\n]+`), TokenNewline},
	{regexp.MustCompile(`\s+`), TokenWhitespace},
	{regexp.MustCompile(`--.*?(\r\n|\r|\n|$)`), TokenComment},

	{regexp.MustCompile(`\*`), TokenWildcard},
	{regexp.MustCompile(`-?\d+(\.\d+)?E-?\d+`), TokenNumberFloat},
	{regexp.MustCompile(`[^'"()_A-ZÀ-Ü]-?(\d+(\.\d*)|\.\d+)[^'"()_A-ZÀ-Ü]`), TokenNumberFloat},
	{regexp.MustCompile(`[^'"()_A-ZÀ-Ü]-?\d+[^'"()_A-ZÀ-Ü]`), TokenNumberInteger},
	{
		regexp.MustCompile(`'(''|\\'|[^'])*'`),
		TokenString,
	},
	{
		regexp.MustCompile("`(\\`|[^`])*`"),
		TokenString,
	},
	{
		regexp.MustCompile(`((LEFT\s+|RIGHT\s+|FULL\s+)?(INNER\s+|OUTER\s+|STRAIGHT\s+)?|(CROSS\s+|NATURAL\s+)?)?JOIN\b`),
		TokenKeyword,
	},
	{regexp.MustCompile(`ORDER\s+BY\b`), TokenKeyword},
	{regexp.MustCompile(`GROUP\s+BY\b`), TokenKeyword},
	{regexp.MustCompile(`UNION\s+ALL\b`), TokenKeyword},
	{
		regexp.MustCompile(`[<>=~!]+`),
		TokenOperator,
	},
	{
		regexp.MustCompile(`\[+/@#%^&|^-]+\b`),
		TokenOperator,
	},
	{
		regexp.MustCompile(`\w[$#\w]*`),
		TokenUseAsKeyword,
	},
	{
		regexp.MustCompile(`[;()[\],.]`),
		TokenPunctuation,
	},
	{regexp.MustCompile(`[+/@#%^&|-]+`), TokenOperator},
}

var defaultKeywords = []matchInstruction[string]{
	{"SELECT", TokenKeyword},
	{"FROM", TokenKeyword},
	{"WHERE", TokenKeyword},
	{"AS", TokenKeyword},
	{"WITH", TokenKeywordCTE},
	{"AND", TokenKeyword},
	{"OR", TokenKeyword},
	{"NOT", TokenKeyword},
	{"IS", TokenKeyword},
	{"NULL", TokenKeyword},
	{"LIKE", TokenKeyword},
	{"IN", TokenKeyword},
	{"ORDER", TokenKeyword},
	{"BY", TokenKeyword},
	{"UNION", TokenKeyword},
}

type Token struct {
	Value string
	Type  TokenType
}

type Lexer struct {
	regexChecks []matchInstruction[*regexp.Regexp]
	keywords    []matchInstruction[string]
}

func defaultLexer() *Lexer {
	return &Lexer{
		regexChecks: defaultRegexChecks,
		keywords:    defaultKeywords,
	}
}

func findInSlice(needle string, haystack []matchInstruction[string]) *matchInstruction[string] {
	for _, check := range haystack {
		if check.value == needle {
			return &check
		}
	}
	return nil
}

// Clear removes all keywords and regex checks from the lexer, useful to reset the lexer to a clean state
// while adding your own custom keywords and regex checks, without needing to create a new lexer instance.
func (l *Lexer) Clear() {
	l.keywords = []matchInstruction[string]{}
	l.regexChecks = []matchInstruction[*regexp.Regexp]{}
}

func (l *Lexer) AddKeyword(keyword string, keywordType TokenType) {
	if findInSlice(keyword, l.keywords) != nil {
		return
	}
	l.keywords = append(l.keywords, matchInstruction[string]{value: keyword, instructionType: keywordType})
}

func (l *Lexer) AddRegexp(re *regexp.Regexp, tokenType TokenType) {
	l.regexChecks = append(l.regexChecks, matchInstruction[*regexp.Regexp]{value: re, instructionType: tokenType})
}

func (l *Lexer) process(accum string) (t Token) {
	var match *matchInstruction[*regexp.Regexp]
	var strMatch string
	for _, check := range l.regexChecks {
		matchPos := check.value.FindStringIndex(accum)
		if len(matchPos) == 0 {
			continue
		}
		if matchPos[0] != 0 {
			// we only care about matches on the beginning of the string
			continue
		}

		strMatch = accum[matchPos[0]:matchPos[1]]
		if strMatch != "" {
			match = &check
			break
		}
	}

	if match == nil {
		return
	}
	matchType := match.instructionType

	if matchType == TokenUnknown {
		return
	}

	if matchType == TokenUseAsKeyword {
		if l.IsKeyword(strMatch) {
			matchType = TokenKeyword
		} else {
			matchType = TokenName
		}
	}

	t.Value = strMatch
	t.Type = matchType

	return t
}

func (l *Lexer) IsKeyword(s string) bool {
	return findInSlice(strings.ToUpper(s), l.keywords) != nil
}

func (l *Lexer) GetTokens(data string) ([]Token, error) {
	var tokens []Token
	var pos int

	for pos < len(data) {
		token := l.process(data[pos:])
		if token.Value == "" {
			break
		}

		tokens = append(tokens, token)
		pos += len(token.Value)
	}

	return tokens, nil
}

func GetTokens(data string) ([]Token, error) {
	lexer := defaultLexer()
	return lexer.GetTokens(data)
}
