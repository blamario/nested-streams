package pipe

import "testing"

func compareSlices[item comparable](t *testing.T, xs, ys []item) {
	if len(xs) != len(ys) {
		t.Fail()
	}
	for i, n := range xs {
		if n != ys[i] {
			t.Fail()
		}
	}
}

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
	compareSlices(t, testItems, buffer)

	buffer = nil
	runPipe(pipeProducerTransducer(Itemize(testItems), Passthrough[int]), Collect[int](&buffer))
	compareSlices(t, testItems, buffer)

	buffer = nil
	runPipe(pipeProducerTransducer(Itemize(testItems), Filter(func(x int) bool { return x != 6 })), Collect[int](&buffer))
	compareSlices(t, []int{4, 8}, buffer)

	buffer = nil
	runPipe(pipeProducerTransducer(Itemize(testItems), Map(func(x int) float32 { return x / 2 })), Collect[int](&buffer))
	compareSlices(t, []int{2, 3, 4}, buffer)
}

func FuzzPipes(f *testing.F) {
	f.Fuzz(func(t *testing.T, testItems []byte) {
		var buffer []byte
		runPipe(Itemize(testItems), Collect[byte](&buffer))
		compareSlices(t, testItems, buffer)
	})
}
