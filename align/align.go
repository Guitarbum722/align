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

type justification byte

// Left, Right or Center justification options
const (
	JustifyRight justification = iota
	JustifyCenter
	JustifyLeft
)

// Alignable ...
type Alignable interface {
	ColumnCounts() []string
	Export([]string)
	SplitWithQual(string, string, string) []string
	ColumnSize(int) int
	UpdatePadding(PaddingOpts)
}

// TextQualifier is used to configure the scanner to account for a text qualifier
type TextQualifier struct {
	On        bool
	Qualifier string
}

// PaddingOpts provides configurability for left/center/right justification and padding length
type PaddingOpts struct {
	Justification justification
}

// Aligner scans input and writes output
type Aligner struct {
	S            *bufio.Scanner
	W            *bufio.Writer
	sep          string // separator string or delimiter
	columnCounts map[int]int
	txtq         TextQualifier
	padOpts      PaddingOpts
}

// NewAligner creates and initializes a ScanWriter with in and out as its initial Reader and Writer
// and sets del to the desired delimiter to be used for alignment.
// It is meant to read the contents of its io.Reader to determine the length of each field
// and output the results in an aligned format.
// Left justification is used by default.  See UpdatePadding to set the justification.
func NewAligner(in io.Reader, out io.Writer, sep string, qu TextQualifier) Alignable {
	return &Aligner{
		S:            bufio.NewScanner(in),
		W:            bufio.NewWriter(out),
		sep:          sep,
		columnCounts: make(map[int]int),
		txtq:         qu,
		padOpts: PaddingOpts{
			Justification: JustifyLeft,
		},
	}
}

// Init accepts the same arguments as NewAligner.  It simply provides another option
// for initializing an Aligner which is already allocated.
func (a *Aligner) Init(in io.Reader, out io.Writer, sep string, qu TextQualifier) {
	a.S = bufio.NewScanner(in)
	a.W = bufio.NewWriter(out)
	a.sep = sep
	a.columnCounts = make(map[int]int)
	a.txtq = qu
}

// ColumnSize looks up the Aligner's columnCounts key with num and returns the value
// that was set by ColumnCounts().
// If num is not a valid key in Aligner.columnCounts, then -1 is returned.
func (a *Aligner) ColumnSize(num int) int {
	if _, ok := a.columnCounts[num]; !ok {
		return -1
	}
	return a.columnCounts[num]
}

// UpdatePadding uses PaddingOpts p to update the Aligner's padding options.
func (a *Aligner) UpdatePadding(p PaddingOpts) {
	a.padOpts = p
}

// FieldLen works in a similar manner to the standard lib function strings.Index().
// Instead of returning the index of the first instance of sep, it returns the length
// of s before the first index of sep.
func FieldLen(s, sep string) int {
	return genFieldLen(s, sep, "")
}

// FieldLenEscaped works in the same way as FieldLen, but a text qualifer string can
// be provided.  If qual is an empty string, then the behavior will be exactly the same
// as FieldLen.
func FieldLenEscaped(s, sep, qual string) int {
	return genFieldLen(s, sep, qual)
}

func genFieldLen(s, sep, qual string) int {
	i := 0
	if qual == "" || !strings.HasPrefix(s, qual) {
		i = strings.Index(s, sep)
	} else {
		i = strings.Index(s, qual+sep)

		if i == -1 {
			return len(s)
		}
		return len(s[:i+len(qual)])
	}

	if i == -1 {
		return len(s)
	}

	return len(s[:i])
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

		if a.txtq.On {
			for start := 0; start < len(line); {
				temp = FieldLenEscaped(line[start:], a.sep, a.txtq.Qualifier)
				start += temp + len(a.sep)
				if temp > a.columnCounts[columnNum] {
					a.columnCounts[columnNum] = temp
				}
				columnNum++
				temp = 0
			}
		} else {
			for start := 0; start < len(line); {
				temp = FieldLen(line[start:], a.sep)
				start += temp + len(a.sep)
				if temp > a.columnCounts[columnNum] {
					a.columnCounts[columnNum] = temp
				}
				columnNum++
				temp = 0
			}
		}

		lines = append(lines, line)
	}

	return lines
}

// Export will pad each field in lines based on the Aligner's column counts
func (a *Aligner) Export(lines []string) {
	for _, line := range lines {
		words := a.SplitWithQual(line, a.sep, a.txtq.Qualifier)

		var columnNum int
		for _, word := range words {
			word = pad(word, columnNum, a.columnCounts[columnNum], a.padOpts)
			columnNum++

			// Do not add a delimiter to the last field
			// This also properly aligns the output even if there are lines with a different number of fields
			if columnNum == len(words) {
				a.W.WriteString(word + "\n")
				continue
			}
			a.W.WriteString(word + a.sep)
		}
	}
	a.W.Flush()
}

// pad s based on the supplied PaddingOpts
func pad(s string, columnNum, count int, p PaddingOpts) string {
	padLength := countPadding(s, count)

	switch p.Justification {
	case JustifyLeft:
		s = trailingPad(s, padLength)
	case JustifyRight:
		s = leadingPad(s, padLength)
	case JustifyCenter:

	default:
		s = trailingPad(s, padLength)
	}

	// at least one space to pad every field after the delimiter for readability
	if columnNum > 0 {
		s = SingleSpace + s
	}
	s = s + SingleSpace

	return s
}

// determines the length of the padding needed
func countPadding(s string, count int) int {
	padLength := count - len(s)
	rCount, wordLen := utf8.RuneCountInString(s), len(s)
	if rCount < wordLen {
		padLength += wordLen - rCount
	}
	return padLength
}

// prepends padding
func leadingPad(s string, padLen int) string {
	pad := make([]byte, 0, padLen)

	for len(pad) < padLen {
		pad = append(pad, ' ')
	}

	return string(pad) + s
}

// appends padding
func trailingPad(s string, padLen int) string {
	pad := make([]byte, 0, padLen)

	for len(pad) < padLen {
		pad = append(pad, ' ')
	}

	return s + string(pad)
}

// SplitWithQual basically works like the standard strings.Split() func, but will consider a text qualifier if set.
func (a *Aligner) SplitWithQual(s, sep, qual string) []string {

	if !a.txtq.On {
		return strings.Split(s, sep) // use standard Split() method if no qualifier is considered
	}
	var words = make([]string, 0, strings.Count(s, sep))

	for start := 0; start < len(s); {
		count := genFieldLen(s[start:], sep, qual)
		words = append(words, s[start:start+count])
		start += count + len(sep)
	}

	return words
}
