package vipertemplate_test

import (
	"fmt"
	"testing"
	"text/template"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"

	viperTemplate "github.com/sv-tools/viper-template"
)

func TestGetWithViper(t *testing.T) {
	t.Cleanup(func() {
		viper.Reset()
	})
	viper.Set("foo", 43)

	v := viper.New()
	v.Set("foo", 42)

	val, err := viperTemplate.Get("foo", viperTemplate.WithViper(v))
	require.NoError(t, err)
	require.Equal(t, 42, val)
}

func TestGetWithData(t *testing.T) {
	t.Cleanup(func() {
		viper.Reset()
	})
	viper.Set("foo", "{{ .Bar }}")

	data := struct {
		Bar int
	}{
		Bar: 42,
	}

	val, err := viperTemplate.Get("foo", viperTemplate.WithData(&data))
	require.NoError(t, err)
	require.Equal(t, "42", val)
}

func TestGetWithFuncs(t *testing.T) {
	t.Cleanup(func() {
		viper.Reset()
	})
	viper.Set("foo", "{{ Bar }}")

	funcs := template.FuncMap{
		"Bar": func() int {
			return 42
		},
	}

	val, err := viperTemplate.Get("foo", viperTemplate.WithFuncs(funcs))
	require.NoError(t, err)
	require.Equal(t, "42", val)
}

func ExampleGet_with_options() {
	v := viper.New()
	v.Set("foo", `{{ Get "bar" }}`)
	v.Set("bar", `{{ Mul . 2 }}`)

	type Data struct {
		Bar int
	}
	data := Data{
		Bar: 42,
	}

	funcs := template.FuncMap{
		"Mul": func(d *Data, v int) int {
			return d.Bar * v
		},
	}

	val, err := viperTemplate.Get(
		"foo",
		viperTemplate.WithViper(v),
		viperTemplate.WithData(&data),
		viperTemplate.WithFuncs(funcs),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println(val)
	// Output: 84
}
