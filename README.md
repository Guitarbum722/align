# true-up
_A general purpose package that aligns text_

The focus of this package is to provide a fast, efficient, and useful library for aligning text.

### Install

```sh
$ go get github.com/Guitarbum722/true-up
```

### Included

* Fully capable library, waiting to be consume by your program.
* A simple CLI with options to specify your delimiter, input and output files, etc.

_Why?_

Sometimes, it's just easier to align a CSV (or delimited file) by its delimiter and view the columns in your plain text editor ( which saves you from opening Excel!).

Another use is to align blocks of code by `=` or `->`, etc.


### Usage - CLI examples
_Specify your input file, output file, delimiter._
*You can also pipe input to Stdin (if the `-file` option is provided, it will overwrite Stdin)*
If no `-output` option is provided, Stdout will be used.

```sh
$ true-up -file input_file.csv -output output_file.csv

$ true-up -file input_file.csv -output 
```

Do you have rows with a different number of fields?  This might be more common with code, but `true-up` doesn't care!

```
$ echo "field1|field2\nValue1|Value2\nCoolValue1|CoolValue2|CoolValue3" | true-up -sep \|
field1     | field2
Value1     | Value2
CoolValue1 | CoolValue2| CoolValue3
```