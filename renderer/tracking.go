package renderer

import (
	"io"
	"time"

	"github.com/lvhungdev/tt/tracker"
)

func RenderRecord(writer io.Writer, record tracker.Record) (int, error) {
	var content string

	if record.End.IsZero() {
		content += "tracking \"" + record.Name + "\"\n"
		content += "  started: " + record.Start.Local().Format("2006-01-02 15:04:05") + "\n"
		content += "  current: " + timeDiffString(record.Start, time.Now()) + "\n"
	} else {
		content += "recorded \"" + record.Name + "\"\n"
		content += "  started: " + record.Start.Local().Format("2006-01-02 15:04:05") + "\n"
		content += "  ended  : " + timeDiffString(record.Start, record.End) + "\n"
		content += "  total  : " + durationString(record.Start, record.End) + "\n"
	}

	return writer.Write([]byte(content))
}
