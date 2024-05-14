package command

import (
	"fmt"
	"slices"
	"strings"
	"time"
)

type Command interface {
	Time() time.Time
}

type Base struct {
	t time.Time
}

func (b Base) Time() time.Time {
	return b.t
}

type GetCurrent struct {
	Base
}

type StartTracking struct {
	Base
	Name string
}

type StopTracking struct {
	Base
}

type Report struct {
	Base
	From time.Time
	To   time.Time
}

func Get(args []string) (Command, error) {
	if len(args) == 0 {
		return GetCurrent{Base{time.Now()}}, nil
	}

	switch args[0] {
	case "start":
		return parseStartTracking(args[1:])
	case "stop":
		return parseStopTracking(args[1:])
	case "report":
		return parseReport(args[1:])
	default:
		return nil, fmt.Errorf("unknown command: %s", args[0])
	}
}

func parseStartTracking(args []string) (Command, error) {
	now := time.Now()
	at, idx := getOption(args, "at")
	if idx != -1 {
		t, err := parseTime(at)
		if err != nil {
			return nil, err
		}

		now = t
	}

	name := strings.Join(
		slices.DeleteFunc(args, func(s string) bool {
			return s == args[idx]
		}),
		" ")

	return StartTracking{Base{now}, name}, nil
}

func parseStopTracking(args []string) (Command, error) {
	now := time.Now()

	at, idx := getOption(args, "at")
	if idx != -1 {
		t, err := parseTime(at)
		if err != nil {
			return nil, err
		}

		now = t
	}

	return StopTracking{Base{now}}, nil
}

func parseReport(args []string) (Command, error) {
	now := time.Now()
	from := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	to := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.Local)

	d, idx := getOption(args, "dur")
	if idx != -1 {
		start, end, err := parseDuration(d)
		if err != nil {
			return nil, err
		}

		from = start
		to = end
	}

	return Report{
		Base: Base{now},
		From: from,
		To:   to,
	}, nil
}

func getOption(args []string, name string) (string, int) {
	for i, arg := range args {
		option, value, found := strings.Cut(arg, ":")
		if !found {
			continue
		}

		if option != name {
			continue
		}

		return value, i
	}

	return "", -1
}
