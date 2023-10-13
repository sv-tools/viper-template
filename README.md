# viper-template
An extension for Golang [Viper](https://github.com/spf13/viper) library to use and execute Go templates in the viper configs.

[![Code Analysis](https://github.com/sv-tools/viper-template/actions/workflows/checks.yaml/badge.svg)](https://github.com/sv-tools/viper-template/actions/workflows/checks.yaml)
[![Go Reference](https://pkg.go.dev/badge/github.com/sv-tools/viper-template.svg)](https://pkg.go.dev/github.com/sv-tools/viper-template)
[![codecov](https://codecov.io/gh/sv-tools/viper-template/branch/main/graph/badge.svg?token=0XVOTDR1CW)](https://codecov.io/gh/sv-tools/viper-template)
[![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/sv-tools/viper-template?style=flat-square)](https://github.com/sv-tools/viper-template/releases)

## Usage

```go
package main

import (
	"fmt"
	"text/template"

	"github.com/spf13/viper"
	vipertemplate "github.com/sv-tools/viper-template"
)

func main() {
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
		"foo",
		vipertemplate.WithViper(v),
		vipertemplate.WithData(&data),
		vipertemplate.WithFuncs(funcs),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println(val)
	// Output: 84
}
```

## Benchmarks

### `v1.5.0`
using go `v1.21.3`

```
goos: darwin
goarch: arm64
pkg: github.com/sv-tools/viper-template
BenchmarkGetParallel-8            807338              1510 ns/op            4329 B/op         43 allocs/op
BenchmarkGetSequential-8          431893              2716 ns/op            4327 B/op         43 allocs/op
```
