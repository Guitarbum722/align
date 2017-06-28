package main

import (
	"log"
	"os"

	"github.com/Guitarbum722/true-up/align"
)

func main() {
	delimiter := align.Comma // default
	var sw align.Alignable

	fi, _ := os.Stdin.Stat()
	if (fi.Mode() & os.ModeCharDevice) == 0 {
		switch len(os.Args) {
		case 2:
			f, err := os.Create(os.Args[1])
			if err != nil {
				log.Fatalln(err)
			}
			defer f.Close()
			sw = align.NewAligner(os.Stdin, f, delimiter)
		default:
			sw = align.NewAligner(os.Stdin, os.Stdout, delimiter)
		}
	} else {
		switch len(os.Args) {
		case 2:
			f, err := os.Open(os.Args[1])
			if err != nil {
				log.Fatalln(err)
			}
			defer f.Close()
			sw = align.NewAligner(f, os.Stdout, delimiter)
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
			sw = align.NewAligner(f, out, delimiter)
		default:
			sw = align.NewAligner(os.Stdin, os.Stdout, delimiter)
		}
	}

	lines := sw.ColumnCounts()
	sw.Export(lines)
}
