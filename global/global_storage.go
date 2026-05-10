package global

import "errors"

var storage map[string]any
var (
	ErrNotFound      = errors.New("По этому ключу ничего не найдено")
	ErrTypeAssertion = errors.New("Ошибка приведения типов")
)

func Init() {
	storage = map[string]any{}
}

func Add(name string, value any) {
	storage[name] = value
}

func Get(name string) (any, error) {
	v, ok := storage[name]
	if !ok {
		return nil, ErrNotFound
	}

	return v, nil
}
