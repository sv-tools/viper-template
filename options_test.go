package vipertemplate_test

import (
	"fmt"
	"testing"
	"text/template"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"

	vipertemplate "github.com/sv-tools/viper-template/v2"
)

func TestGetWithData(t *testing.T) {
	v := viper.New()
	v.Set("foo", "{{ .Bar }}")

	data := struct {
		Bar int
	}{
		Bar: 42,
	}

	val, err := vipertemplate.Get(v, "foo", vipertemplate.WithData(&data))
	require.NoError(t, err)
	require.Equal(t, "42", val)
}

func TestGetWithFuncs(t *testing.T) {
	v := viper.New()
	v.Set("foo", "{{ Bar }}")

	funcs := template.FuncMap{
		"Bar": func() int {
			return 42
		},
	}

	val, err := vipertemplate.Get(v, "foo", vipertemplate.WithFuncs(funcs))
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

	val, err := vipertemplate.Get(
		v,
		"foo",
		vipertemplate.WithData(&data),
		vipertemplate.WithFuncs(funcs),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println(val)
	// Output: 84
}
