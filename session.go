package vslparser

import (
	"bufio"
	"fmt"
	"io"
)

// SessionParser implements varnishlog session grouped log (produced by
// "varnishlog -g session" command) parsing functionality.
type SessionParser struct {
	scanner *bufio.Scanner
}

// NewSessionParser creates a new SessionParser reading & parsing r.
func NewSessionParser(r io.Reader) *SessionParser {
	return &SessionParser{
		scanner: bufio.NewScanner(r),
	}
}

// Parse parses log stream produced by varnishlog with enabled grouping.
// Presence of End tag is required.
//
// Example expected of input:
//      *   << Session  >> 413073608
//      -   Begin          sess 0 HTTP/1
//      -   Link           req 413073609 rxreq
//      -   End
//      **  << Request  >> 413073609
//      --  Begin          req 413073608 rxreq
//      --  ReqURL         /healthz
//      --  End
func (p *SessionParser) Parse() ([]Entry, error) {
	var entries []Entry
	for i := 0; p.scanner.Scan(); i++ {
		// Empty line '\n\n' is session log group delimiter.
		if len(p.scanner.Bytes()) == 0 {
			return entries, nil
		}

		e, err := parseEntry(p.scanner)
		if err != nil {
			return nil, fmt.Errorf("cannot parse entry %d: %w", i, err)
		}

		entries = append(entries, e)
	}

	if err := p.scanner.Err(); err != nil {
		return nil, fmt.Errorf("group scanning failed: %w", err)
	}

	// We have reached EOF in the scanner internal reader. We might still
	// have some entries already parsed, but as we haven't seen an empty
	// line delimiter, we are almost certain that entries doesn't represent
	// a complete session log. So we drop it and return EOF as we will not
	// see more full session logs.
	return nil, io.EOF
}
