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

// Tags work as ordered dictionary for faster search in tags. They're meant
// to be read-only.
type Tags struct {
	lookup map[string][]*Tag
	list   []Tag
}

// NewTags creates a Tags structure and does the expensive allocation.
func NewTags(tags []Tag) Tags {
	t := Tags{
		lookup: make(map[string][]*Tag),
		list:   tags,
	}
	for i := range t.list {
		// https://github.com/golang/go/wiki/CommonMistakes#using-reference-to-loop-iterator-variable
		tag := tags[i]
		t.lookup[tag.Key] = append(t.lookup[tag.Key], &tag)
	}
	return t
}

func (t *Tags) FirstWithKey(key string) (Tag, bool) {
	return t.NthWithKey(1, key)
}

func (t *Tags) NthWithKey(n int, key string) (Tag, bool) {
	tags, ok := t.lookup[key]
	if !ok {
		return Tag{}, false
	}
	if len(tags) < n {
		return Tag{}, false
	}
	return *tags[n-1], true
}

func (t *Tags) LastWithKey(key string) (Tag, bool) {
	tags, ok := t.lookup[key]
	if !ok {
		return Tag{}, false
	}
	if len(tags) == 0 {
		return Tag{}, false
	}
	return *tags[len(tags)-1], true
}

func (t *Tags) AllWithKey(key string) []Tag {
	tags := t.lookup[key]
	out := make([]Tag, len(tags))
	for i, t := range tags {
		out[i] = *t
	}
	return out
}

func (t *Tags) All() []Tag {
	out := make([]Tag, len(t.list))
	copy(out, t.list)
	return out
}
