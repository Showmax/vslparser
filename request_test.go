package vslparser

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRequestParser_ParseRealExample(t *testing.T) {
	r := require.New(t)

	file, err := os.Open("testdata/varnishlog_request.txt")
	r.NoError(err)
	defer file.Close()
	parser := NewRequestParser(file)

	expected := []Entry{
		{
			Level: 1,
			Kind:  KindRequest,
			VXID:  2,
			Tags: []Tag{
				{"Begin", "req 1 rxreq"},
				{"Timestamp", "Start: 1646693481.899847 0.000000 0.000000"},
				{"Timestamp", "Req: 1646693481.899847 0.000000 0.000000"},
				{"VCL_use", "boot"},
				{"ReqStart", "127.0.0.1 37976 a0"},
				{"ReqMethod", "GET"},
				{"ReqURL", "/"},
				{"ReqProtocol", "HTTP/1.1"},
				{"ReqHeader", "Host: localhost:6081"},
				{"ReqHeader", "User-Agent: curl/7.82.0"},
				{"ReqHeader", "Accept: */*"},
				{"ReqHeader", "X-Forwarded-For: 127.0.0.1"},
				{"VCL_call", "RECV"},
				{"VCL_return", "hash"},
				{"VCL_call", "HASH"},
				{"VCL_return", "lookup"},
				{"VCL_call", "MISS"},
				{"VCL_return", "fetch"},
				{"Link", "bereq 3 fetch"},
				{"Timestamp", "Fetch: 1646693481.900397 0.000550 0.000550"},
				{"RespProtocol", "HTTP/1.1"},
				{"RespStatus", "503"},
				{"RespReason", "Backend fetch failed"},
				{"RespHeader", "Date: Mon, 07 Mar 2022 22:51:21 GMT"},
				{"RespHeader", "Server: Varnish"},
				{"RespHeader", "Content-Type: text/html; charset=utf-8"},
				{"RespHeader", "Retry-After: 5"},
				{"RespHeader", "X-Varnish: 2"},
				{"RespHeader", "Age: 0"},
				{"RespHeader", "Via: 1.1 varnish (Varnish/7.0)"},
				{"VCL_call", "DELIVER"},
				{"VCL_return", "deliver"},
				{"Timestamp", "Process: 1646693481.900439 0.000591 0.000041"},
				{"Filters", ""},
				{"RespHeader", "Content-Length: 278"},
				{"RespHeader", "Connection: keep-alive"},
				{"Timestamp", "Resp: 1646693481.900513 0.000665 0.000074"},
				{"ReqAcct", "78 0 78 246 278 524"},
				{"End", ""},
			},
		},
		{
			Level: 2,
			Kind:  KindBeReq,
			VXID:  3,
			Tags: []Tag{
				{"Begin", "bereq 2 fetch"},
				{"VCL_use", "boot"},
				{"Timestamp", "Start: 1646693481.900025 0.000000 0.000000"},
				{"BereqMethod", "GET"},
				{"BereqURL", "/"},
				{"BereqProtocol", "HTTP/1.1"},
				{"BereqHeader", "Host: localhost:6081"},
				{"BereqHeader", "User-Agent: curl/7.82.0"},
				{"BereqHeader", "Accept: */*"},
				{"BereqHeader", "X-Forwarded-For: 127.0.0.1"},
				{"BereqHeader", "Accept-Encoding: gzip"},
				{"BereqHeader", "X-Varnish: 3"},
				{"VCL_call", "BACKEND_FETCH"},
				{"VCL_return", "fetch"},
				{"Timestamp", "Fetch: 1646693481.900075 0.000050 0.000050"},
				{"FetchError", "backend default: fail errno 111 (Connection refused)"},
				{"Timestamp", "Beresp: 1646693481.900243 0.000217 0.000167"},
				{"Timestamp", "Error: 1646693481.900247 0.000221 0.000004"},
				{"BerespProtocol", "HTTP/1.1"},
				{"BerespStatus", "503"},
				{"BerespReason", "Backend fetch failed"},
				{"BerespHeader", "Date: Mon, 07 Mar 2022 22:51:21 GMT"},
				{"BerespHeader", "Server: Varnish"},
				{"VCL_call", "BACKEND_ERROR"},
				{"BerespHeader", "Content-Type: text/html; charset=utf-8"},
				{"BerespHeader", "Retry-After: 5"},
				{"VCL_return", "deliver"},
				{"Storage", "malloc Transient"},
				{"Length", "278"},
				{"BereqAcct", "0 0 0 0 0 0"},
				{"End", ""},
			},
		},
	}
	entries, err := parser.Parse()
	r.NoError(err)
	r.Equal(expected, entries)

	expected = []Entry{
		{
			Level: 1,
			Kind:  KindRequest,
			VXID:  5,
			Tags: []Tag{
				{"Begin", "req 4 rxreq"},
				{"Timestamp", "Start: 1646693489.711023 0.000000 0.000000"},
				{"Timestamp", "Req: 1646693489.711023 0.000000 0.000000"},
				{"VCL_use", "boot"},
				{"ReqStart", "127.0.0.1 37978 a0"},
				{"ReqMethod", "POST"},
				{"ReqURL", "/post"},
				{"ReqProtocol", "HTTP/1.1"},
				{"ReqHeader", "Host: localhost:6081"},
				{"ReqHeader", "User-Agent: curl/7.82.0"},
				{"ReqHeader", "Accept: */*"},
				{"ReqHeader", "X-Forwarded-For: 127.0.0.1"},
				{"VCL_call", "RECV"},
				{"VCL_return", "pass"},
				{"VCL_call", "HASH"},
				{"VCL_return", "lookup"},
				{"VCL_call", "PASS"},
				{"VCL_return", "fetch"},
				{"Link", "bereq 6 pass"},
				{"Timestamp", "Fetch: 1646693489.711278 0.000254 0.000254"},
				{"RespProtocol", "HTTP/1.1"},
				{"RespStatus", "503"},
				{"RespReason", "Backend fetch failed"},
				{"RespHeader", "Date: Mon, 07 Mar 2022 22:51:29 GMT"},
				{"RespHeader", "Server: Varnish"},
				{"RespHeader", "Content-Type: text/html; charset=utf-8"},
				{"RespHeader", "Retry-After: 5"},
				{"RespHeader", "X-Varnish: 5"},
				{"RespHeader", "Age: 0"},
				{"RespHeader", "Via: 1.1 varnish (Varnish/7.0)"},
				{"VCL_call", "DELIVER"},
				{"VCL_return", "deliver"},
				{"Timestamp", "Process: 1646693489.711294 0.000270 0.000015"},
				{"Filters", ""},
				{"RespHeader", "Content-Length: 278"},
				{"RespHeader", "Connection: keep-alive"},
				{"Timestamp", "Resp: 1646693489.711344 0.000320 0.000050"},
				{"ReqAcct", "83 0 83 246 278 524"},
				{"End", ""},
			},
		},
		{
			Level: 2,
			Kind:  KindBeReq,
			VXID:  6,
			Tags: []Tag{
				{"Begin", "bereq 5 pass"},
				{"VCL_use", "boot"},
				{"Timestamp", "Start: 1646693489.711099 0.000000 0.000000"},
				{"BereqMethod", "POST"},
				{"BereqURL", "/post"},
				{"BereqProtocol", "HTTP/1.1"},
				{"BereqHeader", "Host: localhost:6081"},
				{"BereqHeader", "User-Agent: curl/7.82.0"},
				{"BereqHeader", "Accept: */*"},
				{"BereqHeader", "X-Forwarded-For: 127.0.0.1"},
				{"BereqHeader", "X-Varnish: 6"},
				{"VCL_call", "BACKEND_FETCH"},
				{"VCL_return", "fetch"},
				{"Timestamp", "Fetch: 1646693489.711118 0.000019 0.000019"},
				{"FetchError", "backend default: fail errno 111 (Connection refused)"},
				{"Timestamp", "Beresp: 1646693489.711199 0.000100 0.000081"},
				{"Timestamp", "Error: 1646693489.711202 0.000103 0.000002"},
				{"BerespProtocol", "HTTP/1.1"},
				{"BerespStatus", "503"},
				{"BerespReason", "Backend fetch failed"},
				{"BerespHeader", "Date: Mon, 07 Mar 2022 22:51:29 GMT"},
				{"BerespHeader", "Server: Varnish"},
				{"VCL_call", "BACKEND_ERROR"},
				{"BerespHeader", "Content-Type: text/html; charset=utf-8"},
				{"BerespHeader", "Retry-After: 5"},
				{"VCL_return", "deliver"},
				{"Storage", "malloc Transient"},
				{"Length", "278"},
				{"BereqAcct", "0 0 0 0 0 0"},
				{"End", ""},
			},
		},
	}
	entries, err = parser.Parse()
	r.NoError(err)
	r.Equal(expected, entries)

	expected = []Entry{
		{
			Level: 1,
			Kind:  KindRequest,
			VXID:  32770,
			Tags: []Tag{
				{"Begin", "req 32769 rxreq"},
				{"Timestamp", "Start: 1646693544.293284 0.000000 0.000000"},
				{"Timestamp", "Req: 1646693544.293284 0.000000 0.000000"},
				{"VCL_use", "boot"},
				{"ReqStart", "127.0.0.1 37980 a0"},
				{"ReqMethod", "PUT"},
				{"ReqURL", "/foo?param=val"},
				{"ReqProtocol", "HTTP/1.1"},
				{"ReqHeader", "Host: localhost:6081"},
				{"ReqHeader", "User-Agent: curl/7.82.0"},
				{"ReqHeader", "Accept: */*"},
				{"ReqHeader", "magic: aloha"},
				{"ReqHeader", "greeting: traveler"},
				{"ReqHeader", "X-Forwarded-For: 127.0.0.1"},
				{"VCL_call", "RECV"},
				{"VCL_return", "pass"},
				{"VCL_call", "HASH"},
				{"VCL_return", "lookup"},
				{"VCL_call", "PASS"},
				{"VCL_return", "fetch"},
				{"Link", "bereq 32771 pass"},
				{"Timestamp", "Fetch: 1646693544.294199 0.000915 0.000915"},
				{"RespProtocol", "HTTP/1.1"},
				{"RespStatus", "503"},
				{"RespReason", "Backend fetch failed"},
				{"RespHeader", "Date: Mon, 07 Mar 2022 22:52:24 GMT"},
				{"RespHeader", "Server: Varnish"},
				{"RespHeader", "Content-Type: text/html; charset=utf-8"},
				{"RespHeader", "Retry-After: 5"},
				{"RespHeader", "X-Varnish: 32770"},
				{"RespHeader", "Age: 0"},
				{"RespHeader", "Via: 1.1 varnish (Varnish/7.0)"},
				{"VCL_call", "DELIVER"},
				{"VCL_return", "deliver"},
				{"Timestamp", "Process: 1646693544.294221 0.000937 0.000021"},
				{"Filters", ""},
				{"RespHeader", "Content-Length: 282"},
				{"RespHeader", "Connection: keep-alive"},
				{"Timestamp", "Resp: 1646693544.294304 0.001020 0.000083"},
				{"ReqAcct", "125 0 125 250 282 532"},
				{"End", ""},
			},
		},
		{
			Level: 2,
			Kind:  KindBeReq,
			VXID:  32771,
			Tags: []Tag{
				{"Begin", "bereq 32770 pass"},
				{"VCL_use", "boot"},
				{"Timestamp", "Start: 1646693544.293767 0.000000 0.000000"},
				{"BereqMethod", "PUT"},
				{"BereqURL", "/foo?param=val"},
				{"BereqProtocol", "HTTP/1.1"},
				{"BereqHeader", "Host: localhost:6081"},
				{"BereqHeader", "User-Agent: curl/7.82.0"},
				{"BereqHeader", "Accept: */*"},
				{"BereqHeader", "magic: aloha"},
				{"BereqHeader", "greeting: traveler"},
				{"BereqHeader", "X-Forwarded-For: 127.0.0.1"},
				{"BereqHeader", "X-Varnish: 32771"},
				{"VCL_call", "BACKEND_FETCH"},
				{"VCL_return", "fetch"},
				{"Timestamp", "Fetch: 1646693544.293823 0.000055 0.000055"},
				{"FetchError", "backend default: fail errno 111 (Connection refused)"},
				{"Timestamp", "Beresp: 1646693544.294076 0.000308 0.000253"},
				{"Timestamp", "Error: 1646693544.294080 0.000313 0.000004"},
				{"BerespProtocol", "HTTP/1.1"},
				{"BerespStatus", "503"},
				{"BerespReason", "Backend fetch failed"},
				{"BerespHeader", "Date: Mon, 07 Mar 2022 22:52:24 GMT"},
				{"BerespHeader", "Server: Varnish"},
				{"VCL_call", "BACKEND_ERROR"},
				{"BerespHeader", "Content-Type: text/html; charset=utf-8"},
				{"BerespHeader", "Retry-After: 5"},
				{"VCL_return", "deliver"},
				{"Storage", "malloc Transient"},
				{"Length", "282"},
				{"BereqAcct", "0 0 0 0 0 0"},
				{"End", ""},
			},
		},
	}
	entries, err = parser.Parse()
	r.NoError(err)
	r.Equal(expected, entries)

	entries, err = parser.Parse()
	r.Equal(io.EOF, err)
	r.Nil(entries)
}

func BenchmarkRequestParser_Parse(b *testing.B) {
	r := require.New(b)

	file, err := os.Open("testdata/varnishlog_request.txt")
	r.NoError(err)

	fileBytes, err := io.ReadAll(file)
	r.NoError(err)
	r.NoError(file.Close())

	parsers := make([]*RequestParser, b.N)
	for i := range parsers {
		parsers[i] = NewRequestParser(bytes.NewReader(fileBytes))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = parsers[i].Parse()
	}
}
