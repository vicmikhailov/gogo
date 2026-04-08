// Package commonlibs showcases popular third-party libraries frequently used in the Go ecosystem. // Package description for third-party libraries showcase.
// For a Java developer: // Guidance for Java background readers.
// - zap: Equivalent to SLF4J + Logback/Log4j2. High-performance, structured logging. // Explaining zap library analogy.
// - uuid: Equivalent to java.util.UUID. // Explaining uuid analogy.
// - gin: Equivalent to Spring Boot (Web) or JAX-RS. The most popular web framework in Go. // Explaining gin analogy.
// - testify: Equivalent to JUnit assertions (AssertJ/Hamcrest) and Mockito. // Explaining testify analogy.
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

// RunCommonLibsDemo demonstrates how to use popular Go libraries. // Comment for the demo orchestrator.
func RunCommonLibsDemo() { // Entry point for the common libraries showcase.
	fmt.Println("--- Common Libraries Demo ---") // Printing the section header.

	// 1. Uber-zap: Structured, High-Performance Logging // Comment for zap demo section.
	// Java comparison: Similar to using SLF4J with a JSON appender, but much faster. // Explaining Java analogy.
	runZapDemo() // Executing the zap logging demo.

	// 2. Google UUID: Generating unique identifiers // Comment for UUID demo section.
	// Java comparison: Exactly like java.util.UUID.randomUUID(). // Explaining Java analogy.
	runUUIDDemo() // Executing the UUID demo.

	// 3. Gin Gonic: Modern Web Framework // Comment for Gin demo section.
	// Java comparison: Similar to Spring Boot or Micronaut, but much lighter. // Explaining Java analogy.
	// We'll just show the setup, not block the whole process. // Noting that server isn't started.
	runGinDemo() // Executing the Gin setup demo.

	// 4. Testify Mock: Mocking dependencies // Comment for mocking demo section.
	// Java comparison: Exactly like Mockito. // Explaining Java analogy.
	runMockDemo() // Executing the mocking demo.

	// 5. Errgroup: Structured concurrency and error propagation // Comment for errgroup demo section.
	// Java comparison: Similar to CompletableFuture.allOf() with error handling. // Explaining Java analogy.
	runErrgroupDemo() // Executing the errgroup concurrency demo.

	fmt.Println("--- Common Libraries Demo End ---") // Printing the section footer.
}

// For a Java developer: // Guidance for Java background readers.
// - Similar to SLF4J + Logback, but designed for zero-allocations in hot paths. // Explaining performance focus.
// - Structured logging (JSON) is first-class, making it easier to parse logs in ELK/Splunk. // Explaining structured logging benefits.
func runZapDemo() { // Function to showcase zap logging.
	fmt.Println("1. Logging with zap (Structured & Fast):") // Printing the zap demo header.

	// In Go, we often use 'Sugar' for a more familiar, printf-like API, // Comment about zap SugaredLogger.
	// or the 'Logger' for maximum performance with type-safe fields. // Comment about core Logger.
	logger, _ := zap.NewProduction() // Creating a production-grade zap logger (JSON output).
	defer logger.Sync()              // Ensuring any buffered logs are flushed before function returns.

	sugar := logger.Sugar() // Getting a SugaredLogger for printf-style logging.

	// Structured logging allows you to add context as key-value pairs // Comment about key-value structured logging.
	// instead of just formatting strings. // Reinforcing the difference vs printf.
	sugar.Infow("Failed to fetch URL", // Logging an info message with additional context fields.
		"url", "http://example.com", // Adding the URL field.
		"attempt", 3, // Adding the attempt count.
		"backoff", time.Second, // Adding the backoff duration.
	)

	// You can also use printf style // Comment about printf-style convenience.
	sugar.Infof("Hello from zap sugar logger!") // Logging using formatted string.

	fmt.Println("   Check the console output above (it's JSON in production mode!)") // Noting JSON output in production mode.
}

// For a Java developer: // Guidance for Java background readers.
// - Java equivalent: `java.util.UUID`. // Comparing to Java UUID API.
func runUUIDDemo() { // Function to showcase UUID generation and parsing.
	fmt.Println("\n2. UUID Generation (google/uuid):") // Printing the UUID demo header.

	// Generating a version 4 (random) UUID // Comment for UUID generation step.
	id := uuid.New()                                   // Generating a random UUID (v4).
	fmt.Printf("   Generated UUID: %s\n", id.String()) // Printing the generated UUID as a string.

	// Parsing a UUID from string // Comment for parsing step.
	parsed, err := uuid.Parse(id.String()) // Parsing the UUID string back to a UUID object.
	if err == nil {                        // Checking if parsing succeeded.
		fmt.Printf("   Parsed UUID back successfully: %v\n", parsed.Version()) // Printing the UUID version.
	}
}

