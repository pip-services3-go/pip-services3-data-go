package persistence

import (
	"reflect"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/jinzhu/copier"
	"github.com/pip-services3-go/pip-services3-commons-go/v3/convert"
	cdata "github.com/pip-services3-go/pip-services3-commons-go/v3/data"
	refl "github.com/pip-services3-go/pip-services3-commons-go/v3/reflect"
)

func toFieldType(obj interface{}) reflect.Type {
	// Unwrap value
	wrap, ok := obj.(refl.IValueWrapper)
	if ok {
		obj = wrap.InnerValue()
	}

	// Move from pointer to real struct
	typ := reflect.TypeOf(obj)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	return typ
}

func matchField(field reflect.StructField, name string) bool {
	// Field must be public and match to name as case insensitive
	r, _ := utf8.DecodeRuneInString(field.Name)
	return unicode.IsUpper(r) &&
		strings.ToLower(field.Name) == strings.ToLower(name)
}

func getValue(obj interface{}) interface{} {
	wrap, ok := obj.(refl.IValueWrapper)
	if ok {
		obj = wrap.InnerValue()
	}

	return obj
}

// Gets value of object property specified by its name.
// Parameters:
// 			- obj interface{}
// 			an object to read property from.
// 			- name string
// 			a name of the property to get.
// Returns interface{}
// the property value or null if property doesn't exist or introspection failed.
func GetProperty(obj interface{}, name string) interface{} {
	if obj == nil || name == "" {
		return nil
	}

	obj = getValue(obj)
	val := reflect.ValueOf(obj)

	if val.Kind() == reflect.Map {
		name = strings.ToLower(name)
		for _, v := range val.MapKeys() {
			key := convert.StringConverter.ToString(v.Interface())
			key = strings.ToLower(key)
			if name == key {
				return val.MapIndex(v).Interface()
			}
		}
		return nil
	}

	defer func() {
		// Do nothing and return nil
		recover()
	}()

	fieldType := toFieldType(obj)
	if fieldType.Kind() == reflect.Struct {
		for index := 0; index < fieldType.NumField(); index++ {
			field := fieldType.Field(index)
			if matchField(field, name) {
				val := reflect.ValueOf(obj)
				if val.Kind() == reflect.Ptr {
					val = val.Elem()
				}
				return val.Field(index).Interface()
			}
		}
	}

	return nil
}

// Sets value of object property specified by its name.
// If the property does not exist or introspection fails this method doesn't do anything and doesn't any throw errors.
// Parameters:
// 			- obj interface{}
// 			an object to write property to.
// 			name string
// 			a name of the property to set.
// 			- value interface{}
// 			a new value for the property to set.
func SetProperty(obj interface{}, name string, value interface{}) {
	if obj == nil || name == "" {
		return
	}

	obj = getValue(obj)
	val := reflect.ValueOf(obj)

	if val.Kind() == reflect.Map {
		name = strings.ToLower(name)
		for _, v := range val.MapKeys() {
			key := convert.StringConverter.ToString(v.Interface())
			key = strings.ToLower(key)
			if name == key {
				val.SetMapIndex(v, reflect.ValueOf(value))
				return
			}
		}
		val.SetMapIndex(reflect.ValueOf(name), reflect.ValueOf(value))
		return
	}

	defer func() {
		// Do nothing and return nil
		recover()
	}()

	fieldType := toFieldType(obj)
	if fieldType.Kind() == reflect.Struct {
		for index := 0; index < fieldType.NumField(); index++ {
			field := fieldType.Field(index)
			if matchField(field, name) {
				val := reflect.ValueOf(obj)
				if val.Kind() == reflect.Ptr {
					val = val.Elem()
				}
				val.Field(index).Set(reflect.ValueOf(value))
				return
			}
		}
	}
}

// Get object Id value
// Parameters:
// 			- item interface{}
// 			an object to read property from.
// Returns interface{}
// the property value or nil if property doesn't exist or introspection failed.
func GetObjectId(item interface{}) interface{} {
	return GetProperty(item, "Id")

}

// SetObjectId is set object Id value
// Parameters:
// 			- item *interface{}
// 			an pointer on object to set id property
// 			- id interface{}
//			id value for set
// Results saved in input object
func SetObjectId(item *interface{}, id interface{}) {
	value := *item
	if reflect.ValueOf(value).Kind() == reflect.Map {
		//refl.ObjectWriter.SetProperty(value, "Id", id)
		SetProperty(value, "Id", id)
	} else {
		typePointer := reflect.New(reflect.TypeOf(value))
		typePointer.Elem().Set(reflect.ValueOf(value))
		typeInterface := typePointer.Interface()
		//refl.ObjectWriter.SetProperty(typeInterface, "Id", id)
		SetProperty(typeInterface, "Id", id)
		*item = reflect.ValueOf(typeInterface).Elem().Interface()
	}
}

// GenerateObjectId is generates a new id value when it's empty
// Parameters:
// 			- item *interface{}
// 			an pointer on object to set id property
// Results saved in input object
func GenerateObjectId(item *interface{}) {
	value := *item
	idField := GetProperty(value, "Id")
	if idField != nil {
		if reflect.ValueOf(idField).IsZero() {
			SetObjectId(item, cdata.IdGenerator.NextLong())
		}
	} else {
		panic("Id field doesn't exist")
	}
}

// CloneObject is clones object function
// Parameters:
// 			- item interface{}
// 			an object to clone
// Return interface{}
// copy of input item
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

// CompareValues are ompares two values
// Parameters:
// 			- value1 interface{}
// 			an object one for compare
// 			- value2 interface{}
// 			an object two for compare
// Return bool
// true if value1 equal value2 and false otherwise
func CompareValues(value1 interface{}, value2 interface{}) bool {
	// Todo: Implement proper comparison
	return value1 == value2
}