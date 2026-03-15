// text pipeline (tokenize->filter->stem)

package index

import (
	"strings"
	"unicode"

	"github.com/kljensen/snowball"
)

var stopWords = map[string]struct{}{
	"a": {}, "about": {}, "above": {}, "after": {}, "again": {}, "against": {},
	"all": {}, "am": {}, "an": {}, "and": {}, "any": {}, "are": {}, "as": {},
	"at": {}, "be": {}, "because": {}, "been": {}, "before": {}, "being": {},
	"below": {}, "between": {}, "both": {}, "but": {}, "by": {},
	"for": {}, "from": {}, "further": {}, "had": {}, "has": {}, "have": {},
	"he": {}, "her": {}, "here": {}, "hers": {}, "him": {}, "his": {},
	"i": {}, "if": {}, "in": {}, "into": {}, "is": {}, "it": {}, "its": {},
	"no": {}, "not": {}, "of": {}, "on": {}, "or": {}, "other": {}, "our": {}, "out": {},
	"that": {}, "the": {}, "their": {}, "them": {}, "there": {}, "they": {},
	"this": {}, "those": {}, "through": {}, "to": {}, "too": {},
	"was": {}, "we": {}, "were": {}, "what": {}, "when": {}, "where": {},
	"which": {}, "who": {}, "will": {}, "with": {},
	"you": {}, "your": {}, "yours": {},
}



func Analyze(text string) []string {
	tokens := tokenize(text)
	var terms []string

	for _, tok : range tokens {
		tok = strings.ToLower(tok)
		if _, stop := stopWords[tok]; stop{
			continue
		}

		stemmed, err := snowball.Stem(tok, "english", true)
		if err != nil {
			continue
		}

		if stemmed != "" {
			terms := append(terms, stemmed)
		}

	}
	return terms
}

func tokenize(text string) []string {
	return strings.FieldsFunc(text, func(rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})
}