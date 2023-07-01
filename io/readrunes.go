package io

import (
	"bufio"
	"os"
)

func ReadRunes(filePath string) (iter FileIto[rune], err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	return &runeReader{file: file, reader: bufio.NewReader(file)}, nil
}

type runeReader struct {
	file    *os.File
	reader  *bufio.Reader
	current rune
}

func (r *runeReader) Next() bool {
	var err error
	r.current, _, err = r.reader.ReadRune()
	return err == nil
}

func (r *runeReader) Item() rune {
	return r.current
}

func (r *runeReader) Close() error {
	return r.file.Close()
}
