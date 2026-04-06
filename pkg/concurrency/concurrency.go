// Package concurrency showcases Go's concurrency model.
//
// For a Java developer:
//   - Goroutines (`go func()`) are NOT OS threads. They are "green threads" multiplexed
//     onto a small number of OS threads. They use ~2KB of stack space vs ~1MB in Java.
//   - Channels (`chan`) are the primary way to communicate between goroutines.
//     "Don't communicate by sharing memory; share memory by communicating."
//   - `sync.WaitGroup` is equivalent to Java's `CountDownLatch`.
//   - `context` is used for cancellation and timeouts, similar to `Future.cancel()`
//     or passing a `CancellationToken` in other languages.
package concurrency

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// contextKey is a custom type for context keys to avoid collisions.
type contextKey string

const requestIDKey contextKey = "request-id"

// RunConcurrencyDemo showcases Go's concurrency model.
func RunConcurrencyDemo() {
	fmt.Println("--- Concurrency Demo ---")

	// 1. Simple Goroutine with Channels
	// Java equivalent: Creating a thread and using a SynchronousQueue.
	fmt.Println("1. Simple Goroutine and Channels:")
	ch := make(chan string) // Unbuffered channel (blocking until both sides are ready)
	go func() {             // The 'go' keyword starts a new goroutine asynchronously
		time.Sleep(100 * time.Millisecond)
		ch <- "Hello from a goroutine!" // Send data into the channel
	}()
	msg := <-ch // Receive data from the channel (blocks until data is available)
	fmt.Println("   Received:", msg)

	// 2. sync.WaitGroup
	// Java equivalent: java.util.concurrent.CountDownLatch(3)
	fmt.Println("2. sync.WaitGroup for multiple workers:")
	var wg sync.WaitGroup
	for i := 1; i <= 3; i++ {
		wg.Add(1) // Increment the counter
		go func(id int) {
			defer wg.Done() // Decrement the counter when the function exits (like finally block)
			fmt.Printf("   Worker %d is working...\n", id)
			time.Sleep(50 * time.Millisecond)
		}(i)
	}
	wg.Wait() // Block until the counter reaches zero
	fmt.Println("   All workers finished.")

	// 3. Select statement and Timers
	// Java equivalent: Complex logic with `Selector` or polling multiple `BlockingQueue`s.
	// `select` lets a goroutine wait on multiple communication operations.
	fmt.Println("3. Select statement with timeout:")
	c1 := make(chan string)
	go func() {
		time.Sleep(200 * time.Millisecond)
		c1 <- "Result 1"
	}()

	select {
	case res := <-c1: // Triggers if c1 receives a value
		fmt.Println("   Received:", res)
	case <-time.After(100 * time.Millisecond): // Triggers if 100ms passes first
		fmt.Println("   Timeout reached (as expected)!")
	}

	// 4. Context for cancellation
	// Java equivalent: `Thread.interrupt()` or `ExecutorService.shutdownNow()`.
	// Go uses `context.Context` to propagate cancellation signals down the call tree.
	fmt.Println("4. Context for cancellation:")
	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel() // Good practice: always call cancel to release resources

	finished := make(chan bool)
	go func(ctx context.Context) {
		select {
		case <-time.After(500 * time.Millisecond):
			fmt.Println("   Worker finished on its own (should not happen)")
			finished <- true
		case <-ctx.Done(): // Triggers when the context is cancelled or times out
			fmt.Println("   Worker received cancellation signal:", ctx.Err())
			finished <- false
		}
	}(ctx)

	<-finished

	// 5. Context Value Propagation
	// Java equivalent: ThreadLocal.
	// Context can also carry request-scoped values (like IDs, tokens) through the stack.
	fmt.Println("5. Context Value Propagation:")
	ctxValue := context.WithValue(context.Background(), requestIDKey, "req-12345")
	processRequest(ctxValue)

	// 6. Worker Pool (fan-out pattern)
	// Java equivalent: ExecutorService with a fixed thread pool.
	fmt.Println("6. Worker Pool (fan-out):")
	jobs := make(chan int, 5) // Buffered channel with capacity 5
	results := make(chan int, 5)

	// Start 3 workers (goroutines)
	for w := 1; w <= 3; w++ {
		go worker(w, jobs, results)
	}

	// Send 5 jobs into the jobs channel
	for j := 1; j <= 5; j++ {
		jobs <- j
	}
	close(jobs) // Closing a channel indicates no more values will be sent

	// Collect results from the results channel
	for a := 1; a <= 5; a++ {
		<-results
	}
	fmt.Println("   All jobs completed via worker pool.")

	// 7. Pipeline Pattern
	// Java comparison: Java Streams or Reactor/RxJava.
	// Go idiom: Connect stages using channels. Each stage is a goroutine.
	fmt.Println("7. Pipeline Pattern (generator -> square -> print):")
	nums := gen(2, 3)
	sq := square(nums)

	// Consume the final stage of the pipeline
	fmt.Print("   Pipeline output: ")
	for n := range sq {
		fmt.Printf("%d ", n)
	}
	fmt.Println("\n   Pipeline completed.")

	fmt.Println("--- Concurrency Demo End ---")
}

// processRequest demonstrates extracting values from a Context.
//
// For a Java developer:
//   - This is similar to extracting values from a `ThreadLocal` or a Request Attribute in Spring.
//   - However, context is passed explicitly rather than being stored in thread-local storage.
//   - `ctx.Value` returns `any`; we use a type assertion `v.(string)` to cast it.
func processRequest(ctx context.Context) {
	// ctx.Value returns any; we use a type assertion (v.(string)) to cast it.
	if reqID, ok := ctx.Value(requestIDKey).(string); ok {
		fmt.Printf("   Processing request with ID: %s\n", reqID)
	} else {
		fmt.Println("   No request ID found in context.")
	}
}

// worker is a standard worker function that consumes from one channel and sends to another.
//
// For a Java developer:
//   - `jobs <-chan int` is a receive-only channel (input).
//   - `results chan<- int` is a send-only channel (output).
//   - This provides compile-time safety for channel usage.
//   - 'range' on a channel continues until the channel is closed.
func worker(id int, jobs <-chan int, results chan<- int) {
	// 'range' on a channel continues until the channel is closed.
	for j := range jobs {
		fmt.Printf("   Worker %d started job %d\n", id, j)
		time.Sleep(10 * time.Millisecond)
		results <- j * 2
	}
}

// gen converts a list of integers to a channel that emits them.
// This is the first stage of the pipeline.
//
// For a Java developer:
// - Java equivalent: `Stream.of(nums)`.
// - Go idiom: Functions that return a receive-only channel (`<-chan`) are often used as "generators".
func gen(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

// square receives integers from a channel and emits their squares.
// This is the second stage of the pipeline.
//
// For a Java developer:
// - Java equivalent: `.map(n -> n * n)`.
func square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()
	return out
}
