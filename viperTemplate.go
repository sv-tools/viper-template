package vipertemplate

import (
	"sync"
	"text/template"

	"github.com/spf13/viper"
)

var parserPool = sync.Pool{
	New: func() any {
		return &parser{}
	},
}

// Get tries to parse a go template stored in a Viper by a given key.
// If the type of the value is not string then just returns,
// otherwise it parses and executes the stored text as a go template.
// The `Get` function added by default. It supports the recursive calls.
func Get(key string, opts ...Option) (any, error) {
	p := parserPool.Get().(*parser)
	defer func() {
		parserPool.Put(p)
	}()
	p.viper = viper.GetViper()
	p.visited = make(map[string]struct{})
	p.data = nil
	p.funcs = template.FuncMap{"Get": p.get}
	for _, opt := range opts {
		opt(p)
	}

	return p.parse(key)
}

// GetString tries to parse a go template stored in a Viper by a given key.
// It fails if the stored value is not string.
func GetString(key string, opts ...Option) (string, error) {
	val, err := Get(key, opts...)
	if err != nil {
		return "", err
	}

	text, ok := val.(string)
	if !ok {
		return "", newError(key, ErrNonStringValue)
	}

	return text, nil
}
