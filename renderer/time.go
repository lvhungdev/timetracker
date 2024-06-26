package renderer

import (
	"fmt"
	"time"
)

func timeDiffString(t1 time.Time, t2 time.Time) string {
	fmtStr := ""

	if t1.Year() != t2.Year() {
		fmtStr = "2006-01-02 15:04:05"
	} else if t1.Month() != t2.Month() {
		fmtStr = "     01-02 15:04:05"
	} else if t1.Day() != t2.Day() {
		fmtStr = "        02 15:04:05"
	} else if t1.Hour() != t2.Hour() {
		fmtStr = "           15:04:05"
	} else if t1.Minute() != t2.Minute() {
		fmtStr = "              04:05"
	} else {
		fmtStr = "                 05"
	}

	return t2.Local().Format(fmtStr)
}

func durationDiffString(t1 time.Time, t2 time.Time) string {
	if t1.IsZero() || t2.IsZero() {
		return ""
	}

	var result string

	d := t2.Sub(t1).Abs().Round(time.Second)
	dInSec := int(d.Seconds())

	if dInSec < 60 {
		result = fmt.Sprintf("%02d", dInSec)
	} else if dInSec < 3600 {
		mins := dInSec / 60
		secs := dInSec % 60
		result = fmt.Sprintf("%02d:%02d", mins, secs)
	} else {
		hours := dInSec / 3600
		mins := (dInSec % 3600) / 60
		secs := dInSec % 60
		result = fmt.Sprintf("%02d:%02d:%02d", hours, mins, secs)
	}

	length := len("2006-01-02 15:04:05")
	spaceToFill := length - len(result)
	for i := 0; i < spaceToFill; i++ {
		result = " " + result
	}

	return result
}

func durationString(t1 time.Time, t2 time.Time) string {
	if t1.IsZero() || t2.IsZero() {
		return ""
	}

	dInSec := int(t2.Sub(t1).Abs().Round(time.Second).Seconds())

	if dInSec < 60 {
		return fmt.Sprintf("%02d:%02d:%02d", 0, 0, dInSec)
	} else if dInSec < 3600 {
		mins := dInSec / 60
		secs := dInSec % 60
		return fmt.Sprintf("%02d:%02d:%02d", 0, mins, secs)
	} else {
		hours := dInSec / 3600
		mins := (dInSec % 3600) / 60
		secs := dInSec % 60
		return fmt.Sprintf("%02d:%02d:%02d", hours, mins, secs)
	}
}

func timeDateString(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	return t.Format("2006-01-02")
}

func timeHourString(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	return t.Format("15:04:05")
}
