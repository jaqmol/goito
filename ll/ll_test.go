package ll_test

import (
	"testing"

	"github.com/jaqmol/goito/ll"
)

func TestParallel(t *testing.T) {
	start := ll.Start(func(size int, next ll.Sink[int]) {
		for i := 0; i < size; i++ {
			t.Errorf("Writing %d\n", i)
			next.Write(i)
		}
		next.Done()
	})
	end := ll.End(func(num int, next ll.Term) {
		t.Errorf("NUM: %d\n", num)
	})
	start.Next(end)
	start.Start()
	start.Write(100)
	t.Errorf("INITIALIZING %v, %v\n", start, end)
	err := end.Wait()
	if err != nil {
		t.Error(err)
	}
}

// func TestParallel(t *testing.T) {
// 	start := ll.Start(func(size int, next ll.Sink[int]) {
// 		for i := 0; i < size; i++ {
// 			t.Errorf("Writing %d\n", i)
// 			next.Write(i)
// 		}
// 		next.Done()
// 	})
// 	trans := ll.Pipe(func(num int, next ll.Sink[string]) {
// 		value := fmt.Sprintf("NUM: %d", num)
// 		next.Write(value)
// 	})
// 	count := 0
// 	end := ll.End(func(str string, next ll.Term) {
// 		expt := fmt.Sprintf("NUM: %d", count)
// 		if str == expt {
// 			t.Errorf("Expected %s, got %s", expt, str)
// 		}
// 		count++
// 	})
// 	start.Next(trans)
// 	trans.Next(end)
// 	// t.Error("STARTING")
// 	start.Start()
// 	start.Write(100)
// 	t.Errorf("INITIALIZING %v, %v, %v\n", start, trans, end)
// 	err := end.Wait()
// 	if err != nil {
// 		t.Error(err)
// 	}
// }
