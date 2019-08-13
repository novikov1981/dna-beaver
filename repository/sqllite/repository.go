package sqllite

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/novikov1981/dna-beaver"
	"github.com/satori/go.uuid"
	"os"
	"time"
)

type Repository struct {
	database *sqlx.DB
}

func NewRepository(dbPath string) (*Repository, error) {
	dbFileExist := checkFileExists(dbPath)
	// connect db
	database, err := sqlx.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	repo := &Repository{database}
	// create tables if file did not exist before
	if !dbFileExist {
		err := repo.create()
		if err != nil {
			return nil, err
		}
	}
	return &Repository{database}, nil
}

func checkFileExists(filePath string) bool {
	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			return false
		}
		fmt.Printf("error: file stat problems: %s\n", err.Error())
		os.Exit(1)
	}
	return true
}

func (r *Repository) create() error {
	// SYNTHESIS
	statement, err := r.database.Prepare(`
	CREATE TABLE IF NOT EXISTS synthesis (
		uuid TEXT PRIMARY KEY,
		name TEXT UNIQUE,
		scale INTEGER CHECK(scale>0),
		created_at TEXT
	);`)
	if err != nil {
		return err
	}
	_, err = statement.Exec()
	if err != nil {
		return err
	}

	// OLIGS
	statement, err = r.database.Prepare(`
		CREATE TABLE IF NOT EXISTS oligs (
			uuid TEXT PRIMARY KEY,
			synthesis_uuid TEXT,
			content TEXT,
			position NUMBER CHECK(position>0),
			FOREIGN KEY(synthesis_uuid) REFERENCES synthesis(uuid)
				);`)
	if err != nil {
		return err
	}

	_, err = statement.Exec()
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) InsertSynthesis(name string, scale int64, oo []string) error {
	createdAt := time.Now().Format("2006-01-02 15:04:05")
	synthesisUUID := uuid.NewV4().String()
	tx, err := r.database.Beginx()
	if err != nil {
		return err
	}
	_, err = tx.Exec(`INSERT INTO synthesis (uuid, name, scale, created_at) VALUES (?, ?, ?, ?)`,
		synthesisUUID, name, scale, createdAt)
	if err != nil {
		tx.Rollback()
		return err
	}
	for i, o := range oo {
		oligUUID := uuid.NewV4().String()
		_, err = tx.Exec(`INSERT INTO oligs (uuid, synthesis_uuid, content, position) VALUES (?, ?, ?, ?)`,
			oligUUID, synthesisUUID, o, i+1)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) FindSynthesis(pattern string) ([]dna_beaver.Synthesis, error) {
	var ss []dna_beaver.Synthesis
	tx, err := r.database.Beginx()
	if err != nil {
		return nil, err
	}
	err = tx.Select(&ss, `
		SELECT * 
		FROM synthesis 
		WHERE EXISTS (
				SELECT 
					1 
				FROM 
					oligs
				WHERE 
					oligs.content LIKE "%" || ? || "%"
					AND 
					oligs.synthesis_uuid = synthesis.uuid
			) OR 
			synthesis.name LIKE "%" || ? || "%"`, pattern, pattern)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	for i, s := range ss {
		var oo []dna_beaver.Olig
		err = tx.Select(&oo, `
			SELECT * 
			FROM oligs 
			WHERE oligs.synthesis_uuid = ?
				AND 
				  oligs.content LIKE "%" || ? || "%"
			ORDER BY oligs.position`, s.Uuid, pattern)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		ss[i].Oligs = oo
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return ss, nil
}

func (r *Repository) GetSynthesis(name string) (*dna_beaver.Synthesis, error) {
	var ss []dna_beaver.Synthesis
	tx, err := r.database.Beginx()
	if err != nil {
		return nil, err
	}
	err = tx.Select(&ss, `
		SELECT * 
		FROM synthesis 
		WHERE synthesis.name = ?`, name)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if len(ss) == 0 {
		return nil, nil
	}
	if len(ss) > 1 {
		return nil, fmt.Errorf("number of results by search on synthesis name is more than 1, synthesis found: %+v", ss)
	}
	synt := ss[0]
	var oo []dna_beaver.Olig
	err = tx.Select(&oo, `
		SELECT * 
		FROM oligs 
		WHERE oligs.synthesis_uuid = ?
    	ORDER BY oligs.position`, synt.Uuid)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	synt.Oligs = oo
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return &synt, nil
}
