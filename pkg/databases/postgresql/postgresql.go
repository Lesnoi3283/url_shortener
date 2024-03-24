package postgresql

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Postgresql struct {
	store *sql.DB
}

func (p *Postgresql) Ping() error {
	return p.store.Ping()
}

func NewPostgresql(connStr string) (Postgresql, error) {
	db, err := sql.Open("pgx", connStr)

	toRet := Postgresql{
		store: db,
	}

	return toRet, err
}
