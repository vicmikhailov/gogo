# gogo: Go Fast Track for Senior Java Developers

`gogo` is a hands-on Go 1.24 training project for experienced Java developers. It focuses on the differences that matter when moving from class-oriented, garbage-collected Java to Go's value semantics, interfaces, explicit errors, built-in collections, and lightweight concurrency. Each package has runnable examples and tests; Java comparisons are guides, not claims that the languages behave identically.

## Learning path

Run the complete showcase with `make run`. It executes the demos in `cmd/gogo/main.go` and then starts the HTTP server; press Enter or Ctrl+C to stop it.

| Step | Package | Focus |
|---|---|---|
| 1 | `syntax` | Values, pointers, zero values, multiple returns, `defer` |
| 2 | `concurrency` | Goroutines, channels, cancellation, worker pools, shared mutable state |
| 3 | `interfaces` | Implicit interface satisfaction and consumer-sized contracts |
| 4 | `generics` | Type parameters, constraints, reusable helpers |
| 5 | `errors` | Error values, wrapping, `errors.Is` and `errors.As` |
| 6 | `collections` | Slices, sets, queues, ordered maps, grouping and stream-like functions |
| 7 | `basictypes` | Maps, strings, JSON, slice/map aliasing, value and pointer mutation |
| 8 | `patterns` | GoF patterns, functional options, iterators, and Java-to-Go antipatterns |
| 9 | `advanced` | Embedding, reflection, closures, atomics, and `embed` |
| 10 | `iosystem` | Files, environment, processes, networking, signals |
| 11 | `commonlibs` | Examples using zap, uuid, testify, and gin |
| 12 | `web` | HTTP handlers, middleware, JSON requests |

For a focused session, call an individual package's `Run*Demo` function or run its tests (for example, `go test ./pkg/concurrency`).

## Getting started

Install the Go version declared in `go.mod` (Go 1.24 or newer), then:

```sh
make run       # build and run all demos; the web server remains active
go test ./...  # run all package tests
go test -race ./... # detect data races in tests
go vet ./...   # run Go's static checks
```

The Makefile also provides `build`, `test`, `fmt`, `vet`, and `clean`. `make run` builds `gogo_binary` before launching it.

## Java-to-Go mental model

| Java starting point | Go perspective |
|---|---|
| Classes and object references | Structs are values. Assignment copies a struct; pointers are explicit aliases. |
| `List<T>` / `ArrayList<T>` | Slices are descriptors over arrays. Element writes can be shared; `append` returns the descriptor to use next. |
| `HashMap<K,V>` / `HashSet<T>` | Built-in maps and `map[T]struct{}` sets. Map iteration order is unspecified; maps need synchronization for concurrent mutation. |
| `implements` interfaces | Interfaces are satisfied implicitly. Define the smallest interface where the consumer needs it. |
| Exceptions | Return `error` values and handle them explicitly; wrap with `%w` when adding context. |
| Threads and executors | Goroutines are runtime-scheduled tasks. Use channels for communication, `sync` primitives for shared state, and `context` for cancellation. |
| `null` | Only certain Go types can be `nil`; zero values and typed-nil interface values have distinct behavior. |

Go slices, maps, and channels are passed by value, but copying their descriptors does not deep-copy the data they refer to. A slice append may allocate a new backing array, so a function that appends should return the resulting slice. Maps share entries across copied map values. These semantics do not make concurrent access safe.

## Patterns worth carrying over—and those to reconsider

Idiomatic patterns shown in the code include narrow consumer-side interfaces, explicit error wrapping, functional options for extensible configuration, composition through embedding, `sync.Once` for one-time initialization, and context-aware pipelines. A GoF pattern is still a tool, not a goal: prefer the simplest code that makes the contract clear.

The antipattern demos cover typed nil errors, copying synchronization primitives, goroutine leaks, producer-side interface overuse, unnecessary pointers to slices/maps, panic for ordinary errors, retaining large arrays through small slices, avoidable slice growth, unsynchronized map access, shadowed errors, and assuming a struct assignment deep-copies slices and maps. Each example is paired with a safer alternative where applicable. Some antipattern functions intentionally demonstrate incorrect behavior; do not copy them into application code.

For collections, use built-in slices and maps by default. The generic `Set`, `Stack`, `Queue`, and `OrderedMap` are teaching examples; an insertion-ordered map is only worth maintaining when order is part of the required behavior. The queue clears consumed entries and compacts its storage, and the set and ordered map support useful zero values.

## Build, dependencies, and deployment

Go Modules (`go.mod`) declare dependencies, roughly analogous to Maven's `pom.xml`. The `go` command provides builds, tests, formatting, vetting, benchmarks, and fuzzing without requiring a separate build framework. This repository also uses third-party libraries in `pkg/commonlibs`.

`go build` creates a native executable. Whether it is fully self-contained depends on build settings and platform dependencies (for example, cgo may link system libraries); test the artifact on the intended target. Cross-compilation is available for supported targets, for example:

```sh
GOOS=linux GOARCH=amd64 go build -o gogo-linux ./cmd/gogo
```

