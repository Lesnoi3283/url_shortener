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

func (p *Postgresql) Close() error {
	return p.store.Close()
}

func NewPostgresql(connStr string) (Postgresql, error) {
	db, err := sql.Open("pgx", connStr)

	toRet := Postgresql{
		store: db,
	}

	if err != nil {
		return toRet, err
	}

	_, err = toRet.store.Exec("CREATE TABLE IF NOT EXISTS urls (id SERIAL PRIMARY KEY, long VARCHAR(2048), short VARCHAR(255));")

	return toRet, err
}
