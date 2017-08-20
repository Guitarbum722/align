package main

import (
	"bufio"
	"io"
	"strings"
	"unicode/utf8"
)

const singleSpace = " "

type justification byte

// Left, Right or Center justification options
const (
	JustifyRight justification = iota
	JustifyCenter
	JustifyLeft
)

// Aligner aligns text based on configuration and options
type Aligner interface {
	Align()
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

// Align scans input and writes output with aligned text
type Align struct {
	S            *bufio.Scanner
	W            *bufio.Writer
	sep          string // separator string or delimiter
	sepOut       string
	columnCounts map[int]int
	txtq         TextQualifier
	padOpts      PaddingOpts
	filter       []int
	filterLen    int
}

// NewAlign creates and initializes a ScanWriter with in and out as its initial Reader and Writer
// and sets del to the desired delimiter to be used for alignment.
// It is meant to read the contents of its io.Reader to determine the length of each field
// and output the results in an aligned format.
// Left justification is used by default.  See UpdatePadding to set the justification.
func newAlign(in io.Reader, out io.Writer, sep string, qu TextQualifier) *Align {
	return &Align{
		S:            bufio.NewScanner(in),
		W:            bufio.NewWriter(out),
		sep:          sep,
		sepOut:       sep,
		columnCounts: make(map[int]int),
		txtq:         qu,
		padOpts: PaddingOpts{
			Justification: JustifyLeft, // default
		},
	}
}

// outputSep sets the output separator string with outsep if a different value from the input sep is desired.
func (a *Align) outputSep(outsep string) {
	a.sepOut = outsep
}

// Align determines the length of each field of text around the configured delimiter and aligns all of the
// text by the delimiter.
func (a *Align) Align() {
	lines := a.columnLength()
	a.export(lines)
}

// columnSize looks up the Aligner's columnCounts key with num and returns the value
// that was set by ColumnCounts().
// If num is not a valid key in Aligner.columnCounts, then -1 is returned.
func (a *Align) columnSize(num int) int {
	if _, ok := a.columnCounts[num]; !ok {
		return -1
	}
	return a.columnCounts[num]
}

// updatePadding uses PaddingOpts p to update the Aligner's padding options.
func (a *Align) updatePadding(p PaddingOpts) {
	a.padOpts = p
}

// fieldLen works in a similar manner to the standard lib function strings.Index().
// Instead of returning the index of the first instance of sep, it returns the length
// of s before the first index of sep.
func fieldLen(s, sep string) int {
	return genFieldLen(s, sep, "")
}

// fieldLenEscaped works in the same way as FieldLen, but a text qualifer string can
// be provided.  If qual is an empty string, then the behavior will be exactly the same
// as FieldLen.
func fieldLenEscaped(s, sep, qual string) int {
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

// columnLength scans the input and determines the maximum length of each field based on
// the longest value for each field in all of the pertaining lines.
// All of the lines of the io.Reader are returned as a string slice.
func (a *Align) columnLength() []string {
	var lines []string
	for a.S.Scan() {
		var columnNum int
		var temp int

		line := a.S.Text()

		if a.txtq.On {
			for start := 0; start < len(line); {
				temp = fieldLenEscaped(line[start:], a.sep, a.txtq.Qualifier)
				start += temp + len(a.sep)
				if temp > a.columnCounts[columnNum] {
					a.columnCounts[columnNum] = temp
				}
				columnNum++
				temp = 0
			}
		} else {
			for start := 0; start < len(line); {
				temp = fieldLen(line[start:], a.sep)
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

// export will pad each field in lines based on the Aligner's column counts
func (a *Align) export(lines []string) {
	for _, line := range lines {
		words := a.splitWithQual(line, a.sep, a.txtq.Qualifier)

		var columnNum int
		var tempColumn int // used for call to pad() to incorporate column filtering
		for _, word := range words {
			if a.filterLen > 0 {
				if !contains(a.filter, columnNum+1) {
					columnNum++
					if columnNum == len(words) {
						a.W.WriteString("\n")
					}
					continue
				}
			}

			word = pad(word, tempColumn, a.columnCounts[columnNum], a.padOpts)
			columnNum++
			tempColumn++

			// Do not add a delimiter to the last field
			// This also properly aligns the output even if there are lines with a different number of fields
			if a.filterLen > 0 && a.filter[a.filterLen-1] == columnNum {
				a.W.WriteString(word + "\n")
				break
			} else if columnNum == len(words) {
				a.W.WriteString(word + "\n")
				break
			}
			a.W.WriteString(word + a.sepOut)
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
		if padLength > 2 {
			s = trailingPad(s, padLength/2)
			s = leadingPad(s, padLength-(padLength/2))
		} else {
			s = trailingPad(s, padLength)
		}
	default:
		s = trailingPad(s, padLength)
	}

	// at least one space to pad every field after the delimiter for readability
	if columnNum > 0 {
		s = singleSpace + s
	}
	s = s + singleSpace

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

// splitWithQual basically works like the standard strings.Split() func, but will consider a text qualifier if set.
func (a *Align) splitWithQual(s, sep, qual string) []string {
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

func (a *Align) filterColumns(c []int) {
	a.filter = c
	a.filterLen = len(c)
}

func contains(nums []int, num int) bool {
	for _, v := range nums {
		if v == num {
			return true
		}
	}
	return false
}
