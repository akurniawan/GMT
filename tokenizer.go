package gmt

import (
	"fmt"
	"regexp"
	"strings"
)

var isN = PerlPropsLoader("isN")
var isAlnum = PerlPropsLoader("isAlnum")
var isSc = PerlPropsLoader("isSc")
var isSo = PerlPropsLoader("isSo")
var isAlpha = PerlPropsLoader("IsAlpha")
var isLower = PerlPropsLoader("IsLower")

// Remove ASCII junk.
var deduplicateSpace = NewReplacement(`\s+`, ` `)
var asciiJunk = NewReplacement(`[\000-\037]`, ``)

// Neurotic Perl heading space, multi-space and trailing space chomp.
// These regexes are kept for reference purposes and shouldn't be used!!
var midStrip = NewReplacement(` +`, ` `)
var leftStrip = NewReplacement(`^ `, ``)
var rightStrip = NewReplacement(` $`, ``)

// Pad all "other" special characters not in isAlnum.
var padNotisAlnum = NewReplacement(fmt.Sprintf("([^%s\\s\\.'\\`\\,\\-])", isAlnum), ` $1 `)

// Splits all hyphens (regardless of circumstances), e.g.
// 'foo-bar' -> 'foo @-@ bar'
var aggressiveHyphenSplit = NewReplacement(fmt.Sprintf(`([%s])\-([%s])`, isAlnum, isAlnum), "$1 @-@ $2")

// Make multi-dots stay together.
var replaceDotWithLiteralstring1 = NewReplacement(`\.([\.]+)`, ` DOTMULTI$1`)
var replaceDotWithLiteralstring2 = NewReplacement(`DOTMULTI\.([^\.])`, `DOTDOTMULTI $1`)
var replaceDotWithLiteralstring3 = NewReplacement(`DOTMULTI\.`, `DOTDOTMULTI`)

// Separate out "," except if within numbers (5,300)
// e.g.  A,B,C,D,E > A , B,C , D,E
// First application uses up B so rule can't see B,C
// two-step version here may create extra spaces but these are removed later
// will also space digit,letter or letter,digit forms (redundant with next section)
var commaSeparate1 = NewReplacement(fmt.Sprintf(`([^%s])[,]`, isN), `$1 , `)
var commaSeparate2 = NewReplacement(fmt.Sprintf(`[,]([^%s])`, isN), ` , $1`)
var commaSeparate3 = NewReplacement(fmt.Sprintf(`([%s])[,]$`, isN), `$1 , `)
var commaSeparate = []Replacement{
	commaSeparate1,
	commaSeparate2,
	commaSeparate3,
}

// Attempt to get correct directional quotes.
var directionalQuote1 = NewReplacement("^``", "`` ")
var directionalQuote2 = NewReplacement(`^"`, "`` ")
var directionalQuote3 = NewReplacement("^\\`([^\\`])", "` $1")
var directionalQuote4 = NewReplacement("^'", "`  ")
var directionalQuote5 = NewReplacement("([ ([{<])\"", "$1 `` ")
var directionalQuote6 = NewReplacement("([ ([{<])``", "$1 `` ")
var directionalQuote7 = NewReplacement("([ ([{<])`([^`])", "$1 ` $2")
var directionalQuote8 = NewReplacement("([ ([{<])'", "$1 ` ")

// Replace ... with _ELLIPSIS_
var replaceEllipsis = NewReplacement(`\.\.\.`, ` _ELLIPSIS_ `)

// Restore _ELLIPSIS_ with ...
var restoreEllipsis = NewReplacement(`_ELLIPSIS_`, `\.\.\.`)

// Pad , with tailing space except if within numbers, e.g. 5,300
var comma1 = NewReplacement(fmt.Sprintf(`([^%s])[,]([^%s])`, isN, isN), `$1 , $2`)
var comma2 = NewReplacement(fmt.Sprintf(`([%s])[,]([^%s])`, isN, isN), `$1 , $2`)
var comma3 = NewReplacement(fmt.Sprintf(`([^%s])[,]([%s])`, isN, isN), `$1 , $2`)

// Pad unicode symbols with spaces.
var symbols = NewReplacement(fmt.Sprintf("([;:@#$%%&%s%s])", isSc, isSo), ` $1 `)

// Separate out intra-token slashes.  PTB tokenization doesn't do this, so
// the tokens should be merged prior to parsing with a PTB-trained parser.
// e.g. "and/or" -> "and @/@ or"
var intratokenSlashes = NewReplacement(fmt.Sprintf(`([%s])\/([%s])`, isAlnum, isAlnum), `$1 \@\/\@ $2`)

// Splits final period at end of string.
var finalPeriod = NewReplacement(`([^.])([.])([\]\)}>"']*) ?$`, `$1 $2$3`)

