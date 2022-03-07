package vslparser

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// Note: Session must end with an empty line '\n\n'.
const sessionExample = `*   << Session  >> 413073608
-   Begin          sess 0 HTTP/1
-   Link           req 413073609 rxreq
-   End
**  << Request  >> 413073609
--  Begin          req 413073608 rxreq
--  ReqURL         /healthz
--  End

`

func TestSessionParser_Parse(t *testing.T) {
	result := []Entry{
		{
			Kind:  "Session",
			VXID:  413073608,
			Level: 1,
			Tags: []Tag{
				{
					Key:   "Begin",
					Value: "sess 0 HTTP/1",
				},
				{
					Key:   "Link",
					Value: "req 413073609 rxreq",
				},
				{Key: "End"},
			},
		}, {
			Kind:  "Request",
			Level: 2,
			VXID:  413073609,
			Tags: []Tag{
				{
					Key:   "Begin",
					Value: "req 413073608 rxreq",
				}, {
					Key:   "ReqURL",
					Value: "/healthz",
				},
				{Key: "End"},
			},
		},
	}

	tests := []struct {
		name        string
		varnishlogs string
		want        [][]Entry
	}{
		{
			name:        "empty_string",
			varnishlogs: "",
		}, {
			name:        "empty_lines",
			varnishlogs: "\n\n\n",
			want:        [][]Entry{nil, nil, nil},
		}, {
			name:        "no_empty_line_at_the_beginning_or_end",
			varnishlogs: sessionExample,
			want:        [][]Entry{result},
		}, {
			name:        "no_empty_line_at_the_beginning",
			varnishlogs: sessionExample + "\n",
			want:        [][]Entry{result, nil},
		}, {
			name:        "no_empty_line_at_the_end",
			varnishlogs: "\n" + sessionExample,
			want:        [][]Entry{nil, result},
		}, {
			name:        "multiple_empty_lines_at_the_beginning",
			varnishlogs: "\n\n\n\n" + sessionExample,
			want:        [][]Entry{nil, nil, nil, nil, result},
		}, {
			name:        "multiple_empty_lines_at_the_end",
			varnishlogs: sessionExample + "\n\n\n\n",
			want:        [][]Entry{result, nil, nil, nil, nil},
		}, {
			name:        "two_groups",
			varnishlogs: sessionExample + sessionExample,
			want:        [][]Entry{result, result},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)
			parser := NewSessionParser(strings.NewReader(tt.varnishlogs))

			for i, expected := range tt.want {
				got, err := parser.Parse()
				r.NoError(err, "index: %d/%d\n", i, len(tt.want))
				r.Equal(expected, got, "index: %d/%d\n", i, len(tt.want))
			}

			_, err := parser.Parse()
			r.Equal(io.EOF, err)
		})
	}
}

func BenchmarkSessionParser_Parse(b *testing.B) {
	readers := make([]io.Reader, b.N)
	for i := range readers {
		readers[i] = strings.NewReader(sessionExample)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parser := NewSessionParser(readers[i])
		if _, err := parser.Parse(); err != nil {
			b.Fatalf("unexpected error: %s", err.Error())
		}
	}
}
