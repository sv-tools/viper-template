package vipertemplate_test

import (
	"fmt"
	"os"
	"testing"
	"text/template"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"

	vipertemplate "github.com/sv-tools/viper-template/v2"
)

func TestGet(t *testing.T) {
	v := viper.New()
	v.Set("bar", 42)
	v.Set("foo", `{{ Get "bar" }}`)

	val, err := vipertemplate.Get(v, "foo")
	require.NoError(t, err)
	require.Equal(t, "42", val)
}

func TestGetIncorrectTemplate(t *testing.T) {
	v := viper.New()
	v.Set("foo", `{{ Get "bar"`)

	val, err := vipertemplate.Get(v, "foo")
	require.EqualError(t, err, "template: foo:1: unclosed action")
	require.Empty(t, val)
}

func TestGetNoKey(t *testing.T) {
	v := viper.New()
	val, err := vipertemplate.Get(v, "foo")
	require.NoError(t, err)
	require.Nil(t, val)
}

func TestGetString(t *testing.T) {
	v := viper.New()
	v.Set("bar", 42)
	v.Set("foo", `{{ Get "bar" }}`)

	val, err := vipertemplate.GetString(v, "foo")
	require.NoError(t, err)
	require.Equal(t, "42", val)
}

func TestGetNonStringValue(t *testing.T) {
	v := viper.New()
	v.Set("foo", 42)

	val, err := vipertemplate.Get(v, "foo")
	require.NoError(t, err)
	require.Equal(t, 42, val)
}

func TestGetStringNonStringValue(t *testing.T) {
	v := viper.New()
	v.Set("foo", 42)

	val, err := vipertemplate.GetString(v, "foo")
	require.ErrorIs(t, err, vipertemplate.ErrNonStringValue)
	require.Empty(t, val)
}

func TestGetCircularDependency(t *testing.T) {
	v := viper.New()
	v.Set("foo", `{{ Get "foo" }}`)

	val, err := vipertemplate.Get(v, "foo")
	require.Error(t, err)
	require.Contains(t, err.Error(), vipertemplate.ErrCircularDependency.Error())
	require.Empty(t, val)
}

func TestGetStringCircularDependency(t *testing.T) {
	v := viper.New()
	v.Set("foo", `{{ Get "foo" }}`)

	val, err := vipertemplate.GetString(v, "foo")
	require.Error(t, err)
	require.Contains(t, err.Error(), vipertemplate.ErrCircularDependency.Error())
	require.Empty(t, val)
}

func ExampleGet_first() {
	v := viper.New()
	v.Set("bar", 42)
	v.Set("foo", `{{ Get "bar" }}`)

	val, err := vipertemplate.Get(v, "foo")
	if err != nil {
		panic(err)
	}
	fmt.Println(val)
	// Output: 42
}

func ExampleGet_second() {
	v := viper.New()
	v.Set("bar", 42)
	v.Set("foo", `{{ Get "bar" }}`)

	val, err := vipertemplate.Get(v, "bar")
	if err != nil {
		panic(err)
	}
	fmt.Println(val)
	// Output: 42
}

func ExampleGet_readme() {
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

func ExampleGet_environmentVariable() {
	if err := os.Setenv("bar", "42"); err != nil {
		panic(err)
	}
	defer func() {
		if err := os.Unsetenv("bar"); err != nil {
			panic(err)
		}
	}()
	v := viper.New()
	v.Set("foo", `{{ Getenv "bar" }}`)

	funcs := template.FuncMap{
		"Getenv": os.Getenv,
	}
	val, err := vipertemplate.Get(
		v,
		"foo",
		vipertemplate.WithFuncs(funcs),
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(val)
	// Output: 42
}

func ExampleGetString_first() {
	v := viper.New()
	v.Set("bar", 42)
	v.Set("foo", `{{ Get "bar" }}`)

	val, err := vipertemplate.GetString(v, "foo")
	if err != nil {
		panic(err)
	}
	fmt.Println(val)

	_, err = vipertemplate.GetString(v, "bar")
	fmt.Println(err)
	// Output: 42
	// non-parsable template for the key 'bar': non-string value
}

func ExampleGetString_second() {
	v := viper.New()
	v.Set("bar", 42)
	v.Set("foo", `{{ Get "bar" }}`)

	_, err := vipertemplate.GetString(v, "bar")
	fmt.Println(err)
	// Output: non-parsable template for the key 'bar': non-string value
}

var benchmarkGetResult any

func BenchmarkGetParallel(b *testing.B) {
	v := viper.New()
	v.Set("bar", 42)
	v.Set("foo", `{{ Get "bar" }}`)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		var r any
		for pb.Next() {
			r, _ = vipertemplate.Get(v, "foo")
		}
		benchmarkGetResult = r
	})
}

func BenchmarkGetSequential(b *testing.B) {
	v := viper.New()
	v.Set("bar", 42)
	v.Set("foo", `{{ Get "bar" }}`)
	var r any
	b.ResetTimer()
	for range b.N {
		r, _ = vipertemplate.Get(v, "foo")
	}
	benchmarkGetResult = r
}
