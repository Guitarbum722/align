package align

import (
	"bufio"
	"io"
	"strings"

	"github.com/mattn/go-runewidth"
)

// Justification is used to set the alignment of the column
// contents itself along the right, left, or center.
type Justification byte

// Left, Right or Center Justification options.
const (
	JustifyRight Justification = iota + 1
	JustifyCenter
	JustifyLeft
)

// TextQualifier is used to configure the scanner to account for a text qualifier.
type TextQualifier struct {
	On        bool
	Qualifier string
}

// PaddingOpts provides configurability for left/center/right Justification and padding length.
type PaddingOpts struct {
	Justification  Justification
	ColumnOverride map[int]Justification //override the Justification of specified columns
	Pad            int                   // padding surrounding the separator
}

// Align scans input and writes output with aligned text.
type Align struct {
	scanner      *bufio.Scanner
	writer       *bufio.Writer
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
// Left Justification is used by default.  See UpdatePadding to set the Justification.
func NewAlign(in io.Reader, out io.Writer, sep string, qu TextQualifier) *Align {
	return &Align{
		scanner:      bufio.NewScanner(in),
		writer:       bufio.NewWriter(out),
		sep:          sep,
		sepOut:       sep,
		columnCounts: make(map[int]int),
		txtq:         qu,
		padOpts: PaddingOpts{
			//defaults
			Justification: JustifyLeft,
			Pad:           1,
		},
	}
}

// OutputSep sets the output separator string with outsep if a different value from the input sep is desired.
func (a *Align) OutputSep(outsep string) {
	a.sepOut = outsep
}

// Align determines the length of each field of text around the configured delimiter and aligns all of the
// text by the delimiter.
func (a *Align) Align() {
	lines := a.columnLength()
	a.export(lines)
}

// columnSize looks up the Align's columnCounts key with num and returns the value
// that was set by ColumnCounts().
// If num is not a valid key in Align.columnCounts, then -1 is returned.
func (a *Align) columnSize(num int) int {
	if _, ok := a.columnCounts[num]; !ok {
		return -1
	}
	return a.columnCounts[num]
}

// UpdatePadding uses PaddingOpts p to update the Align's padding options.
func (a *Align) UpdatePadding(p PaddingOpts) {
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
	for a.scanner.Scan() {
		var columnNum int
		var temp int

		line := a.scanner.Text()

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

// export will pad each field in lines based on the Align's column counts.
func (a *Align) export(lines []string) {
	if a.padOpts.Pad < 0 {
		a.padOpts.Pad = 0
	}

	surroundingPad := make([]byte, 0, a.padOpts.Pad)

	for i := 0; i < a.padOpts.Pad; i++ {
		surroundingPad = append(surroundingPad, ' ')
	}

	for _, line := range lines {
		words := a.splitWithQual(line, a.sep, a.txtq.Qualifier)

		var columnNum int
		var tempColumn int // used for call to pad() to incorporate column filtering
		for _, word := range words {
			if a.filterLen > 0 {
				if !contains(a.filter, columnNum+1) {
					columnNum++
					if columnNum == len(words) {
						a.writer.WriteString("\n")
					}
					continue
				}
			}

			j := a.padOpts.Justification

			// override Justification for the specified columnNum in the key for the PaddingOpts.columnOverride map
			if len(a.padOpts.ColumnOverride) > 0 {
				for k, v := range a.padOpts.ColumnOverride {
					if k == columnNum+1 {
						j = v
					}
				}
			}

			word = applyPadding(word, tempColumn, a.columnCounts[columnNum], j, string(surroundingPad))
			columnNum++
			tempColumn++

			// Do not add a delimiter to the last field
			// This also properly aligns the output even if there are lines with a different number of fields
			if a.filterLen > 0 && a.filter[a.filterLen-1] == columnNum {
				a.writer.WriteString(word + "\n")
				break
			} else if columnNum == len(words) {
				a.writer.WriteString(word + "\n")
				break
			}
			a.writer.WriteString(word + a.sepOut)
		}
	}
	a.writer.Flush()
}

// applyPadding rebuilds word by adding padding appropriately based on the
// desired justification, the overall column count and the supplied surrounding
// padding string.
func applyPadding(word string, columnNum, count int, just Justification, surroundingPad string) string {
	padLength := countPadding(word, count)

	var sb strings.Builder
	sb.Grow(padLength + len(word) + (len(surroundingPad) * 2))

	// add surrounding pad to beginning of column (except for the 1st column)
	if len(surroundingPad) > 0 {
		if columnNum > 0 {
			sb.WriteString(surroundingPad)
		}
	}

	switch just {
	case JustifyLeft:
		sb.WriteString(word)
		for i := 0; i < padLength; i++ {
			sb.WriteByte(' ')
		}
	case JustifyRight:
		for i := 0; i < padLength; i++ {
			sb.WriteByte(' ')
		}
		sb.WriteString(word)
	case JustifyCenter:
		// not much of a point to 'center' justification with such a small padding; default it if <= 2.
		if padLength > 2 {
			for i := 0; i < (padLength - (padLength / 2)); i++ {
				sb.WriteByte(' ')
			}
			sb.WriteString(word)
			for i := 0; i < (padLength / 2); i++ {
				sb.WriteByte(' ')
			}
			// trailingPad(&sb, padLength/2)
			// sb.WriteString(word)
			// leadingPad(&sb, padLength-(padLength/2))
		} else {
			sb.WriteString(word)
			for i := 0; i < padLength; i++ {
				sb.WriteByte(' ')
			}
		}
	}

	// add surrounding pad to end of column
	if len(surroundingPad) > 0 {
		sb.WriteString(surroundingPad)
	}
	return sb.String()
}

// determines the length of the padding needed.
func countPadding(s string, count int) int {
	padLength := count - len(s)
	rCount, wordLen := runewidth.StringWidth(s), len(s)
	if rCount < wordLen {
		padLength += wordLen - rCount
	}
	return padLength
}

// prepends padding.
func leadingPad(sb *strings.Builder, padLen int) {
	for i := 0; i < padLen; i++ {
		sb.WriteByte(' ')
	}
}

// appends padding.
func trailingPad(sb *strings.Builder, padLen int) {
	for i := 0; i < padLen; i++ {
		sb.WriteByte(' ')
	}
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

// FilterColumns sets which column numbers should be output.
func (a *Align) FilterColumns(c []int) {
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
