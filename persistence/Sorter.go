package persistence

/*
Helper type for sorting data in memory persistence
*/
import "reflect"

//------------- Sorter -----------------------

type sorter struct {
	items    []interface{}
	compFunc interface{}
}

func (s sorter) Len() int {
	return len(s.items)
}

func (s sorter) Swap(i, j int) {
	s.items[i], s.items[j] = s.items[j], s.items[i]
}

func (s sorter) Less(i, j int) bool {
	if s.compFunc == nil {
		panic("Sort.Less Error compare function is nil!")
	}
	compFuncType := reflect.TypeOf(s.compFunc)
	if compFuncType.NumIn() != 2 {
		panic("Sort.Less Error compare function must recive 2 arguments")
	}
	_compFunc := reflect.ValueOf(s.compFunc)
	result := _compFunc.Call([]reflect.Value{reflect.ValueOf(s.items[i]), reflect.ValueOf(s.items[j])})
	return result[0].Bool()
}
