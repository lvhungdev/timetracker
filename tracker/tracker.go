package tracker

import (
	"errors"
	"sort"
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

	sort.Slice(records, func(i, j int) bool {
		return records[i].Start.Before(records[j].Start)
	})

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

func (t *Tracker) GetAll(from, to time.Time) []Record {
	results := []Record{}

	for curr := from; curr.Sub(to).Hours() <= 24; curr = curr.AddDate(0, 0, 1) {
		records := t.repo.GetFromDate(curr)
		for _, r := range records {
			if (r.Start.After(from) || r.Start == from) && (r.End.Before(to) || r.End == to) {
				results = append(results, r)
			}
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Start.Before(results[j].Start)
	})

	t.records = results
	return t.records
}

func (t *Tracker) StartTracking(name string) (old *Record, curr *Record, err error) {
	if name == "" {
		err = errors.New("name cannot be empty")
		return
	}

	now := time.Now()

	old = t.GetCurrent()
	if old != nil {
		old.End = now
	}

	t.records = append(t.records, newRecord(name, now))

	curr = &t.records[len(t.records)-1]

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
