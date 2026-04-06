package concurrency

import (
	"context"
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
