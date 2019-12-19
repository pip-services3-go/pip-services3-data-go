package persistence

import (
	"reflect"

	"github.com/jinzhu/copier"
	cdata "github.com/pip-services3-go/pip-services3-commons-go/data"
	refl "github.com/pip-services3-go/pip-services3-commons-go/reflect"
)

// Get object Id value
func GetObjectId(item interface{}) interface{} {
	// Todo: Optimize implementation
	return refl.ObjectReader.GetProperty(item, "Id")
}

// Set object Id value
func SetObjectId(item interface{}, id interface{}) {
	if reflect.ValueOf(item).Kind() == reflect.Map {
		refl.ObjectWriter.SetProperty(item, "Id", cdata.IdGenerator.NextLong())
	} else {
		// Todo: Why do we clone the object here??
		typePointer := reflect.New(reflect.TypeOf(item))
		// Todo: Handler pointers to object
		typePointer.Elem().Set(reflect.ValueOf(item))
		interfacePointer := typePointer.Interface()
		refl.ObjectWriter.SetProperty(interfacePointer, "Id", id)
	}
}

// Generates a new id value when it's empty
func GenerateObjectId(item *interface{}) {
	value := *item
	idField := refl.ObjectReader.GetProperty(value, "Id")
	// Todo: Process maps
	if idField != nil && reflect.ValueOf(idField).IsZero() {
		SetObjectId(item, cdata.IdGenerator.NextLong())
	} else {
		panic("Id field doesn't exist")
	}
}

// Clones object
func CloneObject(item interface{}) interface{} {
	var dest interface{}
	var src = item
	if reflect.TypeOf(src).Kind() == reflect.Ptr {
		src = reflect.ValueOf(src).Elem().Interface()
	}
	if reflect.ValueOf(src).Kind() == reflect.Map {
		itemType := reflect.TypeOf(src)
		itemValue := reflect.ValueOf(src)
		mapType := reflect.MapOf(itemType.Key(), itemType.Elem())
		newMap := reflect.MakeMap(mapType)
		for _, k := range itemValue.MapKeys() {
			v := itemValue.MapIndex(k)
			newMap.SetMapIndex(k, v)
		}
		// inflate map
		typePointer := reflect.New(itemType)
		typePointer.Elem().Set(newMap)
		interfacePointer := typePointer.Interface()
		dest = reflect.ValueOf(interfacePointer).Elem().Interface()

	} else {
		copier.Copy(&dest, &src)
	}
	return dest
}

// Compares two values
func CompareValues(value1 interface{}, value2 interface{}) bool {
	// Todo: Implement proper comparison
	return value1 == value2
}
