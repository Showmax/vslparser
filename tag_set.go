package vslparser

// TagSet work as ordered dictionary for faster search in tags. They're meant
// to be read-only.
type TagSet struct {
	lookup map[string][]Tag
	list   []Tag
}

// NewTagSet creates a TagSet structure and does the expensive allocation.
func NewTagSet(tags []Tag) TagSet {
	lookup := make(map[string][]Tag)
	for _, tag := range tags {
		lookup[tag.Key] = append(lookup[tag.Key], tag)
	}

	return TagSet{
		lookup: lookup,
		list:   tags,
	}
}

func (t *TagSet) FirstWithKey(key string) (Tag, bool) { return t.NthWithKey(1, key) }

func (t *TagSet) NthWithKey(n int, key string) (Tag, bool) {
	tags, ok := t.lookup[key]
	if !ok || len(tags) < n {
		return Tag{}, false
	}

	return tags[n-1], true
}

func (t *TagSet) LastWithKey(key string) (Tag, bool) {
	tags, ok := t.lookup[key]
	if !ok || len(tags) < 1 {
		return Tag{}, false
	}

	return tags[len(tags)-1], true
}

// AllWithKey returns a readonly slice of all tags with a given key.
func (t *TagSet) AllWithKey(key string) []Tag { return t.lookup[key] }

// All returns a readonly slice of all tags in t.
func (t *TagSet) All() []Tag { return t.list }
