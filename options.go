package vipertemplate

import (
	"text/template"
)

// Option is an interface for the additional options for Get or GetString functions
type Option func(p *parser)

// WithFuncs is an option to extend the predefined functions
func WithFuncs(funcs template.FuncMap) Option {
	return func(p *parser) {
		for name, value := range funcs {
			p.funcs[name] = value
		}
	}
}

// WithData is an option to use a given data as an input for the templates
func WithData(data any) Option {
	return func(p *parser) {
		p.data = data
	}
}
