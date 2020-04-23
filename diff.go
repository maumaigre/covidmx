package main

import (
	"bytes"
	"fmt"

	"github.com/sergi/go-diff/diffmatchpatch"
)

// Operation defines the operation of a diff item.
type Operation int8

// DiffCSV converts a []Diff into a colored text report.
func DiffCSV(diffs []diffmatchpatch.Diff) string {
	var buff bytes.Buffer
	for _, diff := range diffs {
		text := diff.Text

		switch diff.Type {
		case 1:
			_, _ = buff.WriteString(fmt.Sprintf("%s\n", text))

		}
	}

	return buff.String()
}
