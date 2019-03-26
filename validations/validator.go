package validation

import (
	"log"
	"strings"
)

type Validator struct {
}

func NewValidator() (*Validator, error) {
	return &Validator{}, nil
}

func (r *Validator) Validate(oo []string) {
	for i, o := range oo {
		if wronSimbol := r.ValidateOne(o); len(wronSimbol) != 0 {
			for _, s := range wronSimbol {
				log.Printf("olig %d '%s' contains invalid character '%s' ", i+1, oo[i], s)
			}
		}
	}
}

func (r *Validator) ValidateOne(o string) (n []string) {
	// make all validations for olig
	wronSimbol := make([]string, 0)
	dna := strings.Split(o, ",")[1]
	dnaU := strings.ToUpper(dna)
	voc := []string{"A", "C", "G", "T", "R", "Y", "K", "M", "S", "W", "B", "D", "H", "V", "N"} // voc - слайс содержащий допустимые симолы он может быть константой и он используется в нескольких функциях
	for _, r := range dnaU {
		count := 0
		s := string(r)
		for _, o := range voc { // voc - слайс содержащий допустимые симолы он может быть константой и он используется в нескольких функциях
			c := strings.Count(s, o)
			count += c
		}
		if count == 0 {
			wronSimbol = append(wronSimbol, string(r))
			//return fmt.Errorf("incorrect content for olig %s", o)
			// В этом месте должна быть ошибка НЕДОПУСТИМЫЙ СИМВОЛ с указанием недопустимого символа, его номера в последовательности
		}
	}

	return wronSimbol
}

//func (r *Validator) Measure(oo []string, ignoreMode bool) (dA, dC, dG, dT float32, err error) {
//if !ignoreMode {
//	if err, _ := r.Validate(oo); err != nil {

//var dA, dC, dG, dT float32 = 0, 0, 0, 0
//voc := []string{"A", "C", "G", "T", "R", "Y", "K", "M", "S", "W", "B", "D", "H", "V", "N"}
//for _, o := range oo {
//	var oSeq = ""
//	oSeq = strings.Split(o, ",")[1]
//
//	dnaU := strings.ToUpper(oSeq)
//	count := 0
//	if count < len(dnaU) {
//		for _, o := range voc {
//			c := strings.Count(dnaU, string(o))
//			if c > 0 {
//				cf := float32(c)
//				switch o {
//
//				case "A":
//					dA += cf
//				case "C":
//					dC += cf
//				case "G":
//					dG += cf
//				case "T":
//					dT += cf
//				case "R":
//					dA += cf / 2
//					dG += cf / 2
//				case "Y":
//					dC += cf / 2
//					dT += cf / 2
//				case "K":
//					dG += cf / 2
//					dT += cf / 2
//				case "M":
//					dA += cf / 2
//					dC += cf / 2
//				case "S":
//					dG += cf / 2
//					dC += cf / 2
//				case "W":
//					dA += cf / 2
//					dT += cf / 2
//				case "B":
//					dC += cf / 3
//					dG += cf / 3
//					dT += cf / 3
//				case "D":
//					dA += cf / 3
//					dG += cf / 3
//					dT += cf / 3
//				case "H":
//					dA += cf / 3
//					dC += cf / 3
//					dG += cf / 3
//				case "V":
//					dA += cf / 3
//					dC += cf / 3
//					dG += cf / 3
//				case "N":
//					dA += cf / 4
//					dC += cf / 4
//					dG += cf / 4
//					dT += cf / 4
//				}
//			}
//			count += c
//		}
//		return 0, 0, 0, 0, err
//	}
//}
//}
///
///
///
//return 0, 0, 0, 0, nil
//}
//return dA, dC, dG, dT, nil
//}
