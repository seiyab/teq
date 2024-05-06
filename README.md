# teq

![CI](https://github.com/seiyab/teq/actions/workflows/go.yml/badge.svg)

teq is a Go library designed to enhance your testing experience by providing a flexible and customizable way to perform deep equality checks. It's especially useful when you need to compare complex objects or types that have specific equality conditions.

## Features

- Transforms: Register a "transform" function that can modify objects before comparison. This allows you to control how equality is determined. For example, by transforming time.Time objects to UTC, you can make your equality checks timezone-insensitive.
- Formats: Register a "format" function that defines how objects are displayed when they are not equal. This is useful for types like time.Time and decimal.Decimal that may not be human-readable in their default format. By registering your own format, you can make the output of your tests more understandable.

## Installation

```sh
go get github.com/seiyab/teq@latest
```

## Usage

To use teq, you first need to create a new instance:

```go
tq := teq.New()
```

Then, you can add your transforms and formats:

```go
// time.Time will be transformed into UTC time. So equality check with `tq` will be timezone-insensitive.
tq.AddTransform(func(d time.Time) time.Time {
    return d.In(d.UTC)
})

// time.Time will be shown in RFC3339 format when it appear in diff.
tq.AddFormat(func(d time.Time) string {
    return d.Format(time.RFC3339)
})
```

Finally, you can use teq to perform deep equality checks in your tests:

```go
tq.Equal(t, expected, actual)
```

If you need "common" equality across your project, we recommend to define function that returns your customized teq.

```go
func NewTeq() teq.Teq {
    tq := teq.New()
    tq.AddTransform(/* ... */)
    tq.AddTransform(/* ... */)
    // :
    return tq
}

// then you can use easily it everywhere
tq := NewTeq()
```

## Prior works

- [testify](https://github.com/stretchr/testify)
- [deepequal.go](https://github.com/weaveworks/scope/blob/12175b96a3456f1c2b050f1e1d6432498ed64d95/test/reflect/deepequal.go) in "github.com/weaveworks/scope/test/reflect" package authored by Weaveworks Ltd
- [go-cmp](https://github.com/google/go-cmp)
