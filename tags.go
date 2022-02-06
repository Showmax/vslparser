package vslparser

// Tags work as ordered dictionary for faster search in tags. They're meant
// to be read-only.
type Tags struct {
	lookup  map[string][]*Tag
	allTags []Tag
}

// NewTags creates a Tags structure and does the expensive allocation.
func NewTags(tags []Tag) Tags {
	//sort.SliceStable(tags, func(i, j int) bool { return tags[i].Key < tags[j].Key })

	t := Tags{
		lookup:  make(map[string][]*Tag),
		allTags: tags,
	}
	for i := range t.allTags {
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
	out := make([]Tag, len(t.allTags))
	copy(out, t.allTags)
	return out
}
