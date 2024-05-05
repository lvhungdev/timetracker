package tracker

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Record struct {
	Id    string
	Name  string
	Start time.Time
	End   time.Time
}

func newRecord(name string, start time.Time) Record {
	return Record{
		Id:    uuid.NewString(),
		Name:  name,
		Start: start,
	}
}

type Tracker struct {
	repo    Repo
	records []Record
}

func New(repo Repo) *Tracker {
	records := repo.GetToday()

	return &Tracker{
		repo:    repo,
		records: records,
	}
}

func (t *Tracker) GetCurrent() *Record {
	for i := 0; i < len(t.records); i++ {
		if t.records[i].End.IsZero() {
			return &t.records[i]
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

	t.records = append(t.records, newRecord(name, now))

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

func (t *Tracker) Save() error {
	return t.repo.SaveAll(t.records)
}
