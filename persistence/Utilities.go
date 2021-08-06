package persistence

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/jinzhu/copier"
	"github.com/pip-services3-go/pip-services3-commons-go/convert"
	cdata "github.com/pip-services3-go/pip-services3-commons-go/data"
	refl "github.com/pip-services3-go/pip-services3-commons-go/reflect"
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
//   - obj interface{}
//   an object to read property from.
//   - name string
//   a name of the property to get.
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
		return getPropertyRecursive(fieldType, obj, name)
	}

	return nil
}

func getPropertyRecursive(fieldType reflect.Type, obj interface{}, name string) interface{} {
	for index := 0; index < fieldType.NumField(); index++ {
		field := fieldType.Field(index)
		val := reflect.ValueOf(obj)
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
		switch field.Type.Kind() {
		default:
			if matchField(field, name) {
				return val.Field(index).Interface()
			}
		case reflect.Struct:
			if item := getPropertyRecursive(field.Type, val.Field(index).Interface(), name); item != nil {
				return item
			}
		}
	}

	return nil
}

// Sets value of object property specified by its name.
// If the property does not exist or introspection fails this method doesn't do anything and doesn't any throw errors.
// Parameters:
//   - obj interface{}
//   an object to write property to.
//   name string
//   a name of the property to set.
//   - value interface{}
//   a new value for the property to set.
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
		if err := recover(); err != nil {
			fmt.Println("Error while set property %v", err)
		}
	}()

	fieldType := toFieldType(obj)
	if fieldType.Kind() == reflect.Struct {
		setPropertyRecursive(fieldType, obj, name, value)
	}

}

func setPropertyRecursive(fieldType reflect.Type, obj interface{}, name string, value interface{}) {
	for index := 0; index < fieldType.NumField(); index++ {
		field := fieldType.Field(index)
		val := reflect.ValueOf(obj)
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
		switch field.Type.Kind() {
		default:
			if matchField(field, name) {
				val.Field(index).Set(reflect.ValueOf(value))
				return
			}
		case reflect.Struct:
			setPropertyRecursive(field.Type, val.Field(index).Addr().Interface(), name, value)
		}
	}
}

// Get object Id value
// Parameters:
//   - item interface{}
//   an object to read property from.
// Returns interface{}
// the property value or nil if property doesn't exist or introspection failed.
func GetObjectId(item interface{}) interface{} {
	return GetProperty(item, "Id")
}

// SetObjectId is set object Id value
// Parameters:
//   - item *interface{}
//   an pointer on object to set id property
//   - id interface{}
//   id value for set
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
//   - item *interface{}
//   an pointer on object to set id property
// Results saved in input object
func GenerateObjectId(item *interface{}) {
	value := *item
	idField := GetProperty(value, "Id")
	if idField == nil {
		idField = GetProperty(value, "ID")
	}
	if idField == nil {
		panic("'Id' or 'ID' field doesn't exist")
	}
	if reflect.ValueOf(idField).IsZero() {
		SetObjectId(item, cdata.IdGenerator.NextLong())
	}
}

// CloneObject is clones object function
// Parameters:
//   - item interface{}
//   an object to clone
// Return interface{}
// copy of input item
func CloneObject(item interface{}, proto reflect.Type) interface{} {
	var dest interface{}
	var src = item

	if reflect.ValueOf(src).Kind() == reflect.Map {
		itemType := reflect.TypeOf(src)
		mapType := reflect.MapOf(itemType.Key(), itemType.Elem())
		newMap := reflect.MakeMap(mapType)
		dest = newMap.Interface()
		err := copier.CopyWithOption(&dest, src, copier.Option{DeepCopy: false, IgnoreEmpty: false})
		if err != nil {
			return nil
		}

	} else {
		var destPtr reflect.Value
		if proto.Kind() == reflect.Ptr {
			destPtr = reflect.New(proto.Elem())
		} else {
			destPtr = reflect.New(proto)
		}
		if reflect.TypeOf(src).Kind() == reflect.Ptr {
			src = reflect.ValueOf(src).Elem().Interface()
		}
		err := copier.CopyWithOption(destPtr.Interface(), src, copier.Option{DeepCopy: false, IgnoreEmpty: false})
		if err != nil {
			return nil
		}

		dest = destPtr.Elem().Interface()
	}
	return dest
}

// CloneObjectForResult is clones object for result function
// Parameters:
// 	  - item interface{}
// 	  an object to clone
//	  -proto reflect.Type
//	  type of returned value, need for detect object or pointer returned type
// Return interface{}
// copy of input item
func CloneObjectForResult(src interface{}, proto reflect.Type) interface{} {
	var dest interface{}

	if reflect.ValueOf(src).Kind() == reflect.Map {
		itemType := reflect.TypeOf(src)
		mapType := reflect.MapOf(itemType.Key(), itemType.Elem())
		newMap := reflect.MakeMap(mapType)
		dest = newMap.Interface()
		err := copier.CopyWithOption(&dest, src, copier.Option{DeepCopy: false, IgnoreEmpty: false})
		if err != nil {
			return nil
		}
	} else {
		var destPtr reflect.Value
		if proto.Kind() == reflect.Ptr {
			destPtr = reflect.New(proto.Elem())
		} else {
			destPtr = reflect.New(proto)
		}
		err := copier.CopyWithOption(destPtr.Interface(), src, copier.Option{DeepCopy: false, IgnoreEmpty: false})
		if err != nil {
			return nil
		}
		// make pointer on clone object, if proto is ptr
		dest = destPtr.Elem().Interface()
		if proto.Kind() == reflect.Ptr {
			dest = destPtr.Interface()
		}
	}

	return dest
}

// CompareValues are ompares two values
// Parameters:
//   - value1 interface{}
//   an object one for compare
//   - value2 interface{}
//   an object two for compare
// Return bool
// true if value1 equal value2 and false otherwise
func CompareValues(value1 interface{}, value2 interface{}) bool {
	// Todo: Implement proper comparison
	return value1 == value2
}

// Convert methods

// FromIds method convert ids string array to array of interface{} object
// Parameters:
//   - ids - []string
//   array of ids
// Return []interface{}
// array of ids
func FromIds(ids []string) []interface{} {
	result := make([]interface{}, len(ids))
	for i, v := range ids {
		result[i] = v
	}
	return result
}

// ToPublicMap method convert interface{} object to map[string]interface{}
// Parameters
//   - value - interface{}
//   input object to convert
// Return map[string]interface{}
// converted object to map
func ToPublicMap(value interface{}) map[string]interface{} {
	if value != nil {
		result, _ := value.(map[string]interface{})
		return result
	}
	return nil
}

// ToPublicArray method convert array of interface{} object to array of map[string]interface{}
// Parameters
//   - value - []interface{}
//   input object to convert
// Return []map[string]interface{}
// converted map array
func ToPublicArray(values []interface{}) []map[string]interface{} {
	if values == nil {
		return nil
	}

	result := make([]map[string]interface{}, len(values))
	for i, v := range values {
		result[i] = ToPublicMap(v)
	}
	return result
}
