package main

import (
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
		_ = a.columnLength()
	}
}

func BenchmarkSplitWithQual(b *testing.B) {
	input := "First,\"Middle, name\",Last,Email,Region,City,Zip,Full_Name"

	a := newAlign(strings.NewReader(input), os.Stdout, comma, TextQualifier{On: true, Qualifier: "\""})

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = a.splitWithQual(input, comma, "\"")
	}
}

func BenchmarkSplitWithQualNoQual(b *testing.B) {
	input := "First,Middle,Last,Email,Region,City,Zip,Full_Name"

	a := newAlign(strings.NewReader(input), os.Stdout, comma, TextQualifier{On: false})

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = a.splitWithQual(input, comma, "\"")
	}
}
