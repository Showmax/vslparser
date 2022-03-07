package vslparser

import (
	"bufio"
	"fmt"
	"io"
)

// GroupParsing implements varnishlog session (produced by "varnishlog -g
// session" command) parsing functionality.
type SessionParser struct {
	scanner *bufio.Scanner
}

// NewSessionParser creates a new SessionParser reading & parsing r.
func NewSessionParser(r io.Reader) *SessionParser {
	return &SessionParser{
		scanner: bufio.NewScanner(r),
	}
}

// Parse parses log stream produced by varnishlog with enabled
// grouping. Presence of End tag is required.
//
// Example expected of input:
//	*   << Session  >> 413073608
//	-   Begin          sess 0 HTTP/1
//	-   Link           req 413073609 rxreq
//	-   End
//	**  << Request  >> 413073609
//	--  Begin          req 413073608 rxreq
//	--  ReqURL         /healthz
//	--  End
func (p *SessionParser) Parse() ([]Entry, error) {
	var entries []Entry

	for i := 0; p.scanner.Scan(); i++ {
		// Groups are separated with an empty line (\n\n).
		if len(p.scanner.Bytes()) == 0 {
			break
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
	if len(entries) == 0 {
		return nil, io.EOF
	}
	return entries, nil
}
