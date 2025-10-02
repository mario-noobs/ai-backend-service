package helper

import (
	"testing"
	"time"
)

func TestTimer_Start(t *testing.T) {
	timer := &Timer{}

	// Record time before starting
	beforeStart := time.Now()

	// Start the timer
	timer.Start()

	// Record time after starting
	afterStart := time.Now()

	// Verify that the start time is set within a reasonable range
	if timer.start.Before(beforeStart) || timer.start.After(afterStart) {
		t.Errorf("Timer start time not set correctly. Expected between %v and %v, got %v",
			beforeStart, afterStart, timer.start)
	}
}

func TestTimer_End(t *testing.T) {
	timer := &Timer{}

	// Start the timer
	timer.Start()

	// Sleep for a known duration
	sleepDuration := 10 * time.Millisecond
	time.Sleep(sleepDuration)

	// End the timer and get elapsed time
	elapsed := timer.End()

	// Verify that elapsed time is at least the sleep duration
	// We allow some tolerance for execution time
	minExpected := sleepDuration.Milliseconds()
	maxExpected := minExpected + 10 // 10ms tolerance

	if elapsed < minExpected {
		t.Errorf("Expected elapsed time to be at least %d ms, got %d ms", minExpected, elapsed)
	}

	if elapsed > maxExpected {
		t.Errorf("Expected elapsed time to be at most %d ms, got %d ms", maxExpected, elapsed)
	}
}

func TestTimer_EndWithoutStart(t *testing.T) {
	timer := &Timer{}

	// Call End without calling Start first
	elapsed := timer.End()

	// Since start time is zero, elapsed should be a large positive number
	// representing time since epoch
	if elapsed <= 0 {
		t.Errorf("Expected positive elapsed time even without Start(), got %d", elapsed)
	}
}

func TestTimer_MultipleStartEnd(t *testing.T) {
	timer := &Timer{}

	// First measurement
	timer.Start()
	time.Sleep(5 * time.Millisecond)
	elapsed1 := timer.End()

	// Second measurement (should reset the start time)
	timer.Start()
	time.Sleep(15 * time.Millisecond)
	elapsed2 := timer.End()

	// Second measurement should be longer
	if elapsed2 <= elapsed1 {
		t.Errorf("Expected second measurement (%d ms) to be longer than first (%d ms)",
			elapsed2, elapsed1)
	}

	// Verify both measurements are reasonable
	if elapsed1 < 5 || elapsed1 > 15 {
		t.Errorf("First measurement out of expected range: %d ms", elapsed1)
	}

	if elapsed2 < 15 || elapsed2 > 25 {
		t.Errorf("Second measurement out of expected range: %d ms", elapsed2)
	}
}

func TestTimer_ConcurrentAccess(t *testing.T) {
	timer := &Timer{}

	// Start timer in main goroutine
	timer.Start()

	// Create a channel to receive results from goroutines
	results := make(chan int64, 2)

	// Start two goroutines that will call End()
	go func() {
		time.Sleep(10 * time.Millisecond)
		results <- timer.End()
	}()

	go func() {
		time.Sleep(20 * time.Millisecond)
		results <- timer.End()
	}()

	// Collect results
	elapsed1 := <-results
	elapsed2 := <-results

	// Both should be positive and reasonable
	if elapsed1 <= 0 || elapsed2 <= 0 {
		t.Errorf("Expected positive elapsed times, got %d and %d", elapsed1, elapsed2)
	}

	// The difference should be reasonable (around 10ms)
	diff := elapsed2 - elapsed1
	if diff < 5 || diff > 15 {
		t.Errorf("Expected difference around 10ms, got %d ms", diff)
	}
}
