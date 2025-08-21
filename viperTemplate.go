package vipertemplate

// Viper is a minimal adapter interface used by this package.
// It decouples us from a specific github.com/spf13/viper version,
// allowing callers to use any Viper version (or a compatible implementation).
type Viper interface {
	Get(key string) any
}

// Get tries to parse a go template stored in a Viper by a given key.
// If the type of the value is not string then just returns,
// otherwise it parses and executes the stored text as a go template.
// The `Get` function added by default. It supports the recursive calls.
func Get(viper Viper, key string, opts ...Option) (any, error) {
	p := parserPool.Get().(*parser)
	defer func() {
		p.reset()
		parserPool.Put(p)
	}()
	p.viper = viper
	for _, opt := range opts {
		opt(p)
	}

	return p.parse(key)
}

// GetString tries to parse a go template stored in a Viper by a given key.
// It fails if the stored value is not string.
func GetString(viper Viper, key string, opts ...Option) (string, error) {
	val, err := Get(viper, key, opts...)
	if err != nil {
		return "", err
	}

	text, ok := val.(string)
	if !ok {
		return "", newError(key, ErrNonStringValue)
	}

	return text, nil
}
