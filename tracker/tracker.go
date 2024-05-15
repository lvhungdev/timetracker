package tracker

import (
	"fmt"
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

	tkr := &Tracker{
		repo:    repo,
		records: records,
	}
	tkr.sortAndRepairRecords()

	return tkr
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
	records := []Record{}
	for curr := from; curr.Sub(to).Hours() <= 24; curr = curr.AddDate(0, 0, 1) {
		r := t.repo.GetFromDate(curr)
		//for _, r := range records {
		//	if (r.Start.After(from) || r.Start == from) && (r.End.Before(to) || r.End == to) {
		//		results = append(results, r)
		//	}
		//}
		records = append(records, r...)
	}

	t.records = records
	t.sortAndRepairRecords()

	results := []Record{}
	for _, r := range t.records {
		if r.End.IsZero() && (r.Start.Before(to)) && (r.Start.After(from) || r.Start == from) {
			results = append(results, r)
		} else if !r.End.IsZero() && (r.End.Before(to) || r.End == to) && (r.Start.After(from) || r.Start == from) {
			results = append(results, r)
		}
	}

	return results
}

func (t *Tracker) StartTracking(name string, at time.Time) (old *Record, curr *Record, err error) {
	if name == "" {
		err = fmt.Errorf("name cannot be empty")
		return
	}

	if at.IsZero() {
		at = time.Now()
	}

	old = t.GetCurrent()
	if old != nil {
		if old.Start.After(at) {
			err = fmt.Errorf("%s is after current's start time %s", at, old.Start)
			return
		}

		old.End = at
	}

	t.records = append(t.records, newRecord(name, at))

	curr = &t.records[len(t.records)-1]

	return
}

func (t *Tracker) StopTracking(at time.Time) (*Record, error) {
	curr := t.GetCurrent()
	if curr == nil {
		return nil, fmt.Errorf("no active record found")
	}

	if !at.IsZero() && at.Before(curr.Start) {
		return nil, fmt.Errorf("%s is before the current end time %s", at, curr.End)
	}

	if at.IsZero() {
		curr.End = time.Now()
	} else {
		curr.End = at
	}

	return curr, nil
}

func (t *Tracker) Save() error {
	return t.repo.SaveAll(t.records)
}

func (t *Tracker) sortAndRepairRecords() {
	sort.Slice(t.records, func(i, j int) bool {
		return t.records[i].Start.Before(t.records[j].Start)
	})

	now := time.Now()
	begin := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	for i := 0; i < len(t.records); i++ {
		r := t.records[i]
		if r.End.IsZero() && r.Start.Before(begin) {
			t.records[i].End = time.Date(r.Start.Year(), r.Start.Month(), r.Start.Day(), 23, 59, 59, 0, time.Local)
		}
	}
}
