package vslparser

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// white returns whether the byte b is considered a whitespace character for
// the purpose of parsing of the log.
func white(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n'
}

// splitLine splits the log line s into a key and value component efficiently
// on white-space boundaries.
func splitLine(s string) (string, string) {
	l := len(s)
	ks := 0
	for ks < l {
		if !white(s[ks]) {
			break
		}
		ks++
	}
	ke := ks
	for ke < l {
		if white(s[ke]) {
			break
		}
		ke++
	}
	vs := ke
	for vs < l {
		if !white(s[vs]) {
			break
		}
		vs++
	}
	return s[ks:ke], s[vs:]
}

// Parse will attempt to produce a single Entry from the log, which it reads
// using the given scanner.
//
// The process is fairly efficient, as the individual fields (e.g. time-stamps)
// aren't converted to any special representation. Instead, the parsed entry
// is kept mostly in its textual form. Only basic processing, such as splitting
// lines into fields with a key and a value, are performed. The Entry struct
// provides various convenience methods which perform the subsequent parsing.
func Parse(scanner *bufio.Scanner) (*Entry, error) {
	if err := skipEmptyLines(scanner); err != nil {
		return nil, err
	}
	return parseEntry(scanner)
}

func parseEntry(scanner *bufio.Scanner) (*Entry, error) {
	var e Entry
	// Parse Parselog entry header, e.g.:
	// *   << BeReq    >> 32086823
	// *   << Request  >> 32742536
	// *   << Session  >> 29236595
	header := strings.Fields(scanner.Text())
	e.Level = len(header[0]) // number of asterisks
	if len(header) != 5 || header[0] != asteriskRepeat(e.Level) {
		return nil, errors.New("header line was expected")
	}
	var err error
	e.Kind = header[2]
	if e.VXID, err = strconv.Atoi(header[4]); err != nil {
		return nil, fmt.Errorf("failed to parse VXID: %w", err)
	}
	// Parse log entries, e.g.:
	// -   ReqStart       136.243.103.218 53602
	// -   ReqURL         /health
	// -   Timestamp      Process: 1545037998.759333 0.000031 0.000031
	foundEnd := false
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			return nil, errors.New("parse error: unexpected empty line")
		}
		dashes := dashRepeat(e.Level)
		if !strings.HasPrefix(line, dashes) {
			return nil, fmt.Errorf("parse error on line %q: does not start with '%s'", line, dashes)
		}
		k, v := splitLine(line[e.Level:])
		if k == "" {
			return nil, fmt.Errorf("parse error on line %q: empty key", line)
		}
		if k == "End" {
			foundEnd = true
			break
		}
		e.Tags = append(e.Tags,
			Tag{Key: k, Value: v},
		)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if !foundEnd {
		return nil, errors.New("unexpected EOF in the middle of a log entry")
	}
	return &e, nil
}

func skipEmptyLines(scanner *bufio.Scanner) error {
	eof := true
	for scanner.Scan() {
		if len(scanner.Bytes()) > 0 {
			eof = false
			break
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	if eof {
		return io.EOF
	}
	return nil
}

// slashRepeat is performance optimization.
func dashRepeat(i int) string {
	switch i {
	case 1:
		return "-"
	case 2:
		return "--"
	case 3:
		return "---"
	case 4:
		return "----"
	default:
		// Unlikely.
		return strings.Repeat("-", i)
	}
}

// slashRepeat is performance optimization.
func asteriskRepeat(i int) string {
	switch i {
	case 1:
		return "*"
	case 2:
		return "**"
	case 3:
		return "***"
	case 4:
		return "****"
	default:
		// Unlikely.
		return strings.Repeat("*", i)
	}
}
