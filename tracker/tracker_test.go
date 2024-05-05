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
	curr, new, err := tracker.StartTracking("task 1")

	if err != nil {
		t.Fatalf("expect error nil, got %v", err)
	}
	if curr != nil {
		t.Fatalf("expect curr nil, got not nil")
	}
	if new == nil {
		t.Fatalf("expect new not nil, got nil")
	}
	if new.Name != "task 1" {
		t.Fatalf("expect new.Name 'task 1', got '%v'", new.Name)
	}
}

func TestStartTrackingWithCurr(t *testing.T) {
	repo := newMockRepo([]Record{newRecord("task 1", time.Now())})
	tracker := New(repo)

	curr, new, err := tracker.StartTracking("task 2")

	if err != nil {
		t.Fatalf("expect error nil, got %v", err)
	}
	if curr == nil {
		t.Fatalf("expect curr not nil, got nil")
	}
	if curr.Name != "task 1" {
		t.Fatalf("expect curr.Name 'task 1', got '%v'", curr.Name)
	}
    if curr.End.IsZero() {
        t.Fatalf("expect curr.End not zero, got zero")
    }
	if new == nil {
		t.Fatalf("expect new not nil, got nil")
	}
	if new.Name != "task 2" {
		t.Fatalf("expect new.name 'task 2', got '%v'", new.Name)
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
