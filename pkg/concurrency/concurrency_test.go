package concurrency

import (
	"context"
	"sync"
	"testing"
	"time"
)

/**
 * ===========================================================================
 * Concurrency Tests
 * ===========================================================================
 *
 * For a Java developer:
 * - Testing concurrency in Go often involves using channels for synchronization.
 * - Use `time.After` to prevent tests from hanging indefinitely.
 */

func TestWorkerPool(t *testing.T) {
	jobs := make(chan int, 5)
	results := make(chan int, 5)

	// Start workers
	for w := 1; w <= 3; w++ {
		go worker(w, jobs, results)
	}

	// Send 3 jobs
	for j := 1; j <= 3; j++ {
		jobs <- j
	}
	close(jobs)

	// Collect results with a timeout to avoid hanging if there's a bug
	timeout := time.After(500 * time.Millisecond)
	for a := 1; a <= 3; a++ {
		select {
		case res := <-results:
			if res%2 != 0 {
				t.Errorf("Expected even result from worker, got %d", res)
			}
		case <-timeout:
			t.Fatal("Test timed out - potential deadlock in worker pool")
		}
	}
}

func TestContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	// Start a goroutine that waits for context cancellation
	done := make(chan bool)
	go func() {
		<-ctx.Done()
		done <- true
	}()

	// Cancel and check if it propagated
	cancel()

	select {
	case <-done:
		// success
	case <-time.After(100 * time.Millisecond):
		t.Error("Context cancellation signal not received")
	}
}

func TestPipelineCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	output := square(ctx, gen(ctx, 1, 2, 3, 4, 5))

	select {
	case got := <-output:
		if got != 1 {
			t.Fatalf("first squared value = %d, want 1", got)
		}
	case <-time.After(500 * time.Millisecond):
		cancel()
		t.Fatal("pipeline did not produce a value")
	}

	cancel()
	done := make(chan struct{})
	go func() {
		for range output {
		}
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("pipeline stages did not stop after cancellation")
	}
}

func TestSafeCounterConcurrentUpdates(t *testing.T) {
	var counter SafeCounter
	const workers = 20
	const increments = 100

	var wg sync.WaitGroup
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < increments; j++ {
				counter.Add(1)
			}
		}()
	}
	wg.Wait()
	if got, want := counter.Value(), workers*increments; got != want {
		t.Fatalf("counter value = %d, want %d", got, want)
	}
}
