package databases

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"github.com/Lesnoi3283/url_shortener/internal/app/entities"
	"time"
)

type JustAMap struct {
	Store     map[string]string
	UserStore map[string]int
}

func NewJustAMap() *JustAMap {
	jm := &JustAMap{Store: make(map[string]string), UserStore: make(map[string]int)}
	return jm
}

func (j *JustAMap) SaveWithUserID(ctx context.Context, userID int, url entities.URL) error {
	j.Store[url.Short] = url.Long
	j.UserStore[url.Short] = userID
	return nil
}

func (j *JustAMap) SaveBatchWithUserID(ctx context.Context, userID int, urls []entities.URL) error {
	for _, el := range urls {
		err := j.SaveWithUserID(ctx, userID, el)
		if err != nil {
			return err
		}
	}
	return nil
}

func (j *JustAMap) DeleteBatchWithUserID(userID int) (urlsChan chan string, err error) {
	urlsChan = make(chan string)
	//go func() {
	//	//to avoid deadlock
	//	for el := range urlsChan {
	//		el = el
	//	}
	//}()
	return urlsChan, ErrThisFuncIsNotSupported()
}

func (j *JustAMap) GetUserUrls(ctx context.Context, userID int) ([]entities.URL, error) {
	toRet := make([]entities.URL, 0)

	for short, full := range j.Store {
		uID := j.UserStore[short]
		if (uID == userID) && (uID != 0) {
			toRet = append(toRet, entities.URL{Long: full, Short: short})
		}
	}

	return toRet, nil
}

func (j *JustAMap) Ping() error {
	return nil
}

// because it actually a session id, not a user id
func (j *JustAMap) CreateUser(ctx context.Context) (int, error) {
	t := time.Now()
	timeBytes := []byte(t.Format(time.RFC3339Nano))
	hasher := sha256.New()
	hasher.Write(timeBytes)
	hashSum := hasher.Sum(nil)

	userID := int(binary.BigEndian.Uint64(hashSum[:8]))
	return userID, nil
}

func (j *JustAMap) Save(ctx context.Context, url entities.URL) error {
	j.Store[url.Short] = url.Long
	return nil
}

func (j *JustAMap) SaveBatch(ctx context.Context, urls []entities.URL) error {
	for _, url := range urls {
		err := j.Save(ctx, url)
		if err != nil {
			return err
		}
	}
	return nil
}

func (j *JustAMap) Get(ctx context.Context, key string) (toRet string, err error) {
	toRet, ok := j.Store[key]
	if !ok {
		err = fmt.Errorf("key doesnt exist")
	}
	return toRet, err
}
