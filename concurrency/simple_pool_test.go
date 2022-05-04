package concurrency

import (
	"sync"
	"testing"
	"time"
)

func TestSimplePool(t *testing.T) {
	pool := NewSimplePool(8)

	result := make(chan int, 1)

	result <- 0

	var taskGroup sync.WaitGroup

	for i := 1; i <= 32; i++ {
		i := i
		taskGroup.Add(1)

		pool.Submit(func() {
			currentResult := <-result

			t.Logf("task #%v => result = %v + %v", i, i, currentResult)

			result <- (i + currentResult)

			time.Sleep(5 * time.Second)

			taskGroup.Done()
		})
	}

	taskGroup.Wait()

	expectedResult := 0

	for i := 1; i <= 32; i++ {
		expectedResult += i
	}

	actualResult := <-result

	t.Logf("result = %v", actualResult)

	if actualResult != expectedResult {
		t.Fatalf("expected %d but got %d", expectedResult, actualResult)
	}
}
