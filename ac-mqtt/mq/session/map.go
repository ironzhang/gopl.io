package session

import (
	"fmt"
	"io"
)

type Map interface {
	Set(s *Session) error
	Get(id string) (*Session, bool)
	Del(id string)
	Count() int
}

type mapFactoryFunc func(param string) (Map, error)

var factorys = make(map[string]mapFactoryFunc)

func Register(name string, factFunc mapFactoryFunc) {
	if factFunc == nil {
		panic(fmt.Sprintf("%q session map factory func is nil", name))
	}

	if _, dup := factorys[name]; dup {
		panic(fmt.Sprintf("%q session map factory func duplicate", name))
	}

	factorys[name] = factFunc
}

func Unregister(name string) {
	delete(factorys, name)
}

func NewMap(name string, param string) (Map, error) {
	f, ok := factorys[name]
	if !ok {
		return nil, fmt.Errorf("unsupport %q manager", name)
	}
	return f(param)
}

func Close(m Map) error {
	if c, ok := m.(io.Closer); ok {
		return c.Close()
	}
	return nil
}
