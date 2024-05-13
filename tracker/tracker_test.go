package tracker

import (
	"testing"
	"time"
)

type mockRepo struct {
	records []Record
}

func newMockRepo(records []Record) *mockRepo {
	if records == nil {
		records = []Record{}
	}

	return &mockRepo{
		records: records,
	}
}

func (r *mockRepo) GetToday() []Record {
	return r.GetFromDate(time.Now())
}

func (r *mockRepo) GetFromDate(date time.Time) []Record {
	results := []Record{}
	from := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local)
	to := time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 0, time.Local)

	for _, record := range r.records {
		if (record.Start.After(from) || record.Start == from) && (record.End.Before(to) || record.End == to) {
			results = append(results, record)
		}
	}

	return results
}

func (r *mockRepo) SaveAll(records []Record) error {
	r.records = records
	return nil
}

func TestGetCurrent(t *testing.T) {
	repo := newMockRepo([]Record{newRecord("task 1", time.Now())})
	tracker := New(repo)

	curr := tracker.GetCurrent()

	if curr == nil {
		t.Fatalf("expect curr not nil, got nil")
	}
	if curr.Name != "task 1" {
		t.Fatalf("expect curr.Name 'task 1', got '%v'", curr.Name)
	}
}

func TestStartTracking(t *testing.T) {
	tracker := New(newMockRepo(nil))
	old, curr, err := tracker.StartTracking("task 1", time.Time{})

	if err != nil {
		t.Fatalf("expect error nil, got %v", err)
	}
	if old != nil {
		t.Fatalf("expect old nil, got not nil")
	}
	if curr == nil {
		t.Fatalf("expect curr not nil, got nil")
	}
	if curr.Name != "task 1" {
		t.Fatalf("expect curr.Name 'task 1', got '%v'", curr.Name)
	}
	if curr.Start.IsZero() {
		t.Fatalf("expect curr.Start to have value, got zero value")
	}
}

func TestStartTrackingWithCurr(t *testing.T) {
	repo := newMockRepo([]Record{newRecord("task 1", time.Now())})
	tracker := New(repo)

	old, curr, err := tracker.StartTracking("task 2", time.Time{})

	if err != nil {
		t.Fatalf("expect error nil, got %v", err)
	}
	if old == nil {
		t.Fatalf("expect old not nil, got nil")
	}
	if old.Name != "task 1" {
		t.Fatalf("expect old.Name 'task 1', got '%v'", old.Name)
	}
	if old.End.IsZero() {
		t.Fatalf("expect old.End not zero, got zero")
	}
	if curr == nil {
		t.Fatalf("expect curr not nil, got nil")
	}
	if curr.Name != "task 2" {
		t.Fatalf("expect curr.name 'task 2', got '%v'", curr.Name)
	}
}

func TestStartTrackingWithTimeSpecified(t *testing.T) {
	at := time.Now().Add(-time.Hour)
	tracker := New(newMockRepo(nil))

	old, curr, err := tracker.StartTracking("task 1", at)
	if err != nil {
		t.Fatalf("expect error nil, got %v", err)
	}
	if old != nil {
		t.Fatalf("expect old nil, got not nil")
	}
	if curr == nil {
		t.Fatalf("expect curr not nil, got nil")
	}
	if curr.Name != "task 1" {
		t.Fatalf("expect curr.Name 'task 1', got '%v'", curr.Name)
	}
	if curr.Start != at {
		t.Fatalf("expect curr.Start '%s 1', got '%s'", at, curr.Start)
	}
}

func TestStartTrackingWithTimeSpecifiedBeforeOldStart(t *testing.T) {
	now := time.Now()
	tracker := New(newMockRepo([]Record{newRecord("task 1", now)}))

	_, _, err := tracker.StartTracking("task 1", now.Add(-time.Hour))
	if err == nil {
		t.Fatalf("expect error not nil, got nil")
	}
}

func TestStartTrackingEmptyName(t *testing.T) {
	tracker := New(newMockRepo(nil))

	_, _, err := tracker.StartTracking("", time.Time{})

	if err == nil {
		t.Fatalf("expect error not nil, got error nil")
	}
}

func TestStopTracking(t *testing.T) {
	repo := newMockRepo([]Record{newRecord("task 1", time.Now())})
	tracker := New(repo)

	curr, err := tracker.StopTracking(time.Time{})
	if err != nil {
		t.Fatalf("expect err nil, got not nil")
	}
	if curr == nil {
		t.Fatalf("expect curr not new, got nil")
	}
	if curr.Name != "task 1" {
		t.Fatalf("expect curr.Name 'task 1', got '%v'", curr.Name)
	}
}

func TestStopTrackingWithoutCurrent(t *testing.T) {
	tracker := New(newMockRepo(nil))

	_, err := tracker.StopTracking(time.Time{})
	if err == nil {
		t.Fatalf("expect err not nil, got nil")
	}
}

func TestStopTrackingWithTimeSpecified(t *testing.T) {
	now := time.Now()
	repo := newMockRepo([]Record{newRecord("task 1", now.Add(-time.Hour*2))})
	tracker := New(repo)

	curr, err := tracker.StopTracking(now.Add(-time.Hour))
	if err != nil {
		t.Fatalf("expect err nil, got not nil")
	}
	if curr.End != now.Add(-time.Hour) {
		t.Fatalf("expect curr.End '%s', got '%s'", now.Add(-time.Hour), curr.End)
	}
}

func TestStopTrackingWithTimeSpecifiedBeforeStart(t *testing.T) {
	now := time.Now()
	repo := newMockRepo([]Record{newRecord("task 1", now)})
	tracker := New(repo)

	_, err := tracker.StopTracking(now.Add(-time.Hour))
	if err == nil {
		t.Fatalf("expect err not nil, got nil")
	}
}

func TestGetAllToday(t *testing.T) {
	now := time.Now()

	r1 := newRecord("task 1", time.Date(now.Year(), now.Month(), now.Day(), 14, 0, 0, 0, time.Local))
	r1.End = r1.Start.Add(time.Hour)
	r2 := newRecord("task 2", time.Date(now.Year(), now.Month(), now.Day(), 4, 0, 0, 0, time.Local))
	r2.End = r2.Start.Add(time.Hour * 2)
	repo := newMockRepo([]Record{r1, r2})
	tracker := New(repo)

	records := tracker.GetAll(
		time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local),
		time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.Local))

	if len(records) != 2 {
		t.Fatalf("expect 2 records, got %v", len(records))
	}
	if records[0].Name != "task 2" {
		t.Fatalf("expect records[0].Name 'task 2', got '%v'", records[0].Name)
	}
	if records[1].Name != "task 1" {
		t.Fatalf("expect records[1].Name 'task 1', got '%v'", records[1].Name)
	}
}
