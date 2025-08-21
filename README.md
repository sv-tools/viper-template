# viper-template

An add-on for the Go Viper library that evaluates Go text/template expressions inside configuration values. Use it to reference other keys, inject custom data and functions, and compute settings dynamically at read time.

[![Code Analysis](https://github.com/sv-tools/viper-template/actions/workflows/checks.yaml/badge.svg)](https://github.com/sv-tools/viper-template/actions/workflows/checks.yaml)
[![Go Reference](https://pkg.go.dev/badge/github.com/sv-tools/viper-template.svg)](https://pkg.go.dev/github.com/sv-tools/viper-template)
[![codecov](https://codecov.io/gh/sv-tools/viper-template/branch/main/graph/badge.svg?token=0XVOTDR1CW)](https://codecov.io/gh/sv-tools/viper-template)
[![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/sv-tools/viper-template?style=flat-square)](https://github.com/sv-tools/viper-template/releases)

## v2 Changes

In v2, Get and GetString accept a Viper instance as the first parameter. The Viper interface abstracts the dependency, allowing this package to use one Viper version (e.g., for testing) while callers use another.

The Viper version pinned in go.mod is for internal testing only and does not constrain the Viper version you use in your project.

## Usage

```shell
go get github.com/sv-tools/viper-template/v2@latest
```

```go
package main

import (
  "fmt"
  "text/template"

  "github.com/spf13/viper"
  vipertemplate "github.com/sv-tools/viper-template/v2"
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
```

## Benchmarks

### `v1.9.0`

using go `v1.23.9`

```shell
% go test -bench=. -benchmem

goos: darwin
goarch: arm64
pkg: github.com/sv-tools/viper-template
cpu: Apple M1
BenchmarkGetParallel-8            757263              1582 ns/op            4330 B/op         43 allocs/op
BenchmarkGetSequential-8          417817              2780 ns/op            4328 B/op         43 allocs/op
```
