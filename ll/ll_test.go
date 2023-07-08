package ll_test

import (
	"fmt"
	"testing"

	"github.com/jaqmol/goito/ll"
)

func TestParallel(t *testing.T) {
	start := ll.Start(func(size int, sink ll.Sink[int]) {
		for i := 0; i < size; i++ {
			sink.Write(i)
		}
		sink.Done()
	})
	cast := ll.Pipe(func(num int, sink ll.Sink[string]) {
		value := fmt.Sprintf("NUM:%d", num)
		sink.Write(value)
	})
	coll := make([]string, 0)
	end := ll.End(func(val string, term ll.Term) {
		coll = append(coll, val)
	})
	start.Next(cast)
	cast.Next(end)
	start.Start()
	start.Write(100)
	err := end.Wait()
	if err != nil {
		t.Error(err)
	} else {
		for i, val := range coll {
			expt := fmt.Sprintf("NUM:%d", i)
			if val != expt {
				t.Errorf("Expected %s, got %s", expt, val)
			}
		}
	}
}
