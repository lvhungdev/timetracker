package main

import (
	"testing"
	"time"
)

func TestGetCommandStart(t *testing.T) {
	cmd, err := getCommand([]string{"start", "task 1"})
	if err != nil {
		t.Fatalf("expected err to be nil, got '%v'", err)
	}
	start, ok := cmd.(cmdStartTracking)
	if !ok {
		t.Fatalf("expected start to be a startTracking command")
	}
	if start.name != "task 1" {
		t.Fatalf("expected start.name to be 'task 1', got '%s'", start.name)
	}
}

func TestGetCommandStop(t *testing.T) {
	cmd, err := getCommand([]string{"stop"})
	if err != nil {
		t.Fatalf("expected err to be nil, got '%v'", err)
	}
	_, ok := cmd.(cmdStopTracking)
	if !ok {
		t.Fatalf("expected stop to be a stopTracking command")
	}
}

func TestGetCommandReport(t *testing.T) {
	cmd, err := getCommand([]string{"report"})
	if err != nil {
		t.Fatalf("expected err to be nil, got '%v'", err)
	}
	_, ok := cmd.(cmdReport)
	if !ok {
		t.Fatalf("expected report to be a report command")
	}
}

func TestGetCommandInvalid(t *testing.T) {
	_, err := getCommand([]string{"invalid"})
	if err == nil {
		t.Fatalf("expected err to be not nil, got nil")
	}
}

func TestParseCmdReportDay(t *testing.T) {
	now := time.Now()
	from := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	to := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.Local)
	args := []string{"report", "day"}

	cmd, err := getCommand(args)

	if err != nil {
		t.Fatalf("expected err to be nil, got '%v'", err)
	}
	report, ok := cmd.(cmdReport)
	if !ok {
		t.Fatalf("expected report to be a report command")
	}
	if report.from != from {
		t.Fatalf("expected report from to be %v, got %v", from, report.from)
	}
	if report.to != to {
		t.Fatalf("expected report to to be %v, got %v", to, report.to)
	}
}

func TestParseCmdReportWithNDay(t *testing.T) {
	now := time.Now()
	from := time.Date(now.Year(), now.Month(), now.Day()-2, 0, 0, 0, 0, time.Local)
	to := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.Local)
	args := []string{"report", "3", "day"}

	cmd, err := getCommand(args)

	if err != nil {
		t.Fatalf("expected err to be nil, got '%v'", err)
	}
	report, ok := cmd.(cmdReport)
	if !ok {
		t.Fatalf("expected report to be a report command")
	}
	if report.from != from {
		t.Fatalf("expected report from to be %v, got %v", from, report.from)
	}
	if report.to != to {
		t.Fatalf("expected report to to be %v, got %v", to, report.to)
	}
}

func TestParseCmdReportWeek(t *testing.T) {
	now := time.Now()
	from := time.Date(now.Year(), now.Month(), now.Day()-int(now.Weekday())+1, 0, 0, 0, 0, time.Local)
	to := time.Date(now.Year(), now.Month(), now.Day()-int(now.Weekday())+7, 23, 59, 59, 0, time.Local)
	args := []string{"report", "week"}

	cmd, err := getCommand(args)

	if err != nil {
		t.Fatalf("expected err to be nil, got '%v'", err)
	}
	report, ok := cmd.(cmdReport)
	if !ok {
		t.Fatalf("expected report to be a report command")
	}
	if report.from != from {
		t.Fatalf("expected report from to be %v, got %v", from, report.from)
	}
	if report.to != to {
		t.Fatalf("expected report to to be %v, got %v", to, report.to)
	}
}

func TestParseCmdReportWithNWeek(t *testing.T) {
	now := time.Now()
	from := time.Date(now.Year(), now.Month(), now.Day()-int(now.Weekday())+1-7*4, 0, 0, 0, 0, time.Local)
	to := time.Date(now.Year(), now.Month(), now.Day()-int(now.Weekday())+7, 23, 59, 59, 0, time.Local)
	args := []string{"report", "5", "week"}

	cmd, err := getCommand(args)

	if err != nil {
		t.Fatalf("expected err to be nil, got '%v'", err)
	}
	report, ok := cmd.(cmdReport)
	if !ok {
		t.Fatalf("expected report to be a report command")
	}
	if report.from != from {
		t.Fatalf("expected report from to be %v, got %v", from, report.from)
	}
	if report.to != to {
		t.Fatalf("expected report to to be %v, got %v", to, report.to)
	}
}

func TestParseCmdReportInvalid(t *testing.T) {
	args := []string{"report", "invalid"}

	_, err := getCommand(args)

	if err == nil {
		t.Fatalf("expected err to be not nil, got nil")
	}
}
