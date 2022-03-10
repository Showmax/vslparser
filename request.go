package vslparser

import (
	"bufio"
	"fmt"
	"io"
)

// RequestParser implements varnishlog request grouped log (produced by
// "varnishlog -g request" command) parsing functionality.
type RequestParser struct {
	scanner *bufio.Scanner
}

// NewRequestParser creates a new RequestParser reading & parsing r.
func NewRequestParser(r io.Reader) *RequestParser {
	return &RequestParser{
		scanner: bufio.NewScanner(r),
	}
}

func (p *RequestParser) Parse() ([]Entry, error) {
	var entries []Entry
	for i := 0; p.scanner.Scan(); i++ {
		// Empty line '\n\n' is request log group delimiter.
		if len(p.scanner.Bytes()) == 0 {
			return entries, nil
		}

		entry, err := parseEntry(p.scanner)
		if err != nil {
			return nil, fmt.Errorf("cannot parse entry %d: %w", i, err)
		}

		entries = append(entries, entry)
	}

	if err := p.scanner.Err(); err != nil {
		return nil, fmt.Errorf("request scanning failed: %w", err)
	}

	// We have reached EOF in the scanner internal reader. We might still
	// have some entries already parsed, but as we haven't seen an empty
	// line delimiter, we are almost certain that entries doesn't represent
	// a complete request log. So we drop it and return EOF as we will not
	// see more full request logs.
	return nil, io.EOF
}
