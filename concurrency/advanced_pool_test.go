package concurrency

import (
	"context"
	"testing"
	"time"
)

func TestAdvancedPool(t *testing.T) {
	pool, err := NewAdvancedPool(12, 4)

	if err != nil {
		t.Fatal(err)
	}

	result := make(chan int, 1)

	result <- 0

	requestContext := context.Background()

	for i := 1; i <= 32; i++ {
		i := i

		pool.Submit(requestContext, func(ctx context.Context) {
			time.Sleep(5 * time.Second)
			result <- i + <-result
		})
	}
	err = pool.Close(requestContext)

	if err != nil {
		t.Log(err)
	}

	expectedResult := 0

	for i := 1; i <= 32; i++ {
		expectedResult += i
	}

	actualResult := <-result

	if actualResult != expectedResult {
		t.Fatalf("expected %d but got %d", expectedResult, actualResult)
	}
}
