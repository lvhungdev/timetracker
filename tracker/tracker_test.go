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
	return r.records
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
	old, curr, err := tracker.StartTracking("task 1")

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
}

func TestStartTrackingWithCurr(t *testing.T) {
	repo := newMockRepo([]Record{newRecord("task 1", time.Now())})
	tracker := New(repo)

	old, curr, err := tracker.StartTracking("task 2")

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

func TestStartTrackingEmptyName(t *testing.T) {
	tracker := New(newMockRepo(nil))

	_, _, err := tracker.StartTracking("")

	if err == nil {
		t.Fatalf("expect error not nil, got error nil")
	}
}

func TestStopTracking(t *testing.T) {
	repo := newMockRepo([]Record{newRecord("task 1", time.Now())})
	tracker := New(repo)

	curr, err := tracker.StopTracking()
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

func TestStopTrackingError(t *testing.T) {
	tracker := New(newMockRepo(nil))

	_, err := tracker.StopTracking()
	if err == nil {
		t.Fatalf("expect err not nil, got nil")
	}
}
