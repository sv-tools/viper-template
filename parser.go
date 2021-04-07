package vipertemplate

import (
	"text/template"

	"github.com/spf13/viper"
)

type parser struct {
	viper   *viper.Viper
	visited map[string]bool
	funcs   template.FuncMap
	data    interface{}
}

func newParser(opts ...Option) *parser {
	p := &parser{
		viper:   viper.GetViper(),
		visited: map[string]bool{},
		funcs:   template.FuncMap{},
		data:    nil,
	}
	p.funcs["Get"] = p.get

	for _, opt := range opts {
		opt(p)
	}

	return p
}

func (p *parser) parse(key string) (interface{}, error) {
	if _, ok := p.visited[key]; ok {
		return nil, newError(key, ErrCircularDependency)
	}

	val := p.viper.Get(key)
	if val == nil {
		return nil, nil
	}

	text, ok := val.(string)
	if !ok {
		return val, nil
	}

	tmpl, err := template.New(key).Funcs(p.funcs).Parse(text)
	if err != nil {
		return "", err
	}

	p.visited[key] = true

	buf := getBuffer()
	defer putBuffer(buf)

	if err := tmpl.Execute(buf, p.data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (p *parser) get(key string) (interface{}, error) {
	return p.parse(key)
}
