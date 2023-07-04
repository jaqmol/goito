package ll_test

import (
	"fmt"
	"testing"

	"github.com/jaqmol/goito/ll"
)

func TestPipe(t *testing.T) {
	gen := ll.Gen(func(in int, out chan<- int, errCh chan<- error) {
		out <- in * 2
		out <- in * 3
	})
	step1 := ll.New(func(in int, out chan<- int, errCh chan<- error) {
		out <- in + 1
	})
	step2 := ll.New(func(in int, out chan<- int, errCh chan<- error) {
		out <- in * 10
	})
	step3 := ll.New(func(in int, out chan<- string, errCh chan<- error) {
		out <- fmt.Sprintf("%d", in)
	})
	count := 0
	end := ll.End(func(in string, errCh chan<- error) {
		if count == 0 && in != "210" {
			errCh <- fmt.Errorf("Invalid end input: %s", in)
		} else if count == 1 && in != "310" {
			errCh <- fmt.Errorf("Invalid end input: %s", in)
		}
		count++
	})
	gen.Next(step1)
	step1.Next(step2)
	step2.Next(step3)
	step3.Next(end)
	err := gen.Start(10)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}
