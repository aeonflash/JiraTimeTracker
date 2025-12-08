package main

import (
	"testing"
	"time"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

// **Feature: jira-reporting-window, Property 36: Request queue ordering**
// For any sequence of report requests, they should be processed in the order they were received.
// **Validates: Requirements 9.5**
func TestPropertyRequestQueueOrdering(t *testing.T) {
	properties := gopter.NewProperties(nil)

	properties.Property("requests are processed in the order they were received", prop.ForAll(
		func(numRequests int) bool {
			// Ensure we have at least 1 request and cap at reasonable number
			if numRequests < 1 {
				numRequests = 1
			}
			if numRequests > 50 {
				numRequests = 50
			}

			// Create a new request queue
			queue := NewRequestQueue()

			// Create and enqueue requests with sequential IDs
			requests := make([]ReportRequest, numRequests)
			for i := 0; i < numRequests; i++ {
				request := ReportRequest{
					ID: i + 1,
					DateRange: DateRange{
						StartDate: time.Now().AddDate(0, -1, 0),
						EndDate:   time.Now(),
					},
					Timestamp: time.Now().Add(time.Duration(i) * time.Millisecond),
				}
				requests[i] = request
				queue.Enqueue(request)
			}

			// Property 1: Queue size should match number of enqueued requests
			if queue.Size() != numRequests {
				return false
			}

			// Property 2: Dequeue all requests and verify they come out in order
			for i := 0; i < numRequests; i++ {
				dequeuedRequest := queue.Dequeue()
				
				// Check that we got a request
				if dequeuedRequest == nil {
					return false
				}

				// Check that the request ID matches the expected order
				expectedID := i + 1
				if dequeuedRequest.ID != expectedID {
					return false
				}

				// Check that the request matches the original
				if dequeuedRequest.ID != requests[i].ID {
					return false
				}
			}

			// Property 3: Queue should be empty after dequeuing all requests
			if !queue.IsEmpty() {
				return false
			}

			// Property 4: Dequeuing from empty queue should return nil
			if queue.Dequeue() != nil {
				return false
			}

			return true
		},
		gen.IntRange(1, 50), // numRequests
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// TestPropertyRequestQueueConcurrency tests thread-safe queue operations
func TestPropertyRequestQueueConcurrency(t *testing.T) {
	properties := gopter.NewProperties(nil)

	properties.Property("queue operations are thread-safe", prop.ForAll(
		func(numRequests int) bool {
			// Ensure we have at least 2 requests and cap at reasonable number
			if numRequests < 2 {
				numRequests = 2
			}
			if numRequests > 30 {
				numRequests = 30
			}

			// Create a new request queue
			queue := NewRequestQueue()

			// Enqueue requests concurrently
			done := make(chan bool, numRequests)
			for i := 0; i < numRequests; i++ {
				go func(id int) {
					request := ReportRequest{
						ID: id,
						DateRange: DateRange{
							StartDate: time.Now().AddDate(0, -1, 0),
							EndDate:   time.Now(),
						},
						Timestamp: time.Now(),
					}
					queue.Enqueue(request)
					done <- true
				}(i + 1)
			}

			// Wait for all enqueues to complete
			for i := 0; i < numRequests; i++ {
				<-done
			}

			// Property 1: Queue size should match number of enqueued requests
			if queue.Size() != numRequests {
				return false
			}

			// Property 2: All requests should be retrievable
			retrievedIDs := make(map[int]bool)
			for i := 0; i < numRequests; i++ {
				request := queue.Dequeue()
				if request == nil {
					return false
				}
				retrievedIDs[request.ID] = true
			}

			// Property 3: All unique IDs should be present
			if len(retrievedIDs) != numRequests {
				return false
			}

			// Property 4: Queue should be empty after dequeuing all
			if !queue.IsEmpty() {
				return false
			}

			return true
		},
		gen.IntRange(2, 30), // numRequests
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// TestPropertyRequestQueuePeek tests that Peek doesn't modify the queue
func TestPropertyRequestQueuePeek(t *testing.T) {
	properties := gopter.NewProperties(nil)

	properties.Property("peek does not modify queue", prop.ForAll(
		func(numRequests int) bool {
			// Ensure we have at least 1 request
			if numRequests < 1 {
				numRequests = 1
			}
			if numRequests > 20 {
				numRequests = 20
			}

			// Create a new request queue
			queue := NewRequestQueue()

			// Enqueue requests
			for i := 0; i < numRequests; i++ {
				request := ReportRequest{
					ID: i + 1,
					DateRange: DateRange{
						StartDate: time.Now().AddDate(0, -1, 0),
						EndDate:   time.Now(),
					},
					Timestamp: time.Now(),
				}
				queue.Enqueue(request)
			}

			initialSize := queue.Size()

			// Property 1: Peek should return the first request
			peekedRequest := queue.Peek()
			if peekedRequest == nil {
				return false
			}
			if peekedRequest.ID != 1 {
				return false
			}

			// Property 2: Queue size should not change after peek
			if queue.Size() != initialSize {
				return false
			}

			// Property 3: Multiple peeks should return the same request
			for i := 0; i < 5; i++ {
				peek := queue.Peek()
				if peek == nil || peek.ID != 1 {
					return false
				}
			}

			// Property 4: Size should still be unchanged
			if queue.Size() != initialSize {
				return false
			}

			// Property 5: Dequeue should return the same request as peek
			dequeuedRequest := queue.Dequeue()
			if dequeuedRequest == nil || dequeuedRequest.ID != peekedRequest.ID {
				return false
			}

			// Property 6: Size should now be one less
			if queue.Size() != initialSize-1 {
				return false
			}

			return true
		},
		gen.IntRange(1, 20), // numRequests
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// TestPropertyRequestQueueClear tests that Clear removes all requests
func TestPropertyRequestQueueClear(t *testing.T) {
	properties := gopter.NewProperties(nil)

	properties.Property("clear removes all requests", prop.ForAll(
		func(numRequests int) bool {
			// Ensure we have at least 1 request
			if numRequests < 1 {
				numRequests = 1
			}
			if numRequests > 20 {
				numRequests = 20
			}

			// Create a new request queue
			queue := NewRequestQueue()

			// Enqueue requests
			for i := 0; i < numRequests; i++ {
				request := ReportRequest{
					ID: i + 1,
					DateRange: DateRange{
						StartDate: time.Now().AddDate(0, -1, 0),
						EndDate:   time.Now(),
					},
					Timestamp: time.Now(),
				}
				queue.Enqueue(request)
			}

			// Property 1: Queue should have requests before clear
			if queue.IsEmpty() {
				return false
			}
			if queue.Size() != numRequests {
				return false
			}

			// Clear the queue
			queue.Clear()

			// Property 2: Queue should be empty after clear
			if !queue.IsEmpty() {
				return false
			}

			// Property 3: Size should be 0
			if queue.Size() != 0 {
				return false
			}

			// Property 4: Dequeue should return nil
			if queue.Dequeue() != nil {
				return false
			}

			// Property 5: Peek should return nil
			if queue.Peek() != nil {
				return false
			}

			return true
		},
		gen.IntRange(1, 20), // numRequests
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}
