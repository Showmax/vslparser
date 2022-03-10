package vslparser

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
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
		"":               {},
		"    foo    ":    {k: "foo"},
		"foo bar":        {k: "foo", v: "bar"},
		"  foo  bar":     {k: "foo", v: "bar"},
		" foo    bar   ": {k: "foo", v: "bar   "},
		"				 foo	bar	 ": {k: "foo", v: "bar	 "},
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

// testEntryParseOK tests that input is entry-parsed as expected without errors.
func testEntryParseOK(t testing.TB, expected Entry, input string) {
	got, err := NewEntryParser(strings.NewReader(input)).Parse()
	require.NoError(t, err)
	require.Equal(t, expected, got)
}

// testEntryParseMultipleOK tests that input is entry-parsed as a chain of entries
// from expected without errors.
func testEntryParseMultipleOK(t testing.TB, expected []Entry, input string) {
	r := require.New(t)
	parser := NewEntryParser(strings.NewReader(input))

	for _, ent := range expected {
		got, err := parser.Parse()
		r.NoError(err)
		r.Equal(ent, got)
	}
}

// testEntryParseError tests that entry-parsing of input produces an error.
func testEntryParseError(t testing.TB, input string) {
	_, err := NewEntryParser(strings.NewReader(input)).Parse()
	require.Error(t, err)
	t.Logf("parsing %q gives: %s\n", input, err.Error())
}

// TestEntryParser_Parse tests that various entry inputs are either parsed
// correctly or produce errors.
func TestEntryParser_Parse(t *testing.T) {
	testEntryParseOK(t, Entry{
		Level: 1,
		Kind:  "BeReq",
		VXID:  123,
		Tags: []Tag{
			{"End", ""},
		},
	}, "* << BeReq >> 123\n- End")

	testEntryParseOK(t, Entry{
		Level: 1,
		Kind:  "Request",
		VXID:  40000000,
		Tags: []Tag{
			{"Foo", "Bar"},
			{"Foo", "Baz"},
			{"Bar", "Foo  Bar    Baz	"}, // Trailing tab valid.
			{"End", ""},
		},
	}, "*   <<  Request >> 40000000\n- Foo Bar\n-Foo Baz\n- Bar     Foo  Bar    Baz	\n- End")
	testEntryParseMultipleOK(t, []Entry{
		{
			Level: 1,
			Kind:  "BeReq",
			VXID:  123,
			Tags: []Tag{
				{"End", ""},
			},
		},
		{
			Level: 1,
			Kind:  "BeReq",
			VXID:  124,
			Tags: []Tag{
				{"End", ""},
			},
		},
	}, "* << BeReq >> 123\n- End\n\n* << BeReq >> 124\n- End")

	testEntryParseError(t, "")
	testEntryParseError(t, "- ")
	testEntryParseError(t, "* << Request >> 1\n - Foo Bar\n- End")
	testEntryParseError(t, "* << Request >> Foo")
	testEntryParseError(t, "* << Request >> 1")
}

func TestEOF(t *testing.T) {
	f, err := os.Open(os.DevNull)
	require.NoError(t, err)
	defer f.Close()

	_, err = NewEntryParser(f).Parse()
	if err != io.EOF {
		t.Errorf("parsing should result in an EOF error. Got error: '%s'", err)
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

func BenchmarkEntryParser_Parse(b *testing.B) {
	readers := make([]io.Reader, b.N)
	for i := range readers {
		readers[i] = strings.NewReader(entryExample)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parser := NewEntryParser(readers[i])
		if _, err := parser.Parse(); err != nil {
			b.Fatal(err)
		}
	}
}
