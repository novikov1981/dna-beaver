package validation

import "strings"

func ExtractDna(olig string) (dna string) {
	return strings.Split(olig, ",")[1]
}
