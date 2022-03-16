package vslparser

// Tags is an array of VSL tags with added functionality.
//
// Contrary to TagSet, Tags are much cheaper to create as it doesn't perform any
// allocations, nor any smart handling. On the other hand searches in plain Tags
// take linear time.
type Tags []Tag

// FirstWithKey finds first tag with key key in t. This method takes O(n).
func (t Tags) FirstWithKey(key string) (Tag, bool) {
	return t.NthWithKey(key, 1)
}

// NthWithKey finds Nth tag with key key in t. This method takes O(n).
func (t Tags) NthWithKey(key string, n int) (Tag, bool) {
	var cnt int
	for _, tag := range t {
		if tag.Key == key {
			cnt++
		}

		if cnt == n {
			return tag, true
		}
	}

	return Tag{}, false
}

// LastWithKey finds last tag with key key in t. This method takes O(n).
func (t Tags) LastWithKey(key string) (Tag, bool) {
	for i := len(t) - 1; i >= 0; i-- {
		if t[i].Key == key {
			return t[i], true
		}
	}

	return Tag{}, false
}
