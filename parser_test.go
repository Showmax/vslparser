package vslparser

import (
	"bufio"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"
)

type kv struct {
	k string
	v string
}

// TestSplitLine tests that splitLine correctly handles white-space and missing
// keys/values in the input line and that white-space is preserved at the tail
// of the value.
func TestSplitLine(t *testing.T) {
	samples := map[string]kv{
		"":               kv{},
		"    foo    ":    kv{k: "foo"},
		"foo bar":        kv{k: "foo", v: "bar"},
		"  foo  bar":     kv{k: "foo", v: "bar"},
		" foo    bar   ": kv{k: "foo", v: "bar   "},
		"				 foo	bar	 ": kv{k: "foo", v: "bar	 "},
	}
	for line, kv := range samples {
		gk, gv := splitLine(line)
		if gk != kv.k {
			t.Errorf("parsing %q should give %q as key, got %q", line, kv.k, gk)
		}
		if gv != kv.v {
			t.Errorf("parsing %q should give %q as value, got %q", line, kv.v, gv)
		}
	}
}

// stringScanner is a helper which provides suitable bufio.Scanner object that
// can be passed to Parse so that the given string is processed by the parser.
func stringScanner(s string) *bufio.Scanner {
	return bufio.NewScanner(strings.NewReader(s))
}

// testParseOK tests that s is parsed as e without errors.
func testParseOK(t *testing.T, e *Entry, s string) {
	got, err := Parse(stringScanner(s))
	if err != nil {
		t.Errorf("failed to parse %q: %v", s, err)
		return
	}
	if !reflect.DeepEqual(e, got) {
		t.Errorf("parsing %q should give %v, got %v", s, e, got)
	}
}

// testParseMultipleOK tests that s is parsed as a chain of entries from e
// without errors.
func testParseMultipleOK(t *testing.T, e []Entry, s string) {
	scanner := stringScanner(s)
	for _, ent := range e {
		got, err := Parse(scanner)
		if err != nil {
			t.Errorf("failed to parse %q: %v", s, err)
			return
		}
		if !reflect.DeepEqual(ent, got) {
			t.Errorf("parsing %q should give %v, got %v", s, ent, got)
		}
	}
}

// testParseError tests that parsing s produces an error.
func testParseError(t *testing.T, s string) {
	_, err := Parse(stringScanner(s))
	if err == nil {
		t.Errorf("parsing %q should be a parse error", s)
	} else {
		t.Logf("parsing %q gives: %v", s, err)
	}
}

// TestParse tests that various inputs are either parsed correctly or produce
// errors.
func TestParse(t *testing.T) {
	testParseOK(t, &Entry{
		Level: 1,
		Kind:  "BeReq",
		VXID:  123,
	}, "* << BeReq >> 123\n- End")

	testParseOK(t, &Entry{
		Level: 1,
		Kind:  "Request",
		VXID:  40000000,
		Tags: []Tag{
			{"Foo", "Bar"},
			{"Foo", "Baz"},
			{"Bar", "Foo  Bar    Baz	"}, // Trailing tab valid.
		},
	}, "*   <<  Request >> 40000000\n- Foo Bar\n-Foo Baz\n- Bar     Foo  Bar    Baz	\n- End")
	testParseMultipleOK(t, []Entry{
		{
			Level: 1,
			Kind:  "BeReq",
			VXID:  123,
		},
		{
			Level: 1,
			Kind:  "BeReq",
			VXID:  124,
		},
	}, "* << BeReq >> 123\n- End\n\n* << BeReq >> 124\n- End")

	testParseError(t, "")
	testParseError(t, "- ")
	testParseError(t, "* << Request >> 1\n - Foo Bar\n- End")
	testParseError(t, "* << Request >> Foo")
	testParseError(t, "* << Request >> 1")
}

func TestEOF(t *testing.T) {
	f, _ := os.Open(os.DevNull)
	defer f.Close()

	scanner := bufio.NewScanner(f)
	_, err := Parse(scanner)
	if err != io.EOF {
		t.Errorf("parsing should result in an EOF error. Got error: '%s'", err)
	} else {
		t.Logf("parsing properly returned EOF")
	}
}

const entryExample = `
*   << Request  >> 29236596
-   Begin          req 29236595 rxreq
-   Timestamp      Start: 1545037998.267746 9.124000 18.152000
-   Timestamp      Bad1: 1545037998.267746 foo 37.1248520
-   Timestamp      Bad2: 1545037998.267746 22.111
-   ReqStart       127.0.0.1 44876
-   ReqMethod      GET
-   ReqURL         /health
-   ReqProtocol    HTTP/1.0
-   ReqHeader      X-Forwarded-For: 192.168.1.1
-   VCL_call       RECV
-   VCL_return     synth
-   VCL_call       HASH
-   VCL_return     lookup
-   Timestamp      Process: 1545037998.267784 0.000038 0.000038
-   RespHeader     Date: Mon, 17 Dec 2018 09:13:18 GMT
-   RespHeader     Server: Varnish
-   RespHeader     X-Varnish: 29236596
-   RespHeader     GoWithout:Spaces
-   Empty
-   EmptyTwice
-   RespProtocol   HTTP/1.1
-   SomeFloat      0.1
-   RespStatus     200
-   RespReason     OK
-   RespReason     OK
-   VCL_call       SYNTH
-   RespHeader     Access-Control-Allow-Origin: *
-   RespHeader     Content-Type: application/json; charset=utf-8
-   EmptyTwice
-   VCL_return     deliver
-   RespHeader     Content-Length: 2
-   Storage        malloc Transient
-   RespHeader     Accept-Ranges: bytes
-   Debug          "RES_MODE 2"
-   RespHeader     Connection: close
-   Timestamp      Resp: 1545037998.267831 0.000085 0.000047
-   ReqAcct        24 0 24 233 2 235
-   Foo            Bar Not a named field because there's no ':' after 'Key'
-   End`

func BenchmarkParse(b *testing.B) {
	scanners := make([]*bufio.Scanner, b.N)
	for i := range scanners {
		r := strings.NewReader(entryExample)
		s := bufio.NewScanner(r)
		scanners[i] = s
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := scanners[i]
		if _, err := Parse(s); err != nil {
			b.Fatal(err)
		}
	}
}
