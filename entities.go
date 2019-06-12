package dna_beaver

type Synthesis struct {
	Uuid      string `db:"uuid"`
	Name      string `db:"name"`
	Scale     int64  `db:"scale"`
	CreatedAt string `db:"created_at"`
	Oligs     []Olig
}

type Olig struct {
	Uuid          string `db:"uuid"`
	SynthesisUuid string `db:"synthesis_uuid"`
	Content       string `db:"content"`
	Position      int64  `db:"position"`
}
