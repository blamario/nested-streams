package pipe

import "testing"

func TestPipes(t *testing.T) {
	runPipe(
		func(sink chan<- int) {
			sink <- 1
			sink <- 2
		},
		func(source <-chan int) {
			if <-source != 1 {
				t.Fail()
			}
			if <-source != 2 {
				t.Fail()
			}
		})
	var buffer []int
	testItems := []int{4, 6, 8}
	runPipe(Itemize(testItems), Collect[int](&buffer))
	if len(buffer) != len(testItems) {
		t.Fail()
	}
	for i, n := range buffer {
		if n != testItems[i] {
			t.Fail()
		}
	}
}

func FuzzPipes(f *testing.F) {
	f.Fuzz(func(t *testing.T, testItems []byte) {
		var buffer []byte
		runPipe(Itemize(testItems), Collect[byte](&buffer))
		if len(buffer) != len(testItems) {
			t.Fail()
		}
		for i, n := range buffer {
			if n != testItems[i] {
				t.Fail()
			}
		}
	})
}