## Testing, race detection, and diagnostics

Go's testing tools can check behavior, concurrency, allocations, and many failure modes without a separate test framework. Use deterministic synchronization in concurrency tests (channels, `sync.WaitGroup`, and contexts) instead of relying on sleeps. Put deadlines around tests that could block, and make each goroutine signal completion so the test can verify shutdown. Test boundary cases too: empty and nil inputs, cancellation before and during work, errors, repeated calls, and concurrent access.

### Tests, coverage, and repeatability

```sh
go test ./...                                      # all package tests
go test ./pkg/collections -run 'TestQueue'         # focused tests
go test ./... -shuffle=on -count=10                # varied order and repeated runs
go test ./... -cover                               # per-package statement coverage
go test ./... -coverprofile=coverage.out           # write a coverage profile
go tool cover -func=coverage.out                   # inspect coverage by function
go test -run '^$' -bench=. -benchmem ./pkg/...     # benchmarks and allocation counts
```

`-count` helps expose flaky timing or state assumptions; it does not prove a test is race-free. Coverage shows which statements ran, not whether assertions were meaningful. Add table-driven tests for input variations, and test the externally visible behavior rather than implementation details where possible.

### Detecting data races and concurrency bugs

Run race-enabled tests regularly, especially after changing goroutines, shared state, maps, slices, or callbacks:

```sh
go test -race ./...                                # instrument test binaries
go test -race -count=20 ./pkg/concurrency          # repeat concurrency-focused tests
go run -race ./cmd/gogo                            # instrument the full showcase
GOMAXPROCS=1 go test -count=20 ./pkg/concurrency   # exercise another scheduling configuration
```

The race detector reports conflicting unsynchronized memory accesses that occur in the code paths actually executed. It is not a proof of race freedom, and it does not detect every concurrency problem: deadlocks, incorrect locking discipline, lost updates implemented with individually race-free operations, and goroutine leaks need explicit tests and review. Keep `-race` in CI where supported; instrumentation increases runtime and memory use.

For concurrent code, verify both results and lifecycle: all workers finish, cancellation unblocks senders and receivers, channels are closed by the sender, and shared mutable state has a clear lock/ownership rule. Use timeouts only as a test safety net, not as the synchronization mechanism. Run high-contention scenarios repeatedly, and consider varying `GOMAXPROCS`.

### Finding goroutine leaks and hangs

Prefer tests that wait for a worker's `done` signal after cancellation and fail if it does not exit by a deadline. Avoid asserting an exact `runtime.NumGoroutine` count; the runtime and test runner have their own goroutines, making exact counts flaky. If a test suite needs systematic leak checks, a dedicated helper such as `go.uber.org/goleak` can verify that goroutines created by a test have exited.

For a running process, capture a goroutine profile to see where goroutines are blocked:

```sh
go tool pprof http://localhost:6060/debug/pprof/goroutine
```

This URL is available only when the application explicitly exposes `net/http/pprof` (typically on a protected diagnostic listener); the showcase does not expose it by default. `go test` also supports block and mutex profiles for locating synchronization contention:

```sh
go test ./pkg/concurrency -blockprofile=block.out -mutexprofile=mutex.out
go tool pprof block.out
go tool pprof mutex.out
```

### Investigating memory growth and performance

First reproduce the growth with a representative workload. A high allocation rate is not automatically a leak: compare live heap (`inuse_space`) with total allocation (`alloc_space`) over time. Look for retained references (for example, a small subslice holding a large backing array), unbounded caches, queues retaining consumed elements, and goroutines blocked while holding data.

```sh
go test -run '^$' -bench=. -benchmem -memprofile=mem.out ./pkg/...
go tool pprof -alloc_space mem.out                   # cumulative allocation hotspots
go tool pprof -inuse_space mem.out                    # retained heap in the profile
GODEBUG=gctrace=1 go test ./pkg/collections          # inspect garbage-collection activity
```

Use benchmarks and profiles to guide optimization rather than assuming preallocation or a custom collection is faster. For a long-running service, compare heap profiles after similar workload periods and allow for garbage collection before concluding memory is retained. Profiles can contain sensitive data; store and share them accordingly.

### Fuzzing, static checks, and build checks

Fuzz parsers and other functions with broad input spaces; Go saves failing inputs as regression corpus entries:

```sh
go test -fuzz=FuzzReverse -fuzztime=30s ./pkg/advanced
```

Also run the built-in analyzer and check every supported target in CI:

```sh
go fmt ./...                                    # format Go code
go vet ./...                                    # suspicious constructs and common mistakes
go test ./...                                   # correctness
go test -race ./...                             # exercised data races
go build ./cmd/gogo                             # compile the executable
```

Optional tools such as `staticcheck` and `govulncheck` can find additional code and dependency issues; install and run them in CI using their official Go tool instructions. For production incidents, add structured logs and metrics with request or operation identifiers, but avoid logging secrets or unbounded/high-cardinality values.

The examples are print-driven for easy exploration; tests assert behavior and guard against regressions. Start with one package, compare its Java analogy with the actual Go implementation, then run the tests and experiment with the examples.
