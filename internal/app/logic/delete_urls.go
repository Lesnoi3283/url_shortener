package logic

import (
	"fmt"
)

//TODO: Принимать интерфейс бд и логера. Передавать их сюда будем из хендлеров HTTP и gRPC.
// В структуру gRPC сервера также можнжо засунуть базу данных и прочие структуры.

//TODO: Перенести интерфейс стораджа в этот пакет.

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
