package tsv

import (
	"encoding/csv"
	"errors"
	"fmt"
	"golang.org/x/text/unicode/norm"
	"io"
	"reflect"
	"strconv"
	"strings"
)

// Parser has information for parser
type Parser struct {
	Headers          []string
	Reader           *csv.Reader
	Data             interface{}
	ref              reflect.Value
	indices          []int // indices is field index list of header array
	structMode       bool
	normalize        norm.Form
	arrayDeliminator string
	nilSign          string
}

// NewStructModeParser creates new TSV parser with given io.Reader as struct mode
func NewParser(reader io.Reader, data interface{}, arrayDeliminator string, nilSign string) (*Parser, error) {
	r := csv.NewReader(reader)
	r.Comma = '\t'

	// first line should be fields
	headers, err := r.Read()

	if err != nil {
		return nil, err
	}

	for i, header := range headers {
		headers[i] = header
	}

	p := &Parser{
		Reader:           r,
		Headers:          headers,
		Data:             data,
		ref:              reflect.ValueOf(data).Elem(),
		indices:          make([]int, len(headers)),
		structMode:       false,
		normalize:        -1,
		arrayDeliminator: arrayDeliminator,
		nilSign:          nilSign,
	}

	// get type information
	t := p.ref.Type()

	for i := 0; i < t.NumField(); i++ {
		// get TSV tag
		tsvtag := t.Field(i).Tag.Get("tsv")
		if tsvtag != "" {
			// find tsv position by header
			for j := 0; j < len(headers); j++ {
				if headers[j] == tsvtag {
					// indices are 1 start
					p.indices[j] = i + 1
					p.structMode = true
				}
			}
		}
	}

	if !p.structMode {
		for i := 0; i < len(headers); i++ {
			p.indices[i] = i + 1
		}
	}

	return p, nil
}

// NewParserWithoutHeader creates new TSV parser with given io.Reader
func NewParserWithoutHeader(reader io.Reader, data interface{}, arrayDeliminator string, nilSign string) *Parser {
	r := csv.NewReader(reader)
	r.Comma = '\t'

	p := &Parser{
		Reader:           r,
		Data:             data,
		ref:              reflect.ValueOf(data).Elem(),
		normalize:        -1,
		arrayDeliminator: arrayDeliminator,
		nilSign:          nilSign,
	}

	return p
}

// Next puts reader forward by a line
func (p *Parser) Next() (eof bool, err error) {

	// Get next record
	var records []string

	for {
		// read until valid record
		records, err = p.Reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				return true, nil
			}
			return false, err
		}
		if len(records) > 0 {
			break
		}
	}

	if len(p.indices) == 0 {
		p.indices = make([]int, len(records))
		// mapping simple index
		for i := 0; i < len(records); i++ {
			p.indices[i] = i + 1
		}
	}

	// record should be a pointer
	for i, record := range records {
		idx := p.indices[i]
		if idx == 0 {
			// skip empty index
			continue
		}
		fmt.Println(record)
		record = strings.TrimSpace(record)
		//Account for other ways of denoting null
		if p.nilSign != "" {
			if record == p.nilSign {
				record = ""
			}
		}
		// get target field
		field := p.ref.Field(idx - 1)
		switch field.Type().String() {
		case "string":
			// Normalize text
			if p.normalize >= 0 {
				record = p.normalize.String(record)
			}
			field.SetString(record)
		case "bool":
			if record == "" {
				field.SetBool(false)
			} else {
				col, err := strconv.ParseBool(record)
				if err != nil {
					return false, err
				}
				field.SetBool(col)
			}
		case "int":
			if record == "" {
				field.SetInt(0)
			} else {
				col, err := strconv.ParseInt(record, 10, 0)
				if err != nil {
					return false, err
				}
				field.SetInt(col)
			}
		case "[]string":
			if p.arrayDeliminator != "" {
				subRecords := strings.Split(record, p.arrayDeliminator)
				var cleanedUpRecords []string
				for _, subRecord := range subRecords {
					subRecord = strings.TrimSpace(subRecord)
					if subRecord == "" {
						continue
					}
					cleanedUpRecords = append(cleanedUpRecords, subRecord)
				}
				slice := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(record)), len(cleanedUpRecords), cap(cleanedUpRecords))
				field.Set(slice)
				for index, cleanedUpRecord := range cleanedUpRecords {
					if p.normalize >= 0 {
						cleanedUpRecords[index] = p.normalize.String(cleanedUpRecord)
					}
					slice.Index(index).Set(reflect.ValueOf(cleanedUpRecords[index]))
				}
			}
		case "float64":
			if record == "" {
				field.SetFloat(0)
			} else {
				col, err := strconv.ParseFloat(record, 64)
				if err != nil {
					return false, err
				}
				field.SetFloat(col)
			}
		case "float32":
			if record == "" {
				field.SetFloat(0)
			} else {
				col, err := strconv.ParseFloat(record, 32)
				if err != nil {
					return false, err
				}
				field.SetFloat(col)
			}
		default:
			return false, errors.New("Unsupported field type:" + field.Type().String())
		}
	}
	return false, nil
}
