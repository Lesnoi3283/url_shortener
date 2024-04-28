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

	_, err = toRet.store.Exec(`
        CREATE TABLE IF NOT EXISTS urls_table (
            id SERIAL PRIMARY KEY,
            long VARCHAR(2048) UNIQUE,
            short VARCHAR(255),
            user_id INT
        );`)
	if err != nil {
		return nil, fmt.Errorf("postgres exec (create urls_table): %w", err)
	}

	_, err = toRet.store.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY
        );`)
	if err != nil {
		return nil, fmt.Errorf("postgres exec (create users): %w", err)
	}

	return toRet, nil
}

// Хорошая ли идея использовать тут скомпилированные запросы? В NewPostgresql их создать,
// в структуре постгрес сохранить и использовать в Save и Get. Потокобезопасно ли это?
func (p *Postgresql) Save(ctx context.Context, short string, full string) error {
	query := "INSERT INTO urls_table (long, short) VALUES ($1, $2) ON CONFLICT (long) DO NOTHING;"

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
		query2 := "SELECT short FROM urls_table WHERE long = $1;"
		row := p.store.QueryRowContext(ctx, query2, full)

		err = row.Scan(&shortURL)
		if err != nil {
			return fmt.Errorf("postgres query: %w", err)
		}
		return NewAlreadyExistsError(shortURL)
	}

	return nil
}

func (p *Postgresql) SaveWithUserID(ctx context.Context, userID int, short string, full string) error {
	query := "INSERT INTO urls_table (user_id, long, short) VALUES ($1, $2, $3) ON CONFLICT (long) DO NOTHING;"

	result, err := p.store.ExecContext(ctx, query, userID, full, short)
	if err != nil {
		return fmt.Errorf("postgres execute: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		shortURL := ""
		query2 := "SELECT short FROM urls_table WHERE long = $1;"
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
	query := "INSERT INTO urls_table (long, short) VALUES ($1, $2);"

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

func (p *Postgresql) SaveBatchWithUserID(ctx context.Context, userID int, urls []entities.URL) error {
	tx, err := p.store.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("postgres transaction start: %w", err)
	}
	query := "INSERT INTO urls_table (user_id, long, short) VALUES ($1, $2, $3);"

	for _, url := range urls {
		_, err = tx.ExecContext(ctx, query, userID, url.Long, url.Short)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("postgres, transaction error: %w", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("postgres, transaction commit: %w", err)
	}

	return nil
}

func (p *Postgresql) Get(ctx context.Context, short string) (full string, err error) {

	query := "SELECT long FROM urls_table WHERE short = $1;"
	row := p.store.QueryRowContext(ctx, query, short)

	err = row.Scan(&full)
	if err != nil {
		return "", fmt.Errorf("postgres query: %w", err)
	}
	return full, nil
}

func (p *Postgresql) GetUserUrls(ctx context.Context, userID int) ([]struct {
	Long  string
	Short string
}, error) {
	query := "SELECT long, short FROM urls_table WHERE user_id = $1;"

	var urls []struct {
		Long  string
		Short string
	}

	rows, err := p.store.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("postgres query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var url struct {
			Long  string
			Short string
		}
		if err := rows.Scan(&url.Long, &url.Short); err != nil {
			return nil, fmt.Errorf("postgres row scan: %w", err)
		}
		urls = append(urls, url)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("postgres rows iteration: %w", err)
	}

	return urls, nil
}

func (p *Postgresql) Ping() error {
	return p.store.Ping()
}

func (p *Postgresql) Close() error {
	return p.store.Close()
}

func (p *Postgresql) CreateUser(ctx context.Context) (int, error) {
	query := "INSERT INTO users DEFAULT VALUES RETURNING id;"

	var userID int

	err := p.store.QueryRowContext(ctx, query).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("postgres create user: %w", err)
	}

	return userID, nil
}
