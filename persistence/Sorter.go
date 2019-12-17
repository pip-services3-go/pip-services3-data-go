package persistence

/*
Helper type for sorting data in memory persistence
*/
//------------- Sorter -----------------------
type sorter struct {
	items    []interface{}
	compFunc func(a, b interface{}) bool
}

// Return length of items array
func (s sorter) Len() int {
	return len(s.items)
}

// Make swap two items in array
func (s sorter) Swap(i, j int) {
	s.items[i], s.items[j] = s.items[j], s.items[i]
}

// Compare less function
func (s sorter) Less(i, j int) bool {
	if s.compFunc == nil {
		panic("Sort.Less Error compare function is nil!")
	}
	return s.compFunc(s.items[i], s.items[j])
}
