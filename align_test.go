package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

const comma = ","

var columnCountCases = []struct {
	input  string
	sep    string
	isQual bool
	qual   string
	counts map[int]int
}{
	{
		"John,Doe,Henry\nMichael,Douglas,F",
		comma,
		false,
		"",
		map[int]int{
			0: 7,
			1: 7,
			2: 5,
		},
	},
	{
		"John,Doe,\"Henry, Mellencamp\",jm@nothing.com\nMichael,Douglas,F",
		comma,
		true,
		"\"",
		map[int]int{
			0: 7,
			1: 7,
			2: 19,
		},
	}, {
		"John,Doe,\"Henry, Mellencamp\",jm@nothing.com\nMichael,Douglas,F",
		comma,
		true,
		"\"",
		map[int]int{
			0: 7,
			1: 7,
			2: 19,
		},
	},
	{
		"one||two||three||four\nuno||dos||tres||quatro",
		"||",
		true,
		"",
		map[int]int{
			0: 3,
			1: 3,
			2: 5,
			3: 6,
		},
	},
	{
		"one||\"two|| number\"||three||four\nuno||dos||\"tres|\"||quatro",
		"||",
		true,
		"\"",
		map[int]int{
			0: 3,
			1: 14,
			2: 7,
			3: 6,
		},
	},
	{
		"one,tisß\nseven,two", // with byte count > 1
		comma,
		false,
		"",
		map[int]int{
			0: 5,
			1: 5,
		},
	},
}

var fieldLenCases = []struct {
	input    string
	sep      string
	expected int
}{
	{
		"first,last",
		",",
		5,
	},
	{
		"phone-number||email",
		"||",
		12,
	},
}

var fieldLenEscapedCases = []struct {
	input    string
	sep      string
	qual     string
	expected int
}{
	{
		"\"address1, address2\",last",
		",",
		"\"",
		20,
	},
	{
		"'expenseline2|expenseline2'||email",
		"||",
		"'",
		27,
	},
}

var countPaddingCases = []struct {
	input    string
	fieldLen int
	expected int
}{
	{
		"Roy",
		10,
		7,
	},
	{
		"S",
		20,
		19,
	},
	{
		"Luü",
		5,
		2,
	},
}

var paddingCases = []struct {
	input       string
	columnCount int
	po          PaddingOpts
	expected    int
}{
	{
		"Go",
		8,
		PaddingOpts{Justification: JustifyLeft},
		10,
	},
	{
		"Go",
		8,
		PaddingOpts{Justification: JustifyCenter},
		10,
	},
	{
		"Go",
		4,
		PaddingOpts{Justification: JustifyCenter},
		6,
	},
	{
		"Go",
		8,
		PaddingOpts{Justification: JustifyRight},
		10,
	},
}

var qualifiedSplitCases = []struct {
	input    string
	sep      string
	qual     string
	expected int
}{
	{
		"First,\"Middle, Nickname\",Last",
		",",
		"\"",
		3,
	},
	{
		"First||\"Middle|| Nickname\"||Last",
		"||",
		"\"",
		3,
	},
	{
		"First,'Middle Nickname',Last",
		",",
		"'",
		3,
	},
	{
		"First,Middle Nickname,Last",
		",",
		"",
		3,
	},
	{
		"First",
		",",
		"\"",
		1,
	},
}

var exportCases = []struct {
	input          io.Reader
	output         io.Writer
	outColumns     []int
	numOfDelimiter int
}{
	{
		strings.NewReader("one,two,three\nfour,five,six"),
		bytes.NewBufferString(""),
		[]int{1},
		1,
	},
	{
		strings.NewReader("first,middle,last"),
		bytes.NewBufferString(""),
		[]int{1, 3},
		0,
	},
	{
		strings.NewReader("first,middle,last"),
		bytes.NewBufferString(""),
		[]int{1, 4},
		0,
	},
	{
		strings.NewReader("first,middle,last"),
		bytes.NewBufferString(""),
		nil,
		0,
	},
}

