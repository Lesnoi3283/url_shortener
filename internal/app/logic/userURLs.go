package logic

import (
	"context"
	"github.com/Lesnoi3283/url_shortener/internal/app/entities"
)

// GetUsersURLs return all user`s URLs with full and short versions.
// Short version includes the base address.
func GetUsersURLs(ctx context.Context, storage URLStorageInterface, baseAddress string, userID int) ([]entities.URL, error) {
	usersURLs, err := storage.GetUserUrls(ctx, userID)
	if err != nil {
		return nil, err
	}
	for i := range usersURLs {
		usersURLs[i].ShortURL = baseAddress + "/" + usersURLs[i].ShortURL
	}
	return usersURLs, nil
}
