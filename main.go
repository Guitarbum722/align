package main

import (
	"bufio"
	"log"
	"os"
	"strings"
	"unicode/utf8"
)

const (
	delimiter = ','
	fn        = "things.csv"
)

func main() {
	var columnNum int
	var columnCounts = make(map[int]int)
	var lines []string
	f, err := os.Open(fn)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()
	out, err := os.Create("out" + fn)
	if err != nil {
		log.Fatalln(err)
	}
	defer out.Close()
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		temp := 0
		columnNum = 0
		line := scanner.Text()
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

	w := bufio.NewWriter(out)
	// w := bufio.NewWriter(os.Stdout)
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
				w.WriteString(word + string(delimiter))
				continue
			}
			w.WriteString(word)
		}
		w.WriteByte('\n')
	}
	w.Flush()
}
