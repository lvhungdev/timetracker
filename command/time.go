package command

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func parseTime(s string) (time.Time, error) {
	now := time.Now()

	if len(s) <= 1 {
		return time.Time{}, fmt.Errorf("invalid time format: %s", s)
	}

	t, err := time.ParseInLocation("2006-01-02T15:04:05", s, time.Local)
	if err == nil {
		return t, nil
	}

	t, err = time.ParseInLocation("15:04:05", s, time.Local)
	if err == nil {
		return time.Date(now.Year(), now.Month(), now.Day(), t.Hour(), t.Minute(), t.Second(), 0, time.Local), nil
	}

	n, err := strconv.Atoi(s[:len(s)-1])
	d := time.Duration(n)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid time format: %s", s)
	}

	switch s[len(s)-1] {
	case 's':
		return now.Add(-time.Second * d), nil
	case 'm':
		return now.Add(-time.Minute * d), nil
	case 'h':
		return now.Add(-time.Hour * d), nil
	case 'd':
		return now.Add(-time.Hour * 24 * d), nil
	case 'w':
		return now.Add(-time.Hour * 24 * 7 * d), nil
	default:
		return time.Time{}, fmt.Errorf("invalid time format: %s", s)
	}
}

func parseDuration(s string) (from, to time.Time, err error) {
	now := time.Now()

	if s == "today" {
		from = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
		to = time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.Local)
		return
	}

	if s == "yesterday" {
		from = time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, time.Local)
		to = time.Date(now.Year(), now.Month(), now.Day()-1, 23, 59, 59, 0, time.Local)
		return
	}

	if strings.HasSuffix(s, "d") {
		amount, e := strconv.Atoi(s[:len(s)-1])
		if e != nil {
			err = fmt.Errorf("invalid time format: %s", s)
			return
		}

		from = time.Date(now.Year(), now.Month(), now.Day()-amount+1, 0, 0, 0, 0, time.Local)
		to = time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.Local)
		return
	}

	if strings.HasSuffix(s, "w") {
		amount, e := strconv.Atoi(s[:len(s)-1])
		if e != nil {
			err = fmt.Errorf("invalid time format: %s", s)
			return
		}

		monday := now.AddDate(0, 0, -int(now.Weekday())+1-7*(amount-1))
		sunday := now.AddDate(0, 0, -int(now.Weekday())+7)
		from = time.Date(monday.Year(), monday.Month(), monday.Day(), 0, 0, 0, 0, time.Local)
		to = time.Date(sunday.Year(), sunday.Month(), sunday.Day(), 23, 59, 59, 0, time.Local)
		return
	}

	err = fmt.Errorf("invalid time format: %s", s)
	return
}
