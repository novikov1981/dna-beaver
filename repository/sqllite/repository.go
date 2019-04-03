package sqllite

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/novikov1981/dna-beaver"
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
		name TEXT PRIMARY KEY,
		scale INTEGER,
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
			synthesis_name TEXT,
			content TEXT,
			FOREIGN KEY(synthesis_name) REFERENCES synthesis(name)
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
	//createdAt := ....
	// INSERT to synthesis
	//for _, o := range oo {
	//	oligUUID := generate()
	//	// INSER to olig
	//}
	//
	tx, err := r.database.Beginx()
	if err != nil {
		return err
	}
	_, err = tx.Exec(`INSERT INTO synthesis (name, scale, created_at) VALUES (?, ?, ?)`, name, scale, time.Now())
	if err != nil {
		tx.Rollback()
		return err
	}
	//for _,o := range oo {
	//	_, err = tx.Exec(`INSERT INTO oligs (name, scale, created_at) VALUES (?, ?, ?)`,name,scale,time.Now())
	//	if err != nil {
	//		tx.Rollback()
	//		return err
	//	}
	//}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) FindSynthesis(pattern string) ([]dna_beaver.Synthesis, error) {
	return nil, nil
}
