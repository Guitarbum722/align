package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
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
  -c           output specific fields (default: all fields)
  `

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

	hFlag := flag.Bool("h", false, usage)
	helpFlag := flag.Bool("help", false, usage)
	fFlag := flag.String("f", "", "")
	oFlag := flag.String("o", "", "")
	qFlag := flag.String("q", "", "")
	sFlag := flag.String("s", ",", "")
	dFlag := flag.String("d", "", "")
	aFlag := flag.String("a", "left", "")
	cFlag := flag.String("c", "", "")

	flag.Parse()

	if *dFlag == "" {
		*dFlag = *sFlag
	}

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

	var outColumns []int

	if *cFlag != "" {
		c := strings.Split(*cFlag, ",")
		outColumns = make([]int, 0, len(c))

		// validate specified field numbers and sort them
		for _, v := range c {
			num, err := strconv.Atoi(v)
			if err != nil {
				return 1, errors.New("make sure entry for -c are numbers (ie 1,2,5,7)")
			}
			if num > 0 {
				outColumns = append(outColumns, num)
			}
		}
		sort.Ints(outColumns)
	}

	aligner := newAlign(input, output, *sFlag, qu)

	switch *aFlag {
	case "left":
		aligner.updatePadding(PaddingOpts{Justification: JustifyLeft})
	case "right":
		aligner.updatePadding(PaddingOpts{Justification: JustifyRight})
	case "center":
		aligner.updatePadding(PaddingOpts{Justification: JustifyCenter})
	default:
		aligner.updatePadding(PaddingOpts{Justification: JustifyLeft})
	}
	aligner.filterColumns(outColumns)
	aligner.outputSep(*dFlag)

	aligner.Align()

	return 0, nil
}
