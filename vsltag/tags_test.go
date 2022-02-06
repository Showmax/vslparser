package vsltag_test

import (
	"net"
	"testing"
	"time"

	vsl "github.com/Showmax/vslparser"
	"github.com/Showmax/vslparser/vsltag"
)

func TestSessOpen(t *testing.T) {
	type want struct {
		remoteIP   net.IP
		remotePort int
		socketName string
		localIP    net.IP
		localPort  int
		sessStrat  time.Time
		fileDesc   int
	}
	tests := []struct {
		name string
		tag  vsl.Tag
		want want
	}{
		{
			tag: vsl.Tag{
				Key:   "SessOpen",
				Value: "10.46.103.82 5480 a0 10.243.103.218 6081 1604933732.219939 25",
			},
			want: want{
				remoteIP:   net.IPv4(10, 46, 103, 82),
				remotePort: 5480,
				socketName: "a0",
				localIP:    net.IPv4(10, 243, 103, 218),
				localPort:  6081,
				sessStrat:  time.Date(2020, 11, 9, 14, 55, 32, 219939000, time.UTC),
				fileDesc:   25,
			},
		},
		{
			tag: vsl.Tag{
				Key:   "SessOpen",
				Value: "10.243.103.218 2040 a0 10.243.103.218 6081 1607417773.924903 32",
			},
			want: want{
				remoteIP:   net.IPv4(10, 243, 103, 218),
				remotePort: 2040,
				socketName: "a0",
				localIP:    net.IPv4(10, 243, 103, 218),
				localPort:  6081,
				sessStrat:  time.Date(2020, 12, 8, 8, 56, 13, 924903000, time.UTC),
				fileDesc:   32,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sessOpen := vsltag.SessOpen(tt.tag)

			if gotAddr, gotPort := sessOpen.RemoteAddr(); !gotAddr.Equal(tt.want.remoteIP) || gotPort != tt.want.remotePort {
				t.Errorf("RemoteAddr() = %v:%v, want %v:%v",
					gotAddr, gotPort,
					tt.want.remoteIP, tt.want.remotePort)
			}
			if got := sessOpen.SocketName(); got != tt.want.socketName {
				t.Errorf("SocketName() = %v, want %v", got, tt.want.socketName)
			}
			if gotAddr, gotPort := sessOpen.LocalAddr(); !gotAddr.Equal(tt.want.localIP) || gotPort != tt.want.localPort {
				t.Errorf("LocalAddr() = %v:%v, want %v:%v",
					gotAddr, gotPort,
					tt.want.localIP, tt.want.localPort)
			}
			if got, err := sessOpen.SessionStart(); err != nil || !got.Equal(tt.want.sessStrat) {
				t.Errorf("SessionStart() = %v, want %v (error: %v)",
					// Show it in UTC to be it easy to compare for humans.
					got.In(time.UTC), tt.want.sessStrat.In(time.UTC), err)
			}
			if got := sessOpen.FileDescriptor(); got != tt.want.fileDesc {
				t.Errorf("FileDescriptor() = %v, want %v", got, tt.want.fileDesc)
			}
		})
	}
}
