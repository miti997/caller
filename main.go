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

func sendRequest(wg *sync.WaitGroup, id int, apiGatewayURL string) {
	defer wg.Done()
	start := time.Now()

	resp, err := http.Get(apiGatewayURL)
	if err != nil {
		fmt.Printf("Request %d failed: %s\n", id, err)
		return
	}
	defer resp.Body.Close()

	duration := time.Since(start)
	fmt.Printf("Request %d: Status %d, Response Time: %v\n", id, resp.StatusCode, duration)
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter API Gateway URL: ")
	apiGatewayURL, _ := reader.ReadString('\n')
	apiGatewayURL = strings.TrimSpace(apiGatewayURL)

	fmt.Print("Enter number of requests: ")
	var numRequests int
	fmt.Scanln(&numRequests)

	var wg sync.WaitGroup
	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go sendRequest(&wg, i, apiGatewayURL)
	}
	wg.Wait()

	fmt.Println("All requests completed.")
}
