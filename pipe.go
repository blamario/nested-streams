package pipe

import "fmt"
import "sync"

type Producer[item any] func(chan<- item)
type Consumer[item any] func(<-chan item)
type Transducer[in any, out any] func(<-chan in, chan<- out)

func runPipe[of any](p Producer[of], c Consumer[of]) {
	var wg sync.WaitGroup
	wg.Add(1)
	pipe := make(chan of)
	go func() { defer func() { wg.Done(); recover() }(); p(pipe) }()
	go func() { defer func() { wg.Done(); recover() }(); c(pipe) }()
	wg.Wait() // until either is done
	wg.Add(1)
	close(pipe)
	wg.Wait() // until the other is done
}

func pipeTransducers[in, mid, out any](up Transducer[in, mid], down Transducer[mid, out]) Transducer[in, out] {
	return func(source <-chan in, sink chan<- out) {
		runPipe[mid](func(m chan<- mid) { up(source, m) },
			func(m <-chan mid) { down(m, sink) })
	}
}

func pipeTransducerConsumer[in, mid any](t Transducer[in, mid], c Consumer[mid]) Consumer[in] {
	return func(source <-chan in) {
		runPipe[mid](func(m chan<- mid) { t(source, m) },
			func(m <-chan mid) { c(m) })
	}
}

func pipeTransducerProducer[mid, out any](p Producer[mid], t Transducer[mid, out]) Producer[out] {
	return func(sink chan<- out) {
		runPipe[mid](func(m chan<- mid) { p(m) },
			func(m <-chan mid) { t(m, sink) })
	}
}

func Print[of any](source <-chan of) {
	for item := range source {
		fmt.Println(item)
	}
}

func Itemize[of any](from []of) func(chan<- of) {
	return func(sink chan<- of) {
		for _, item := range from {
			sink <- item
		}
	}
}

func Collect[of any](buffer *[]of) func(<-chan of) {
	return func(source <-chan of) {
		for item := range source {
			*buffer = append(*buffer, item)
		}
	}
}
