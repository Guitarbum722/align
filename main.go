package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
)

const usage = `Usage: align [-h] [-f] [-o] [-q] [-s] [-d] [-a]
Options:
  -h | --help  help
  -f           input file.  If not specified, pipe input to stdin
  -o           output file. (default: stdout)
  -q           text qualifier (if applicable)
  -s           delimiter (default: ',')
  -d           output delimiter (defaults to the value of sep)
  -a           <left>, <right>, <center> justification (default: left)
`

var hFlag *bool
var helpFlag *bool
var fFlag *string
var oFlag *string
var qFlag *string
var sFlag *string
var dFlag *string
var aFlag *string

func main() {
	if retval, err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(retval)
	}
}

func run() (int, error) {

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usage)
	}

	parseFlags() // parse command line args

	// check for piped input, but use specified input file if supplied
	fi, _ := os.Stdin.Stat()
	isPiped := (fi.Mode() & os.ModeCharDevice) == 0

	if *hFlag == true || *helpFlag == true {
		return 1, errors.New(usage)
	}
	if !isPiped {
		if len(os.Args[1:]) == 0 {
			return 1, errors.New(usage)
		}
	}

	var input io.Reader
	var output io.Writer
	var qu TextQualifier

	if isPiped {
		if *fFlag != "" {
			f, err := os.Open(*fFlag)
			if err != nil {
				return 1, err
			}
			defer f.Close()
			input = f
		} else {
			input = os.Stdin
		}
	} else {
		if *fFlag != "" {
			f, err := os.Open(*fFlag)
			if err != nil {
				return 1, err
			}
			defer f.Close()
			input = f
		} else {
			return 1, errors.New("no input provided \n" + usage)
		}
	}

	if *oFlag != "" {
		f, err := os.Create(*oFlag)
		if err != nil {
			return 1, err
		}
		defer f.Close()
		output = f
	} else {
		output = os.Stdout
	}

	if *qFlag != "" {
		qu = TextQualifier{
			On:        true,
			Qualifier: *qFlag,
		}
	}

	sw := NewAligner(input, output, *sFlag, qu)

	switch *aFlag {
	case "left":
		sw.UpdatePadding(PaddingOpts{Justification: JustifyLeft})
	case "right":
		sw.UpdatePadding(PaddingOpts{Justification: JustifyRight})
	case "center":
		sw.UpdatePadding(PaddingOpts{Justification: JustifyCenter})
	default:
		sw.UpdatePadding(PaddingOpts{Justification: JustifyLeft})
	}

	lines := sw.ColumnCounts()

	sw.OutputSep(*dFlag)
	sw.Export(lines)

	return 0, nil
}

func parseFlags() {
	hFlag = flag.Bool("h", false, usage)
	helpFlag = flag.Bool("help", false, usage)
	fFlag = flag.String("f", "", "")
	oFlag = flag.String("o", "", "")
	qFlag = flag.String("q", "", "")
	sFlag = flag.String("s", ",", "")
	dFlag = flag.String("d", "", "")
	aFlag = flag.String("a", "left", "")

	flag.Parse()

	if *dFlag == "" {
		*dFlag = *sFlag
	}
}
