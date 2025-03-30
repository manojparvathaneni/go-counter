package counter

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"sync/atomic"
	"syscall"
)

type VisitCounter struct {
	Visits atomic.Int64
}

type CounterData struct {
	Visits int64 `json:"visits"`
}

var (
	ErrDecode = errors.New("could not decode counter data")
)

// LoadCounter loads the counter from a JSON file (locked for read)
func LoadCounter(path string) (*VisitCounter, error) {
	counter := &VisitCounter{}

	f, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_SH); err != nil {
		return nil, err
	}
	defer syscall.Flock(int(f.Fd()), syscall.LOCK_UN)

	var data CounterData
	if err := json.NewDecoder(f).Decode(&data); err != nil {
		log.Println("No previous data or invalid format. Starting fresh.")
		return counter, nil
	}

	counter.Visits.Store(data.Visits)
	return counter, nil
}

// SaveCounter stores the counter data into JSON file (locked for write)
func SaveCounter(counter *VisitCounter, path string) error {
	data := CounterData{Visits: counter.Visits.Load()}

	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX); err != nil {
		return err
	}
	defer syscall.Flock(int(f.Fd()), syscall.LOCK_UN)

	_, err = f.Write(bytes)
	return err
}
