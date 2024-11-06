package logic

import (
	"fmt"
)

// DeleteURLs deletes URLs from database. It creates a new goroutine witch deletes URLs.
func DeleteURLs(userID int, shortURLs []string, storage URLStorageInterface) error {
	inputCh, err := storage.DeleteBatchWithUserID(userID)
	if err != nil {
		return fmt.Errorf("error while deleting URLs: %w", err)
	}
	//fan-out
	go func() {
		defer close(inputCh)
		for _, URL := range shortURLs {
			inputCh <- URL
		}
	}()
	return nil
}
