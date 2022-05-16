package vslparser

// Tags is an array of VSL tags with added functionality.
//
// Contrary to TagSet, Tags are much cheaper to create as it doesn't perform any
// allocations, nor any smart handling. On the other hand searches in plain Tags
// take linear time.
type Tags []Tag

// FirstWithKey finds first tag with key k in t. This method takes O(n).
func (t Tags) FirstWithKey(k string) (Tag, bool) {
	return t.NthWithKey(k, 0)
}

// NthWithKey finds Nth (indexed from zero) tag with key k in t. This method
// takes O(n).
func (t Tags) NthWithKey(k string, n uint) (Tag, bool) {
	var cnt uint
	for _, tag := range t {
		if tag.Key != k {
			continue
		}

		if cnt == n {
			return tag, true
		}
		cnt++
	}

	return Tag{}, false
}

// LastWithKey finds last tag with key k in t. This method takes O(n).
func (t Tags) LastWithKey(k string) (Tag, bool) {
	for i := len(t) - 1; i >= 0; i-- {
		if t[i].Key == k {
			return t[i], true
		}
	}

	return Tag{}, false
}
