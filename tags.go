package vslparser

// Tags provides a set of utility functions on an array of Tags.
type Tags []Tag

func (t Tags) FirstWithKey(key string) (Tag, bool) {
	return t.NthWithKey(key, 1)
}

func (t Tags) NthWithKey(key string, n int) (Tag, bool) {
	var cnt int
	for _, tag := range t {
		if tag.Key == key {
			cnt++
			if cnt >= n {
				return tag, true
			}
		}
	}
	return Tag{}, false
}

func (t Tags) LastWithKey(key string) (Tag, bool) {
	for i := len(t) - 1; i > -1; i-- {
		if t[i].Key == key {
			return t[i], true
		}
	}

	return Tag{}, false
}

func (t Tags) AllWithKey(key string) []Tag {
	tags := make([]Tag, 0, 1)
	for _, tag := range t {
		if tag.Key == key {
			tags = append(tags, tag)
		}
	}
	return tags
}
