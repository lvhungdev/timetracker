package tracker

import (
	"errors"
	"time"
)

type Record struct {
	Name  string
	Start time.Time
	End   time.Time
}

type Tracker struct {
	records []Record
}

func New() *Tracker {
	return &Tracker{
		records: []Record{},
	}
}

func (t *Tracker) GetCurrent() *Record {
	for _, r := range t.records {
		if r.End.IsZero() {
			return &r
		}
	}

	return nil
}

func (t *Tracker) StartTracking(name string) (curr *Record, new *Record, err error) {
	if name == "" {
		err = errors.New("name cannot be empty")
		return
	}

	now := time.Now()

	curr = t.GetCurrent()
	if curr != nil {
		curr.End = now
	}

	t.records = append(t.records, Record{
		Name:  name,
		Start: now,
	})

	new = &t.records[len(t.records)-1]

	return
}

func (t *Tracker) StopTracking() (*Record, error) {
	curr := t.GetCurrent()
	if curr == nil {
		return nil, errors.New("no active record found")
	}

	curr.End = time.Now()

	return curr, nil
}
