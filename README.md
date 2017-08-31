# align
_A general purpose application that aligns text_

[![GoDoc](https://img.shields.io/badge/api-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/Guitarbum722/align) 
[![Build Status](https://travis-ci.org/Guitarbum722/align.svg?branch=master)](https://travis-ci.org/Guitarbum722/align)
[![Go Report Card](https://goreportcard.com/badge/github.com/Guitarbum722/align)](https://goreportcard.com/report/github.com/Guitarbum722/align)
[![Coverage Status](https://coveralls.io/repos/github/Guitarbum722/align/badge.svg?branch=master)](https://coveralls.io/github/Guitarbum722/align?branch=master)

The focus of this application is to provide a fast, efficient, and useful tool for aligning text.
Its creation is the result of inspiration from several other amazing alignment tools, namely [`column`](http://man7.org/linux/man-pages/man1/column.1.html) or the Sublime Text plugin [`AlignTab`](https://github.com/randy3k/AlignTab).

### Included

* A simple yet useful CLI with options to specify your delimiter, input and output files, etc.
* Align by any string as your delimiter or separator, not just a single character.
* If your separator string is contained within the data itself, it can be escaped by specifying a text qualifier.
* Right, Center, or Left justification of each field.

_Why?_

Sometimes, it's just easier to align a CSV (or delimited file) by its delimiter and view the columns in your plain text editor (which saves you from opening Excel!).

Another use is to align blocks of code by `=` or `=>`, etc.

### Install

```sh
$ go get github.com/Guitarbum722/align
$ make install

$ # build all binaries
$ make release
```

### Usage - CLI examples

```
Usage: align [-h] [-f] [-o] [-q] [-s] [-d] [-a]
Options:
  -h | --help  help
  -f           input file.  If not specified, pipe input to stdin
  -o           output file. (default: stdout)
  -q           text qualifier (if applicable)
  -s           delimiter (default: ',')
  -d           output delimiter (defaults to the value of sep)
  -a           <left>, <right>, <center> justification (default: left)
  -c           output specific fields (default: all fields)
  -i           override justification by column number (e.g. 2:center,5:right)
```

_Specify your input file, output file, delimiter._
*You can also pipe input to stdin (if the `-f` option is provided, it will take precedence over Stdin)*
If no `-o` option is provided, stdout will be used.

```sh
$ align -f input_file.csv -o output_file.csv

$ align -f input_file.csv -o 

$ cat awesome.csv | align
```

Do you have rows with a different number of fields?  This might be more common with code, but `align` doesn't care!

```
$ echo "field1|field2\nValue1|Value2\nCoolValue1|CoolValue2|CoolValue3" | align -s \|
field1     | field2
Value1     | Value2
CoolValue1 | CoolValue2 | CoolValue3
```

Column filtering (specifiy output fields and optionally override the justification of the output fields).  This might be useful if you would like to display a dollar amount or number field differently.  The specified fields are indexed at 1.

```sh
# output fields 1,3,5 justified 'right'
$ cat file.csv | align -a right -c 1,3,5

# output fields 1,2,3,7,8 with default justification (left) except for field 7, which is right justified
$ cat file.csv | align -c 1,2,3,7,8 -v 1-right,7-center

#output all fields by default, with right justification, with overridden justification on certain columns
$ cat file.csv | align -a right -v 1-center,5-left
```

Support for worldwide characters.
```
first          , last              , middle  , email
paul           , danny             ,  かど    , や製油
```

It is perfectly acceptable to even use emojis as your input/output delimiters.
```
first  😮 last     😮 email
Hector 😮 Gonzalez 😮 h.g@nothing.com
```

### Contributions

If you have suggestions or discover a bug, please open an issue.  If you think you can make the fix, please use the Fork / Pull Request on your feature branch approach.
