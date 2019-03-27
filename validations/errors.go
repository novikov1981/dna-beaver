package validation

import (
	"fmt"
	"strings"
)

type IncorrectOligSymbolsError struct {
	oligNumber int
	olig       string
	positions  []int
}

func NewIncorrectOligSymbolsError(olig string, oligNumber int) IncorrectOligSymbolsError {
	err := IncorrectOligSymbolsError{olig: olig, oligNumber: oligNumber}
	return err
}

func (e *IncorrectOligSymbolsError) AddErrorPosition(n int) {
	e.positions = append(e.positions, n)
}

func (e *IncorrectOligSymbolsError) Error() string {
	sb := strings.Builder{}
	sb.WriteString(
		fmt.Sprintf("Validation error: olig number %d '%s' contains wrong symbols\n", e.oligNumber+1, e.olig))
	for _, pos := range e.positions {
		sb.WriteString(fmt.Sprintf("position %d character '%s'\n", pos+1, string(ExtractDna(e.olig)[pos])))
	}
	return sb.String()
}

func (e *IncorrectOligSymbolsError) Empty() bool {
	return len(e.positions) == 0
}

type IncorrectSynthesisError struct {
	errs []error
}

func NewIncorrectSynthesisError() IncorrectSynthesisError {
	err := IncorrectSynthesisError{errs: []error{}}
	return err
}

func (e *IncorrectSynthesisError) AddOligError(err error) {
	e.errs = append(e.errs, err)
}

func (e *IncorrectSynthesisError) Error() string {
	sb := strings.Builder{}
	sb.WriteString("Synthesis contains errors:\n")
	for _, err := range e.errs {
		sb.WriteString(err.Error())
	}
	return sb.String()
}

func (e *IncorrectSynthesisError) Empty() bool {
	return len(e.errs) == 0
}
