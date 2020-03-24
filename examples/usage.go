package main

import (
	"fmt"

	moses "github.com/akurniawan/GMT"
)

func main() {
	normalizer := moses.NewNormalizer("en", true, true, true, false, false)
	fmt.Println(normalizer.Normalize("adit ganteng"))
}
