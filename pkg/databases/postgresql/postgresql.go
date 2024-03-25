package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Postgresql struct {
	store *sql.DB
}

func NewPostgresql(connStr string) (*Postgresql, error) {

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, err
	}

	toRet := &Postgresql{
		store: db,
	}

	_, err = toRet.store.Exec("CREATE TABLE IF NOT EXISTS urls (id SERIAL PRIMARY KEY, long VARCHAR(2048), short VARCHAR(255));")

	return toRet, err
}

func (p *Postgresql) Save(ctx context.Context, short string, full string) error {
	query := "INSERT INTO urls (long, short) VALUES ($1, $2);"

	ping := p.store == nil
	fmt.Print(ping)
	_, err := p.store.ExecContext(ctx, query, full, short)
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgresql) Get(ctx context.Context, short string) (full string, err error) {

	query := "SELECT long FROM urls WHERE short = $1;"
	row := p.store.QueryRowContext(ctx, query, short)

	err = row.Scan(&full)
	if err != nil {
		return "", err
	}

	return full, nil
}

func (p *Postgresql) Ping() error {
	return p.store.Ping()
}

func (p *Postgresql) Close() error {
	return p.store.Close()
}
