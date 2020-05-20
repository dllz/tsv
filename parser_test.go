package tsv

import (
	"fmt"
	"golang.org/x/text/unicode/norm"
	"os"
	"testing"
)

type TestRow struct {
	Name        string
	Age         int
	Gender      string
	Active      bool
	MiddleNames []string
	BigNumber   float64
	SmallNumber float32
}

type TestTaggedRow struct {
	Age         int      `tsv:"age"`
	Active      bool     `tsv:"active"`
	Gender      string   `tsv:"gender"`
	Name        string   `tsv:"name"`
	MiddleNames []string `tsv:"middleNames"`
}

func TestParserWithoutHeader(t *testing.T) {

	file, err := os.Open("example_simple.tsv")
	if err != nil {
		t.Error(err)
		return
	}
	defer file.Close()

	data := TestRow{}
	parser := NewParserWithoutHeader(file, &data, ",", "\\N")

	i := 0

	for {
		eof, err := parser.Next()
		if eof {
			return
		}
		if i == 0 {
			if data.Name != "alex" ||
				data.Age != 10 ||
				data.Gender != "male" ||
				data.Active != true ||
				len(data.MiddleNames) != 3 {
				fmt.Println(data)
				t.Error("Record does not match index:0")
				if err != nil {
					t.Error(err)
				}
			}
		}
		if i == 1 {
			if data.Name != "john" ||
				data.Age != 24 ||
				data.Gender != "male" ||
				data.Active != false ||
				len(data.MiddleNames) != 2 {
				fmt.Println(data)
				t.Error("Record does not match index:1")
				if err != nil {
					t.Error(err)
				}
			}
		}
		if i == 2 {
			if data.Name != "sara" ||
				data.Age != 30 ||
				data.Gender != "female" ||
				data.Active != true ||
				len(data.MiddleNames) != 1 {
				t.Error("Record does not match index:2")
				if err != nil {
					t.Error(err)
				}
			}
		}
		if i == 3 {
			if err == nil {
				t.Error("Error should be caused")
				return
			}
		}
		if i == 4 {
			if err == nil {
				t.Error("Error should be caused")
				return
			}
		}
		if i == 5 {
			if data.Name != "mike" ||
				data.Age != 55 ||
				data.Gender != "male" ||
				data.Active != false ||
				len(data.MiddleNames) != 0 {
				fmt.Println(data)
				t.Error("Record does not match index:5")
				if err != nil {
					t.Error(err)
				}
			}
		}
		if i == 6 {
			if data.Name != "" ||
				data.Age != 55 ||
				data.Gender != "male" ||
				data.Active != false ||
				len(data.MiddleNames) != 0 {
				t.Error("Record does not match index:6")
				if err != nil {
					t.Error(err)
				}
			}
		}
		if i == 7 {
			if data.Name != "big" ||
				data.Age != 69 ||
				data.Gender != "female" ||
				data.Active != true ||
				len(data.MiddleNames) != 0 ||
				data.BigNumber != 1231231231231212312321354123312312312 ||
				data.SmallNumber != 12312312312541234123 {
				fmt.Println(data)
				t.Error("Record does not match index:7")
				if err != nil {
					t.Error(err)
				}
			}
		}
		i++
	}

}

func TestParserTaggedStructure(t *testing.T) {

	file, err := os.Open("example.tsv")
	if err != nil {
		t.Error(err)
		return
	}
	defer file.Close()

	data := TestTaggedRow{}
	parser, err := NewParser(file, &data, ",", "\\N")
	if err != nil {
		t.Error(err)
		return
	}

	i := 0

	for {
		eof, err := parser.Next()
		if eof {
			return
		}
		if i == 0 {
			if err != nil {
				t.Error(err)
			}
			if data.Name != "alex" ||
				data.Age != 10 ||
				data.Gender != "male" ||
				data.Active != true ||
				len(data.MiddleNames) != 2 {
				t.Error("Record does not match index:0")
			}
		}
		if i == 1 {
			if err != nil {
				t.Error(err)
			}
			if data.Name != "john" ||
				data.Age != 24 ||
				data.Gender != "male" ||
				data.Active != false ||
				len(data.MiddleNames) != 1 {
				t.Error("Record does not match index:1")
			}
		}
		if i == 2 {
			if err != nil {
				t.Error(err)
			}
			if data.Name != "sara" ||
				data.Age != 30 ||
				data.Gender != "female" ||
				data.Active != true ||
				len(data.MiddleNames) != 0 {
				t.Error("Record does not match index:2")
			}
		}
		i++
	}

}

func TestParserNormalize(t *testing.T) {

	file, err := os.Open("example_norm.tsv")
	if err != nil {
		t.Error(err)
		return
	}
	defer file.Close()

	data := TestRow{}
	parser, err := NewParser(file, &data, ",", "\\N")
	if err != nil {
		t.Error(err)
		return
	}
	// Use NFC as normalization
	parser.normalize = norm.NFKC

	i := 0

	for {
		eof, err := parser.Next()
		if eof {
			return
		}
		if err != nil {
			t.Error(err)
		}
		if i == 0 && data.Name != "アレックス" {
			t.Errorf("name is not normalized %v", data.Name)
		}
		if i == 1 && data.Name != "デボラ" {
			t.Errorf("name is not normalized %v", data.Name)
		}
		if i == 2 && data.Name != "デボラ" {
			t.Errorf("name is not normalized %v", data.Name)
		}
		if i == 3 && data.Name != "(テスト)" {
			t.Errorf("name is not normalized %v", data.Name)
		}
		if i == 4 && data.Name != "/" {
			t.Errorf("name is not normalized %v", data.Name)
		}
		i++
	}

}
