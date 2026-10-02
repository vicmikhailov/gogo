// Package concurrency showcases Go's concurrency model.
//
// For a Java developer:
//   - Goroutines (`go func()`) are NOT OS threads. They are "green threads" multiplexed
//     onto a small number of OS threads. They use ~2KB of stack space vs ~1MB in Java.
//   - Channels (`chan`) are the primary way to communicate between goroutines.
//     "Don't communicate by sharing memory; share memory by communicating."
//   - Java-ism to avoid: Using `sync.Mutex` for everything. While Go has Mutexes,
//     channels are often a cleaner way to coordinate state.
//   - Java-ism to avoid: Thinking `go` is like `new Thread().start()`. It's much cheaper;
//     don't be afraid to spawn thousands of them.
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

// contextKey is a custom type for context keys to prevent namespace collisions across packages.
type contextKey string

const requestIDKey contextKey = "request-id"

// RunConcurrencyDemo showcases Go's core concurrency primitives and patterns.
func RunConcurrencyDemo() {
	fmt.Println("--- Concurrency Demo ---")

	// 1. Simple Goroutine with Channels
	// Java equivalent: Creating a thread and using a SynchronousQueue.
	// - make(chan string): creates an unbuffered channel (synchronous; sender blocks until receiver is ready).
	// - go func(): starts a new concurrent goroutine.
	// - ch <- data: sends data into the channel.
	// - <-ch: receives data from the channel (blocks until data arrives).
	fmt.Println("1. Simple Goroutine and Channels:")
	ch := make(chan string)
	go func() {
		time.Sleep(100 * time.Millisecond)
		ch <- "Hello from a goroutine!"
	}()
	msg := <-ch
	fmt.Println("   Received:", msg)

	// 2. sync.WaitGroup for multiple workers
	// Java equivalent: java.util.concurrent.CountDownLatch(3).
	// - wg.Add(1): increments the task counter before launching each goroutine.
	// - defer wg.Done(): decrements the counter when the worker completes (like a finally block).
	// - wg.Wait(): blocks until the counter reaches zero.
	fmt.Println("2. sync.WaitGroup for multiple workers:")
	var wg sync.WaitGroup
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			fmt.Printf("   Worker %d is working...\n", id)
			time.Sleep(50 * time.Millisecond)
		}(i)
	}
	wg.Wait()
	fmt.Println("   All workers finished.")

	// 3. Select statement with timeout
	// Java equivalent: Complex polling with `Selector` or multiple `BlockingQueue.poll(timeout)`.
	// - `select` lets a goroutine wait simultaneously on multiple channel operations.
	// - case res := <-c1: triggers if c1 receives data.
	// - case <-time.After(...): triggers if the timeout duration expires first.
	fmt.Println("3. Select statement with timeout:")
	c1 := make(chan string)
	go func() {
		time.Sleep(200 * time.Millisecond)
		c1 <- "Result 1"
	}()

	select {
	case res := <-c1:
		fmt.Println("   Received:", res)
	case <-time.After(100 * time.Millisecond):
		fmt.Println("   Timeout reached (as expected)!")
	}

	// 4. Context for cancellation
	// Java equivalent: `Thread.interrupt()` or `ExecutorService.shutdownNow()`.
	// Go uses `context.Context` to propagate cancellation signals down the call tree.
	// - context.WithTimeout: creates a context with an automatic deadline.
	// - defer cancel(): releases timer resources associated with the context.
	// - <-ctx.Done(): receives a signal when cancelled or timed out; ctx.Err() provides the reason.
	fmt.Println("4. Context for cancellation:")
	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()

	finished := make(chan bool)
	go func(ctx context.Context) {
		select {
		case <-time.After(500 * time.Millisecond):
			fmt.Println("   Worker finished on its own (should not happen)")
			finished <- true
		case <-ctx.Done():
			fmt.Println("   Worker received cancellation signal:", ctx.Err())
			finished <- false
		}
	}(ctx)

	<-finished

	// 5. Context Value Propagation
	// Java equivalent: ThreadLocal or Spring RequestContextHolder.
	// Context can also carry request-scoped metadata (like request IDs or auth tokens) down the stack.
	fmt.Println("5. Context Value Propagation:")
	ctxValue := context.WithValue(context.Background(), requestIDKey, "req-12345")
	processRequest(ctxValue)

	// 6. Worker Pool (fan-out pattern)
	// Java equivalent: ExecutorService with a fixed thread pool (Executors.newFixedThreadPool(3)).
	// - jobs := make(chan int, 5): buffered channel holding up to 5 pending tasks without blocking.
	// - close(jobs): signals to workers that no more jobs will be produced.
	fmt.Println("6. Worker Pool (fan-out):")
	jobs := make(chan int, 5)
	results := make(chan int, 5)

	for w := 1; w <= 3; w++ {
		go worker(w, jobs, results)
	}

	for j := 1; j <= 5; j++ {
		jobs <- j
	}
	close(jobs)

	for a := 1; a <= 5; a++ {
		<-results
	}
	fmt.Println("   All jobs completed via worker pool.")

	// 7. Pipeline Pattern (generator -> square -> print)
	// Java comparison: Java Streams (.map()) or Reactive Streams (RxJava/Project Reactor).
	// Go idiom: Connect stages using channels where each stage runs in its own goroutine.
	fmt.Println("7. Pipeline Pattern (generator -> square -> print):")
	nums := gen(2, 3)
	sq := square(nums)

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
//   - Similar to extracting values from a `ThreadLocal` or a Request Attribute in Spring.
//   - Unlike ThreadLocal, Context is passed explicitly down the call hierarchy.
//   - `ctx.Value` returns `any`; type assertion `v.(string)` casts it safely.
func processRequest(ctx context.Context) {
	if reqID, ok := ctx.Value(requestIDKey).(string); ok {
		fmt.Printf("   Processing request with ID: %s\n", reqID)
	} else {
		fmt.Println("   No request ID found in context.")
	}
}

// worker consumes tasks from one channel and sends results to another.
//
// For a Java developer:
//   - Directional channels provide compile-time safety:
//   - `jobs <-chan int`: receive-only channel (input).
//   - `results chan<- int`: send-only channel (output).
//   - `for j := range jobs` iterates until the channel is closed.
func worker(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		fmt.Printf("   Worker %d started job %d\n", id, j)
		time.Sleep(10 * time.Millisecond)
		results <- j * 2
	}
}

// gen converts a variable number of integers into a channel stream (generator stage).
//
// For a Java developer:
//   - Java equivalent: `Stream.of(nums)`.
//   - Go idiom: Functions returning `<-chan T` (receive-only) act as concurrent generators.
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

// square receives integers from an input channel, squares them, and emits to an output channel.
//
// For a Java developer:
//   - Java equivalent: `.map(n -> n * n)`.
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
