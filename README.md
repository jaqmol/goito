# goito

**Go Iterator and parallel Pipeline Library**

## Package io

With this package you can list directories, read lines of files and read runes of files. All with one unified generic iterator API:

```go
type Ito[T any] interface {
	Next() bool
	Item() T
}
```

## Package ito

This package only defines the iterator interface, plus a sorting function for said iterators.

## Package ll

This package implements a pipelining library in which every step is executed in a different goroutine, while actual execution is sequential. See test for example.