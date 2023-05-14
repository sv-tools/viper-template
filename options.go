package vipertemplate

import (
	"text/template"

	"github.com/spf13/viper"
)

// Option is an interface for the additional options for Get or GetString functions
type Option func(p *parser)

// WithViper is an option to use a custom viper object
func WithViper(v *viper.Viper) Option {
	return func(p *parser) {
		p.viper = v
	}
}

// WithFuncs is an option to extend the predefined functions
func WithFuncs(funcs template.FuncMap) Option {
	return func(p *parser) {
		for name, value := range funcs {
			p.funcs[name] = value
		}
	}
}

// WithData is an option to use a given data as an input for the templates
func WithData(data interface{}) Option {
	return func(p *parser) {
		p.data = data
	}
}
