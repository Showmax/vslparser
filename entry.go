package vslparser

// Entry holds a single log entry. An entry consists mostly of a collection of
// log fields, called Tags.
type Entry struct {
	Level int
	Kind  string
	VXID  int
	Tags  []Tag
}

// Tag is the key/value pair making up a VSL tag.
type Tag struct {
	Key   string
	Value string
}
