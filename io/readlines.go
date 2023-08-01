package io

import (
	"bufio"
	"os"
)

func ReadLines(filePath string) (iter FileIto[string], err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	return &lineReader{file: file, scanner: bufio.NewScanner(file)}, nil
}

type lineReader struct {
	file    *os.File
	scanner *bufio.Scanner
	line    string
}

func (l *lineReader) Next() bool {
	next := l.scanner.Scan()
	if next {
		l.line = l.scanner.Text()
	}
	return next
}

func (l *lineReader) Item() string {
	return l.line
}

func (l *lineReader) Close() error {
	return l.file.Close()
}
