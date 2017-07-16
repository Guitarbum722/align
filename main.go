package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/Guitarbum722/true-up/align"
	"github.com/fatih/flags"
)

func main() {
	if retval, err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(retval)
	}
}

func run() (int, error) {
	args := os.Args[1:]

	// defaults
	sep := ','
	var input io.Reader
	var output io.Writer
	var qu align.TextQualifier

	if flags.Has("sep", args) {
		delimiter, err := flags.Value("sep", args)
		if err != nil {
			return 1, err
		}
		sep = []rune(delimiter)[0]
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
			Qualifier: []rune(q)[0],
		}
	}

	sw := align.NewAligner(input, output, sep, qu)

	lines := sw.ColumnCounts()
	sw.Export(lines)

	return 0, nil
}
