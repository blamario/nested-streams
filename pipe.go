package pipe

import "fmt"
import "sync"

// A producer feeds items into the provided sink channel
type Producer[item any] func(chan<- item)

// A consumer receives items from the provided source channel
type Consumer[item any] func(<-chan item)

// A transducer receives items from one channel (the source) and feeds items into another channel (the sink)
type Transducer[in any, out any] func(<-chan in, chan<- out)

// Primitive operation to fork a producer-consumer pair together with a channel through which they communicate. When
// either the producer or the consumer ends, the channel is closed. The function runPipe returns only when both
// goroutine siblings are done.
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

// Pipe two transducers together into a single transducer
func pipeTransducers[in, mid, out any](up Transducer[in, mid], down Transducer[mid, out]) Transducer[in, out] {
	return func(source <-chan in, sink chan<- out) {
		runPipe[mid](func(m chan<- mid) { up(source, m) },
			func(m <-chan mid) { down(m, sink) })
	}
}

// Pipe a transducer and a consumer together into a consumer
func pipeTransducerConsumer[in, mid any](t Transducer[in, mid], c Consumer[mid]) Consumer[in] {
	return func(source <-chan in) {
		runPipe[mid](func(m chan<- mid) { t(source, m) },
			func(m <-chan mid) { c(m) })
	}
}

// Pipe a producer and a transducer together into a producer
func pipeProducerTransducer[mid, out any](p Producer[mid], t Transducer[mid, out]) Producer[out] {
	return func(sink chan<- out) {
		runPipe[mid](func(m chan<- mid) { p(m) },
			func(m <-chan mid) { t(m, sink) })
	}
}

// Producers

// Produce each item in the given slice.
func Itemize[of any](from []of) func(chan<- of) {
	return func(sink chan<- of) {
		for _, item := range from {
			sink <- item
		}
	}
}

// Consumers

// Print each item consumed from the source.
func Print[of any](source <-chan of) {
	for item := range source {
		fmt.Println(item)
	}
}

// Collect into the given buffer each item consumed from the source.
func Collect[of any](buffer *[]of) func(<-chan of) {
	return func(source <-chan of) {
		for item := range source {
			*buffer = append(*buffer, item)
		}
	}
}

// Suppress all items from the source.
func Suppress[of any](source <-chan of) {
	for _ = range source {
	}
}

// Transducers

// Pass-through transducer with no effect.
func Passthrough[of any](source <-chan of, sink chan<- of) {
	for item := range source {
		sink <- item
	}
}

// Filter each consumed item, re-produce it only if the given predicate returns true.
func Filter[of any](pred func(of) bool) func(<-chan of, chan<- of) {
	return func(source <-chan of, sink chan<- of) {
		for item := range source {
			if pred(item) {
				sink <- item
			}
		}
	}
}

// Map each source item using the given function
func Map[in, out any](f func(in) out) func(<-chan in, chan<- out) {
	return func(source <-chan in, sink chan<- out) {
		for item := range source {
			sink <- f(item)
		}
	}
}
