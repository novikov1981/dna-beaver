package sqllite

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/novikov1981/dna-beaver"
	"github.com/satori/go.uuid"
	"log"
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
	if _, err := os.Stat("./synthesis.db"); err != nil {
		if os.IsNotExist(err) {
			return false
		}
		log.Fatalf("error: file stat problems: %s", err.Error())
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
	// TODO: implement this
	return nil, nil
}
