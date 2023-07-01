package io

import (
	"io/fs"
	"io/ioutil"

	"github.com/jaqmol/goito/ito"
)

func ListDir(dirPath string) (iter ito.Ito[fs.FileInfo], err error) {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}
	return &lister{files: files, index: -1}, nil
}

type lister struct {
	files []fs.FileInfo
	index int
}

func (l *lister) Next() bool {
	l.index++
	return l.index < len(l.files)
}

func (l *lister) Item() fs.FileInfo {
	return l.files[l.index]
}
