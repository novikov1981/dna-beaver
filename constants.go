package dna_beaver

const (
	SqlLiteFile  = "./synthesis.db"
	ValidateMode = "validate"
	SaveMode     = "save"
	SearchMode   = "search"
)

const (
	ValidNotations = "ACGTRYKMSWBDHVN"
	DA             = "A"
	DC             = "C"
	DG             = "G"
	DT             = "T"
)

var (
	Amedits = []string{DA, DC, DG, DT}
)
