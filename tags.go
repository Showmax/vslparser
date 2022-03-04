package vslparser

// Tags work as ordered dictionary for faster search in tags. They're meant
// to be read-only.
type Tags struct {
	lookup map[string][]Tag
	list   []Tag
}

// NewTags creates a Tags structure and does the expensive allocation.
func NewTags(tags []Tag) Tags {
	lookup := make(map[string][]Tag)
	for _, tag := range tags {
		lookup[tag.Key] = append(lookup[tag.Key], tag)
	}

	return Tags{
		lookup: lookup,
		list:   tags,
	}
}

func (t *Tags) FirstWithKey(key string) (Tag, bool) { return t.NthWithKey(1, key) }

func (t *Tags) NthWithKey(n int, key string) (Tag, bool) {
	tags, ok := t.lookup[key]
	if !ok || len(tags) < n {
		return Tag{}, false
	}

	return tags[n-1], true
}

func (t *Tags) LastWithKey(key string) (Tag, bool) {
	tags, ok := t.lookup[key]
	if !ok || len(tags) < 1 {
		return Tag{}, false
	}

	return tags[len(tags)-1], true
}

// AllWithKey returns a readonly slice of all tags with a given key.
func (t *Tags) AllWithKey(key string) []Tag { return t.lookup[key] }

// All returns a readonly slice of all tags in t.
func (t *Tags) All() []Tag { return t.list }
