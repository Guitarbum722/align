package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"
	"unicode/utf8"
)

const (
	delimiter = ','
	fn        = "things.csv"
)

// ScanWriter scans input and writes output
type ScanWriter struct {
	s *bufio.Scanner
	w *bufio.Writer
}

func newScanWriter(in io.Reader, out io.Writer) *ScanWriter {
	return &ScanWriter{
		s: bufio.NewScanner(in),
		w: bufio.NewWriter(out),
	}
}

func main() {
	var columnNum int
	var columnCounts = make(map[int]int)
	var lines []string
	var sw *ScanWriter

	fi, _ := os.Stdin.Stat()
	if (fi.Mode() & os.ModeCharDevice) == 0 {
		switch len(os.Args) {
		case 2:
			f, err := os.Create(os.Args[1])
			if err != nil {
				log.Fatalln(err)
			}
			defer f.Close()
			sw = newScanWriter(os.Stdin, f)
		default:
			sw = newScanWriter(os.Stdin, os.Stdout)
		}
	} else {
		switch len(os.Args) {
		case 2:
			f, err := os.Open(os.Args[1])
			if err != nil {
				log.Fatalln(err)
			}
			defer f.Close()
			sw = newScanWriter(f, os.Stdout)
		case 3:
			f, err := os.Open(os.Args[1])
			if err != nil {
				log.Fatalln(err)
			}
			defer f.Close()
			out, err := os.Create(os.Args[2])
			if err != nil {
				log.Fatalln(err)
			}
			defer out.Close()
			sw = newScanWriter(f, out)
		default:
			sw = newScanWriter(os.Stdin, os.Stdout)
		}
	}

	for sw.s.Scan() {
		temp := 0
		columnNum = 0
		line := sw.s.Text()
		for i, v := range line {
			temp += utf8.RuneLen(v)
			if v != delimiter && i < len(line)-1 {
				continue
			}
			if temp > columnCounts[columnNum] {
				columnCounts[columnNum] = temp
			}
			columnNum++
			temp = 0
		}
		lines = append(lines, line)
	}

	for _, line := range lines {
		words := strings.Split(line, string(delimiter))
		columnNum = 0
		for _, word := range words {
			for len(word) < columnCounts[columnNum] {
				word += " "
			}
			rCount, wordLen := utf8.RuneCountInString(word), len(word)
			if rCount < wordLen {
				for i := 0; i < wordLen-rCount; i++ {
					word += " "
				}
			}
			columnNum++
			// since columnNum was just incremented, do not add a comma to the last field
			if _, ok := columnCounts[columnNum]; ok {
				sw.w.WriteString(word + string(delimiter))
				continue
			}
			sw.w.WriteString(word)
		}
		sw.w.WriteByte('\n')
	}
	sw.w.Flush()
}
