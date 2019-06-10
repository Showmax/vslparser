package vslparser

import (
	"bufio"
	"reflect"
	"strings"
	"testing"
	"time"
)

func example() *Entry {
	s := `
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
	e, err := Parse(bufio.NewScanner(strings.NewReader(s)))
	if err != nil {
		panic(err)
	}
	return e
}

func TestTimestamp(t *testing.T) {
	e := example()
	stamps := map[string]*Timestamp{
		"Start": &Timestamp{
			AbsTime:     time.Unix(0, 1545037998267746000),
			UsSinceUnit: 9124000,
			UsSincePrev: 18152000,
		},
	}
	badStamps := []string{
		"Missing",
		"Bad1",
		"Bad2",
	}
	for name, stamp := range stamps {
		got, err := e.Timestamp(name)
		if err != nil {
			t.Errorf("parsing timestamp %q should not fail, got: %v", name, err)
			continue
		}
		if !stamp.AbsTime.Equal(got.AbsTime) {
			t.Errorf("timestamp %q should have AbsTime %v, has %v",
				name, stamp.AbsTime, got.AbsTime)
		}
		if stamp.UsSinceUnit != got.UsSinceUnit {
			t.Errorf("timestamp %q should have UsSinceUnit %v, has %v",
				name, stamp.UsSinceUnit, got.UsSinceUnit)
		}
		if stamp.UsSincePrev != got.UsSincePrev {
			t.Errorf("timestamp %q should have UsSincePrev %v, has %v",
				name, stamp.UsSincePrev, got.UsSincePrev)
		}
	}
	for _, name := range badStamps {
		_, err := e.Timestamp(name)
		if err == nil {
			t.Errorf("parsing stamp %q should fail", name)
		} else {
			t.Logf("parsing timestamp %q gives: %v", name, err)
		}
	}
}

func TestNamedField(t *testing.T) {
	e := example()
	respHeaders := map[string]string{
		"Content-Length": "2",
		"Date":           "Mon, 17 Dec 2018 09:13:18 GMT",
		"dAtE":           "Mon, 17 Dec 2018 09:13:18 GMT", // mixed case
		"X-Varnish":      "29236596",
		"GoWithout":      "Spaces", // no space after :
	}
	for name, value := range respHeaders {
		got, err := e.NamedField("RespHeader", name)
		if err != nil {
			t.Errorf("getting RespHeader %q should produce no error, got: %v", name, err)
			continue
		}
		if got != value {
			t.Errorf("header %q expected to be %q, got %q", name, value, got)
		}
	}
	badHeaders := map[string]string{
		"ReqMethod": "",
		"":          "Date",
		"Missing":   "Name",
		"ReqHeader": "X-Forwarded-For:", // ':' not part of name
		"Foo":       "Bar",
		"Begin":     "re",
	}
	for key, name := range badHeaders {
		if v, err := e.NamedField(key, name); err == nil {
			t.Errorf("e.NamedField(%q, %q) should produce an error, got: %v", key, name, v)
		} else {
			t.Logf("e.NamedField(%q, %q) gives: %v", key, name, err)
		}
	}
}

func TestField(t *testing.T) {
	e := example()
	samples := map[string][]string{
		"ReqAcct":  []string{"24 0 24 233 2 235"},
		"Begin":    []string{"req 29236595 rxreq"},
		"VCL_call": []string{"RECV", "HASH", "SYNTH"},
		"RespHeader": []string{
			"Date: Mon, 17 Dec 2018 09:13:18 GMT",
			"Server: Varnish",
			"X-Varnish: 29236596",
			"GoWithout:Spaces",
			"Access-Control-Allow-Origin: *",
			"Content-Type: application/json; charset=utf-8",
			"Content-Length: 2",
			"Accept-Ranges: bytes",
			"Connection: close",
		},
		"Empty":      []string{""},
		"EmptyTwice": []string{"", ""},
	}
	for k, vs := range samples {
		gvs, err := e.Field(k)
		if err != nil {
			t.Errorf("e.Field(%q) should produce no error, got: %v", k, err)
			continue
		}
		if !reflect.DeepEqual(gvs, vs) {
			t.Errorf("e.Field(%q) should return %v, got %v", k, vs, gvs)
		}
	}
	bad := []string{
		"Missing",
	}
	for _, k := range bad {
		if _, err := e.Field(k); err == nil {
			t.Errorf("e.Field(%q) should fail", k)
		} else {
			t.Logf("e.Field(%q) gives: %v", k, err)
		}
	}
}

func TestIntField(t *testing.T) {
	e := example()
	samples := map[string]int{
		"RespStatus": 200,
	}
	for name, i := range samples {
		gi, err := e.IntField("RespStatus")
		if err != nil {
			t.Errorf("getting %q as int should not fail, got: %v", name, err)
			continue
		}
		if gi != i {
			t.Errorf("getting %q should return %d, got %d", name, i, gi)
		}
	}
	bad := []string{
		"VCL_call",
		"RespReason",
		"MissingField",
		"SomeFloat",
	}
	for _, name := range bad {
		if _, err := e.IntField(name); err == nil {
			t.Errorf("getting %q as int should fail", name)
		} else {
			t.Logf("getting %q as int gives: %v", name, err)
		}
	}
}
