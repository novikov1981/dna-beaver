package utils

import "strings"

func ExtractDna(olig string) (dna string) {
	return strings.ToUpper(strings.Split(olig, ",")[1])
}
