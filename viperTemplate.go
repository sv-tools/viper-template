package vipertemplate

// Get tries to parse a go template stored in a Viper by a given key.
// If the type of the value is not string then just returns,
// otherwise it parses and executes the stored text as a go template.
// The `Get` function added by default. It supports the recursive calls.
func Get(key string, opts ...Option) (interface{}, error) {
	return newParser(opts...).parse(key)
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
