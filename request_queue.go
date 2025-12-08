package main

// NewRequestQueue creates a new RequestQueue instance
func NewRequestQueue() *RequestQueue {
	return &RequestQueue{
		queue:      make([]ReportRequest, 0),
		processing: false,
	}
}

// Enqueue adds a new request to the queue
func (rq *RequestQueue) Enqueue(request ReportRequest) {
	rq.mutex.Lock()
	defer rq.mutex.Unlock()
	
	rq.queue = append(rq.queue, request)
}

// Dequeue removes and returns the next request from the queue
// Returns nil if queue is empty
func (rq *RequestQueue) Dequeue() *ReportRequest {
	rq.mutex.Lock()
	defer rq.mutex.Unlock()
	
	if len(rq.queue) == 0 {
		return nil
	}
	
	request := rq.queue[0]
	rq.queue = rq.queue[1:]
	
	return &request
}

// Peek returns the next request without removing it
// Returns nil if queue is empty
func (rq *RequestQueue) Peek() *ReportRequest {
	rq.mutex.Lock()
	defer rq.mutex.Unlock()
	
	if len(rq.queue) == 0 {
		return nil
	}
	
	return &rq.queue[0]
}

// Size returns the current number of requests in the queue
func (rq *RequestQueue) Size() int {
	rq.mutex.Lock()
	defer rq.mutex.Unlock()
	
	return len(rq.queue)
}

// IsEmpty returns true if the queue has no requests
func (rq *RequestQueue) IsEmpty() bool {
	rq.mutex.Lock()
	defer rq.mutex.Unlock()
	
	return len(rq.queue) == 0
}

// Clear removes all requests from the queue
func (rq *RequestQueue) Clear() {
	rq.mutex.Lock()
	defer rq.mutex.Unlock()
	
	rq.queue = make([]ReportRequest, 0)
}

// SetProcessing marks the queue as processing or not processing
func (rq *RequestQueue) SetProcessing(processing bool) {
	rq.mutex.Lock()
	defer rq.mutex.Unlock()
	
	rq.processing = processing
}

// IsProcessing returns true if a request is currently being processed
func (rq *RequestQueue) IsProcessing() bool {
	rq.mutex.Lock()
	defer rq.mutex.Unlock()
	
	return rq.processing
}

// GetAllRequests returns a copy of all requests in the queue (for testing)
func (rq *RequestQueue) GetAllRequests() []ReportRequest {
	rq.mutex.Lock()
	defer rq.mutex.Unlock()
	
	// Return a copy to prevent external modification
	result := make([]ReportRequest, len(rq.queue))
	copy(result, rq.queue)
	return result
}
