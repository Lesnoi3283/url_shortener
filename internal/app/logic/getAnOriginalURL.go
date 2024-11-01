package logic

import (
	"context"
	"fmt"
)

// GetOriginalURL is just a wrapper for a storages func.
// It returns an original URL for given short URL.
// Can return the databases.ErrURLWasDeleted error.
func GetOriginalURL(ctx context.Context, shortURL string, storage URLStorageInterface) (string, error) {
	//reading from DB
	fullURL, err := storage.Get(ctx, shortURL)
	if err != nil {
		return "", fmt.Errorf("error while getting full url from db: %w", err)
	} else {
		return fullURL, nil
	}
}
