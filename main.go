package main

import (
	"bufio"
	"fmt"
	_ "os"
	"strings"
)

const (
	delimiter = ','
)

func main() {
	s := `FirstName,MiddleName,LastName,Email
John,K,Moore,jkdoe@nothing.com
Winston,Moby,Kazowski,Winston@nothing.net
Dilbert,Sylvester,AlecBaldwin,dsa@silly.com`
	var columnNum int
	var columnCounts = make(map[int]int)
	scanner := bufio.NewScanner(strings.NewReader(s))

	for scanner.Scan() {
		temp := 0
		columnNum = 0
		line := scanner.Text()
		for i, v := range line {
			temp++
			if v != delimiter && i < len(line)-1 {
				continue
			}
			if temp > columnCounts[columnNum] {
				columnCounts[columnNum] = temp
			}
			columnNum++
			temp = 0
		}
		fmt.Println(columnCounts)
	}
}
