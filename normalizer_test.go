package gmt

import (
	"fmt"
	"reflect"
	"testing"
)

func newNormalizerDefault(lang string) *Normalizer {
	return NewNormalizer(lang, true, true, true, false, false)
}

func newNormalizerWithoutQuoteNorm(lang string) *Normalizer {
	return NewNormalizer(lang, true, false, true, false, false)
}

func newNormalizerWithoutNumbNorm(lang string) *Normalizer {
	return NewNormalizer(lang, true, true, false, false, false)
}

func newNormalizerPrePost(lang string) *Normalizer {
	return NewNormalizer(lang, true, true, true, true, true)
}

func TestNormalizer_Normalize_Default(t *testing.T) {
	normalizer := newNormalizerDefault("en")
	type args struct {
		text string
	}
	tests := []struct {
		name               string
		args               args
		wantNormalizedText string
	}{
		{"testNormalizeDocuments1", args{"The United States in 1805 (color map)                 _Facing_     193"},
			"The United States in 1805 (color map) _Facing_ 193"},
		{"testNormalizeDocuments2", args{"=Formation of the Constitution.=--(1) The plans before the convention,"},
			"=Formation of the Constitution.=-- (1) The plans before the convention,"},
		{"testNormalizeDocuments3", args{"directions--(1) The infective element must be eliminated. When the ulcer"},
			"directions-- (1) The infective element must be eliminated. When the ulcer"},
		{"testNormalizeDocuments4", args{"College of Surgeons, Edinburgh.)]"},
			"College of Surgeons, Edinburgh.) ]"},
		{"testSingleAposthrophe", args{"yesterday ’s reception"},
			"yesterday 's reception"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotNormalizedText := normalizer.Normalize(tt.args.text); gotNormalizedText != tt.wantNormalizedText {
				t.Errorf("Normalizer.Normalize() = %v, want %v", gotNormalizedText, tt.wantNormalizedText)
			}
		})
	}
}

func TestNormalizer_NormalizeQuoteComma(t *testing.T) {
	normalizers := []Normalizer{
		*newNormalizerWithoutQuoteNorm("en"),
		*newNormalizerDefault("en"),
	}

	type args struct {
		text string
	}
	tests := []struct {
		name               string
		args               args
		wantNormalizedText string
		idx                int
	}{
		{"testWithoutQuote", args{"THIS EBOOK IS OTHERWISE PROVIDED TO YOU \"AS-IS\"."},
			"THIS EBOOK IS OTHERWISE PROVIDED TO YOU \"AS-IS\".", 0},
		{"testWithQuote", args{"THIS EBOOK IS OTHERWISE PROVIDED TO YOU \"AS-IS\"."},
			"THIS EBOOK IS OTHERWISE PROVIDED TO YOU \"AS-IS.\"", 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotNormalizedText := normalizers[tt.idx].Normalize(tt.args.text); gotNormalizedText != tt.wantNormalizedText {
				t.Errorf("Normalizer.Normalize() = %v, want %v", gotNormalizedText, tt.wantNormalizedText)
			}
		})
	}
}

func TestNormalizer_NormalizeNumbers(t *testing.T) {
	normalizers := []Normalizer{
		*newNormalizerWithoutNumbNorm("en"),
		*newNormalizerDefault("en"),
	}

	type args struct {
		text string
	}
	tests := []struct {
		name               string
		args               args
		wantNormalizedText string
		idx                int
	}{
		{"testWithoutNumberNormalization", args{"12\u00A0123"}, "12\u00A0123", 0},
		{"testWithNumberNormalization", args{"12\u00A0123"}, "12.123", 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotNormalizedText := normalizers[tt.idx].Normalize(tt.args.text); gotNormalizedText != tt.wantNormalizedText {
				fmt.Println(reflect.TypeOf(gotNormalizedText), reflect.TypeOf(tt.wantNormalizedText))
				t.Errorf("Normalizer.Normalize() = %v, want %v", gotNormalizedText, tt.wantNormalizedText)
			}
		})
	}
}

func TestNormalizer_NormalizationUnicode(t *testing.T) {
	normalizer := newNormalizerPrePost("en")

	type args struct {
		text string
	}
	tests := []struct {
		name               string
		args               args
		wantNormalizedText string
	}{
		{"testReplaceUnicodePunct", args{"０《１２３》      ４５６％  '' 【７８９】"},
			`0"123" 456% " [789]`},
		{"testReplaceUnicodePunct", args{"０《１２３》 ４５６％ 【７８９】"},
			`0"123" 456% [789]`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotNormalizedText := normalizer.Normalize(tt.args.text); gotNormalizedText != tt.wantNormalizedText {
				fmt.Println(reflect.TypeOf(gotNormalizedText), reflect.TypeOf(tt.wantNormalizedText))
				t.Errorf("Normalizer.Normalize() = %v, want %v", gotNormalizedText, tt.wantNormalizedText)
			}
		})
	}
}
