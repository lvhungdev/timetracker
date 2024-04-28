package tracker

import (
	"testing"
	"time"
)

func TestGetCurrent(t *testing.T) {
	tracker := New()
	tracker.records = append(tracker.records, Record{
		Name:  "task 1",
		Start: time.Now(),
	})

	curr := tracker.GetCurrent()

	if curr == nil {
		t.Fatalf("expect curr not nil, got nil")
	}
	if curr.Name != "task 1" {
		t.Fatalf("expect curr.Name 'task 1', got '%v'", curr.Name)
	}
}

func TestStartTracking(t *testing.T) {
	tracker := New()
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
	tracker := New()
	tracker.records = append(tracker.records, Record{
		Name:  "task 1",
		Start: time.Now(),
	})

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
	if new == nil {
		t.Fatalf("expect new not nil, got nil")
	}
	if new.Name != "task 2" {
		t.Fatalf("expect new.name 'task 2', got '%v'", new.Name)
	}
}

func TestStartTrackingEmptyName(t *testing.T) {
	tracker := New()

	_, _, err := tracker.StartTracking("")

	if err == nil {
		t.Fatalf("expect error not nil, got error nil")
	}
}

func TestStopTracking(t *testing.T) {
	tracker := New()
	tracker.records = append(tracker.records, Record{
		Name:  "task 1",
		Start: time.Now(),
	})

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
	tracker := New()

	_, err := tracker.StopTracking()
	if err == nil {
		t.Fatalf("expect err not nil, got nil")
	}
}
