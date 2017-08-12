package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/fatih/flags"
)

const usage = `Usage: align [-h] [-f] [-o] [-q] [-s] [-d] [-a]
Options:
  -h | --help  : help
  -f           : input file.  If not specified, pipe input to stdin
  -o           : output file. (defaults to stdout)
  -q           : text qualifier (if applicable)
  -s           : delimiter. (defaults to ',')
  -d           : output delimiter (defaults to the value of sep)
  -a           : <left>, <right>, <center> justification (default: left)
`

func main() {
	if retval, err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(retval)
	}
}

func run() (int, error) {
	var args []string
	if len(args) == 0 {
		return 1, errors.New(usage)
	}
	args = os.Args[1:]

	sep := "," // default
	var outSep string
	var input io.Reader
	var output io.Writer
	var qu TextQualifier

	if flags.Has("h", args) || flags.Has("help", args) {
		return 1, errors.New(usage)
	}

	if flags.Has("s", args) {
		val, err := flags.Value("s", args)
		if !validArg(err, val) {
			return 1, errors.New("invalid entry for -s \n" + usage)
		}
		sep = val
	}

	if flags.Has("d", args) {
		outSep = sep
		val, err := flags.Value("d", args)
		if !validArg(err, val) {
			return 1, errors.New("invalid entry for -d \n" + usage)
		}
		outSep = val
	} else {
		outSep = sep
	}

	// check for piped input, but use specified input file if supplied
	fi, _ := os.Stdin.Stat()
	if (fi.Mode() & os.ModeCharDevice) == 0 {
		if flags.Has("f", args) {
			fn, err := flags.Value("f", args)
			if !validArg(err, fn) {
				return 1, errors.New("invalid entry for -f \n" + usage)
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
		if flags.Has("f", args) {
			fn, err := flags.Value("f", args)
			if !validArg(err, fn) {
				return 1, errors.New("invalid entry for -f \n" + usage)
			}

			f, err := os.Open(fn)
			if err != nil {
				return 1, err
			}
			defer f.Close()
			input = f
		} else {
			return 1, errors.New("no input provided \n" + usage)
		}
	}

	// if --output flag is not provided with a file name, then use Stdout
	if flags.Has("o", args) {
		fn, err := flags.Value("o", args)
		if !validArg(err, fn) {
			return 1, errors.New("invalid entry for -o \n" + usage)
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

	if flags.Has("q", args) {
		q, err := flags.Value("q", args)
		if !validArg(err, q) {
			return 1, errors.New("invalid entry for -q \n" + usage)
		}

		qu = TextQualifier{
			On:        true,
			Qualifier: q,
		}
	}

	sw := NewAligner(input, output, sep, qu)

	if flags.Has("a", args) {
		val, err := flags.Value("a", args)
		if !validArg(err, val) {
			return 1, errors.New("invalid entry for -a \n" + usage)
		}
		switch val {
		case "left":
			sw.UpdatePadding(PaddingOpts{Justification: JustifyLeft})
		case "right":
			sw.UpdatePadding(PaddingOpts{Justification: JustifyRight})
		case "center":
			sw.UpdatePadding(PaddingOpts{Justification: JustifyCenter})
		default:
			sw.UpdatePadding(PaddingOpts{Justification: JustifyLeft})
		}
	}

	lines := sw.ColumnCounts()

	sw.OutputSep(outSep)
	sw.Export(lines)

	return 0, nil
}

func validArg(err error, arg string) bool {
	if err != nil || arg == "" {
		return false
	}
	return true
}
