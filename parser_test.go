package vipertemplate_test

import (
	"errors"
	"testing"
	"text/template"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"

	vipertemplate "github.com/sv-tools/viper-template/v2"
)

// Test mutual (indirect) circular dependency detection between two keys.
func TestGetMutualCircularDependency(t *testing.T) {
	v := viper.New()
	v.Set("foo", `{{ Get "bar" }}`)
	v.Set("bar", `{{ Get "foo" }}`)

	val, err := vipertemplate.GetString(v, "foo")
	require.Error(t, err)
	require.Contains(t, err.Error(), vipertemplate.ErrCircularDependency.Error())
	require.Empty(t, val)
}

// Test that custom functions don't leak across pooled parser instances.
func TestGetFunctionMapReset(t *testing.T) {
	v := viper.New()
	v.Set("foo", `{{ Custom }}`)

	// First call with the custom function should work.
	funcs := template.FuncMap{
		"Custom": func() int { return 7 },
	}
	val, err := vipertemplate.GetString(v, "foo", vipertemplate.WithFuncs(funcs))
	require.NoError(t, err)
	require.Equal(t, "7", val)

	// Second call without providing funcs again should fail because function is not defined.
	val, err = vipertemplate.GetString(v, "foo")
	require.Error(t, err)
	require.Contains(t, err.Error(), `function "Custom" not defined`)
	require.Empty(t, val)
}

var errBoom = errors.New("boom")

// Test execution error (template parses but function returns an execution error).
func TestGetExecutionError(t *testing.T) {
	v := viper.New()
	v.Set("foo", `{{ Fail }}`)

	funcs := template.FuncMap{
		"Fail": func() (string, error) { return "", errBoom },
	}
	val, err := vipertemplate.Get(v, "foo", vipertemplate.WithFuncs(funcs))
	require.Error(t, err)
	require.Contains(t, err.Error(), "boom")
	require.Empty(t, val)
}

// Test that non-string nested values inside another template render correctly with additional static text.
func TestGetNestedNonStringValueFormatting(t *testing.T) {
	v := viper.New()
	v.Set("bar", 42) // non-string
	v.Set("foo", `Value: {{ Get "bar" }}`)

	val, err := vipertemplate.GetString(v, "foo")
	require.NoError(t, err)
	require.Equal(t, "Value: 42", val)
}
