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

// FirstWithKey returns first tag with key k in t. This method takes O(1).
func (t *TagSet) FirstWithKey(k string) (Tag, bool) { return t.NthWithKey(k, 0) }

// NthWithKey returns Nth (indexed from zero) tag with key k in t. This method
// takes O(1).
func (t *TagSet) NthWithKey(k string, n uint) (Tag, bool) {
	tags, ok := t.lookup[k]
	if !ok || uint(len(tags)) <= n {
		return Tag{}, false
	}

	return tags[n], true
}

// LastWithKey returns last tag with key k in t. This method takes O(1).
func (t *TagSet) LastWithKey(k string) (Tag, bool) {
	tags, ok := t.lookup[k]
	if !ok || len(tags) < 1 {
		return Tag{}, false
	}

	return tags[len(tags)-1], true
}

// AllWithKey returns a readonly slice of all tags with a given key.
func (t *TagSet) AllWithKey(k string) []Tag { return t.lookup[k] }

// All returns a readonly slice of all tags in t.
func (t *TagSet) All() []Tag { return t.list }
