# true-up
_A general purpose package that aligns text_

[![GoDoc](https://img.shields.io/badge/api-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/Guitarbum722/true-up/align) 
[![Build Status](https://travis-ci.org/Guitarbum722/true-up.svg?branch=master)](https://travis-ci.org/Guitarbum722/true-up)

The focus of this package is to provide a fast, efficient, and useful library for aligning text.

### Included

* Fully capable library, waiting to be consumed by your program.
* A simple yet usefule CLI with options to specify your delimiter, input and output files, etc.
* Align by any string as your delimiter or separator, not just a single character.
* If your separator string is contained within the data itself, it can be escaped by specifying
a text qualifier.
* Right, Center or Left justification of each field.

_Why?_

Sometimes, it's just easier to align a CSV (or delimited file) by its delimiter and view the columns in your plain text editor ( which saves you from opening Excel!).

Another use is to align blocks of code by `=` or `=>`, etc.

### Install

```sh
$ go get github.com/Guitarbum722/true-up
```

### Usage - CLI examples

```
Usage: true-up [-sep] [-output] [-file] [-qual]
Options:
  -h | --help  : help
  -file        : input file.  If not specified, pipe input to stdin
  -output      : output file. (defaults to stdout)
  -qual        : text qualifier (if applicable)
  -sep         : delimiter. (defaults to ',')
  -left        : left justification. (default)
  -center      : center justification
  -right       : right justification
```

_Specify your input file, output file, delimiter._
*You can also pipe input to Stdin (if the `-file` option is provided, it will take precedence over Stdin)*
If no `-output` option is provided, Stdout will be used.

```sh
$ true-up -file input_file.csv -output output_file.csv

$ true-up -file input_file.csv -output 

$ cat awesome.csv | true-up
```

Do you have rows with a different number of fields?  This might be more common with code, but `true-up` doesn't care!

```
$ echo "field1|field2\nValue1|Value2\nCoolValue1|CoolValue2|CoolValue3" | true-up -sep \|
field1     | field2
Value1     | Value2
CoolValue1 | CoolValue2 | CoolValue3
```

### Usage - The True-Up library

Initialize your `Aligner` which returns an `Alignable`.
```go
func main() {
	aligner := align.NewAligner(strings.NewReader("one,two,three\nfour,five,six\nseven,eight,nine"),
		os.Stdout,
		",",
		align.TextQualifier{On: false})

	// update justification (default is JustifyLeft)
	aligner.UpdatePadding(align.PaddingOpts{Justification: align.JustifyCenter})

	lines := aligner.ColumnCounts()

	aligner.Export(lines)
}

Output:
one   , two   , three
four  , five  , six
seven , eight , nine

```

or create an `Aligner` and call `Init()`

```go
    s := &align.Aligner{}
    s.Init(input, output, sep)
```

```go
    lines := a.ColumnCounts()
    a.Export(lines)
```

### Contributions

If you have suggestions or discover a bug, please open an issue.  If you think you can make the fix,
please use the Fork / Pull Request on your feature branch approach.