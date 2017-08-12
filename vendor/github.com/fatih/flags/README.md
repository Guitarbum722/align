# Flags [![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://godoc.org/github.com/fatih/flags) [![Build Status](http://img.shields.io/travis/fatih/flags.svg?style=flat-square)](https://travis-ci.org/fatih/flags)


Flags is a low level package for parsing or managing single flag arguments and
their associated values from a list of arguments. It's useful for CLI
applications or creating logic for parsing arguments(custom or os.Args)
manually. 

Note that there is no context available for flags. You need to know upfront how
flags are supposed to be parsed.

Checkout the usage below for examples:

## Install

```bash
go get github.com/fatih/flags
```

## Usage and examples

Let us define three flags. Flags needs to be compatible with the
[flag](https://golang.org/pkg/flag/) package.

```go
args := []string{"--key", "123", "--name=example", "--debug"}
```

Check if a flag exists in the argument list

```go
flags.Has("key", args)    // true
flags.Has("--key", args)  // true
flags.Has("secret", args) // false
```

Get the value for from a flag name

```go
val, _ := flags.Value("--key", args) // val -> "123"
val, _ := flags.Value("name", args)  // val -> "example"
val, _ := flags.Value("debug", args) // val -> "" (means true boolean)
```

Exclude a flag and it's value from the argument list

```go
rArgs := flags.Exclude("key", args)  // rArgs -> ["--name=example", "--debug"]
rArgs := flags.Exclude("name", args) // rArgs -> ["--key", "123", "--debug"]
rArgs := flags.Exclude("foo", args)  // rArgs -> ["--key", "123", "--name=example "--debug"]
```

Is a flag in its valid representation (compatible with the flag package)?

```go
flags.Valid("foo")      // false
flags.Valid("--foo")    // true
flags.Valid("-key=val") // true
flags.Valid("-name=")   // true
```

Parse a flag and return the name

```go
name, _ := flags.Parse("foo")        // returns error, because foo is invalid
name, _ := flags.Parse("--foo")      // name -> "foo
name, _ := flags.Parse("-foo")       // name -> "foo
name, _ := flags.Parse("-foo=value") // name -> "foo
name, _ := flags.Parse("-foo=")      // name -> "foo
```

## flag.Value implementations (StringSlice and IntSlice)

Parse into a `[]string` or `[]int` variable

```go
os.Args = []string{"cmd", "--key", "123,456", "--regions", "us-east-1,eu-west-1"}

var regions []string
var ids []int

flags.StringSliceVar(&regions, []string{}, "to", "Regions to be used")
flags.IntSliceVar(&ids, []int{678}, "ids", "Servers to be used")
flag.Parse()

fmt.Println(regions) // prints: ["us-east-1", "eu-west-1"]
fmt.Println(ids)     // prints: [123,456]
```

Or plug it into a `flag.FlagSet` instance:

```go
args := []string{"--key", "123,456", "--regions", "us-east-1,eu-west-1"}

var regions []string
var ids []int

f := flag.NewFlagSet()
f.Var(flags.NewStringSlice(nil, &regions), "to", "Regions to be used")
f.Var(flags.NewIntSlice(nil, &ids), "to", "Regions to be used")
f.Parse(args)

fmt.Println(regions) // prints: ["us-east-1", "eu-west-1"]
fmt.Println(ids)     // prints: [123,456]
```