var runCases = []struct {
	hValue    bool
	helpValue bool
	fValue    string
	oValue    string
	qValue    string
	sValue    string
	dValue    string
	aValue    string
	cValue    string
	iValue    string
	shouldErr bool
}{
	{
		hValue:    true,
		shouldErr: true,
	},
	{
		helpValue: true,
		shouldErr: true,
	},
	{
		iValue:    "1:right,2:center",
		shouldErr: true,
	},
	{
		iValue:    "1:left,2:,NaN:left",
		shouldErr: true,
	},
	{
		iValue:    "3:noexist",
		shouldErr: true,
	},
	{
		cValue:    "a2-right,NaN",
		shouldErr: true,
	},
	{
		cValue:    "1,2",
		shouldErr: true,
	},
	{
		cValue:    "1-left,2-right,3-center",
		shouldErr: true,
	},
	{
		cValue:    "1,2,a",
		shouldErr: true,
	},
	{
		qValue:    "\"",
		shouldErr: true,
	},
	{
		oValue:    "out.csv",
		shouldErr: true,
	},
}

var ouputSepCases = []struct {
	input    string
	expected string
}{
	{
		",",
		",",
	},
	{
		"||",
		"||",
	},
}

var columnSizeCases = []struct {
	counts   map[int]int
	cNum     int
	expected int
}{
	{
		map[int]int{
			0: 2,
			1: 2,
			2: 5,
		},
		0,
		2,
	},
	{
		map[int]int{
			0: 2,
			1: 2,
			2: 5,
		},
		3,
		-1,
	},
}

var updatePaddingCases = []struct {
	input    justification
	expected justification
}{
	{
		JustifyCenter,
		JustifyCenter,
	},
	{
		JustifyRight,
		JustifyRight,
	},
}

func TestRun(t *testing.T) {
	// flag.Usage()

	for _, tt := range runCases {
		*dFlag = tt.dValue
		*sFlag = tt.sValue
		*hFlag = tt.hValue
		*helpFlag = tt.helpValue
		*fFlag = tt.fValue
		*oFlag = tt.oValue
		*aFlag = tt.aValue
		*cFlag = tt.cValue
		*qFlag = tt.qValue
		*iFlag = tt.iValue

		code, _ := run()

		if tt.shouldErr != (code == 1) {
			t.Fatalf("run() = %v; want %v", code, tt.shouldErr)
		}
	}
}

func TestUpdatePadding(t *testing.T) {
	for _, tt := range updatePaddingCases {
		a := &Align{}
		a.updatePadding(PaddingOpts{Justification: tt.input})

		if a.padOpts.Justification != tt.expected {
			t.Fatalf("updatePadding(%v) = %v; want %v", tt.input, a.padOpts.Justification, tt.expected)
		}
	}
}

func TestColumnSize(t *testing.T) {
	for _, tt := range columnSizeCases {
		a := &Align{columnCounts: tt.counts}

		got := a.columnSize(tt.cNum)
		if got != tt.expected {
			t.Fatalf("columnSize(%v) = %v; want %v", tt.cNum, got, tt.expected)
		}
	}
}

func TestOutputSep(t *testing.T) {
	for _, tt := range ouputSepCases {
		a := &Align{}
		a.outputSep(tt.input)

		if a.sepOut != tt.expected {
			t.Fatalf("outputSep(%v) = %v; want %v", tt.input, a.sepOut, tt.expected)
		}
	}
}

func TestColumnFilter(t *testing.T) {
	for _, tt := range exportCases {
		a := newAlign(tt.input, tt.output, comma, TextQualifier{})
		a.filterColumns(tt.outColumns)

		a.Align()
	}
}

func TestSplit(t *testing.T) {
	for _, tt := range qualifiedSplitCases {
		a := newAlign(strings.NewReader(tt.input), os.Stdout, comma, TextQualifier{On: true, Qualifier: "\""})
		got := a.splitWithQual(tt.input, tt.sep, tt.qual)

		if len(got) != tt.expected {
			t.Fatalf("splitWithQual(%v, %v, %v) = %v; want %v", tt.input, tt.sep, tt.qual, len(got), tt.expected)
		}
	}
}

func TestPad(t *testing.T) {
	for _, tt := range paddingCases {
		got := pad(tt.input, 1, tt.columnCount, tt.po.Justification)

		if len(got) != tt.expected {
			t.Fatalf("pad(%v) =%v; want %v", tt.input, got, tt.expected)
		}
	}
}

func TestColumnCounts(t *testing.T) {
	for _, tt := range columnCountCases {
		a := newAlign(strings.NewReader(tt.input), os.Stdout, tt.sep, TextQualifier{On: tt.isQual, Qualifier: tt.qual})
		a.columnLength()
		for i := range tt.counts {
			if a.columnSize(i) != tt.counts[i] {
				t.Fatalf("Count for column %v = %v, want %v", i, a.columnSize(i), tt.counts[i])
			}
		}
	}
}

