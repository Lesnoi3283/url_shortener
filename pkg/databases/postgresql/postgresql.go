package postgresql

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Postgresql struct {
	store *sql.DB
}

func (p *Postgresql) Ping() error {
	return p.store.Ping()
}

func NewPostgresql(connStr string) (Postgresql, error) {
	ps := fmt.Sprintf(connStr)
	db, err := sql.Open("pgx", ps)

	toRet := Postgresql{
		store: db,
	}

	return toRet, err
}
