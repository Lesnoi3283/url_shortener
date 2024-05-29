package databases

import (
	"bufio"
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/Lesnoi3283/url_shortener/internal/app/entities"
	"os"
	"sync"
	"time"
)

type data struct {
	ID         int    `json:"id"`
	Key        string `json:"key"`
	Val        string `json:"val"`
	UserID     int    `json:"user_id"`
	WasDeleted bool   `json:"was_deleted"`
}

type JSONFileStorage struct {
	Path   string
	lastID int
	mutex  sync.Mutex
}

func NewJSONFileStorage(path string) *JSONFileStorage {
	toRet := &JSONFileStorage{Path: path}
	return toRet
}

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

func (j *JSONFileStorage) DeleteBatchWithUserID(userID int) (urlsChan chan string, err error) {
	urlsChan = make(chan string)
	//go func() {
	//	//to avoid deadlock
	//	for el := range urlsChan {
	//		el = el
	//	}
	//}()
	return urlsChan, ErrThisFuncIsNotSupported()
}

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

func (j *JSONFileStorage) Ping() error {
	return nil
}

// because it actually a session id, not a user id
func (j *JSONFileStorage) CreateUser(ctx context.Context) (int, error) {
	t := time.Now()
	timeBytes := []byte(t.Format(time.RFC3339Nano))
	hasher := sha256.New()
	hasher.Write(timeBytes)
	hashSum := hasher.Sum(nil)

	userID := int(binary.BigEndian.Uint64(hashSum[:8]))
	return userID, nil
}
