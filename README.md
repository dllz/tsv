Advanced TSV parser for Go 
====
Based on github.com/dogenzaka/tsv but has some improvements

- Support for []string 
- Support for float64 
- Support for float32
- Support for nil values in columns 
- More resilient to empty spaces
- Supports go modules

Note it is not a complete drop in replacement as some of the public API needed to be changed in order to support more functionality

tsv advanced is tab-separated values parser for GO. It will parse lines and insert data into any type of struct. tsv supports both simple structs and structs with tagging.

```
go get github.com/dllz/tsv-advanced
```

Quickstart
--

tsv inserts data into struct by fields order.

```go

import (
    "fmt"
    "os"
    "testing"
    )

type TestRow struct {
  Name   string // 0
  Age    int    // 1
  Gender string // 2
  Active bool   // 3
}

func main() {

  file, _ := os.Open("example.tsv")
  defer file.Close()
  arrayDeliminator := ","
  nilSign := "\\N"
  data := TestRow{}
  parser, err := NewParser(file, &data, arrayDeliminator, nilSign)

  for {
    eof, err := parser.Next()
    if eof {
      return
    }
    if err != nil {
      panic(err)
    }
    fmt.Println(data)
  }

}

```

You can define tags to struct fields to map values.

```go
type TestRow struct {
	Name        string   `tsv:"name"`
	Age         int      `tsv:"age"`
	Active      bool     `tsv:"active"`
	Gender      string   `tsv:"gender"`
	MiddleNames []string `tsv:"middleNames"`
	BigNumber   float64  `tsv:"bigNumber"`
	SmallNumber float32  `tsv:"smallNumber"`
}
```

Supported field types
--

Currently, this library supports limited fields but more can easily be added on request or feel free to open PRs with the added functionality

- int
- string
- bool
- float32
- float64
- []string

