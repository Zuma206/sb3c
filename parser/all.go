package parser

// Parses all steps in sequence
type All []Parse

// Parses all steps in sequence
func (all All) Parse(p *Parser) error {
	for _, parse := range all {
		if err := parse.Parse(p); err != nil {
			return err
		}
	}
	return nil
}
