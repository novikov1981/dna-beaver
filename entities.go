package dna_beaver

type Synthesis struct {
	//these are tags, provide additional information for the field for some library
	Uuid      string `db:"uuid"`
	Name      string `db:"name"`
	Scale     int64  `db:"scale"`
	CreatedAt string `db:"created_at"`
}

type Oligs struct {
	//these are tags, provide additional information for the field for some library
	Uuid          string `db:"uuid"`
	SynthesisUuid string `db:"synthesis_uuid"`
	Content       string `db:"content"`
}
