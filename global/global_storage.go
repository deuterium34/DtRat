package global

import "errors"

var storage map[string]func()
var ErrNotFound = errors.New("Метод по этому ключу не найден")

func Init() {
	storage = map[string]func(){}
}

func Add(name string, value func()) {
	storage[name] = value
}

func Get(name string) (func(), error) {
	v, ok := storage[name]
	if !ok {
		return nil, ErrNotFound
	}

	return v, nil
}
