package main

import (
	"bufio"
	_ "fmt"
	"log"
	"os"
	"strings"
	"unicode/utf8"
)

const (
	delimiter = ','
	fn        = "unicode_things.csv"
)

func main() {
	// 	s := `FirstName,MiddleName,LastName,Email
	// John,K,Moore,jkdoe@nothing.com
	// Winston,Moby,Kazowski,Winston@nothing.net
	// Dilbert,Sylvester,AlecBaldwin,dsa@silly.com`
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

	// w := bufio.NewWriter(out)
	w := bufio.NewWriter(os.Stdout)
	for _, line := range lines {
		words := strings.Split(line, string(delimiter))
		columnNum = 0
		for _, word := range words {
			for len(word) < columnCounts[columnNum] {
				word += " "
			}
			columnNum++
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
