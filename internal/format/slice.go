package format

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/tmknom/tfmod/internal/errlib"
)

type SliceFormatter struct {
	format string
	items  []string
	writer io.Writer
}

func NewSliceFormatter(format string, items []string, writer io.Writer) *SliceFormatter {
	return &SliceFormatter{
		format: format,
		items:  items,
		writer: writer,
	}
}

func (p *SliceFormatter) Print() error {
	switch p.format {
	case TextFormat:
		return p.text()
	case JsonFormat:
		return p.json()
	default:
		return p.text()
	}
}

func (p *SliceFormatter) text() error {
	result := strings.Join(p.items, " ")
	_, err := fmt.Fprintln(p.writer, result)
	return err
}

func (p *SliceFormatter) json() error {
	bytes, err := json.Marshal(p.items)
	if err != nil {
		return errlib.Wrapf(err, "failed to marshal items: %v", p.items)
	}

	_, err = fmt.Fprintln(p.writer, string(bytes))
	return err
}
