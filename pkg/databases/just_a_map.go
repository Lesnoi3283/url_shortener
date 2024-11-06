package databases

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"sync"
	"time"

	"github.com/Lesnoi3283/url_shortener/internal/app/entities"
)

// JustAMap is an in-memory storage.
type JustAMap struct {
	Store     map[string]string
	UserStore map[string]int
	Mutex     sync.RWMutex
}

// NewJustAMap build a new JustAMap.
func NewJustAMap() *JustAMap {
	jm := &JustAMap{Store: make(map[string]string), UserStore: make(map[string]int)}
	return jm
}

// SaveWithUserID saves a URL with userID.
func (j *JustAMap) SaveWithUserID(ctx context.Context, userID int, url entities.URL) error {
	j.Mutex.Lock()
	defer j.Mutex.Unlock()
	j.Store[url.ShortURL] = url.OriginalURL
	j.UserStore[url.ShortURL] = userID
	return nil
}

// SaveBatchWithUserID save a batch of URLs with userID.
func (j *JustAMap) SaveBatchWithUserID(ctx context.Context, userID int, urls []entities.URL) error {
	for _, el := range urls {
		err := j.SaveWithUserID(ctx, userID, el)
		if err != nil {
			return err
		}
	}
	return nil
}

// DeleteBatchWithUserID deletes a batch of URLs (if their userID matches with given one).
func (j *JustAMap) DeleteBatchWithUserID(userID int) (urlsChan chan string, err error) {
	urlsChan = make(chan string)
	return urlsChan, ErrThisFuncIsNotSupported()
}

// GetUserUrls returns all URLs of a user.
func (j *JustAMap) GetUserUrls(ctx context.Context, userID int) ([]entities.URL, error) {
	j.Mutex.RLock()
	defer j.Mutex.RUnlock()

	toRet := make([]entities.URL, 0)
	for short, full := range j.Store {
		uID := j.UserStore[short]
		if (uID == userID) && (uID != 0) {
			toRet = append(toRet, entities.URL{OriginalURL: full, ShortURL: short})
		}
	}

	return toRet, nil
}

// Ping always returns true.
func (j *JustAMap) Ping() error {
	return nil
}

// CreateUser returns just a random int and doesn`t saves anything.
// Because it`s actually a session id, not a user id.
func (j *JustAMap) CreateUser(ctx context.Context) (int, error) {
	t := time.Now()
	timeBytes := []byte(t.Format(time.RFC3339Nano))
	hasher := sha256.New()
	hasher.Write(timeBytes)
	hashSum := hasher.Sum(nil)

	userID := int(binary.BigEndian.Uint64(hashSum[:8]))
	return userID, nil
}

// Save saves a new url to a storage.
func (j *JustAMap) Save(ctx context.Context, url entities.URL) error {
	j.Mutex.Lock()
	defer j.Mutex.Unlock()
	j.Store[url.ShortURL] = url.OriginalURL
	return nil
}

// SaveBatch saves a batch of URLs.
func (j *JustAMap) SaveBatch(ctx context.Context, urls []entities.URL) error {
	for _, url := range urls {
		err := j.Save(ctx, url)
		if err != nil {
			return err
		}
	}
	return nil
}

// Get returns an original URL using it`s short version.
func (j *JustAMap) Get(ctx context.Context, key string) (toRet string, err error) {
	j.Mutex.RLock()
	defer j.Mutex.RUnlock()
	toRet, ok := j.Store[key]
	if !ok {
		err = fmt.Errorf("key doesnt exist")
	}
	return toRet, err
}

// GetUsersCount returns the total number of users in the database.
// JustAMap DOESN`T SUPPORT IT NOW!
func (j *JustAMap) GetUsersCount(ctx context.Context) (int, error) {
	return 0, ErrThisFuncIsNotSupported()
}

// GetShortURLCount returns the total number of short URLs in the map storage.
func (j *JustAMap) GetShortURLCount(ctx context.Context) (int, error) {
	j.Mutex.Lock()
	defer j.Mutex.Unlock()
	return len(j.Store), nil
}