// Pad all question marks and exclamation marks with spaces.
var padQuestionExclamationMark = NewReplacement(`([?!])`, ` $1 `)

// Handles parentheses, brackets and converts them to PTB symbols.
var padParenthesis = NewReplacement(`([\]\[\(\){}<>])`, ` $1 `)
var convertParenthesis1 = NewReplacement(`\(`, `-LRB-`)
var convertParenthesis2 = NewReplacement(`\)`, `-RRB-`)
var convertParenthesis3 = NewReplacement(`\[`, `-LSB-`)
var convertParenthesis4 = NewReplacement(`\]`, `-RSB-`)
var convertParenthesis5 = NewReplacement(`\{`, `-LCB-`)
var convertParenthesis6 = NewReplacement(`\}`, `-RCB-`)

// Pads double dashes with spaces.
var padDoubleDashes = NewReplacement(`--`, ` -- `)

// Adds spaces to start and end of string to simplify further regexps.
var padStartOfStr = NewReplacement(`^`, ` `)
var padEndOfStr = NewReplacement(`$`, ` `)

// Converts double quotes to two single quotes and pad with spaces.
var convertDoubleToSingleQuotes = NewReplacement(`"`, ` '' `)

// Handles single quote in possessives or close-single-quote.
var handlesSingleQuotes = NewReplacement(`([^'])' `, `$1 ' `)

// Pad apostrophe in possessive or close-single-quote.
var apostrophe = NewReplacement(`([^'])'`, `$1 ' `)

// Prepend space on contraction apostrophe.
var contraction1 = NewReplacement(`'([sSmMdD]) `, ` '$1 `)
var contraction2 = NewReplacement(`'ll `, ` 'll `)
var contraction3 = NewReplacement(`'re `, ` 're `)
var contraction4 = NewReplacement(`'ve `, ` 've `)
var contraction5 = NewReplacement(`n't `, ` n't `)
var contraction6 = NewReplacement(`'LL `, ` 'LL `)
var contraction7 = NewReplacement(`'RE `, ` 'RE `)
var contraction8 = NewReplacement(`'VE `, ` 'VE `)
var contraction9 = NewReplacement(`N'T `, ` N'T `)

// Informal Contractions.
var contraction10 = NewReplacement(` ([Cc])annot `, ` $1an not `)
var contraction11 = NewReplacement(` ([Dd])'ye `, ` $1' ye `)
var contraction12 = NewReplacement(` ([Gg])imme `, ` $1im me `)
var contraction13 = NewReplacement(` ([Gg])onna `, ` $1on na `)
var contraction14 = NewReplacement(` ([Gg])otta `, ` $1ot ta `)
var contraction15 = NewReplacement(` ([Ll])emme `, ` $1em me `)
var contraction16 = NewReplacement(` ([Mm])ore'n `, ` $1ore 'n `)
var contraction17 = NewReplacement(` '([Tt])is `, ` '$1 is `)
var contraction18 = NewReplacement(` '([Tt])was `, ` '$1 was `)
var contraction19 = NewReplacement(` ([Ww])anna `, ` $1an na `)

// Clean out extra spaces
var cleanExtraSpace1 = NewReplacement(`  *`, ` `)
var cleanExtraSpace2 = NewReplacement(`^ *`, ``)
var cleanExtraSpace3 = NewReplacement(` *$`, ``)

// Neurotic Perl regexes to escape special characters.
var escapeAmpersand = NewReplacement(`&`, `&amp;`)
var escapePipe = NewReplacement(`\|`, `&#124;`)
var escapeLeftAngleBracket = NewReplacement(`<`, `&lt;`)
var escapeRightAngleBracket = NewReplacement(`>`, `&gt;`)
var escapeSingleQuote = NewReplacement(`\'`, `&apos;`)
var escapeDoubleQuote = NewReplacement(`\"`, `&quot;`)
var escapeLeftSquareBracket = NewReplacement(`\[`, `&#91;`)
var escapeRightSquareBracket = NewReplacement(`]`, `&#93;`)

var enSpecific1 = NewReplacement(fmt.Sprintf(`([^%s])[']([^%s])`, isAlpha, isAlpha), `$1 ' $2`)
var enSpecific2 = NewReplacement(fmt.Sprintf(`([^%s%s])[']([%s])`, isAlpha, isN, isAlpha), `$1 ' $2`)
var enSpecific3 = NewReplacement(fmt.Sprintf(`([%s])[']([^%s])`, isAlpha, isAlpha), `$1 ' $2`)
var enSpecific4 = NewReplacement(fmt.Sprintf(`([%s])[']([%s])`, isAlpha, isAlpha), `$1 '$2`)
var enSpecific5 = NewReplacement(fmt.Sprintf(`([%s])[']([s])`, isN), `$1 '$2`)

