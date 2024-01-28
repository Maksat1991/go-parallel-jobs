package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("success")
}

func run() error {
	numberChan := make(chan int)

	var result []int

	ctx := context.Background()
	ctx, cancelFunc := context.WithTimeout(ctx, time.Second*4) //change to 2
	defer cancelFunc()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		result = printNumbers(ctx, numberChan)
	}()

	generateNumbers(3, numberChan)

	close(numberChan)

	fmt.Println("Waiting for goroutines to finish...")
	wg.Wait()

	fmt.Println(len(result))

	fmt.Println("Done!")

	return nil
}

func generateNumbers(total int, ch chan<- int) {

	for idx := 1; idx <= total; idx++ {
		time.Sleep(time.Second)
		fmt.Printf("sending %d to channel\n", idx)
		ch <- idx
	}
}

func printNumbers(ctx context.Context, ch <-chan int) []int {
	result := make([]int, 0)

	for {
		select {
		case <-ctx.Done():
			return result
		case num, ok := <-ch:
			if !ok {
				return result
			}
			fmt.Printf("read %d from channel\n", num)
			result = append(result, num)
		}

	}

}
