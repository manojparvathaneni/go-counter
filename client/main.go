package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"time"
)

type counterData struct {
	Visits int64 `json:"visits"`
}

func main() {
	host := flag.String("host", "http://localhost:8090", "API base URL")
	times := flag.Int("times", 1, "Number of requests to send")
	delay := flag.Duration("delay", 0, "Delay between requests (e.g., 500ms)")
	flag.Parse()

	for i := 0; i < *times; i++ {
		resp, err := http.Get(*host + "/api/counter")
		if err != nil {
			fmt.Printf("Request failed: %v\n", err)
			return
		}
		defer resp.Body.Close()

		var data counterData
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			fmt.Println("Failed to parse response:", err)
			return
		}

		fmt.Printf("Visit #%d | Total Visits: %d\n", i+1, data.Visits)

		time.Sleep(*delay)
	}
}
