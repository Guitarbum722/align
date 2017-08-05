package align

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

func TestColumnCounts(t *testing.T) {
	for _, tt := range columnCountCases {
		aligner := NewAligner(strings.NewReader(tt.input), os.Stdout, tt.sep, TextQualifier{On: tt.isQual, Qualifier: tt.qual})
		aligner.ColumnCounts()
		for i := range tt.counts {
			if aligner.ColumnSize(i) != tt.counts[i] {
				t.Fatalf("Count for column %v = %v, want %v", i, aligner.ColumnSize(i), tt.counts[i])
			}
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

	sw := NewAligner(strings.NewReader(input), os.Stdout, comma, TextQualifier{On: false})

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = sw.ColumnCounts()
	}
}

func BenchmarkSplitWithQual(b *testing.B) {
	input := "First,\"Middle, name\",Last,Email,Region,City,Zip,Full_Name"

	sw := NewAligner(strings.NewReader(input), os.Stdout, comma, TextQualifier{On: true, Qualifier: "\""})

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = sw.SplitWithQual(input, comma, "\"")
	}
}

func BenchmarkSplitWithQualNoQual(b *testing.B) {
	input := "First,Middle,Last,Email,Region,City,Zip,Full_Name"

	sw := NewAligner(strings.NewReader(input), os.Stdout, comma, TextQualifier{On: false})

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = sw.SplitWithQual(input, comma, "\"")
	}
}
