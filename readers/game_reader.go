package readers

import (
	"bufio"
	"encoding/csv"
	"io"
	"log"
	"strings"
)

type GameReader struct {
	scanner *bufio.Scanner
}

func NewGameReader(r io.Reader) GameReader {
	gr := GameReader{
		scanner: bufio.NewScanner(r),
	}
	return gr
}

func (gr *GameReader) Read() ([]string, error) {
	ok := gr.scanner.Scan()
	if !ok {
		return nil, io.EOF
	}
	t := gr.scanner.Text()
	r := csv.NewReader(strings.NewReader(t))
	r.LazyQuotes = true
	fields, err := r.Read()
	if err != nil {
		log.Println("Line read failed: ", t)
	}
	return fields, err
}
