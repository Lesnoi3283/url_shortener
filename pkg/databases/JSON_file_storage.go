// Package databases contains all supported database adapters in this project.
package databases

import (
	"bufio"
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/Lesnoi3283/url_shortener/internal/app/entities"
)

type data struct {
	ID         int    `json:"id"`
	Key        string `json:"key"`
	Val        string `json:"val"`
	UserID     int    `json:"user_id"`
	WasDeleted bool   `json:"was_deleted"`
}

// JSONFileStorage is storage witch uses a file to store data. It writes a JSON arrays to it. Thread-safe.
type JSONFileStorage struct {
	Path   string
	lastID int
	mutex  sync.Mutex
}

// NewJSONFileStorage build a new JSONFileStorage.
func NewJSONFileStorage(path string) *JSONFileStorage {
	toRet := &JSONFileStorage{Path: path}
	return toRet
}

// Save saves a new url to a storage.
func (j *JSONFileStorage) Save(ctx context.Context, url entities.URL) error {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	// Open file
	file, err := os.OpenFile(j.Path, (os.O_RDWR | os.O_APPEND | os.O_CREATE), 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	//Find last id
	scanner := bufio.NewScanner(file)

	if j.lastID == 0 {
		lastLine := ""
		for scanner.Scan() {
			lastLine = scanner.Text()
		}
		if lastLine != "" {
			lastData := data{}
			err = json.Unmarshal([]byte(lastLine), &lastData)
			if err != nil {
				return err
			}
			j.lastID = lastData.ID
		}
	}

	newData := data{
		ID:  j.lastID + 1,
		Key: url.Short,
		Val: url.Long,
	}
	JSONData, err := json.Marshal(newData)
	JSONData = append(JSONData, '\n')

	if err != nil {
		return err
	}

	wr := bufio.NewWriter(file)
	_, err = wr.Write(JSONData)
	if err != nil {
		return err
	}
	wr.Flush()

	return nil
}

// SaveWithUserID saves a URL with userID.
func (j *JSONFileStorage) SaveWithUserID(ctx context.Context, userID int, url entities.URL) error {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	// Open file
	file, err := os.OpenFile(j.Path, (os.O_RDWR | os.O_APPEND | os.O_CREATE), 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	//Find last id
	scanner := bufio.NewScanner(file)

	if j.lastID == 0 {
		lastLine := ""
		for scanner.Scan() {
			lastLine = scanner.Text()
		}
		if lastLine != "" {
			lastData := data{}
			err = json.Unmarshal([]byte(lastLine), &lastData)
			if err != nil {
				return err
			}
			j.lastID = lastData.ID
		}
	}

	newData := data{
		ID:     j.lastID + 1,
		Key:    url.Short,
		Val:    url.Long,
		UserID: userID,
	}
	JSONData, err := json.Marshal(newData)
	JSONData = append(JSONData, '\n')

	if err != nil {
		return err
	}

	wr := bufio.NewWriter(file)
	_, err = wr.Write(JSONData)
	if err != nil {
		return err
	}
	wr.Flush()

	return nil
}

// SaveBatch saves a batch of URLs.
func (j *JSONFileStorage) SaveBatch(ctx context.Context, urls []entities.URL) error {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	// Open file
	file, err := os.OpenFile(j.Path, (os.O_RDWR | os.O_APPEND | os.O_CREATE), 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	//Find last id
	scanner := bufio.NewScanner(file)

	if j.lastID == 0 {
		lastLine := ""
		for scanner.Scan() {
			lastLine = scanner.Text()
		}
		if lastLine != "" {
			lastData := data{}
			err = json.Unmarshal([]byte(lastLine), &lastData)
			if err != nil {
				return err
			}
			j.lastID = lastData.ID
		}
	}

	//save all
	wr := bufio.NewWriter(file)
	for _, url := range urls {
		newData := data{
			ID:  j.lastID + 1,
			Key: url.Short,
			Val: url.Long,
		}
		j.lastID++

		JSONData, err := json.Marshal(newData)
		if err != nil {
			return err
		}

		JSONData = append(JSONData, '\n')

		_, err = wr.Write(JSONData)
		if err != nil {
			return err
		}
	}

	wr.Flush()

	return nil
}

// SaveBatchWithUserID save a batch of URLs with userID.
func (j *JSONFileStorage) SaveBatchWithUserID(ctx context.Context, userID int, urls []entities.URL) error {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	// Open file
	file, err := os.OpenFile(j.Path, (os.O_RDWR | os.O_APPEND | os.O_CREATE), 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	//Find last id
	scanner := bufio.NewScanner(file)

	if j.lastID == 0 {
		lastLine := ""
		for scanner.Scan() {
			lastLine = scanner.Text()
		}
		if lastLine != "" {
			lastData := data{}
			err = json.Unmarshal([]byte(lastLine), &lastData)
			if err != nil {
				return err
			}
			j.lastID = lastData.ID
		}
	}

	//save all
	wr := bufio.NewWriter(file)
	for _, url := range urls {
		newData := data{
			ID:     j.lastID + 1,
			Key:    url.Short,
			Val:    url.Long,
			UserID: userID,
		}
		j.lastID++

		JSONData, err := json.Marshal(newData)
		if err != nil {
			return err
		}

		JSONData = append(JSONData, '\n')

		_, err = wr.Write(JSONData)
		if err != nil {
			return err
		}
	}

	wr.Flush()

	return nil
}

// DeleteBatchWithUserID deletes a batch of URLs (if their userID matches with given one).
func (j *JSONFileStorage) DeleteBatchWithUserID(userID int) (urlsChan chan string, err error) {
	urlsChan = make(chan string)
	return urlsChan, ErrThisFuncIsNotSupported()
}

// GetUserUrls returns all URLs of a user.
func (j *JSONFileStorage) GetUserUrls(ctx context.Context, userID int) (URLs []entities.URL, err error) {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	file, err := os.Open(j.Path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	lastLine := ""
	URLs = make([]entities.URL, 0)
	for scanner.Scan() {
		lastLine = scanner.Text()
		lastData := data{}
		err = json.Unmarshal([]byte(lastLine), &lastData)
		if err != nil {
			return nil, err
		}

		if lastData.UserID == userID {
			URLs = append(URLs, entities.URL{Long: lastData.Val, Short: lastData.Key})
		}
	}

	return URLs, err
}

// Get returns an original URL using it`s short version.
func (j *JSONFileStorage) Get(ctx context.Context, key string) (string, error) {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	file, err := os.Open(j.Path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	lastLine := ""
	for scanner.Scan() {
		lastLine = scanner.Text()
		lastData := data{}
		err = json.Unmarshal([]byte(lastLine), &lastData)
		if err != nil {
			return "", err
		}

		if lastData.Key == key {
			return lastData.Val, nil
		}
	}

	err = fmt.Errorf("key doesnt exist")
	return "", err
}

// Ping always returns true.
func (j *JSONFileStorage) Ping() error {
	return nil
}

// CreateUser returns just a random int and doesn`t saves anything.
// Because it`s actually a session id, not a user id.
func (j *JSONFileStorage) CreateUser(ctx context.Context) (int, error) {
	t := time.Now()
	timeBytes := []byte(t.Format(time.RFC3339Nano))
	hasher := sha256.New()
	hasher.Write(timeBytes)
	hashSum := hasher.Sum(nil)

	userID := int(binary.BigEndian.Uint64(hashSum[:8]))
	return userID, nil
}
