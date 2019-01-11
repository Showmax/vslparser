package vslparser

import (
	"github.com/pkg/errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Timestamp represents a parsed time-stamp.
type Timestamp struct {
	AbsTime     time.Time // Absolute time component.
	UsSinceUnit int       // Number of microseconds since the start of work unit.
	UsSincePrev int       // Number of microseconds since previous stamp.
}

// Fields are a collection of log fields. A single log field consists of a key
// and a list of values which appeared with this key. Consider this fragment of
// varnishlog output:
//
// - Foo	Bar Baz
// - Foo	Foobar
// - Bar
//
// There are two fields, "Foo" with values "Bar Baz" and "Foobar" and the field
// "Bar" with an empty string as a sole value.
type Fields map[string][]string

// Entry holds a single log entry. An entry consists mostly of a collection
// of log fields.
type Entry struct {
	Kind   string
	VXID   int
	Fields Fields
}

// newEntry returns a new empty log entry.
func newEntry() *Entry {
	return &Entry{
		Fields: Fields{},
	}
}

// parseUs returns the number of microseconds encoded in the given string.
func parseUs(s string) (int, error) {
	seconds, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, errors.Wrap(err, "cannot parse float component")
	}
	usecs := int(1000 * 1000 * seconds)
	return usecs, nil
}

// parseAbsTime returns a time.Time object parsed from the given string.
func parseAbsTime(s string) (time.Time, error) {
	us, err := parseUs(s)
	if err != nil {
		return time.Time{}, errors.Wrap(err, "cannot parse absolute time")
	}
	sec := int64(us / 1e6)
	ns := int64(1000 * (us % 1e6))
	return time.Unix(sec, ns), nil
}

// Field returns a fields of the log field with the given key. For example, in
// "- Foo: bar", "Foo" is the key of the field and "bar" is one of the values
// returned.
func (e *Entry) Field(key string) ([]string, error) {
	fs, ok := e.Fields[key]
	if !ok {
		return nil, errors.Errorf("entry has no %q field", key)
	}
	return fs, nil
}

// TryField returns the first value of a log field with the given key if present,
// or an empty string otherwise.
func (e *Entry) TryField(key string) string {
	fs, err := e.Field(key)
	if err != nil {
		return ""
	}
	return fs[0]
}

// IntField returns a field with the given key, converted into an integer if
// possible.
func (e *Entry) IntField(key string) (int, error) {
	fs, err := e.Field(key)
	if err != nil {
		return 0, err
	}
	i, err := strconv.Atoi(fs[0])
	if err != nil {
		return 0, errors.Wrapf(err, "cannot convert field %q to an int", key)
	}
	return i, nil
}

// URLField attempts to parse a URL from the field with the given key.
func (e *Entry) URLField(key string) (*url.URL, error) {
	fs, err := e.Field(key)
	if err != nil {
		return nil, err
	}
	url, err := url.Parse(fs[0])
	if err != nil {
		return nil, errors.Wrapf(err, "cannot parse URL from field %q", key)
	}
	return url, nil
}

// Headers returns the values of the field with the given key parsed as a list
// of HTTP headers if possible. Each value of the field has to be a non-empty
// header name immediately followed by a colon, followed by the optional value
// of the header. Headers can appear multiple times.
func (e *Entry) HeadersField(key string) (http.Header, error) {
	h := http.Header{}
	for _, f := range e.Fields[key] {
		name, val := splitLine(f)
		if len(name) == 0 || name[len(name)-1] != ':' {
			return nil, errors.Errorf("%q in field %q is not an HTTP header", f, key)
		}
		h.Add(name[:len(name)-1], val)
	}
	return h, nil
}

// NamedField returns a structured field with the given key, whose name
// component is the given name, if it exists. Names are compared
// case-insensitive to accommodate for HTTP headers.
//
// A structured field is any field whose value is made of an additional name
// and value components, for example "- BerespHeader: X-Header: foo", where
// "BerespHeader" is the key, "X-Header" is the name and "foo" is the value
// returned.
func (e *Entry) NamedField(key, name string) (string, error) {
	fs, err := e.Field(key)
	if err != nil {
		return "", err
	}
	for _, v := range fs {
		fn, fv := splitLine(v)
		l := len(fn) - 1
		if len(fn) > 0 && strings.EqualFold(fn[:l], name) && fn[l] == ':' {
			return fv, nil
		}
	}
	return "", errors.Errorf("entry has no %q field named %q", key, name)
}

// NamedFieldParts returns the same value as a corresponding NamedFields call,
// except that the string value is split by whitespace using strings.Fields.
func (e *Entry) NamedFieldParts(key, name string) ([]string, error) {
	fv, err := e.NamedField(key, name)
	if err != nil {
		return nil, err
	}
	return strings.Fields(fv), nil
}

// Timestamp parses and returns a Timestamp with the given name, if the log entry
// contains such a timestamp.
func (e *Entry) Timestamp(name string) (*Timestamp, error) {
	stamp, err := e.NamedFieldParts("Timestamp", name)
	if err != nil {
		return nil, errors.Wrapf(err, "entry has no timestamp %q", name)
	}
	if len(stamp) != 3 {
		return nil, errors.Errorf("timestamp %q is malformed", name)
	}
	ts := &Timestamp{}
	if ts.AbsTime, err = parseAbsTime(stamp[0]); err != nil {
		return nil, errors.Wrap(err, "cannot parse absolute time")
	}
	if ts.UsSinceUnit, err = parseUs(stamp[1]); err != nil {
		return nil, errors.Wrap(err, "cannot parse time since work unit start")
	}
	if ts.UsSincePrev, err = parseUs(stamp[2]); err != nil {
		return nil, errors.Wrap(err, "cannot parse time since previous timestamp")
	}
	return ts, nil
}
