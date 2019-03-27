package validation

import (
	"strings"
)

const (
	validNotations = "ACGTRYKMSWBDHVN"
)

type Validator struct {
}

func NewValidator() (*Validator, error) {
	return &Validator{}, nil
}

func (r *Validator) Validate(oo []string) error {
	verr := NewIncorrectSynthesisError()
	for i, o := range oo {
		if err := r.ValidateOne(o, i); err != nil {
			verr.AddOligError(err)
		}
	}
	if !verr.Empty() {
		return &verr
	}
	return nil
}

func (r *Validator) ValidateOne(o string, p int) error {
	verr := NewIncorrectOligSymbolsError(o, p)
	dna := strings.ToUpper(ExtractDna(o))
	for p, a := range dna {
		if strings.Index(validNotations, string(a)) == -1 {
			verr.AddErrorPosition(p)
		}
	}
	if !verr.Empty() {
		return &verr
	}
	return nil
}

func (r *Validator) Measure(oo []string) (x map[string]int) {

	seqMap := make(map[string]int)
	trueLinks := []string{"A", "C", "G", "T", "R", "Y", "K", "M", "S", "W", "B", "D", "H", "V", "N"}
	for _, o := range oo {
		var dna = ""
		dna = strings.Split(o, ",")[1]
		dnaU := strings.ToUpper(dna)
		count := 0

		for _, o := range trueLinks {
			c := strings.Count(dnaU, string(o))
			if c > 0 {
				switch o {
				case "A":
					seqMap["A"] += c
				case "C":
					seqMap["C"] += c
				case "G":
					seqMap["G"] += c
				case "T":
					seqMap["T"] += c
				case "R":
					seqMap["R"] += c
					//					dA += cf / 2
					//					dG += cf / 2
				case "Y":
					seqMap["Y"] += c
					//					dC += cf / 2
					//					dT += cf / 2
				case "K":
					seqMap["K"] += c
					//					dG += cf / 2
					//					dT += cf / 2
				case "M":
					seqMap["M"] += c
					//					dA += cf / 2
					//					dC += cf / 2
				case "S":
					seqMap["S"] += c
					//					dG += cf / 2
					//					dC += cf / 2
				case "W":
					seqMap["W"] += c
					//					dA += cf / 2
					//					dT += cf / 2
				case "B":
					seqMap["B"] += c
					//					dC += cf / 3
					//					dG += cf / 3
					//					dT += cf / 3
				case "D":
					seqMap["D"] += c
					//					dA += cf / 3
					//					dG += cf / 3
					//					dT += cf / 3
				case "H":
					seqMap["H"] += c
					//					dA += cf / 3
					//					dC += cf / 3
					//					dG += cf / 3
				case "V":
					seqMap["V"] += c
					//					dA += cf / 3
					//					dC += cf / 3
					//					dG += cf / 3
				case "N":
					seqMap["N"] += c
					//					dA += cf / 4
					//					dC += cf / 4
					//					dG += cf / 4
					//					dT += cf / 4
				}
			}
			count += c
		}

		//}
		seqMap["wronSimbol"] += len(dna) - count
		seqMap["allLinks"] += len(dna)
		seqMap["allOligs"] += 1
	}

	return seqMap
}

///
///
///
//return 0, 0, 0, 0, nil
//}
//return dA, dC, dG, dT, nil
//}
