package validation

import (
	"strings"
)

const (
	validNotations = "ACGTRYKMSWBDHVN"
)

type Validator struct {
}

type Statistic struct {
	a, c, g, t, r, y, k, m, s, w, b, d, h, v, n, allLinks, allOligs, wrongSymbols int
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
func (r *Statistic) Measure(oo []string) (synthesisStatistic Statistic) {
	for _, o := range oo {
		statisticOne := r.MeasureOne(o)
		synthesisStatistic += statisticOne
	}
	return synthesisStatistic
}

func (r *Statistic) MeasureOne(o string) (statisticOne Statistic) {
	dna := strings.ToUpper(ExtractDna(o))
	count := 0
	for _, o := range validNotations {
		c := strings.Count(dna, string(o))
		if c > 0 {
			switch string(o) {
			case "A":
				statisticOne.a += c
			case "C":
				statisticOne.c += c
			case "G":
				statisticOne.g += c
			case "T":
				statisticOne.t += c
			case "R":
				statisticOne.r += c
				//					dA += cf / 2
				//					dG += cf / 2
			case "Y":
				statisticOne.y += c
				//					dC += cf / 2
				//					dT += cf / 2
			case "K":
				statisticOne.k += c
				//					dG += cf / 2
				//					dT += cf / 2
			case "M":
				statisticOne.m += c
				//					dA += cf / 2
				//					dC += cf / 2
			case "S":
				statisticOne.s += c
				//					dG += cf / 2
				//					dC += cf / 2
			case "W":
				statisticOne.w += c
				//					dA += cf / 2
				//					dT += cf / 2
			case "B":
				statisticOne.b += c
				//					dC += cf / 3
				//					dG += cf / 3
				//					dT += cf / 3
			case "D":
				statisticOne.d += c
				//					dA += cf / 3
				//					dG += cf / 3
				//					dT += cf / 3
			case "H":
				statisticOne.h += c
				//					dA += cf / 3
				//					dC += cf / 3
				//					dG += cf / 3
			case "V":
				statisticOne.v += c
				//					dA += cf / 3
				//					dC += cf / 3
				//					dG += cf / 3
			case "N":
				statisticOne.n += c
				//					dA += cf / 4
				//					dC += cf / 4
				//					dG += cf / 4
				//					dT += cf / 4
			}
		}
		count += c
	}

	statisticOne.wrongSymbols += len(dna) - count
	statisticOne.allLinks += len(dna)

	return statisticOne
}
