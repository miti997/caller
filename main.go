package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

func sendRequest(wg *sync.WaitGroup, id int, apiGatewayURL string, file *os.File, successCount *int, failureCount *int) {
	defer wg.Done()
	start := time.Now()

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Get(apiGatewayURL)
	if err != nil {
		*failureCount++
		msg := fmt.Sprintf("Request %d failed: %s\n", id, err)
		fmt.Print(msg)
		fmt.Fprint(file, msg)
		return
	}
	defer resp.Body.Close()

	*successCount++
	duration := time.Since(start)
	msg := fmt.Sprintf("Request %d: Status %d, Response Time: %v\n", id, resp.StatusCode, duration)
	fmt.Print(msg)
	fmt.Fprint(file, msg)
}

func main() {
	file, err := os.OpenFile("request_results.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter API Gateway URL: ")
	apiGatewayURL, _ := reader.ReadString('\n')
	apiGatewayURL = strings.TrimSpace(apiGatewayURL)

	fmt.Print("Enter number of requests: ")
	var numRequests int
	fmt.Scanln(&numRequests)

	fmt.Fprintln(file, "Request Results:")
	fmt.Fprintln(file, "-----------------")

	var successCount, failureCount int

	var wg sync.WaitGroup
	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go sendRequest(&wg, i, apiGatewayURL, file, &successCount, &failureCount)
	}
	wg.Wait()

	// Print the totals
	fmt.Printf("\nTotal Successful Requests: %d\n", successCount)
	fmt.Printf("Total Failed Requests: %d\n", failureCount)

	fmt.Fprintf(file, "\nTotal Successful Requests: %d\n", successCount)
	fmt.Fprintf(file, "Total Failed Requests: %d\n", failureCount)

	fmt.Println("All requests completed.")
}
