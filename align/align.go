package align

import (
	"bufio"
	"io"
	"strings"
	"unicode/utf8"
)

// Commonly used file delimiters or alignment characters
const (
	Comma    = ','
	VertPipe = '|'
	Star     = '*'
	Tab      = '\t'
	Equal    = '='
	GThan    = '>'
	LThan    = '<'
	Hyphen   = '-'
	Plus     = '+'
	RParen   = ')'
	LParen   = '('
)

// Various values to be used by consuming programs
const (
	SingleSpace = " "
	DoubleQuote = "\""
	SingleQuote = "'"
)

// Alignable ...
type Alignable interface {
	ColumnCounts() []string
	Export([]string)
}

type columnCount int

// Aligner scans input and writes output
type Aligner struct {
	S            *bufio.Scanner
	W            *bufio.Writer
	del          rune // delimiter
	columnCounts map[int]int
}

// NewAligner creates and initializes a ScanWriter with in and out as its initial Reader and Writer
// and sets del to the desired delimiter to be used for alignment.
// It is meant to read the contents of its io.Reader to determine the length of each field
// and output the results in an aligned format.
func NewAligner(in io.Reader, out io.Writer, delimiter rune) Alignable {
	return &Aligner{
		S:            bufio.NewScanner(in),
		W:            bufio.NewWriter(out),
		del:          delimiter,
		columnCounts: make(map[int]int),
	}
}

// Init accepts the same arguments as NewAligner.  It simply provides another option
// for initializing an Aligner which is already allocated.
func (a *Aligner) Init(in io.Reader, out io.Writer, delimiter rune) {
	a.S = bufio.NewScanner(in)
	a.W = bufio.NewWriter(out)
	a.del = delimiter
	a.columnCounts = make(map[int]int)
}

// ColumnCounts scans the input and determines the maximum length of each field based on
// the longest value for each field in all of the pertaining lines.
// All of the lines of the io.Reader are returned as a string slice.
func (a *Aligner) ColumnCounts() []string {
	var lines []string
	for a.S.Scan() {
		var columnNum int
		var temp int

		line := a.S.Text()

		for i, v := range line {
			temp += utf8.RuneLen(v)
			if v != a.del && i < len(line)-1 {
				continue
			}
			if temp > a.columnCounts[columnNum] {
				a.columnCounts[columnNum] = temp
			}
			columnNum++
			temp = 0
		}
		lines = append(lines, line)
	}

	return lines
}

// Export will pad each field in lines based on the Aligner's column counts
func (a *Aligner) Export(lines []string) {
	for _, line := range lines {
		words := strings.Split(line, string(a.del))

		var columnNum int

		for _, word := range words {
			// leading padding for all fields except for the first
			if columnNum > 0 {
				word = SingleSpace + word
			}
			for len(word) < a.columnCounts[columnNum] {
				word += SingleSpace
			}
			rCount, wordLen := utf8.RuneCountInString(word), len(word)
			if rCount < wordLen {
				for i := 0; i < wordLen-rCount; i++ {
					word += SingleSpace
				}
			}
			columnNum++

			// Do not add a delimiter to the last field
			// This also properly aligns the output even if there are lines with a different number of fields
			if columnNum == len(words) {
				a.W.WriteString(word + "\n")
				continue
			}
			a.W.WriteString(word + string(a.del))
		}
	}
	a.W.Flush()
}
