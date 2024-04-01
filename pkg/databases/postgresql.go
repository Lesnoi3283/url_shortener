package databases

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Lesnoi3283/url_shortener/internal/app/entities"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Postgresql struct {
	store *sql.DB
}

func NewPostgresql(connStr string) (*Postgresql, error) {

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("postgres sql open: %w", err)
	}

	toRet := &Postgresql{
		store: db,
	}

	_, err = toRet.store.Exec("CREATE TABLE IF NOT EXISTS urls (id SERIAL PRIMARY KEY, long VARCHAR(2048) UNIQUE, short VARCHAR(255));")

	if err != nil {
		return nil, fmt.Errorf("postgres exec: %w", err)
	}
	return toRet, err
}

// Хорошая ли идея использовать тут скомпилированные запросы? В NewPostgresql их создать,
// в структуре постгрес сохранить и использовать в Save и Get. Потокобезопасно ли это?
func (p *Postgresql) Save(ctx context.Context, short string, full string) error {
	//query := "INSERT INTO urls (long, short) VALUES ($1, $2) ON CONFLICT (long) DO NOTHING;"
	query := "INSERT INTO urls (long, short) VALUES ($1, $2);"

	result, err := p.store.ExecContext(ctx, query, full, short)
	if err != nil {
		return fmt.Errorf("postgres execute: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		shortURL := ""
		query2 := "SELECT short FROM urls WHERE long = $1;"
		row := p.store.QueryRowContext(ctx, query2, full)

		err = row.Scan(&shortURL)
		if err != nil {
			return fmt.Errorf("postgres query: %w", err)
		}
		return NewAlreadyExistsError(shortURL)
	}

	return nil
}

func (p *Postgresql) SaveBatch(ctx context.Context, urls []entities.URL) error {
	tx, err := p.store.Begin()
	if err != nil {
		return fmt.Errorf("postgres transaction start: %w", err)
	}
	query := "INSERT INTO urls (long, short) VALUES ($1, $2);"

	for _, url := range urls {
		_, err = tx.ExecContext(ctx, query, url.Long, url.Short)
		if err != nil {
			tx.Rollback()
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("postgres, transaction commit: %w", err)
	}

	return nil
}

func (p *Postgresql) Get(ctx context.Context, short string) (full string, err error) {

	query := "SELECT long FROM urls WHERE short = $1;"
	row := p.store.QueryRowContext(ctx, query, short)

	err = row.Scan(&full)
	if err != nil {
		return "", fmt.Errorf("postgres query: %w", err)
	}
	return full, nil
}

func (p *Postgresql) Ping() error {
	return p.store.Ping()
}

func (p *Postgresql) Close() error {
	return p.store.Close()
}
