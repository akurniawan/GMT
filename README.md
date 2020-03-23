# GMT

Golang port of Moses tokenizer and normalizer

You can refer to the following repositories for reference to the original code
1. [Sacremoses](https://github.com/alvations/sacremoses)
2. [mosesdecoder](https://github.com/moses-smt/mosesdecoder)

## Features & Limitation

Currently the port is only for tokenizer and normalizer for english and non-chinese languages. While the original sacremoses has detokenizer and true casing as well, they are not yet currently implemented.


# Install


# Usage

## Tokenizer

```go
tokenizer := NewTokenizer("en")
text := "This, weird\xbb symbols\u2026 appearing everywhere\xbf"
exptected := "This , weird \xbb symbols \u2026 appearing everywhere \xbf"
tokenized := tokenizer.Tokenize(text, false, true)
println(text == expected)
```

## Normalizer

```go
normalizer := NewNormalizer("en", true, true, true, false, false)
text := "12\u00A0123"
exptected := "12.123"
normalized := normalizer.mlizedmmmmmmmalse, true)
println(text == normalized)
```
