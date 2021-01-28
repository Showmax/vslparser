package vsltag

import (
	"github.com/Showmax/vslparser"
	"math"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Based on Varnish docs
// https://varnish-cache.org/docs/trunk/reference/vsl.html

type BackendOpen vslparser.Tag

func (l BackendOpen) FileDescriptor() int {
	sp := strings.SplitN(l.Value, " ", 2)
	fd, _ := strconv.Atoi(sp[0])
	return fd
}

func (l BackendOpen) Name() string {
	sp := strings.SplitN(l.Value, " ", 3)
	return sp[1]
}

func (l BackendOpen) RemoteAddr() (addr string, port string) {
	sp := strings.SplitN(l.Value, " ", 5)
	return sp[2], sp[3]
}

func (l BackendOpen) LocalAddr() (addr string, port string) {
	sp := strings.SplitN(l.Value, " ", 6)
	return sp[4], sp[5]
}

// Begin marks the start of a VXID, the first record of a VXID transaction.
type Begin vslparser.Tag

func (b Begin) Type() string {
	sp := strings.SplitN(b.Value, " ", 2)
	return sp[0]
}

func (b Begin) ParentVXID() int {
	sp := strings.SplitN(b.Value, " ", 3)
	vxid, _ := strconv.Atoi(sp[1])
	return vxid
}

func (b Begin) Reason() string {
	i := strings.LastIndex(b.Value, " ")
	return b.Value[i+1:]
}

type BereqMethod vslparser.Tag

func (b BereqMethod) Method() string {
	return b.Value
}

type BerespProtocol vslparser.Tag

func (b BerespProtocol) Protocol() string {
	return b.Value
}

type BerespStatus vslparser.Tag

func (b BerespStatus) Status() int {
	i, _ := strconv.Atoi(b.Value)
	return i
}

type Link vslparser.Tag

func (l Link) ChildType() string {
	sp := strings.SplitN(l.Value, " ", 2)
	return sp[0]
}

func (l Link) ChildVXID() int {
	sp := strings.SplitN(l.Value, " ", 3)
	vxid, _ := strconv.Atoi(sp[1])
	return vxid
}

func (l Link) Reason() string {
	i := strings.LastIndex(l.Value, " ")
	return l.Value[i+1:]
}

type ReqURL vslparser.Tag

func (r ReqURL) URL() url.URL {
	u, _ := url.Parse(r.Value)
	return *u
}

type SessClose vslparser.Tag

func (s SessClose) Reason() string {
	sp := strings.SplitN(s.Value, " ", 2)
	return sp[0]
}

func (s SessClose) Duration() time.Duration {
	i := strings.LastIndex(s.Value, " ")
	td, _ := time.ParseDuration(s.Value[i+1:] + "s")
	return td
}

type Timestamp vslparser.Tag

func (t Timestamp) Event() string {
	sp := strings.SplitN(t.Value, ": ", 2)
	return sp[0]
}

func (t Timestamp) Time() time.Time {
	sp := strings.SplitN(t.Value, " ", 3)
	unixnano, _ := strconv.ParseFloat(sp[1], 64)
	sec, dec := math.Modf(unixnano)
	return time.Unix(int64(sec), int64(dec*1e9))
}

func (t Timestamp) SinceStart() time.Duration {
	sp := strings.SplitN(t.Value, " ", 4)
	tdNano, _ := strconv.ParseFloat(sp[2], 64)
	return time.Duration(tdNano * 1e9)
}

func (t Timestamp) SinceLast() time.Duration {
	i := strings.LastIndex(t.Value, " ")
	tdNano, _ := strconv.ParseFloat(t.Value[i+1:], 64)
	return time.Duration(tdNano * 1e9)
}

type Hit vslparser.Tag

func (h Hit) VXID() int {
	sp := strings.SplitN(h.Value, " ", 2)
	vxid, _ := strconv.Atoi(sp[0])
	return vxid
}

func (h Hit) TTL() float64 {
	sp := strings.SplitN(h.Value, " ", 3)
	f, _ := strconv.ParseFloat(sp[1], 64)
	return f
}

func (h Hit) Grace() float64 {
	sp := strings.SplitN(h.Value, " ", 4)
	f, _ := strconv.ParseFloat(sp[2], 64)
	return f
}

func (h Hit) Keep() float64 {
	i := strings.LastIndex(h.Value, " ")
	f, _ := strconv.ParseFloat(h.Value[i+1:], 64)
	return f
}
