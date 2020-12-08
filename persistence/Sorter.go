package persistence

/*
Helper class for sorting data in MemoryPersistence
implements sort.Interface
*/

//------------- Sorter -----------------------
type sorter struct {
	items    []interface{}
	compFunc func(a, b interface{}) bool
}

// Calculate lenth
// Return length of items array
func (s sorter) Len() int {
	return len(s.items)
}

// Make swap two items in array
// Parameters:
//	 - i,j int
//	indexes of array for swap
func (s sorter) Swap(i, j int) {
	s.items[i], s.items[j] = s.items[j], s.items[i]
}

// Compare less function
// Parameters:
//	 - i,j int
//	 indexes of array for compare
// Returns bool
// true if items[i] < items[j] and false otherwise
func (s sorter) Less(i, j int) bool {
	if s.compFunc == nil {
		panic("Sort.Less Error compare function is nil!")
	}
	return s.compFunc(s.items[i], s.items[j])
}
