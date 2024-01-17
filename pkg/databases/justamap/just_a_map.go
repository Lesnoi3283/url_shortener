package justamap

import "fmt"

type JustAMap struct {
	Store map[string]string
}

func NewJustAMap() *JustAMap {
	jm := &JustAMap{Store: make(map[string]string)}
	return jm
}

func (r *JustAMap) Save(key string, val string) error {
	r.Store[key] = val
	return nil
}

func (r *JustAMap) Get(key string) (toRet string, err error) {
	toRet, ok := r.Store[key]
	if !ok {
		err = fmt.Errorf("key doesnt exist")
	}
	return toRet, err
}
