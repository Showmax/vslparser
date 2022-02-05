package vslparser

import (
	"bufio"
	"io"
)

// ParseGroup parses log stream produced by varnishlog with enabled
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
func ParseGroup(scanner *bufio.Scanner) ([]Entry, error) {
	if err := skipEmptyLines(scanner); err != nil {
		return nil, err
	}

	var ee []Entry
	// do { } while scanner.Scan()
	for ok := true; ok; ok = scanner.Scan() {
		// Groups are separated with an empty line (\n\n).
		if len(scanner.Bytes()) == 0 {
			break
		}
		e, err := parseEntry(scanner)
		if err != nil {
			return nil, err
		}
		ee = append(ee, *e)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if len(ee) == 0 {
		return nil, io.EOF
	}
	return ee, nil
}
