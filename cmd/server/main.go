package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/manojparvathaneni/go-counter/internal/counter"
)

var (
	counterFile = "counter.json"
	autoSaveSec = 10 // default interval in seconds
	counterObj  *counter.VisitCounter
)

func main() {
	// Read interval from env
	if env := os.Getenv("COUNTER_AUTOSAVE_INTERVAL"); env != "" {
		if val, err := strconv.Atoi(env); err == nil && val > 0 {
			autoSaveSec = val
		}
	}

	// Load counter state
	var err error
	counterObj, err = counter.LoadCounter(counterFile)
	if err != nil {
		log.Printf("Failed to load counter, starting from zero: %v", err)
		counterObj = &counter.VisitCounter{}
	}

	// Set up context for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Start periodic save loop
	go startAutoSaver(ctx, time.Duration(autoSaveSec)*time.Second)

	// Set up HTTP handler
	mux := http.NewServeMux()
	mux.HandleFunc("/api/counter", counterHandler)

	srv := &http.Server{
		Addr:    ":8090",
		Handler: mux,
	}

	log.Printf("Counter server running on http://localhost:8090")
	log.Printf("Auto-save interval: %ds", autoSaveSec)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for shutdown
	<-ctx.Done()
	log.Println("Shutdown signal received. Saving counter...")

	if err := counter.SaveCounter(counterObj, counterFile); err != nil {
		log.Printf("Final save failed: %v", err)
	} else {
		log.Println("Counter saved successfully.")
	}

	ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctxTimeout); err != nil {
		log.Fatalf("Graceful shutdown failed: %v", err)
	}
	log.Println("Server shutdown complete.")
}

func counterHandler(w http.ResponseWriter, r *http.Request) {
	counterObj.Visits.Add(1)

	if err := counter.SaveCounter(counterObj, counterFile); err != nil {
		log.Println("Error saving counter:", err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(counter.CounterData{Visits: counterObj.Visits.Load()})
}

func startAutoSaver(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Auto-saver shutting down.")
			return
		case <-ticker.C:
			if err := counter.SaveCounter(counterObj, counterFile); err != nil {
				log.Println("Auto-save failed:", err)
			} else {
				log.Println("Counter auto-saved.")
			}
		}
	}
}
