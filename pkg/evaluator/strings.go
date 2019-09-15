package evaluator

import (
	"fmt"
	"os"

	"golang.org/x/text/collate"
	"golang.org/x/text/language"
)

var collator *collate.Collator

func init() {
	lang := os.Getenv("LANG")
	if lang == "" {
		lang = "en_US"
	}

	if len(lang) > 5 {
		lang = lang[0:5]
	}

	tag, err := language.Parse(lang)
	if err != nil {
		fmt.Printf("Cannot parse %q using \"en_US\" instead\n", lang)
		tag = language.AmericanEnglish
	}

	collator = collate.New(tag)
}

func compareStrings(a, b []rune) int {
	return collator.CompareString(string(a), string(b))
}

func compareRunes(a, b rune) int {
	return collator.CompareString(string(a), string(b))
}
