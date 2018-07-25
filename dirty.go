// Package dirty tracks changes to struct fields.
package dirty

import (
	"reflect"
)

var (
	tracked = map[interface{}]trackedStruct{}
)

type trackedStruct struct {
	s      interface{}
	fields map[string]trackedField
}

type trackedField struct {
	s   *trackedStruct
	i   int
	ft  reflect.StructField
	typ reflect.Type
	was reflect.Value
}

func (f trackedField) is() reflect.Value {
	rv := reflect.ValueOf(f.s.s).Elem()
	rvf := rv.Field(f.i)
	if f.ft.Type.Kind() == reflect.Ptr && !rvf.IsNil() {
		rvf = rvf.Elem()
	}
	return rvf
}

func (f trackedField) changed() bool {
	if isZero(f.was) && isZero(f.is()) {
		return false
	}

	if !isZero(f.was) && !isZero(f.is()) {
		return f.was.Interface() != f.is().Interface()
	}

	return true
}

func track(o interface{}) trackedStruct {
	s := trackedStruct{
		s:      o,
		fields: map[string]trackedField{},
	}

	rt := reflect.TypeOf(o)
	rv := reflect.ValueOf(o)
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
		rv = rv.Elem()
	}

	for i := 0; i < rt.NumField(); i++ {
		rtf := rt.Field(i)
		switch rtf.Type.Kind() {
		case reflect.Array, reflect.Map, reflect.Slice:
			continue
		}

		rvf := rv.Field(i)
		if rtf.Type.Kind() == reflect.Ptr && !rvf.IsNil() {
			rvf = rvf.Elem()
		}
		s.fields[rtf.Name] = trackedField{
			s:   &s,
			i:   i,
			ft:  rtf,
			typ: rtf.Type,
			was: reflect.ValueOf(rvf.Interface()),
		}
	}

	return s
}

func isZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}

// Track starts tracking `o` for changes. Panics if `o` is not a pointer to a
// struct.
func Track(o interface{}) {
	rt := reflect.TypeOf(o)
	if rt.Kind() != reflect.Ptr {
		panic("can only track pointers")
	}
	rv := reflect.ValueOf(o)
	if rv.Elem().Kind() != reflect.Struct {
		panic("can only track pointers to structs")
	}

	tracked[o] = track(o)
}

// Forget forgets all about `o`.
func Forget(o interface{}) {
	delete(tracked, o)
}

// Changed checks if `o` was changed, i.e. the value of any of its fields was
// changed.
func Changed(o interface{}) bool {
	return len(Changes(o)) > 0
}

// Changes returns the changes, if any, that happened to any of the values of
// the fields of `o`. The changes are in the form of a map, with the key being
// the name of the changed field, and the value being a slice of 2 interfaces:
// the value of the field at the time `o` was tracked, and the current/changed
// value of the field.
func Changes(o interface{}) map[string][]interface{} {
	s, ok := tracked[o]
	if !ok {
		panic("not tracked")
	}

	changes := map[string][]interface{}{}
	for name, f := range s.fields {
		if f.changed() {
			changes[name] = []interface{}{
				f.was.Interface(), f.is().Interface(),
			}
		}
	}
	return changes
}
