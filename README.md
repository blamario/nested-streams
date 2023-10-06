# Nested Streams

This project is a re-design and re-implementation in Go of a dormant [Haskell
package](https://hub.darcs.net/blamario/SCC.wiki) called Streaming Component Combinators or SCC.

The SCC framework is a collection of streaming components and component combinators. The components operate on data streams, and most of them start producing output before consuming their entire input. Multiple components can be grouped together into a compound component using one of the combinators. The resulting framework can be seen as a domain-specific language for stream processing.
