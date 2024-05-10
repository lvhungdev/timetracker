package renderer

import (
	"io"
	"strings"
	"time"

	"github.com/lvhungdev/tt/tracker"
)

func RenderRecord(writer io.Writer, record tracker.Record) {
	var content string

	if record.End.IsZero() {
		content += "tracking \"" + record.Name + "\"\n"
		content += "  started: " + record.Start.Local().Format("2006-01-02 15:04:05") + "\n"
		content += "  current: " + timeDiffString(record.Start, time.Now()) + "\n"
		if record.Start.Sub(time.Now()).Abs().Seconds() > 1 {
			content += "  total  : " + durationDiffString(record.Start, time.Now()) + "\n"
		}
	} else {
		content += "recorded \"" + record.Name + "\"\n"
		content += "  started: " + record.Start.Local().Format("2006-01-02 15:04:05") + "\n"
		content += "  ended  : " + timeDiffString(record.Start, record.End) + "\n"
		content += "  total  : " + durationDiffString(record.Start, record.End) + "\n"
	}

	_, err := writer.Write([]byte(content))
	if err != nil {
		// TODO handle error properly
		panic(err)
	}
}

func RenderRecords(writer io.Writer, records []tracker.Record) {
	serializeRecord := func(record tracker.Record) []string {
		date := timeDateString(record.Start)
		day := strings.ToLower(record.Start.Weekday().String())
		start := timeHourString(record.Start)
		end := timeHourString(record.End)
		total := durationString(record.Start, record.End)

		return []string{date, day, start, end, total}
	}

	header := []string{"date", "day", "start", "end", "total"}
	content := [][]string{}
	for _, r := range records {
		content = append(content, serializeRecord(r))
	}

	table, _ := newTable(header, content)
	_, err := writer.Write([]byte(table.String()))
	if err != nil {
		// TODO handle error properly
		panic(err)
	}
}
