package main

import (
	"bufio"
	"bytes"
	"os"
	"regexp"
)

var (
	patEols  = regexp.MustCompile(`[\r\n]+`)
	pat2Eols = regexp.MustCompile(`[\r\n]{2}`)
)

func ScanTwoConsecutiveNewlines(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	if loc := pat2Eols.FindIndex(data); loc != nil && loc[0] >= 0 {
		// Replace newlines within string with a space
		s := patEols.ReplaceAll(data[0:loc[0]+1], []byte(" "))
		// Trim spaces and newlines from string
		s = bytes.Trim(s, "\n ")
		return loc[1], s, nil
	}

	if atEOF {
		// Replace newlines within string with a space
		s := patEols.ReplaceAll(data, []byte(" "))
		// Trim spaces and newlines from string
		s = bytes.Trim(s, "\n ")
		return len(data), s, nil
	}

	// Request more data
	return 0, nil, nil
}

type WebsterParser struct {
	scanner *bufio.Scanner
}

func NewWebsterParser(file *os.File) *WebsterParser {
	scanner := bufio.NewScanner(file)
	scanner.Split(ScanTwoConsecutiveNewlines)
	return &WebsterParser{scanner: scanner}
}

// Finds the next section. String is nil when there are no more sections left
func (wp *WebsterParser) NextSection() (*string, error) {
	if wp.scanner.Scan() {
		str := wp.scanner.Text()
		err := wp.scanner.Err()
		return &str, err
	}
	return nil, nil
}

// Build the next DictEntry
func (wp *WebsterParser) NextEntry() (DictEntry, error) {
	// TODO
}

// Impl DictionaryParser for WebsterParser
func (wp *WebsterParser) Entries(filePaths []string) <-chan DictEntry {
	ch := make(chan DictEntry)
	go func() {
		defer close(ch)
		for _, filePath := range filePaths {
			file, err := os.Open(filePath)
			if err != nil {
				// Handle error, maybe log it and continue to the next file
				println("Could not open file {}", filePath)
				continue
			}

			wp := NewWebsterParser(file)

			for {
				entry, err := wp.NextSection()
				if err != nil {
					println("Error parsing the next entry: {}", err)
					return
				}
				if entry == nil {
					// no more sections left
					break
				}

				// ch <- entry
			}
		}
	}()
	return ch
}
