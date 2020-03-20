package gmt

var extraWhitespace = []Replacement{
	NewReplacement("\\r", ""),
	NewReplacement("\\(", " ("),
	NewReplacement("\\)", ") "),
	NewReplacement(" +", " "),
	NewReplacement("\\) ([.!:?;,])", "}$1"),
	NewReplacement("\\( ", "("),
	NewReplacement(" \\)", ")"),
	NewReplacement("(\\d) %", "$1%"),
	NewReplacement(" :", ":"),
	NewReplacement(" ;", ";"),
}

var normalizeUnicodeIfNotPenn = []Replacement{
	NewReplacement("`", "'"),
	NewReplacement("''", " \" "),
}

var normalizeUnicode = []Replacement{
	NewReplacement("„", "\""),
	NewReplacement("“", "\""),
	NewReplacement("”", "\""),
	NewReplacement("–", "-"),
	NewReplacement("—", " - "),
	NewReplacement(" +", " "),
	NewReplacement("´", "'"),
	NewReplacement("([a-zA-Z])‘([a-zA-Z])", "$1'$2"),
	NewReplacement("([a-zA-Z])’([a-zA-Z])", "$1'$2"),
	NewReplacement("‘", "'"),
	NewReplacement("‚", "'"),
	NewReplacement("’", "'"),
	NewReplacement("''", "\""),
	NewReplacement("´´", "\""),
	NewReplacement("…", "..."),
}

var frenchQuotes = []Replacement{
	NewReplacement("\u00A0«\u00A0", "\""),
	NewReplacement("«\u00A0", "\""),
	NewReplacement("«", "\""),
	NewReplacement("\u00A0»\u00A0", "\""),
	NewReplacement("\u00A0»", "\""),
	NewReplacement("»", "\""),
}

var pseudoSpaces = []Replacement{
	NewReplacement("\u00A0%", "%"),
	NewReplacement("nº\u00A0", "nº "),
	NewReplacement("\u00A0:", ":"),
	NewReplacement("\u00A0ºC", " ºC"),
	NewReplacement("\u00A0cm", " cm"),
	NewReplacement("\u00A0\\?", "?"),
	NewReplacement("\u00A0\\!", "!"),
	NewReplacement("\u00A0;", ";"),
	NewReplacement(",\u00A0", ", "),
	NewReplacement(" +", " "),
}

var enQuotationFollowedByComma = []Replacement{
	NewReplacement("\"([,.]+)", "$1\""),
}

var deEsFrQuotationFollowedByComma = []Replacement{
	NewReplacement(",\"", "\","),
	NewReplacement("(\\.+)\"(\\s*[^<])", "\"$1$2"),
}

var deEsCzCsFr = []Replacement{
	NewReplacement("(\\d)\u00A0(\\d)", "\\$1,$2"),
}

var other = []Replacement{
	NewReplacement("(\\d)\u00A0(\\d)", `$1.$2`),
}

var replaceUnicodePunctuation = []Replacement{
	NewReplacement(`，`, ","),
	NewReplacement(`。\s*`, ". "),
	NewReplacement(`、`, ","),
	NewReplacement(`”`, `"`),
	NewReplacement(`“`, `"`),
	NewReplacement(`∶`, ":"),
	NewReplacement(`：`, ":"),
	NewReplacement(`？`, "?"),
	NewReplacement(`《`, `"`),
	NewReplacement(`》`, `"`),
	NewReplacement(`）`, ")"),
	NewReplacement(`！`, "!"),
	NewReplacement(`（`, "("),
	NewReplacement(`；`, ";"),
	NewReplacement(`」`, `"`),
	NewReplacement(`「`, `"`),
	NewReplacement(`０`, "0"),
	NewReplacement(`１`, `1`),
	NewReplacement(`２`, "2"),
	NewReplacement(`３`, "3"),
	NewReplacement(`４`, "4"),
	NewReplacement(`５`, "5"),
	NewReplacement(`６`, "6"),
	NewReplacement(`７`, "7"),
	NewReplacement(`８`, "8"),
	NewReplacement(`９`, "9"),
	NewReplacement(`．\s*`, ". "),
	NewReplacement(`～`, "~"),
	NewReplacement(`’`, "'"),
	NewReplacement(`…`, "..."),
	NewReplacement(`━`, "-"),
	NewReplacement(`〈`, "<"),
	NewReplacement(`〉`, ">"),
	NewReplacement(`【`, "["),
	NewReplacement(`】`, "]"),
	NewReplacement(`％`, "%"),
}

var controlChars = []Replacement{
	NewReplacement("\\p{C}", ""),
}

// Normalizer is a golang port of the MOses punctuation normalizer from
// https://github.com/moses-smt/mosesdecoder/blob/master/scripts/tokenizer/normalize-punctuation.perl
// Designs are mostly copied from the python version https://github.com/alvations/sacremoses/blob/master/sacremoses/normalize.py
type Normalizer struct {
	rep []Replacement
}

// NewNormalizer create new instance of normalizer. Several parameters are provided to disable
// specific rules for normalization such as quote normalization, number normalization and
// unicode normalization
func NewNormalizer(lang string, penn bool, normQuoteCommas bool, normNumbers bool, preReplaceUniPunct bool, postRemoveCtrlChars bool) *Normalizer {
	var rules [][]Replacement

	if preReplaceUniPunct {
		rules = append(rules, replaceUnicodePunctuation)
	}

	rules = append(rules, extraWhitespace)
	if penn {
		rules = append(rules, normalizeUnicodeIfNotPenn)
	}
	rules = append(rules, normalizeUnicode)
	rules = append(rules, frenchQuotes)
	rules = append(rules, pseudoSpaces)

	if normQuoteCommas {
		switch lang {
		case "en":
			rules = append(rules, enQuotationFollowedByComma)
		case "de", "es", "fr":
			rules = append(rules, deEsFrQuotationFollowedByComma)
		}
	}

	if normNumbers {
		switch lang {
		case "de", "es", "cz", "cs", "fr":
			rules = append(rules, deEsCzCsFr)
		default:
			rules = append(rules, other)
		}
	}

	if postRemoveCtrlChars {
		rules = append(rules, controlChars)
	}

	normalizer := &Normalizer{Flatten(rules)}

	return normalizer
}

// Normalize the incoming text according to pre-defined rules
func (n Normalizer) Normalize(text string) (normalizedText string) {
	normalizedText = text
	for _, re := range n.rep {
		normalizedText = re.rgx.ReplaceAllString(normalizedText, re.sub)
	}
	return
}
