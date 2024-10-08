package databases

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Lesnoi3283/url_shortener/internal/app/entities"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// Postgresql is a struct witch has some functions to work with PostgreSQL database.
type Postgresql struct {
	store *sql.DB
}

// build a new Postgresql.
func NewPostgresql(connStr string) (*Postgresql, error) {
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("postgres sql open: %w", err)
	}

	toRet := &Postgresql{
		store: db,
	}

	_, err = toRet.store.Exec(`
    CREATE TABLE IF NOT EXISTS user_urls_table (
        id SERIAL PRIMARY KEY,
        long VARCHAR(2048) UNIQUE,
        short VARCHAR(255),
        user_id INT,
        is_deleted BOOLEAN DEFAULT false
    );
`)

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

// Save saves a new url to a storage.
func (p *Postgresql) Save(ctx context.Context, url entities.URL) error {
	query := "INSERT INTO user_urls_table (long, short) VALUES ($1, $2) ON CONFLICT (long) DO NOTHING;"

	result, err := p.store.ExecContext(ctx, query, url.Long, url.Short)
	if err != nil {
		return fmt.Errorf("postgres execute: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		//в случае если ссылка уже была сохранена ранее
		shortURL := ""
		query2 := "SELECT short FROM user_urls_table WHERE long = $1;"
		row := p.store.QueryRowContext(ctx, query2, url.Long)

		err = row.Scan(&shortURL)
		if err != nil {
			return fmt.Errorf("postgres query: %w", err)
		}
		return NewAlreadyExistsError(shortURL)
	}

	return nil
}

// SaveWithUserID saves a URL with userID.
func (p *Postgresql) SaveWithUserID(ctx context.Context, userID int, url entities.URL) error {
	query := "INSERT INTO user_urls_table (user_id, long, short) VALUES ($1, $2, $3) ON CONFLICT (long) DO NOTHING;"

	result, err := p.store.ExecContext(ctx, query, userID, url.Long, url.Short)
	if err != nil {
		return fmt.Errorf("postgres execute: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		shortURL := ""
		query2 := "SELECT short FROM user_urls_table WHERE long = $1;"
		row := p.store.QueryRowContext(ctx, query2, url.Long)

		err = row.Scan(&shortURL)
		if err != nil {
			return fmt.Errorf("postgres query: %w", err)
		}
		return NewAlreadyExistsError(shortURL)
	}

	return nil
}

// SaveBatch saves a batch of URLs.
func (p *Postgresql) SaveBatch(ctx context.Context, urls []entities.URL) error {
	tx, err := p.store.Begin()
	if err != nil {
		return fmt.Errorf("postgres transaction start: %w", err)
	}
	query := "INSERT INTO user_urls_table (long, short) VALUES ($1, $2);"

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

// SaveBatchWithUserID save a batch of URLs with userID.
func (p *Postgresql) SaveBatchWithUserID(ctx context.Context, userID int, urls []entities.URL) error {
	tx, err := p.store.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("postgres transaction start: %w", err)
	}
	query := "INSERT INTO user_urls_table (user_id, long, short) VALUES ($1, $2, $3);"

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

// DeleteBatchWithUserID deletes a batch of URLs (if their userID matches with given one).
func (p *Postgresql) DeleteBatchWithUserID(userID int) (urlsChan chan string, err error) {
	tx, err := p.store.BeginTx(context.TODO(), nil)
	if err != nil {
		return nil, fmt.Errorf("postgres transaction start: %w", err)
	}

	query := "UPDATE user_urls_table SET is_deleted = true WHERE short = $1 AND user_id = $2;"

	go func() {
		for url := range urlsChan {
			_, errLocal := tx.Exec(query, url, userID)
			if errLocal != nil {
				tx.Rollback()
				return
			}
		}
		if errLocal := tx.Commit(); errLocal != nil {
			tx.Rollback()
			return
		}
	}()

	urlsChan = make(chan string)
	return urlsChan, nil
}

// Get returns an original URL using it`s short version.
func (p *Postgresql) Get(ctx context.Context, short string) (full string, err error) {

	query := "SELECT long, is_deleted  FROM user_urls_table WHERE short = $1;"
	row := p.store.QueryRowContext(ctx, query, short)

	var isDeleted bool
	err = row.Scan(&full, &isDeleted)

	if err != nil {
		return "", fmt.Errorf("postgres query: %w", err)
	}
	if isDeleted {
		return "", ErrURLWasDeleted()
	}
	return full, nil
}

// GetUserUrls returns all URLs of a user.
func (p *Postgresql) GetUserUrls(ctx context.Context, userID int) ([]entities.URL, error) {
	query := "SELECT long, short FROM user_urls_table WHERE user_id = $1;"

	var urls []entities.URL

	rows, err := p.store.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("postgres query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var url entities.URL
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

// Ping func pings real database and returns the answer.
func (p *Postgresql) Ping() error {
	return p.store.Ping()
}

// Close closes connection to a database.
func (p *Postgresql) Close() error {
	return p.store.Close()
}

// CreateUser creates a new user and saves it in a database.
func (p *Postgresql) CreateUser(ctx context.Context) (int, error) {
	query := "INSERT INTO users DEFAULT VALUES RETURNING id;"

	var userID int

	err := p.store.QueryRowContext(ctx, query).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("postgres create user: %w", err)
	}

	return userID, nil
}
