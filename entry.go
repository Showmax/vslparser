package vslparser

// VXID is varnish transaction ID.
//
// In Varnish code, VXIDs are represented as uint32_t type and the value wraps
// at 1<<30, so Go uint32 us absolutely sufficient to hold any possible value of
// VXID.
//
// Note: Zero value (vxid == 0) is not reported in any varnishlog mode with an
// exception of 'raw' mode. Consequently, zero value of VXID can be in some
// contexts (but not all!) used to indicate and "invalid VXID" or "none" value.
// https://varnish-cache.org/docs/trunk/reference/vsl-query.html
type VXID uint32

// Entry holds a single log entry. An entry consists mostly of a collection of
// log fields, called Tags.
type Entry struct {
	Level int
	Kind  string
	VXID  VXID
	Tags  []Tag
}

// Tag is the key/value pair making up a VSL tag.
type Tag struct {
	Key   string
	Value string
}
