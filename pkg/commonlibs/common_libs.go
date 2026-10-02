// Package commonlibs showcases popular third-party libraries frequently used in the Go ecosystem.
//
// For a Java developer:
//   - zap: Equivalent to SLF4J + Logback/Log4j2. High-performance, structured logging.
//   - uuid: Equivalent to java.util.UUID.
//   - gin: Equivalent to Spring Boot (Web) or JAX-RS. The most popular web framework in Go.
//   - testify: Equivalent to JUnit assertions (AssertJ/Hamcrest) and Mockito.
//   - errgroup: Equivalent to CompletableFuture.allOf() with structured concurrency and error propagation.
package commonlibs

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

// RunCommonLibsDemo orchestrates demonstrations of popular third-party Go libraries.
func RunCommonLibsDemo() {
	fmt.Println("--- Common Libraries Demo ---")

	// 1. Uber-zap: Structured, High-Performance Logging
	// Java comparison: Similar to using SLF4J with a JSON appender, but much faster.
	runZapDemo()

	// 2. Google UUID: Generating unique identifiers
	// Java comparison: Exactly like java.util.UUID.randomUUID().
	runUUIDDemo()

	// 3. Gin Gonic: Modern Web Framework
	// Java comparison: Similar to Spring Boot or Micronaut, but much lighter.
	// We demonstrate the engine setup and routing without blocking the process.
	runGinDemo()

	// 4. Testify Mock: Mocking dependencies
	// Java comparison: Exactly like Mockito.
	runMockDemo()

	// 5. Errgroup: Structured concurrency and error propagation
	// Java comparison: Similar to CompletableFuture.allOf() with error handling.
	runErrgroupDemo()

	fmt.Println("--- Common Libraries Demo End ---")
}

// runZapDemo showcases structured logging with uber-go/zap.
//
// For a Java developer:
//   - Similar to SLF4J + Logback, but designed for zero-allocations in hot paths.
//   - Structured logging (JSON) is first-class, making it easier to parse in ELK/Splunk.
//   - In Go, we often use `Sugar` for a familiar printf-like API, or `Logger` for maximum type-safe performance.
//   - Structured key-value pairs (e.g. "url", "attempt", "backoff") replace string concatenation.
func runZapDemo() {
	fmt.Println("1. Logging with zap (Structured & Fast):")

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	sugar := logger.Sugar()

	sugar.Infow("Failed to fetch URL",
		"url", "http://example.com",
		"attempt", 3,
		"backoff", time.Second,
	)

	sugar.Infof("Hello from zap sugar logger!")

	fmt.Println("   Check the console output above (it's JSON in production mode!)")
}

// runUUIDDemo showcases generating and parsing UUIDs with google/uuid.
//
// For a Java developer:
//   - Java equivalent: `java.util.UUID.randomUUID()` and `UUID.fromString(str)`.
func runUUIDDemo() {
	fmt.Println("\n2. UUID Generation (google/uuid):")

	// Generate a version 4 (random) UUID
	id := uuid.New()
	fmt.Printf("   Generated UUID: %s\n", id.String())

	// Parse a UUID from a string
	parsed, err := uuid.Parse(id.String())
	if err == nil {
		fmt.Printf("   Parsed UUID back successfully: %v\n", parsed.Version())
	}
}

// runGinDemo showcases setting up Gin routes and custom middleware.
//
// For a Java developer:
//   - Gin is the most popular Go web framework, comparable to Spring Boot or Micronaut but much lighter.
//   - In Spring, you define `@RestController` classes and `@GetMapping` annotations. In Gin, you create an engine and explicitly register routes (`r.GET`).
//   - Middleware is equivalent to Servlet Filters or Spring Interceptors.
//   - Gin handler functions take a `*gin.Context`; `c.JSON()` is equivalent to `@ResponseBody`.
func runGinDemo() {
	fmt.Println("\n3. Web Framework (gin-gonic/gin) with Middleware:")

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	// Custom logging middleware: records start time, invokes c.Next(), then logs latency
	r.Use(func(c *gin.Context) {
		t := time.Now()
		c.Next()
		latency := time.Since(t)
		fmt.Printf("   [Custom Middleware] Request to %s took %v\n", c.Request.URL.Path, latency)
	})
	r.Use(gin.Recovery())

	// Explicit route definition returning JSON
	r.GET("/api/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
			"library": "gin",
			"status":  "awesome",
		})
	})

	fmt.Println("   Gin engine initialized with route: GET /api/ping")
	fmt.Println("   (In a real app, you would call r.Run(':8081'))")
}

// ---------------------------------------------------------------------------
// 4. Mocking with testify/mock
// ---------------------------------------------------------------------------

// Database defines a simple interface for demonstrating mocks.
type Database interface {
	GetUser(id string) string
}

// mockDB embeds testify's `mock.Mock` to record method calls and return stubbed values.
//
// For a Java developer:
//   - Java equivalent: `@Mock private Database db;` with Mockito.
//   - Embedding `mock.Mock` provides `m.Called()`, `m.On()`, and `m.AssertExpectations()`.
type mockDB struct {
	mock.Mock
}

func (m *mockDB) GetUser(id string) string {
	args := m.Called(id)
	return args.String(0)
}

// runMockDemo demonstrates setting expectations and verifying mocks with stretchr/testify/mock.
//
// For a Java developer:
//   - `m.On("GetUser", "123").Return("John Doe")` is like `when(db.getUser("123")).thenReturn("John Doe")`.
//   - `m.AssertExpectations(nil)` is like `verify(db).getUser("123")`.
func runMockDemo() {
	fmt.Println("\n4. Mocking (stretchr/testify/mock):")

	m := new(mockDB)
	m.On("GetUser", "123").Return("John Doe")

	result := m.GetUser("123")
	fmt.Printf("   Mock Result: %s\n", result)

	m.AssertExpectations(nil)
	fmt.Println("   Mock expectations verified.")
}

// runErrgroupDemo showcases structured concurrency and error propagation with golang.org/x/sync/errgroup.
//
// For a Java developer:
//   - Java equivalent: `CompletableFuture.allOf()` with error handling, but with built-in context cancellation.
//   - If any task fails, the shared context is cancelled so other workers can stop early.
//   - `g.Wait()` blocks until all tasks succeed or returns the first encountered error.
func runErrgroupDemo() {
	fmt.Println("\n5. Structured Concurrency (errgroup):")

	g, ctx := errgroup.WithContext(context.Background())
	urls := []string{"http://example.com", "http://google.com", "http://golang.org"}

	for _, url := range urls {
		u := url // Capture loop variable for closure
		g.Go(func() error {
			select {
			case <-time.After(10 * time.Millisecond):
				fmt.Printf("   Successfully fetched %s\n", u)
				return nil
			case <-ctx.Done():
				return ctx.Err()
			}
		})
	}

	if err := g.Wait(); err != nil {
		fmt.Printf("   Error from errgroup: %v\n", err)
	} else {
		fmt.Println("   All tasks in errgroup finished successfully.")
	}
}
