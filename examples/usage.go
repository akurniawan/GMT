package main

import (
	"fmt"

	moses "github.com/akurniawan/GMT"
)

func main() {
	normalizer := moses.NewNormalizer("en", true, true, true, false, false)
	fmt.Println(normalizer.Normalize("The United States in 1805 (color map)                 _Facing_     193"))

	tokenizer := moses.NewTokenizer("en")
	fmt.Println(tokenizer.Tokenize("adit, ganteng", false, true))
}
