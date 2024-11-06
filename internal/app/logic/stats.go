package logic

import (
	"context"
	"fmt"
)

// GetStats returns an amount of users and shorten urls.
func GetStats(ctx context.Context, storage URLStorageInterface) (URLsAmount int, usersAmount int, err error) {
	//get URLs count
	URLsAmount, err = storage.GetShortURLCount(ctx)
	if err != nil {
		return 0, 0, fmt.Errorf("catn get a shorten URLs count: %w", err)
	}

	//get users count
	usersAmount, err = storage.GetUsersCount(ctx)
	if err != nil {
		return 0, 0, fmt.Errorf("catn get a users count: %w", err)
	}

	return usersAmount, URLsAmount, nil
}