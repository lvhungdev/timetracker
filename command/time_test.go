package command

import (
	"testing"
	"time"
)

var now = time.Date(
	time.Now().Year(),
	time.Now().Month(),
	time.Now().Day(),
	time.Now().Hour(),
	time.Now().Minute(),
	time.Now().Second(),
	0,
	time.Local,
)

func TestParseTimeWithDateTime(t *testing.T) {
	arg := now.Format("2006-01-02T15:04:05")

	parsed, err := parseTime(arg)
	if err != nil {
		t.Fatalf("expect err nil, got %v", err)
	}

	if parsed != now {
		t.Fatalf("expect %v, got %v", now, parsed)
	}
}

func TestParseTimeWithTime(t *testing.T) {
	arg := now.Format("15:04:05")

	parsed, err := parseTime(arg)
	if err != nil {
		t.Fatalf("expect err nil, got %v", err)
	}

	if parsed != now {
		t.Fatalf("expect %v, got %v", now, parsed)
	}
}

func TestParseTimeWithRelativeSecond(t *testing.T) {
	arg := "10s"

	parsed, err := parseTime(arg)
	if err != nil {
		t.Fatalf("expect err nil, got %v", err)
	}

	parsed = parsed.Truncate(time.Second)
	if parsed != now.Add(-time.Second*10) {
		t.Fatalf("expect %v, got %v", now, parsed)
	}
}

func TestParseTimeWithRelativeMinute(t *testing.T) {
	arg := "2m"

	parsed, err := parseTime(arg)
	if err != nil {
		t.Fatalf("expect err nil, got %v", err)
	}

	parsed = parsed.Truncate(time.Second)
	if parsed != now.Add(-time.Minute*2) {
		t.Fatalf("expect %v, got %v", now, parsed)
	}
}

func TestParseTimeWithRelativeHour(t *testing.T) {
	arg := "50h"

	parsed, err := parseTime(arg)
	if err != nil {
		t.Fatalf("expect err nil, got %v", err)
	}

	parsed = parsed.Truncate(time.Second)
	if parsed != now.Add(-time.Hour*50) {
		t.Fatalf("expect %v, got %v", now, parsed)
	}
}

func TestParseTimeWithRelativeDay(t *testing.T) {
	arg := "6d"

	parsed, err := parseTime(arg)
	if err != nil {
		t.Fatalf("expect err nil, got %v", err)
	}

	parsed = parsed.Truncate(time.Second)
	if parsed != now.Add(-time.Hour*24*6) {
		t.Fatalf("expect %v, got %v", now, parsed)
	}
}

func TestParseTimeWithRelativeWeek(t *testing.T) {
	arg := "3w"

	parsed, err := parseTime(arg)
	if err != nil {
		t.Fatalf("expect err nil, got %v", err)
	}

	parsed = parsed.Truncate(time.Second)
	if parsed != now.Add(-time.Hour*24*7*3) {
		t.Fatalf("expect %v, got %v", now, parsed)
	}
}

func TestParseTimeWithInvalidInput(t *testing.T) {
	_, err := parseTime("")
	if err == nil {
		t.Fatalf("expect err not nil, got nil")
	}
}

func TestParseDurationToday(t *testing.T) {
	expectedFrom := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	expectedTo := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.Local)
	arg := "today"

	from, to, err := parseDuration(arg)
	if err != nil {
		t.Fatalf("expect err nil, got %v", err)
	}

	if expectedFrom != from {
		t.Fatalf("expect from to be %v, got %v", expectedFrom, from)
	}
	if expectedTo != to {
		t.Fatalf("expect to to be %v, got %v", expectedTo, to)
	}
}

func TestParseDurationYesterday(t *testing.T) {
	expectedFrom := time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, time.Local)
	expectedTo := time.Date(now.Year(), now.Month(), now.Day()-1, 23, 59, 59, 0, time.Local)
	arg := "yesterday"

	from, to, err := parseDuration(arg)
	if err != nil {
		t.Fatalf("expect err nil, got %v", err)
	}

	if expectedFrom != from {
		t.Fatalf("expect from to be %v, got %v", expectedFrom, from)
	}
	if expectedTo != to {
		t.Fatalf("expect to to be %v, got %v", expectedTo, to)
	}
}

func TestParseDurationDay(t *testing.T) {
	expectedFrom := time.Date(now.Year(), now.Month(), now.Day()-2, 0, 0, 0, 0, time.Local)
	expectedTo := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.Local)
	arg := "3d"

	from, to, err := parseDuration(arg)
	if err != nil {
		t.Fatalf("expect err nil, got %v", err)
	}

	if expectedFrom != from {
		t.Fatalf("expect from to be %v, got %v", expectedFrom, from)
	}
	if expectedTo != to {
		t.Fatalf("expect to to be %v, got %v", expectedTo, to)
	}
}

func TestParseDurationWeek(t *testing.T) {
	monday := now.AddDate(0, 0, -int(now.Weekday())+1-7*4)
	sunday := now.AddDate(0, 0, -int(now.Weekday())+7)
	expectedFrom := time.Date(monday.Year(), monday.Month(), monday.Day(), 0, 0, 0, 0, time.Local)
	expectedTo := time.Date(sunday.Year(), sunday.Month(), sunday.Day(), 23, 59, 59, 0, time.Local)
	arg := "5w"

	from, to, err := parseDuration(arg)
	if err != nil {
		t.Fatalf("expect err nil, got %v", err)
	}

	if expectedFrom != from {
		t.Fatalf("expect from to be %v, got %v", expectedFrom, from)
	}
	if expectedTo != to {
		t.Fatalf("expect to to be %v, got %v", expectedTo, to)
	}
}

func TestParseDurationInvalidInput(t *testing.T) {
	arg := "invalid"

	_, err := parseTime(arg)

	if err == nil {
		t.Fatalf("expect err not nil, got nil")
	}
}
