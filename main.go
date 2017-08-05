package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/Guitarbum722/true-up/align"
	"github.com/fatih/flags"
)

const usage = `Usage: true-up [-sep] [-output] [-file] [-qual] [-jstfy]
Options:
  -h | --help  : help
  -file        : input file.  If not specified, pipe input to stdin
  -output      : output file. (defaults to stdout)
  -qual        : text qualifier (if applicable)
  -sep         : delimiter. (defaults to ',')
  -outsep      : output delimiter (defaults to the value of sep)
  -left        : left justification. (default)
  -center      : center justification
  -right       : right justification
`

func main() {
	if retval, err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(retval)
	}
}

func run() (int, error) {
	args := os.Args[1:]

	// defaults
	sep := ","
	var outSep string
	var input io.Reader
	var output io.Writer
	var qu align.TextQualifier

	if flags.Has("h", args) || flags.Has("help", args) {
		return 1, errors.New(usage)
	}

	if flags.Has("sep", args) {
		if len(args) < 2 {
			return 1, errors.New("argument to -sep required")
		}
		delimiter, err := flags.Value("sep", args)
		if err != nil {
			return 1, err
		}
		sep = delimiter
	}

	if flags.Has("outsep", args) {
		if len(args) < 2 {
			return 1, errors.New("argument to -outsep required")
		}
		delimiter, err := flags.Value("outsep", args)
		if err != nil {
			return 1, err
		}
		outSep = delimiter
	} else {
		outSep = sep
	}

	// check for piped input, but use specified input file if supplied
	fi, _ := os.Stdin.Stat()
	if (fi.Mode() & os.ModeCharDevice) == 0 {
		if flags.Has("file", args) {
			fn, err := flags.Value("file", args)
			if err != nil {
				return 1, err
			}
			f, err := os.Open(fn)
			if err != nil {
				return 1, err
			}
			defer f.Close()
			input = f
		} else {
			input = os.Stdin
		}
	} else {
		if flags.Has("file", args) {
			fn, err := flags.Value("file", args)
			if err != nil {
				return 1, err
			}
			f, err := os.Open(fn)
			if err != nil {
				return 1, err
			}
			defer f.Close()
			input = f
		} else {
			return 1, errors.New("no input provided")
		}
	}

	// if --output flag is not provided with a file name, then use Stdout
	if flags.Has("output", args) {
		if len(args) < 2 {
			return 1, errors.New("argument to -output required")
		}
		fn, err := flags.Value("output", args)
		if err != nil {
			return 1, err
		}
		f, err := os.Create(fn)
		if err != nil {
			return 1, err
		}
		defer f.Close()
		output = f
	} else {
		output = os.Stdout
	}

	if flags.Has("qual", args) {
		q, err := flags.Value("qual", args)
		if err != nil {
			return 1, err
		}

		qu = align.TextQualifier{
			On:        true,
			Qualifier: q,
		}
	}

	sw := align.NewAligner(input, output, sep, qu)

	if flags.Has("left", args) {
		sw.UpdatePadding(align.PaddingOpts{Justification: align.JustifyLeft})
	} else if flags.Has("right", args) {
		sw.UpdatePadding(align.PaddingOpts{Justification: align.JustifyRight})
	} else if flags.Has("center", args) {
		sw.UpdatePadding(align.PaddingOpts{Justification: align.JustifyCenter})
	}

	lines := sw.ColumnCounts()

	sw.OutputSep(outSep)
	sw.Export(lines)

	return 0, nil
}
