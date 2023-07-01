package io

import "github.com/jaqmol/goito/ito"

type FileIto[T any] interface {
	ito.Ito[T]
	Close() error
}
