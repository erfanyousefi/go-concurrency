package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

func fetchData(ctx context.Context, id int) (string, error) {
	select {
	case <-time.After(time.Duration(rand.Intn(2000) * int(time.Minute.Milliseconds()))):
		return fmt.Sprintf("Raw data: %d", id), nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}
func processData(ctx context.Context, rawData string) (string, error) {
	select {
	case <-time.After(time.Duration(rand.Intn(2000) * int(time.Millisecond))):
		return fmt.Sprintf("processed data: %s", rawData), nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}
func saveData(ctx context.Context, data string) error {
	select {
	case <-time.After(time.Duration(rand.Intn(2000) * int(time.Millisecond))):
		fmt.Printf("Saved data successfully: %s \n", data)
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func worker(ctx context.Context, jobID int) {
	fmt.Printf("Worker %d started \n", jobID)
	rawData, err1 := fetchData(ctx, jobID)
	if err1 != nil {
		fmt.Printf("Worker %d stopped from fetchData %v \n", jobID, err1)
		return
	}
	processedData, err2 := processData(ctx, rawData)
	if err2 != nil {
		fmt.Printf("Worker %d stopped from ProcessData %v \n", jobID, err2)
		return
	}
	err3 := saveData(ctx, processedData)
	if err3 != nil {
		fmt.Printf("Worker %d stopped from saveData %v \n", jobID, err3)
		return
	}
	fmt.Printf("Worker %d done successfully \n", jobID)
}
func main() {
	rand.Seed(time.Now().UnixNano())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	for i := range 3 {
		go worker(ctx, i)
	}
	<-ctx.Done()
	fmt.Printf("Main finished: %v \n", ctx.Err())

}
