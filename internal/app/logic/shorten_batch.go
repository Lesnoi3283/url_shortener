package logic

import (
	"context"
	"fmt"
	"github.com/Lesnoi3283/url_shortener/internal/app/entities"
)

// ShortenBatch saves a batch of URLs to a storage.
// Returns a slice of URLs with an original URL and a short version (with a base address).
// Use "userID = -1" to save URLs without a userID.
func ShortenBatch(ctx context.Context, URLs []entities.URL, baseAddress string, storage URLStorageInterface, userID int) ([]entities.URL, error) {
	//shorting
	for i, url := range URLs {
		URLs[i].ShortURL = string(ShortenURL([]byte(url.OriginalURL)))
		URLs[i].OriginalURL = ""
	}

	//url saving
	var err error
	if userID != -1 {
		err = storage.SaveBatchWithUserID(ctx, userID, URLs)
	} else {
		err = storage.SaveBatch(ctx, URLs)
	}
	if err != nil {
		return nil, fmt.Errorf("error while saving URLs to a storage: %w", err)
	}

	//adding base address to return
	for i := range URLs {
		URLs[i].ShortURL = baseAddress + "/" + URLs[i].ShortURL
	}
	return URLs, nil
}