func TestFieldLenEscaped(t *testing.T) {
	for _, tt := range fieldLenEscapedCases {
		got := fieldLenEscaped(tt.input, tt.sep, tt.qual)
		if got != tt.expected {
			t.Fatalf("FieldLenEscaped(%v) = %v; want %v", tt.input, got, tt.expected)
		}
	}
}

func TestFieldLen(t *testing.T) {
	for _, tt := range fieldLenCases {
		got := fieldLen(tt.input, tt.sep)
		if got != tt.expected {
			t.Fatalf("FieldLen(%v) = %v; want %v", tt.input, got, tt.expected)
		}
	}
}

func TestCountPadding(t *testing.T) {
	for _, tt := range countPaddingCases {
		got := countPadding(tt.input, tt.fieldLen)
		if got != tt.expected {
			t.Fatalf("countPadding(%v) = %v; want %v", tt.input, got, tt.expected)
		}
	}
}

func BenchmarkColumnCounts(b *testing.B) {
	input := `First,Middle,Last,Email,Region,City,Zip,Full_Name,First,Middle,Last,Email,Region,City,Zip,Full_Name,First,Middle,Last,Email,Region,City,Zip,Full_Name
Karleigh,Destiny,Dean,nunc.In@lorem.edu,Stockholms län,Märsta,9038,Shaine Reilly,Karleigh,Destiny,Dean,nunc.In@lorem.edu,Stockholms län,Märsta,9038,Shaine Reilly,Karleigh,Destiny,Dean,nunc.In@lorem.edu,Stockholms län,Märsta,9038,Shaine Reilly        
Alisa,Walker,Armand,Sed@Nuncmauriselit.com,Himachal Pradesh,Shimla,MZ0 4QS,Olivia Velez ,Alisa,Walker,Armand,Sed@Nuncmauriselit.com,Himachal Pradesh,Shimla,MZ0 4QS,Olivia Velez ,Alisa,Walker,Armand,Sed@Nuncmauriselit.com,Himachal Pradesh,Shimla,MZ0 4QS,Olivia Velez
Karleigh,Destiny,Dean,nunc.In@lorem.edu,Stockholms län,Märsta,9038,Shaine Reilly,Karleigh,Destiny,Dean,nunc.In@lorem.edu,Stockholms län,Märsta,9038,Shaine Reilly,Karleigh,Destiny,Dean,nunc.In@lorem.edu,Stockholms län,Märsta,9038,Shaine Reilly        
Alisa,Walker,Armand,Sed@Nuncmauriselit.com,Himachal Pradesh,Shimla,MZ0 4QS,Olivia Velez ,Alisa,Walker,Armand,Sed@Nuncmauriselit.com,Himachal Pradesh,Shimla,MZ0 4QS,Olivia Velez ,Alisa,Walker,Armand,Sed@Nuncmauriselit.com,Himachal Pradesh,Shimla,MZ0 4QS,Olivia Velez
Karleigh,Destiny,Dean,nunc.In@lorem.edu,Stockholms län,Märsta,9038,Shaine Reilly,Karleigh,Destiny,Dean,nunc.In@lorem.edu,Stockholms län,Märsta,9038,Shaine Reilly,Karleigh,Destiny,Dean,nunc.In@lorem.edu,Stockholms län,Märsta,9038,Shaine Reilly
Alisa,Walker,Armand,Sed@Nuncmauriselit.com,Himachal Pradesh,Shimla,MZ0 4QS,Olivia Velez ,Alisa,Walker,Armand,Sed@Nuncmauriselit.com,Himachal Pradesh,Shimla,MZ0 4QS,Olivia Velez ,Alisa,Walker,Armand,Sed@Nuncmauriselit.com,Himachal Pradesh,Shimla,MZ0 4QS,Olivia Velez
`

	a := newAlign(strings.NewReader(input), os.Stdout, comma, TextQualifier{On: false})

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		a.columnLength()
	}
}

func BenchmarkSplitWithQual(b *testing.B) {
	input := "First,\"Middle, name\",Last,Email,Region,City,Zip,Full_Name"

	a := newAlign(strings.NewReader(input), os.Stdout, comma, TextQualifier{On: true, Qualifier: "\""})

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		a.splitWithQual(input, comma, "\"")
	}
}

func BenchmarkSplitWithQualNoQual(b *testing.B) {
	input := "First,Middle,Last,Email,Region,City,Zip,Full_Name"

	a := newAlign(strings.NewReader(input), os.Stdout, comma, TextQualifier{On: false})

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		a.splitWithQual(input, comma, "\"")
	}
}
