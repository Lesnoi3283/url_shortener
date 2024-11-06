package logic

import (
	"context"
	"fmt"
	"github.com/Lesnoi3283/url_shortener/internal/app/entities"
)

// Shorten saves one URL to a storage.
// Returns a short version with a base address. Even after an error from database it will return a short version.
// Can return a wrapped databases.AlreadyExistsError (in this case use short url value from error).
// Use "userID = -1" to save URLs without a userID.
func Shorten(ctx context.Context, URL []byte, baseAddress string, storage URLStorageInterface, userID int) (string, error) {
	urlShort := string(ShortenURL(URL))

	//url saving
	var err error
	if userID != -1 {
		err = storage.SaveWithUserID(ctx, userID, entities.URL{
			ShortURL:    urlShort,
			OriginalURL: string(URL),
		})
	} else {
		err = storage.Save(ctx, entities.URL{
			ShortURL:    urlShort,
			OriginalURL: string(URL),
		})
	}
	if err != nil {
		return "", fmt.Errorf("error while saving URL into a storage: %w", err)
	}

	//return
	return baseAddress + "/" + urlShort, nil
}
