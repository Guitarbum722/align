package align

import (
	"bufio"
	"io"
)

// ScanWriter scans input and writes output
type ScanWriter struct {
	S *bufio.Scanner
	W *bufio.Writer
}

func NewScanWriter(in io.Reader, out io.Writer) *ScanWriter {
	return &ScanWriter{
		S: bufio.NewScanner(in),
		W: bufio.NewWriter(out),
	}
}
