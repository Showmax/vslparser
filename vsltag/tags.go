package vsltag

import (
	"github.com/Showmax/vslparser"
	"math"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Based on Varnish docs
// https://varnish-cache.org/docs/trunk/reference/vsl.html

// BackendOpen stands for Backend connection opened. Logged when a new backend
// connection is opened.
type BackendOpen vslparser.Tag

func (l BackendOpen) FileDescriptor() int {
	sp := strings.SplitN(l.Value, " ", 2)
	return parseInt(sp[0])
}

func (l BackendOpen) Name() string {
	sp := strings.SplitN(l.Value, " ", 3)
	return sp[1]
}

func (l BackendOpen) RemoteAddr() (addr net.IP, port int) {
	sp := strings.SplitN(l.Value, " ", 5)
	return net.ParseIP(sp[2]), parseInt(sp[3])
}

func (l BackendOpen) LocalAddr() (addr net.IP, port int) {
	sp := strings.SplitN(l.Value, " ", 6)
	return net.ParseIP(sp[4]), parseInt(sp[5])
}

// Begin marks the start of a VXID, the first record of a VXID transaction.
type Begin vslparser.Tag

func (b Begin) Type() string {
	sp := strings.SplitN(b.Value, " ", 2)
	return sp[0]
}

func (b Begin) ParentVXID() int {
	sp := strings.SplitN(b.Value, " ", 3)
	return parseInt(sp[1])
}

func (b Begin) Reason() string {
	i := strings.LastIndex(b.Value, " ")
	return b.Value[i+1:]
}

// BereqMethod stands for Backend request method. The HTTP request method used.
type BereqMethod vslparser.Tag

func (b BereqMethod) Method() string {
	return b.Value
}

// BerespProtocol stands for Backend response protocol. The HTTP protocol version
// information.
type BerespProtocol vslparser.Tag

func (b BerespProtocol) Protocol() string {
	return b.Value
}

// BerespStatus stands for Backend response status. The HTTP response status
// code.
type BerespStatus vslparser.Tag

func (b BerespStatus) Status() int {
	return parseInt(b.Value)
}

// Link links to a child VXID. Links this VXID to any child VXID it initiates.
type Link vslparser.Tag

// ChildType returns "req" or "bereq"
func (l Link) ChildType() string {
	sp := strings.SplitN(l.Value, " ", 2)
	return sp[0]
}

func (l Link) ChildVXID() int {
	sp := strings.SplitN(l.Value, " ", 3)
	return parseInt(sp[1])
}

func (l Link) Reason() string {
	i := strings.LastIndex(l.Value, " ")
	return l.Value[i+1:]
}

// ReqURL contains client request URL. The HTTP request URL.
type ReqURL vslparser.Tag

func (r ReqURL) URL() url.URL {
	u, _ := url.Parse(r.Value)
	return *u
}

// SessClose is the last record for any client connection.
type SessClose vslparser.Tag

func (s SessClose) Reason() string {
	sp := strings.SplitN(s.Value, " ", 2)
	return sp[0]
}

func (s SessClose) Duration() time.Duration {
	i := strings.LastIndex(s.Value, " ")
	return parseDuration(s.Value[i+1:])
}

// SessOpen is the first record for a client connection, with the socket-endpoints of the connection.
type SessOpen vslparser.Tag

func (s SessOpen) RemoteAddr() (addr net.IP, port int) {
	sp := strings.SplitN(s.Value, " ", 3)
	return net.ParseIP(sp[0]), parseInt(sp[1])
}

func (s SessOpen) SocketName() string {
	sp := strings.SplitN(s.Value, " ", 4)
	return sp[2]
}

func (s SessOpen) LocalAddr() (addr net.IP, port int) {
	sp := strings.SplitN(s.Value, " ", 6)
	return net.ParseIP(sp[3]), parseInt(sp[4])
}

func (s SessOpen) SessionStart() time.Time {
	sp := strings.SplitN(s.Value, " ", 7)
	return parseUnixFloat(sp[5])
}

func (s SessOpen) FileDescriptor() int {
	sp := strings.SplitN(s.Value, " ", 8)
	return parseInt(sp[6])
}

// Timestamp contains timing information for the Varnish worker threads.
type Timestamp vslparser.Tag

func (t Timestamp) Event() string {
	sp := strings.SplitN(t.Value, ": ", 2)
	return sp[0]
}

func (t Timestamp) Time() time.Time {
	sp := strings.SplitN(t.Value, " ", 3)
	return parseUnixFloat(sp[1])
}

func (t Timestamp) SinceStart() time.Duration {
	sp := strings.SplitN(t.Value, " ", 4)
	return parseDuration(sp[2])
}

func (t Timestamp) SinceLast() time.Duration {
	i := strings.LastIndex(t.Value, " ")
	return parseDuration(t.Value[i+1:])
}

// Hit object in cache. Object looked up in cache.
type Hit vslparser.Tag

func (h Hit) VXID() int {
	sp := strings.SplitN(h.Value, " ", 2)
	return parseInt(sp[0])
}

func (h Hit) TTL() float64 {
	sp := strings.SplitN(h.Value, " ", 3)
	return parseFloat(sp[1])
}

func (h Hit) Grace() float64 {
	sp := strings.SplitN(h.Value, " ", 4)
	return parseFloat(sp[2])
}

func (h Hit) Keep() float64 {
	i := strings.LastIndex(h.Value, " ")
	return parseFloat(h.Value[i+1:])
}

func parseInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}

func parseFloat(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return f
}

func parseDuration(s string) time.Duration {
	tdNano, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return time.Duration(tdNano * 1e9)
}

func parseUnixFloat(s string) time.Time {
	// TODO consider split by dot as Varnish doesn't try to scam us.
	// It would be faster and resilient to float precision issues.

	unixnano, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return time.Time{}
	}
	sec, dec := math.Modf(unixnano)
	// round to microsecond
	micros := int64(dec*1e6 + 0.5)
	// add trailing zeros, Varnish returns microseconds only.
	// adding zeros later to int type increases accuracy
	nsec := micros * 1e3
	return time.Unix(int64(sec), nsec)
}
