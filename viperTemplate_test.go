package vipertemplate_test

import (
	"fmt"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"

	vipertemplate "github.com/sv-tools/viper-template"
)

func TestGet(t *testing.T) {
	t.Cleanup(func() {
		viper.Reset()
	})

	viper.Set("bar", 42)
	viper.Set("foo", `{{ Get "bar" }}`)

	val, err := vipertemplate.Get("foo")
	require.NoError(t, err)
	require.Equal(t, "42", val)
}

func TestGetIncorrectTemplate(t *testing.T) {
	t.Cleanup(func() {
		viper.Reset()
	})

	viper.Set("foo", `{{ Get "bar"`)

	val, err := vipertemplate.Get("foo")
	require.EqualError(t, err, "template: foo:1: unclosed action")
	require.Empty(t, val)
}

func TestGetNoKey(t *testing.T) {
	t.Cleanup(func() {
		viper.Reset()
	})

	val, err := vipertemplate.Get("foo")
	require.NoError(t, err)
	require.Nil(t, val)
}

func TestGetString(t *testing.T) {
	t.Cleanup(func() {
		viper.Reset()
	})

	viper.Set("bar", 42)
	viper.Set("foo", `{{ Get "bar" }}`)

	val, err := vipertemplate.GetString("foo")
	require.NoError(t, err)
	require.Equal(t, "42", val)
}

func TestGetNonStringValue(t *testing.T) {
	t.Cleanup(func() {
		viper.Reset()
	})

	viper.Set("foo", 42)

	val, err := vipertemplate.Get("foo")
	require.NoError(t, err)
	require.Equal(t, 42, val)
}

func TestGetStringNonStringValue(t *testing.T) {
	t.Cleanup(func() {
		viper.Reset()
	})

	viper.Set("foo", 42)

	val, err := vipertemplate.GetString("foo")
	require.ErrorIs(t, err, vipertemplate.ErrNonStringValue)
	require.Empty(t, val)
}

func TestGetCircularDependency(t *testing.T) {
	t.Cleanup(func() {
		viper.Reset()
	})

	viper.Set("foo", `{{ Get "foo" }}`)

	val, err := vipertemplate.Get("foo")
	require.Error(t, err)
	require.Contains(t, err.Error(), vipertemplate.ErrCircularDependency.Error())
	require.Empty(t, val)
}

func TestGetStringCircularDependency(t *testing.T) {
	t.Cleanup(func() {
		viper.Reset()
	})

	viper.Set("foo", `{{ Get "foo" }}`)

	val, err := vipertemplate.GetString("foo")
	require.Error(t, err)
	require.Contains(t, err.Error(), vipertemplate.ErrCircularDependency.Error())
	require.Empty(t, val)
}

func ExampleGet_first() {
	defer viper.Reset()
	viper.Set("bar", 42)
	viper.Set("foo", `{{ Get "bar" }}`)

	val, err := vipertemplate.Get("foo")
	if err != nil {
		panic(err)
	}
	fmt.Println(val)
	// Output: 42
}

func ExampleGet_second() {
	defer viper.Reset()
	viper.Set("bar", 42)
	viper.Set("foo", `{{ Get "bar" }}`)

	val, err := vipertemplate.Get("bar")
	if err != nil {
		panic(err)
	}
	fmt.Println(val)
	// Output: 42
}

func ExampleGetString_first() {
	defer viper.Reset()
	viper.Set("bar", 42)
	viper.Set("foo", `{{ Get "bar" }}`)

	val, err := vipertemplate.GetString("foo")
	if err != nil {
		panic(err)
	}
	fmt.Println(val)

	_, err = vipertemplate.GetString("bar")
	fmt.Println(err)
	// Output: 42
	// non-parsable template for the key 'bar': non-string value
}

func ExampleGetString_second() {
	defer viper.Reset()
	viper.Set("bar", 42)
	viper.Set("foo", `{{ Get "bar" }}`)

	_, err := vipertemplate.GetString("bar")
	fmt.Println(err)
	// non-parsable template for the key 'bar': non-string value
}
