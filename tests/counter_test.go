package tests

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/manojparvathaneni/go-counter/internal/counter"
)

const testFile = "test_counter.json"

func cleanup() {
	_ = os.Remove(testFile)
}

func TestSaveAndLoadCounter(t *testing.T) {
	cleanup()
	defer cleanup()

	vc := &counter.VisitCounter{}
	vc.Visits.Store(42)

	// Save counter
	if err := counter.SaveCounter(vc, testFile); err != nil {
		t.Fatalf("Failed to save counter: %v", err)
	}

	// Load counter
	loaded, err := counter.LoadCounter(testFile)
	if err != nil {
		t.Fatalf("Failed to load counter: %v", err)
	}

	if loaded.Visits.Load() != 42 {
		t.Errorf("Expected visit count 42, got %d", loaded.Visits.Load())
	}
}

func TestJSONEncodingRoundTrip(t *testing.T) {
	data := counter.CounterData{Visits: 123}
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	var parsed counter.CounterData
	if err := json.Unmarshal(bytes, &parsed); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if parsed.Visits != 123 {
		t.Errorf("Expected 123 visits after round-trip, got %d", parsed.Visits)
	}
}

func TestLoadFromEmptyFile(t *testing.T) {
	cleanup()
	defer cleanup()

	// Create an empty file
	f, err := os.Create(testFile)
	if err != nil {
		t.Fatalf("Could not create test file: %v", err)
	}
	f.Close()

	vc, err := counter.LoadCounter(testFile)
	if err != nil {
		t.Fatalf("Error loading empty file: %v", err)
	}

	if vc.Visits.Load() != 0 {
		t.Errorf("Expected zero visits from empty file, got %d", vc.Visits.Load())
	}
}

func TestSaveCreatesFile(t *testing.T) {
	cleanup()
	defer cleanup()

	vc := &counter.VisitCounter{}
	vc.Visits.Store(999)

	err := counter.SaveCounter(vc, testFile)
	if err != nil {
		t.Fatalf("Failed to save: %v", err)
	}

	// File should exist
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Errorf("Expected file %s to exist, but it does not", filepath.Base(testFile))
	}
}
