package databases

import (
	"context"
	"fmt"
	"github.com/Lesnoi3283/url_shortener/internal/app/entities"
)

type JustAMap struct {
	Store map[string]string
}

func NewJustAMap() *JustAMap {
	jm := &JustAMap{Store: make(map[string]string)}
	return jm
}

func (r *JustAMap) Save(ctx context.Context, key string, val string) error {
	r.Store[key] = val
	return nil
}

func (r *JustAMap) SaveBatch(ctx context.Context, urls []entities.URL) error {
	for _, url := range urls {
		err := r.Save(ctx, url.Short, url.Long)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *JustAMap) Get(ctx context.Context, key string) (toRet string, err error) {
	//Есть идея использовать тут контекст и горутины,
	//в селекте ожидать получение первого из значений
	// - завершения контекста или возвращения значения.
	//Насколько хорошая идея?
	//P.s. Аналогично в JSON_file_storge и других хранилищах,
	//в которых контекст не используется по умолчанию (как в постгрес)
	toRet, ok := r.Store[key]
	if !ok {
		err = fmt.Errorf("key doesnt exist")
	}
	return toRet, err
}
