package dsn

import "errors"

type DSNHelper interface {
	Gen() string
}

var dsnFactory = map[string]func() DSNHelper{}

func RegisterDSN(name string, generator func() DSNHelper) {
	if name == "" {
		panic("disallow empty name")
	}

	dsnFactory[name] = generator
}

func GetGenerator(name string) (DSNHelper, error) {
	g, ok := dsnFactory[name]
	if !ok {
		return nil, errors.New("no generator exists")
	}
	return g(), nil
}
