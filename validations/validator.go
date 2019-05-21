package validation

import (
	. "github.com/novikov1981/dna-beaver"
	"github.com/novikov1981/dna-beaver/utils"
	"strings"
)

func Validate(oo []string) error {
	verr := NewIncorrectSynthesisError()
	for i, o := range oo {
		if err := ValidateOne(o, i); err != nil {
			verr.AddOligError(err)
		}
	}
	if !verr.Empty() {
		return &verr
	}
	return nil
}

func ValidateOne(o string, p int) error {
	verr := NewIncorrectOligSymbolsError(o, p)
	dna := utils.ExtractDna(o)
	for p, a := range dna {
		if strings.Index(ValidNotations, string(a)) == -1 {
			verr.AddErrorPosition(p)
		}
	}
	if !verr.Empty() {
		return &verr
	}
	return nil
}
