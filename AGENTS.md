# AGENTS.md

## What this repo is
- `gogo` is a single-module Go 1.24 project (`go.mod`) that showcases language features, not a layered production service.
- The executable entrypoint is `cmd/gogo/main.go`; reusable demos live in sub-packages within `pkg/`.

## Big-picture architecture
- `cmd/gogo/main.go` is an ordered orchestrator that calls each demo in sequence, ending with a web server:
  `syntax.RunSyntaxDemo` → `concurrency.RunConcurrencyDemo` → `interfaces.RunInterfacesDemo` → `generics.RunGenericsDemo` → `errors.RunErrorsDemo` → `collections.RunCollectionsDemo` → `basictypes.RunBasicTypesDemo` → `patterns.RunPatternsDemo` → `advanced.RunAdvancedDemo` → `iosystem.RunIOSystemDemo` → `commonlibs.RunCommonLibsDemo` → `web.StartWebServer("8080")`.
- Most demos are print-driven and side-effect-based (stdout), so behavior is observed by running the binary.

### File responsibilities
| Package/File | Purpose |
|---|---|
| `pkg/syntax/` | Core syntax: variables, multiple returns, named returns, pointers vs values, defer, zero-values, visibility, iota enums |
| `pkg/concurrency/` | Goroutines, channels, cancellation-aware pipelines, mutex-protected shared state, `sync.WaitGroup`, worker pools, context values |
| `pkg/interfaces/` | `Shape` interface, `Rectangle`/`Circle`, polymorphism, type assertion, type switch |
| `pkg/generics/` | `List[T]`, `MapValues[T,R]` generic helpers |
| `pkg/errors/` | Custom errors (`MyCustomError`), `errors.As` |
| `pkg/collections/` | Generic data structures (`Set`, `Stack`, `Queue`, `OrderedMap`) and functional slice operations (`Filter`, `Reduce`, `FlatMap`, `GroupBy`, `Partition`, `Sorted`, `Distinct`, `Any`, `All`, `Zip`) – comparable to Java Collections + Stream API |
| `pkg/basictypes/` | Basic Go types manipulation: slices (lists), maps, mutability/copy semantics, strings, and JSON |
| `pkg/patterns/` | GoF design patterns (Singleton, Factory Method, Builder, Strategy, Observer, Decorator, Iterator, Adapter, Template Method, Command) and Go Antipatterns (`antipatterns.go`) demonstrating Java pitfalls vs idiomatic Go solutions |
| `pkg/advanced/` | Advanced features: embedding (≈ inheritance), iota enums, bitmasks, closures/currying/memoization, defer/panic/recover, reflection, type constraints (`Number`), concurrent fan-out, struct tags/JSON, `embed` static assets. Also contains benchmarks, fuzzing, and Build System (Go vs Maven) comparison. |
| `pkg/iosystem/` | I/O and System-level programming: file manipulation, directory walking, env vars, CLI flags, exec commands, networking, signal handling |
| `pkg/commonlibs/` | Popular 3rd party libraries: `testify` (assertions), `zap` (logging), `uuid` (unique IDs), `gin` (web framework) |
| `pkg/web/` | HTTP handlers (`/hello`, `/json`, `/echo`) with logging middleware and JSON body parsing |

## Go Build System (for Java Developers)
- **Go tool**: A single executable CLI (`go`) handles build, test, package management, and more.
- **Go Modules (`go.mod`)**: Equivalent to Maven's `pom.xml`.
- **`go build`**: Compiles a native executable; static linking depends on build settings and platform dependencies.
- **`go test`**: Runs unit tests, benchmarks, and fuzzing.

## Common Java-isms to Avoid
| Java-ism | Idiomatic Go | Why? |
|---|---|---|
| Getters/Setters | Direct Field Access | Go prefers simplicity; only use methods for logic/encapsulation. |
| `Panic` as `Exception` | Return `error` | `Panic` is for unrecoverable errors; `error` is a normal value. |
| Typed nil error | `return nil` | Returning typed `*MyError(nil)` as `error` leaves non-nil type; `err != nil` is true! |
| Mutex by value | Pointer receiver `*T` | Go structs are value types; copying a mutex duplicates lock state and breaks synchronization. |
| Goroutine without cancel | `context.Context` / buffered chan | Background goroutines blocked on unbuffered channels leak permanently. |
| Producer-side Interface | Consumer-side Interface | Duck typing allows consumers to define only what they need. |
| Pointers to slice/map | Pass descriptors directly by value in most cases | Copies share referenced backing data; slice append can return a new descriptor. |
| Sub-slicing huge buffers | Explicit `copy()` | Sub-slicing `huge[:4]` pins the entire backing array in RAM (memory leak). |
| `new` everything | Zero-value usability | Many Go types are ready to use when declared (e.g., `sync.Mutex`). |
| Deep package nesting | Flat package structure | Go packages should be broad and meaningful, not deeply hierarchical. |
| Try-Catch blocks | `if err != nil` | Explicit error handling makes control flow obvious. |
| Object inheritance | Composition (Embedding) | Go favors composition over complex class hierarchies. |

## Critical workflows
- Run with Makefile (Recommended): `make run`
- Run full showcase: `go run ./cmd/gogo` (web server remains active until Enter is pressed in stdin).
- Run all tests: `make test` or `go test ./...`
- Build binary: `make build` (creates `gogo_binary`)
- Clean artifacts: `make clean`
- Run benchmarks: `go test -bench=. ./pkg/...`
- Run fuzz test: `go test -fuzz=FuzzReverse ./pkg/advanced`
- Run race detector: `go test -race ./...`
- Validate web endpoints while app is running:
  - `curl "http://localhost:8080/hello?name=GoExpert"`
  - `curl http://localhost:8080/json`
  - `curl -X POST -H 'Content-Type: application/json' -d '{"msg":"hello"}' http://localhost:8080/echo`
- Focus a single test when iterating: `go test ./pkg/generics -run TestMapValues -v`.

## Project-specific coding patterns
- Keep feature demos discoverable by exposing top-level functions named `Run*Demo` in `pkg/*/`.
- Preserve zero-value usability where present (example: `Set[T]{}`, `Stack[T]{}`, `Queue[T]{}`, `OrderedMap[K,V]{}`, and `List[T]{}` are valid before calling methods).
- Match existing testing style: direct `t.Errorf` checks, `reflect.DeepEqual` for slices, `math.Abs(…) > 1e-9` for float comparisons.
- Prefer standard library APIs already used in demos before introducing dependencies.
- GoF patterns use idiomatic Go: interfaces instead of abstract classes, embedding instead of inheritance, `sync.Once` for one-time initialization, and range functions or channels for iterators.

## Integration points and caveats
- No external services; core integrations use Go stdlib (`net/http`, `encoding/json`, `context`, `sync`, `reflect`), with third-party library examples in `pkg/commonlibs/`.
- `StartWebServer` uses global `http.HandleFunc`/default mux; registering the same routes twice in one process will fail.
- `GetAppConfig()` is a true process-wide singleton via `sync.Once` — calling it in tests shares state across the test binary.
- Server lifecycle is intentionally lightweight: launched with `go` routine, readiness is a fixed `100ms` sleep, and shutdown is process-driven.
- If adding new demos, wire them in `main.go` in sequence and keep output clearly sectioned like existing `--- Demo ---` / `--- Demo End ---` markers.