var englishSpecificApostrophe = []Replacement{
	enSpecific1,
	enSpecific2,
	enSpecific3,
	enSpecific4,
	enSpecific5,
}

var frItSpecific1 = NewReplacement(fmt.Sprintf(`([^%s])[']([^%s])`, isAlpha, isAlpha), `$1 ' $2`)
var frItSpecific2 = NewReplacement(fmt.Sprintf(`([^%s])[']([%s])`, isAlpha, isAlpha), `$1 ' $2`)
var frItSpecific3 = NewReplacement(fmt.Sprintf(`([%s])[']([^%s])`, isAlpha, isAlpha), `$1 ' $2`)
var frItSpecific4 = NewReplacement(fmt.Sprintf(`([%s])[']([%s])`, isAlpha, isAlpha), `$1' $2`)

var frItSpecificApostrophe = []Replacement{
	frItSpecific1,
	frItSpecific2,
	frItSpecific3,
	frItSpecific4,
}

var nonSpecificApostrophe = NewReplacement(`\'`, ` ' `)

var trailingDotApostrophe = NewReplacement("\\.' ?$", " . ' ")

var mosesPennRegexes1 = []Replacement{
	deduplicateSpace,
	asciiJunk,
	directionalQuote1,
	directionalQuote2,
	directionalQuote3,
	directionalQuote4,
	directionalQuote5,
	directionalQuote6,
	directionalQuote7,
	directionalQuote8,
	replaceEllipsis,
	comma1,
	comma2,
	comma3,
	symbols,
	intratokenSlashes,
	finalPeriod,
	padQuestionExclamationMark,
	padParenthesis,
	convertParenthesis1,
	convertParenthesis2,
	convertParenthesis3,
	convertParenthesis4,
	convertParenthesis5,
	convertParenthesis6,
	padDoubleDashes,
	padStartOfStr,
	padEndOfStr,
	convertDoubleToSingleQuotes,
	handlesSingleQuotes,
	apostrophe,
	contraction1,
	contraction2,
	contraction3,
	contraction4,
	contraction5,
	contraction6,
	contraction7,
	contraction8,
	contraction9,
	contraction10,
	contraction11,
	contraction12,
	contraction13,
	contraction14,
	contraction15,
	contraction16,
	contraction17,
	contraction18,
	contraction19,
}

var mosesPennRegexes2 = []Replacement{
	restoreEllipsis,
	cleanExtraSpace1,
	cleanExtraSpace2,
	cleanExtraSpace3,
	escapeAmpersand,
	escapePipe,
	escapeLeftAngleBracket,
	escapeRightAngleBracket,
	escapeSingleQuote,
	escapeDoubleQuote,
}

var mosesEscapeXMLRegexes = []Replacement{
	escapeAmpersand,
	escapePipe,
	escapeLeftAngleBracket,
	escapeRightAngleBracket,
	escapeSingleQuote,
	escapeDoubleQuote,
	escapeLeftSquareBracket,
	escapeRightSquareBracket,
}

// Tokenizer is an instance to tokenize text. This is a golang port
// of Moses Tokenizer from https://github.com/moses-smt/mosesdecoder/blob/master/scripts/tokenizer/tokenizer.perl
// Designs are mostly copied from the python version https://github.com/alvations/sacremoses/blob/master/sacremoses/tokenize.py
type Tokenizer struct {
	lang                string
	nonBreakingPrefixes []string
	numericOnlyPrefixes []string
}

func (t Tokenizer) replaceMultidots(text string) string {
	rgx := regexp.MustCompile(`\.([\.]+)`)
	rgx1 := regexp.MustCompile(`DOTMULTI\.([^\.])`)
	rgx2 := regexp.MustCompile(`DOTMULTI\.`)

	text = rgx.ReplaceAllString(text, " DOTMULTI$1")
	for matched, _ := regexp.MatchString(`DOTMULTI\.`, text); matched; {
		text = rgx1.ReplaceAllString(text, "DOTDOTMULTI $1")
		text = rgx2.ReplaceAllString(text, "DOTDOTMULTI")

		matched, _ = regexp.MatchString(`DOTMULTI\.`, text)
	}
	return text
}

func (t Tokenizer) restoreMultidots(text string) string {
	rgx := regexp.MustCompile(`DOTDOTMULTI`)
	rgx1 := regexp.MustCompile(`DOTMULTI`)
	for rgx.MatchString(text) {
		text = rgx.ReplaceAllString(text, `DOTMULTI.`)
	}
	text = rgx1.ReplaceAllString(text, `.`)
	return text
}

