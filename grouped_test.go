package vslparser

import (
	"bufio"
	"reflect"
	"strings"
	"testing"
)

const groupExample = `*   << Session  >> 413073608
-   Begin          sess 0 HTTP/1
-   Link           req 413073609 rxreq
-   End
**  << Request  >> 413073609
--  Begin          req 413073608 rxreq
--  ReqURL         /healthz
--  End`

func TestParseGroup(t *testing.T) {
	result := []Entry{
		{
			Kind: "Session",
			VXID: 413073608,
			Fields: Fields{
				"Begin": {"sess 0 HTTP/1"},
				"Link":  {"req 413073609 rxreq"},
			},
		},
		{
			Kind: "Request",
			VXID: 413073609,
			Fields: Fields{
				"Begin":  {"req 413073608 rxreq"},
				"ReqURL": {"/healthz"},
			},
		},
	}

	tests := []struct {
		name        string
		varnishlogs string
		want        []Entry
		wantErr     bool
	}{
		{
			name:        "empty string",
			varnishlogs: "",
			wantErr:     true,
		}, {
			name:        "empty lines",
			varnishlogs: "\n\n\n",
			wantErr:     true,
		}, {
			name:        "no empty line at the beginning or end",
			varnishlogs: groupExample,
			want:        result,
		}, {
			name:        "no empty line at the beginning",
			varnishlogs: groupExample + "\n",
			want:        result,
		}, {
			name:        "no empty line at the end",
			varnishlogs: "\n" + groupExample,
			want:        result,
		}, {
			name:        "multiple empty lines at the beginning",
			varnishlogs: "\n\n\n\n" + groupExample,
			want:        result,
		}, {
			name:        "multiple empty lines at the end",
			varnishlogs: groupExample + "\n\n\n\n",
			want:        result,
		}, {
			name:        "two groups",
			varnishlogs: groupExample + "\n\n" + groupExample,
			want:        result,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := strings.NewReader(tt.varnishlogs)
			s := bufio.NewScanner(r)
			got, err := ParseGroup(s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseGroup() error = %v, wantErr %v",
					err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseGroup() got = \n%v, want \n%v", got, tt.want)
			}
		})
	}
}

func BenchmarkParseGroup(b *testing.B) {
	scanners := make([]*bufio.Scanner, b.N)
	for i := range scanners {
		r := strings.NewReader(groupExample)
		s := bufio.NewScanner(r)
		scanners[i] = s
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := scanners[i]
		if _, err := ParseGroup(s); err != nil {
			b.Fatal(err)
		}
	}
}