// For a Java developer: // Guidance for Java background readers.
// - Gin is the most popular Go web framework, similar to Spring Boot or Micronaut but much lighter. // Comparing to Java frameworks.
// - Middleware is equivalent to Servlet Filters or Spring Interceptors. // Explaining middleware analogy.
func runGinDemo() { // Function to showcase Gin setup and middleware.
	fmt.Println("\n3. Web Framework (gin-gonic/gin) with Middleware:") // Printing the Gin demo header.

	// Set Gin to release mode to keep stdout clean // Comment about Gin mode.
	gin.SetMode(gin.ReleaseMode) // Switching Gin to release mode.

	// In Java/Spring, you'd use @RestController. In Gin, you create an engine. // Explaining controller vs engine.
	r := gin.New() // Creating a new Gin engine instance.

	// Middleware in Gin is like Spring Interceptors or Servlet Filters. // Explaining where middleware fits.
	// This custom middleware logs the request and handles any panic. // Describing middleware behavior.
	r.Use(func(c *gin.Context) {
		t := time.Now()                                                                           // Recording start time.
		c.Next()                                                                                  // Processing the request (next handler in the chain).
		latency := time.Since(t)                                                                  // Computing latency after request is processed.
		fmt.Printf("   [Custom Middleware] Request to %s took %v\n", c.Request.URL.Path, latency) // Logging route and latency.
	})
	r.Use(gin.Recovery()) // Adding Gin's built-in recovery middleware to handle panics gracefully.

	// Defining a route is very explicit (no annotations) // Comment about explicit routing.
	r.GET("/api/ping", func(c *gin.Context) { // Registering a GET route at /api/ping.
		// c.JSON is equivalent to returning a POJO in Spring @ResponseBody // Explaining JSON response.
		c.JSON(http.StatusOK, gin.H{ // Writing a JSON response with HTTP 200 OK.
			"message": "pong",    // A simple message field.
			"library": "gin",     // Indicating the library.
			"status":  "awesome", // Including a status field.
		})
	})

	fmt.Println("   Gin engine initialized with route: GET /api/ping") // Indicating route initialization.
	fmt.Println("   (In a real app, you would call r.Run(':8081'))")   // Clarifying we aren't starting the server here.
}

// 4. Mocking with testify/mock // Section header for mocking demo.
// Java comparison: Using @Mock and when(mock.method()).thenReturn(value) in Mockito. // Explaining Java analogy.

// Database defines a simple interface for mocking demonstration. // Comment for the Database interface.
type Database interface { // Interface representing a database dependency.
	GetUser(id string) string // Method to fetch a user by ID.
}

type mockDB struct { // Struct embedding testify's Mock type.
	mock.Mock // Embedding Mock provides methods like On, Called, AssertExpectations.
}

// GetUser is a mocked method. // Comment for the mocked method implementation.
func (m *mockDB) GetUser(id string) string { // Method satisfying Database interface using testify.Mock.
	args := m.Called(id) // Recording the call and retrieving configured return arguments.
	return args.String(0)
}

// For a Java developer: // Guidance for Java background readers.
// - Similar to Mockito's `when(...).thenReturn(...)` and `verify(...)`. // Comparing to Mockito usage.
func runMockDemo() { // Function to showcase creating and using a mock.
	fmt.Println("\n4. Mocking (stretchr/testify/mock):") // Printing the mocking demo header.

	// Create an instance of our mock // Comment for mock creation.
	m := new(mockDB) // Instantiating a new mock database.

	// Set expectations (similar to Mockito's when(...).thenReturn(...)) // Comment for expectation setup.
	m.On("GetUser", "123").Return("John Doe") // Configuring GetUser("123") to return "John Doe".

	// Use the mock // Comment for exercising the mock.
	result := m.GetUser("123")                 // Invoking the mocked method.
	fmt.Printf("   Mock Result: %s\n", result) // Printing the mocked result.

	// Assert that expectations were met (similar to Mockito's verify(...)) // Comment for verification step.
	m.AssertExpectations(nil)                     // Verifying that all expectations were satisfied.
	fmt.Println("   Mock expectations verified.") // Printing verification success message.
}

// For a Java developer: // Guidance for Java background readers.
//   - Similar to `CompletableFuture.allOf()` but with built-in error propagation // Comparing to Java CompletableFuture.
//     and cancellation (if one task fails, others are cancelled). // Explaining cancellation.
func runErrgroupDemo() { // Function to showcase errgroup and context cancellation.
	fmt.Println("\n5. Structured Concurrency (errgroup):") // Printing the errgroup demo header.

	// errgroup provides synchronization, error propagation, and context-aware // Comment describing errgroup features.
	// cancellation for groups of goroutines. // Continuation of description.
	g, ctx := errgroup.WithContext(context.Background()) // Creating a new errgroup with a cancellable context.

	urls := []string{"http://example.com", "http://google.com", "http://golang.org"} // Defining URLs to "fetch".

	for _, url := range urls { // Iterating over the URLs to launch a worker per URL.
		// Capture url for the closure // Comment explaining loop variable capture.
		u := url            // Assigning to a new variable to avoid closure capture pitfall.
		g.Go(func() error { // Launching a goroutine managed by errgroup.
			// Simulate an operation that could fail or be cancelled // Comment for simulated workload.
			select { // Selecting between simulated work and context cancellation.
			case <-time.After(10 * time.Millisecond): // Simulating a successful fetch after a brief delay.
				fmt.Printf("   Successfully fetched %s\n", u) // Printing success message.
				return nil                                    // Indicating success.
			case <-ctx.Done(): // If the context is cancelled or deadline exceeded.
				return ctx.Err() // Propagating the context error.
			}
		})
	}

	// Wait blocks until all goroutines have finished, or one returns an error. // Comment for the wait behavior.
	if err := g.Wait(); err != nil { // Waiting for completion and capturing first error.
		fmt.Printf("   Error from errgroup: %v\n", err) // Printing error if occurred.
	} else { // If no errors occurred.
		fmt.Println("   All tasks in errgroup finished successfully.") // Printing success when all tasks complete.
	}
}
