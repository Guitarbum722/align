package align

import (
	"os"
	"strings"
	"testing"
)

func BenchmarkColumnCounts(b *testing.B) {
	input := `First,Middle,Last,Email,Region,City,Zip,Full_Name,First,Middle,Last,Email,Region,City,Zip,Full_Name,First,Middle,Last,Email,Region,City,Zip,Full_Name
Karleigh,Destiny,Dean,nunc.In@lorem.edu,Stockholms län,Märsta,9038,Shaine Reilly,Karleigh,Destiny,Dean,nunc.In@lorem.edu,Stockholms län,Märsta,9038,Shaine Reilly,Karleigh,Destiny,Dean,nunc.In@lorem.edu,Stockholms län,Märsta,9038,Shaine Reilly        
Alisa,Walker,Armand,Sed@Nuncmauriselit.com,Himachal Pradesh,Shimla,MZ0 4QS,Olivia Velez ,Alisa,Walker,Armand,Sed@Nuncmauriselit.com,Himachal Pradesh,Shimla,MZ0 4QS,Olivia Velez ,Alisa,Walker,Armand,Sed@Nuncmauriselit.com,Himachal Pradesh,Shimla,MZ0 4QS,Olivia Velez         
Karleigh,Destiny,Dean,nunc.In@lorem.edu,Stockholms län,Märsta,9038,Shaine Reilly,Karleigh,Destiny,Dean,nunc.In@lorem.edu,Stockholms län,Märsta,9038,Shaine Reilly,Karleigh,Destiny,Dean,nunc.In@lorem.edu,Stockholms län,Märsta,9038,Shaine Reilly        
Alisa,Walker,Armand,Sed@Nuncmauriselit.com,Himachal Pradesh,Shimla,MZ0 4QS,Olivia Velez ,Alisa,Walker,Armand,Sed@Nuncmauriselit.com,Himachal Pradesh,Shimla,MZ0 4QS,Olivia Velez ,Alisa,Walker,Armand,Sed@Nuncmauriselit.com,Himachal Pradesh,Shimla,MZ0 4QS,Olivia Velez         
Karleigh,Destiny,Dean,nunc.In@lorem.edu,Stockholms län,Märsta,9038,Shaine Reilly,Karleigh,Destiny,Dean,nunc.In@lorem.edu,Stockholms län,Märsta,9038,Shaine Reilly,Karleigh,Destiny,Dean,nunc.In@lorem.edu,Stockholms län,Märsta,9038,Shaine Reilly        
Alisa,Walker,Armand,Sed@Nuncmauriselit.com,Himachal Pradesh,Shimla,MZ0 4QS,Olivia Velez ,Alisa,Walker,Armand,Sed@Nuncmauriselit.com,Himachal Pradesh,Shimla,MZ0 4QS,Olivia Velez ,Alisa,Walker,Armand,Sed@Nuncmauriselit.com,Himachal Pradesh,Shimla,MZ0 4QS,Olivia Velez         
`

	sw := NewAligner(strings.NewReader(input), os.Stdout, ',', TextQualifier{On: false})

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = sw.ColumnCounts()
	}
}

func BenchmarkSplitWithQual(b *testing.B) {
	input := "First,\"Middle, name\",Last,Email,Region,City,Zip,Full_Name"

	sw := NewAligner(strings.NewReader(input), os.Stdout, ',', TextQualifier{On: true, Qualifier: '"'})

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = sw.SplitWithQual(input, ",", "\"")
	}
}

func BenchmarkSplitWithQualNoQual(b *testing.B) {
	input := "First,Middle,Last,Email,Region,City,Zip,Full_Name"

	sw := NewAligner(strings.NewReader(input), os.Stdout, ',', TextQualifier{On: false})

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = sw.SplitWithQual(input, ",", "\"")
	}
}
