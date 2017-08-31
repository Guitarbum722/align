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
  -i           override justification by column number (e.g. 2:center,5:right)
  `

func main() {
	if retval, err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(retval)
	}
}

var hFlag *bool
var helpFlag *bool
var fFlag *string
var oFlag *string
var qFlag *string
var sFlag *string
var dFlag *string
var aFlag *string
var cFlag *string
var iFlag *string

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usage)
	}

	hFlag = flag.Bool("h", false, usage)
	helpFlag = flag.Bool("help", false, usage)
	fFlag = flag.String("f", "", "")
	oFlag = flag.String("o", "", "")
	qFlag = flag.String("q", "", "")
	sFlag = flag.String("s", ",", "")
	dFlag = flag.String("d", "", "")
	aFlag = flag.String("a", "left", "")
	cFlag = flag.String("c", "", "")
	iFlag = flag.String("i", "", "")
}

func run() (int, error) {
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
	var outColumns []int
	var justifyOverrides = make(map[int]justification)

	if *iFlag != "" {
		c := strings.Split(*iFlag, ",")

		for _, v := range c {
			if strings.HasSuffix(v, ":right") || strings.HasSuffix(v, ":center") || strings.HasSuffix(v, ":left") {
				overrides := strings.Split(v, ":")
				v = overrides[0]

				num, err := strconv.Atoi(v)
				if err != nil {
					return 1, errors.New("make sure entry for -v are numbers with a justification separated by ':' (ie 1-right,3-center)")
				}

				switch overrides[1] {
				case "left":
					justifyOverrides[num] = JustifyLeft
				case "center":
					justifyOverrides[num] = JustifyCenter
				case "right":
					justifyOverrides[num] = JustifyRight
				}
			}
		}

		if len(justifyOverrides) < 1 {
			return 1, errors.New("make sure entry for -v are numbers with a justification separated by ':' (ie 1:right,3:center)")
		}
	}

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

	if *qFlag != "" {
		qu = TextQualifier{
			On:        true,
			Qualifier: *qFlag,
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

	aligner := newAlign(input, output, *sFlag, qu)

	switch *aFlag {
	case "left":
		aligner.updatePadding(PaddingOpts{Justification: JustifyLeft, columnOverride: justifyOverrides})
	case "right":
		aligner.updatePadding(PaddingOpts{Justification: JustifyRight, columnOverride: justifyOverrides})
	case "center":
		aligner.updatePadding(PaddingOpts{Justification: JustifyCenter, columnOverride: justifyOverrides})
	default:
		aligner.updatePadding(PaddingOpts{Justification: JustifyLeft, columnOverride: justifyOverrides})
	}
	aligner.filterColumns(outColumns)
	aligner.outputSep(*dFlag)

	aligner.Align()

	return 0, nil
}
