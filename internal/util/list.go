package util

// TagsIntersection returns true if there are any tags in common between
// the two lists
func TagsIntersection(l1, l2 []string) bool {
	for _, i1 := range l1 {
		for _, i2 := range l2 {
			if i1 == i2 {
				return true
			}
		}
	}

	return false
}
