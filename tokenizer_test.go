package gmt

import (
	"reflect"
	"testing"
)

func TestTokenizer_Tokenize(t *testing.T) {
	tokenizer := NewTokenizer("en")

	type args struct {
		text                string
		aggresiveDashSplits bool
		escape              bool
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 []string
	}{
		{"testWeirdSymbols", args{"This, weird\xbb symbols\u2026 appearing everywhere\xbf", false, true},
			"This , weird \xbb symbols \u2026 appearing everywhere \xbf", nil},
		{"testNonbreakingFullStop", args{"abc,def.", false, true}, "abc , def .", []string{"abc", ",", "def", "."}},
		{"testNonBreakingNumericOnlyLastToken", args{"2016, pp.", false, true}, "2016 , pp .", []string{"2016", ",", "pp", "."}},
		{"testEscapeXML", args{"This ain't funny. It's actually hillarious, yet double Ls. | [] < > [ ] & You're gonna shake it off? Don't?", false, true},
			"", []string{
				"This",
				"ain",
				"&apos;t",
				"funny",
				".",
				"It",
				"&apos;s",
				"actually",
				"hillarious",
				",",
				"yet",
				"double",
				"Ls",
				".",
				"&#124;",
				"&#91;",
				"&#93;",
				"&lt;",
				"&gt;",
				"&#91;",
				"&#93;",
				"&amp;",
				"You",
				"&apos;re",
				"gonna",
				"shake",
				"it",
				"off",
				"?",
				"Don",
				"&apos;t",
				"?"}},
		{"testNotEscapeXML", args{"This ain't funny. It's actually hillarious, yet double Ls. | [] < > [ ] & You're gonna shake it off? Don't?", false, false},
			"", []string{
				"This",
				"ain",
				"'t",
				"funny",
				".",
				"It",
				"'s",
				"actually",
				"hillarious",
				",",
				"yet",
				"double",
				"Ls",
				".",
				"|",
				"[",
				"]",
				"<",
				">",
				"[",
				"]",
				"&",
				"You",
				"'re",
				"gonna",
				"shake",
				"it",
				"off",
				"?",
				"Don",
				"'t",
				"?"}},
		{"testAposthrophe", args{"this 'is' the thing", false, true}, "this &apos; is &apos; the thing", []string{"this", "&apos;", "is", "&apos;", "the", "thing"}},
		{"testAggresiveSplit", args{"foo-bar", true, true}, "foo @-@ bar", []string{"foo", "@-@", "bar"}},
		{"testOpeningBrackets", args{"By the mid 1990s a version of the game became a Latvian television series (with a parliamentary setting, and played by Latvian celebrities).", false, true},
			"By the mid 1990s a version of the game became a Latvian television series ( with a parliamentary setting , and played by Latvian celebrities ) .",
			[]string{
				"By",
				"the",
				"mid",
				"1990s",
				"a",
				"version",
				"of",
				"the",
				"game",
				"became",
				"a",
				"Latvian",
				"television",
				"series",
				"(",
				"with",
				"a",
				"parliamentary",
				"setting",
				",",
				"and",
				"played",
				"by",
				"Latvian",
				"celebrities",
				")",
				".",
			}},
		{"testDotSplitting", args{"The meeting will take place at 11:00 a.m. Tuesday.", true, true},
			"The meeting will take place at 11 : 00 a.m. Tuesday .",
			[]string{"The", "meeting", "will", "take", "place", "at", "11", ":", "00", "a.m.", "Tuesday", "."}},
		{"testTrainingDotApostrophe", args{"'Hello.'", true, true},
			"&apos;Hello . &apos;",
			[]string{"&apos;Hello", ".", "&apos;"}},
		{"testTrainingDot", args{"'So am I.", true, true},
			"&apos;So am I .",
			[]string{"&apos;So", "am", "I", "."}},
		{"testTrainingDot2", args{"It's 7 p.m.", true, true},
			"It &apos;s 7 p.m .",
			[]string{"It", "&apos;s", "7", "p.m", "."}},
		{"testFinalCommaSplitAfterNumber", args{"Sie sollten vor dem Upgrade eine Sicherung dieser Daten erstellen (wie unter Abschnitt 4.1.1, „Sichern aller Daten und Konfigurationsinformationen“ beschrieben). ", true, true},
			"Sie sollten vor dem Upgrade eine Sicherung dieser Daten erstellen ( wie unter Abschnitt 4.1.1 , „ Sichern aller Daten und Konfigurationsinformationen “ beschrieben ) .",
			[]string{
				"Sie",
				"sollten",
				"vor",
				"dem",
				"Upgrade",
				"eine",
				"Sicherung",
				"dieser",
				"Daten",
				"erstellen",
				"(",
				"wie",
				"unter",
				"Abschnitt",
				"4.1.1",
				",",
				"„",
				"Sichern",
				"aller",
				"Daten",
				"und",
				"Konfigurationsinformationen",
				"“",
				"beschrieben",
				")",
				".",
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tokenizer.Tokenize(tt.args.text, tt.args.aggresiveDashSplits, tt.args.escape)
			if tt.want != "" {
				if got != tt.want {
					t.Errorf("Tokenizer.Tokenize() got = %v, want %v", got, tt.want)
				}
			}
			if tt.want1 != nil {
				if !reflect.DeepEqual(got1, tt.want1) {
					t.Errorf("Tokenizer.Tokenize() got1 = %v, want %v", got1, tt.want1)
				}
			}
		})
	}
}
