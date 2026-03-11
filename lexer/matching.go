package lexer

// A matcher takes in a token and performs validation/matching on it
type MatcherFunc func(token *Token) error

// Allows MatcherFunc to implement the Matcher interface
func (matcher MatcherFunc) MatchLexToken(token *Token) error {
	return matcher(token)
}

// A criteria that can be matched against a token
type Matcher interface {
	// Returns an error if the token does not match the criteria
	MatchLexToken(*Token) error
}
