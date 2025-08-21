package vipertemplate

import (
	"bytes"
	"sync"
	"text/template"
)

var bytesPool = sync.Pool{
	New: func() any {
		return new(bytes.Buffer)
	},
}

type parser struct {
	viper   Viper
	visited map[string]struct{}
	funcs   template.FuncMap
	data    any
}

var parserPool = sync.Pool{
	New: func() any {
		p := new(parser)
		p.reset()
		return p
	},
}

func (p *parser) parse(key string) (any, error) {
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

	p.visited[key] = struct{}{}

	buf := bytesPool.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		bytesPool.Put(buf)
	}()

	if err := tmpl.Execute(buf, p.data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (p *parser) get(key string) (any, error) {
	return p.parse(key)
}

func (p *parser) reset() {
	p.visited = make(map[string]struct{})
	p.data = nil
	p.funcs = template.FuncMap{"Get": p.get}
}
