package renderer

import (
	"fmt"
	"strings"
)

type table struct {
	header  []string
	content [][]string
}

func newTable(header []string, content [][]string) (table, error) {
	rows := append(content, header)
	for _, r := range rows {
		if len(r) != len(header) {
			return table{}, fmt.Errorf("table: len(header) != len(content)")
		}
	}

	return table{
		header:  header,
		content: content,
	}, nil
}

func (t table) String() string {
	header := strings.Join(t.alignedHeader(), "  ")
	divider := strings.Join(t.divider(), "  ")
	content := make([]string, len(t.content))

	aligned := t.alignedContent()
	for i, row := range aligned {
		content[i] = strings.Join(row, "  ")
	}

	return header + "\n" + divider + "\n" + strings.Join(content, "\n") + "\n"
}

func (t table) alignedHeader() []string {
	colWidths := t.colWidths()

	aligned := make([]string, len(t.header))
	for i, col := range t.header {
		aligned[i] = col
		spaceToFill := colWidths[i] - len(col)
		for j := 0; j < spaceToFill; j++ {
			aligned[i] += " "
		}
	}

	return aligned
}

func (t table) alignedContent() [][]string {
	content := [][]string{}
	colWidths := t.colWidths()

	for _, row := range t.content {
		aligned := make([]string, len(row))
		for i, col := range row {
			aligned[i] = col
			spaceToFill := colWidths[i] - len(col)
			for j := 0; j < spaceToFill; j++ {
				aligned[i] += " "
			}
		}

		content = append(content, aligned)
	}

	return content
}

func (t table) divider() []string {
	divider := make([]string, len(t.header))
	colWidths := t.colWidths()

	for i, col := range t.header {
		divider[i] = strings.Repeat("-", len(col))

		spaceToFill := colWidths[i] - len(col)
		for j := 0; j < spaceToFill; j++ {
			divider[i] += " "
		}
	}

	return divider
}

func (t table) colWidths() []int {
	rows := append(t.content, t.header)

	if len(rows) == 0 {
		return []int{}
	}

	result := make([]int, len(rows[0]))

	for _, row := range rows {
		for i, col := range row {
			if len(col) > result[i] {
				result[i] = len(col)
			}
		}
	}

	return result
}
