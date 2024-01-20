package front

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

func Testsse(c *gin.Context) {

	// Use a WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Channel to communicate results
	resultCh := make(chan string)

	// Start multiple asynchronous tasks
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go asyncTask(i, &wg, resultCh)
	}

	// Use a goroutine to collect results
	go func() {
		// Wait for all tasks to finish
		wg.Wait()

		// Close the result channel to signal that no more results will be sent
		close(resultCh)
	}()

	// Read and send results to the client
	for result := range resultCh {
		c.SSEvent("message", result)
	}

	c.Status(http.StatusOK)
}

func asyncTask(id int, wg *sync.WaitGroup, resultCh chan<- string) {
	defer wg.Done()

	// Simulate some asynchronous processing
	time.Sleep(time.Second)

	result := fmt.Sprintf("Task %d completed", id)
	resultCh <- result
}