func (t Tokenizer) handlesNonbreakingPrefixes(text string) string {
	tokens := RemoveEmptyStringFromSlice(strings.Split(text, " "))
	tokenEndsWithPeriodRgx := regexp.MustCompile(`^(\S+)\.$`)
	numTokens := len(tokens)

	for i := 0; i < numTokens; i++ {
		if tokenEndsWithPeriodRgx.MatchString(tokens[i]) {
			prefix := tokenEndsWithPeriodRgx.FindStringSubmatch(tokens[i])[1]
			// 1st condition:
			// Adding unconditional extra split in final dots
			// https://github.com/moses-smt/mosesdecoder/pull/204
			//
			// 2nd condition:
			// Checks for 3 conditions if
			// i.   the prefix contains a fullstop and
			//      any char in the prefix is within the IsAlpha charset
			// ii.  the prefix is in the list of nonbreaking prefixes and
			//      does not contain #NUMERIC_ONLY#
			// iii. the token is not the last token and that the
			//      next token contains all lowercase.
			//
			// 3rd condition:
			// Checks if the prefix is in NUMERIC_ONLY_PREFIXES
			// and ensures that the next word is a digit.
			if i == numTokens-1 {
				tokens[i] = prefix + " ."
			} else if (strings.Contains(prefix, ".") && IsAnyAlphabet(prefix)) ||
				(IsInArray(prefix, t.nonBreakingPrefixes) && !IsInArray(prefix, t.numericOnlyPrefixes)) ||
				(i != numTokens-1 && i+1 < numTokens && IsLower(string(tokens[i+1][0]))) {
				// Do nothing
			} else if IsInArray(prefix, t.numericOnlyPrefixes) && i+1 < numTokens && IsNumber(tokens[i+1]) {
				// Do nothing
			} else {
				tokens[i] = prefix + " ."
			}
		}
	}
	return strings.Join(tokens, " ")
}

// NewTokenizer creates new Tokenizer instance with predefined language
func NewTokenizer(lang string) (tokenizer *Tokenizer) {
	nonBreakingPrefixes := NonBreakingPrefixesLoader(lang)
	var numericOnlyPrefixes []string
	for _, prefix := range nonBreakingPrefixes {
		if matched, _ := regexp.MatchString(`(.*)[\s]+(\#NUMERIC_ONLY\#)`, prefix); matched {
			thisPrefix := strings.Split(prefix, " ")[0]
			numericOnlyPrefixes = append(numericOnlyPrefixes, thisPrefix)
		}
	}
	tokenizer = &Tokenizer{lang, nonBreakingPrefixes, numericOnlyPrefixes}
	return
}

// Tokenize incoming string in accordance to predefined language option.
// We can choose to enable more aggresive dash splitting such as "foo-bar" to "foo @-@ bar"
// and escaping XML tags
func (t Tokenizer) Tokenize(text string, aggresiveDashSplits bool, escapeXML bool) (string, []string) {
	text = deduplicateSpace.rgx.ReplaceAllString(text, deduplicateSpace.sub)
	text = asciiJunk.rgx.ReplaceAllString(text, asciiJunk.sub)

	// TODO: need to implement protected_patterns
	// https://github.com/alvations/sacremoses/blob/ce4703ba4c53f6c53bcfb59bf398de1cfdd827af/sacremoses/tokenize.py#L441

	text = strings.TrimSpace(text)
	text = padNotisAlnum.rgx.ReplaceAllString(text, padNotisAlnum.sub)

	if aggresiveDashSplits {
		text = aggressiveHyphenSplit.rgx.ReplaceAllString(text, aggressiveHyphenSplit.sub)
	}

	text = t.replaceMultidots(text)

	for _, csr := range commaSeparate {
		text = csr.rgx.ReplaceAllString(text, csr.sub)
	}

	if t.lang == "en" {
		for _, re := range englishSpecificApostrophe {
			text = re.rgx.ReplaceAllString(text, re.sub)
		}
	} else if t.lang == "it" || t.lang == "fr" {
		for _, re := range frItSpecificApostrophe {
			text = re.rgx.ReplaceAllString(text, re.sub)
		}
	} else {
		text = nonSpecificApostrophe.rgx.ReplaceAllString(text, nonSpecificApostrophe.sub)
	}
	text = t.handlesNonbreakingPrefixes(text)
	// Cleans up extraneous spaces
	text = deduplicateSpace.rgx.ReplaceAllString(text, deduplicateSpace.sub)
	// Split trailing "."
	text = trailingDotApostrophe.rgx.ReplaceAllString(text, trailingDotApostrophe.sub)

	// TODO: need to implement protected_patterns

	text = t.restoreMultidots(text)

	if escapeXML {
		for _, r := range mosesEscapeXMLRegexes {
			text = r.rgx.ReplaceAllString(text, r.sub)
		}
	}
	text = strings.TrimSpace(text)

	return text, strings.Split(text, " ")
}
