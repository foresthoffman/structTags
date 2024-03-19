// Package struct_tags allows for marshalling non-JSON and third-party struct
// tags. The standard JSON.Marshal() function cannot marshal custom tags, however
// this package can.
//
// e.g. Third-party struct tags could look like the "db" tag below.
//
//   type MyStruct struct {
//     Field string `json:"api_field" db:"db_field"`
//   }
package struct_tags

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"reflect"
)

const (
	ErrNilObject = "object was nil"
)

// fieldMetadata helps maintain the order of a temporary list of reflect.Value field objects.
type fieldMetadata struct {
	TagValue string
	Value    reflect.Value
}

// CustomMarshaller allows for marshalling non-JSON and third-party struct tags.
type CustomMarshaller struct {
	TargetTag          string
	IgnoreTagWithValue string
}

// NewCustomMarshaller creates a new custom-tag marshalling instance.
func NewCustomMarshaller(targetTag, ignoreTag string) *CustomMarshaller {
	return &CustomMarshaller{
		TargetTag:          targetTag,
		IgnoreTagWithValue: ignoreTag,
	}
}

// reflectStructFields recursively searches for struct field tags and writes the
// tagged key-value pair to the provided buffer.
func (m *CustomMarshaller) reflectStructFields(w io.Writer, v reflect.Value) {
	if v == reflect.Zero(reflect.TypeOf(v)).Interface() {
		return
	}

	// Remove ignored fields and ensure that the order of the fields won't change.
	fields := map[int]fieldMetadata{}
	count := 0
	for j := 0; j < v.NumField(); j++ {
		tagValue := v.Type().Field(j).Tag.Get(m.TargetTag)
		if tagValue == m.IgnoreTagWithValue {
			continue
		}
		fields[count] = fieldMetadata{
			TagValue: tagValue,
			Value:    v.Field(j),
		}
		count++
	}

	for i := 0; i < len(fields); i++ {
		switch fields[i].Value.Type().Kind() {
		case reflect.Struct:
			w.Write([]byte(fmt.Sprintf("%q:", fields[i].TagValue)))
			w.Write([]byte("{"))
			m.reflectStructFields(w, fields[i].Value)
			w.Write([]byte("}"))
		case reflect.Ptr:
			w.Write([]byte(fmt.Sprintf("%q:", fields[i].TagValue)))
			w.Write([]byte("{"))
			m.reflectStructFields(w, fields[i].Value.Elem())
			w.Write([]byte("}"))
		case reflect.Slice:
			w.Write([]byte(fmt.Sprintf("%q:", fields[i].TagValue)))
			w.Write([]byte("["))
			numItems := fields[i].Value.Len()
			for x := 0; x < numItems; x++ {
				w.Write([]byte("{"))
				m.reflectStructFields(w, fields[i].Value.Index(x))
				w.Write([]byte("}"))
				if x+1 < numItems {
					w.Write([]byte(","))
				}
			}
			w.Write([]byte("]"))
		default:
			w.Write([]byte(fmt.Sprintf("%q:%q", fields[i].TagValue, fields[i].Value.String())))
		}
		if i+1 < len(fields) {
			w.Write([]byte(","))
		}
	}
}

// Marshal takes the provided object and JSON-marshals it using the
// pre-configured target tag and ignored tag values. Normal JSON struct tags will
// not be used.
func (m *CustomMarshaller) Marshal(obj interface{}) ([]byte, error) {
	if obj == nil {
		return nil, errors.New(ErrNilObject)
	}

	v := reflect.ValueOf(obj)
	if val, ok := obj.(reflect.Value); ok {
		v = val
	}

	w := bytes.NewBuffer([]byte{})
	switch reflect.TypeOf(obj).Kind() {
	case reflect.Struct:
		w.Write([]byte("{"))
		m.reflectStructFields(w, v)
		w.Write([]byte("}\n"))
	case reflect.Ptr:
		w.Write([]byte("{"))
		m.reflectStructFields(w, v.Elem())
		w.Write([]byte("}\n"))
	case reflect.Slice:
		numItems := v.Len()
		for i := 0; i < numItems; i++ {
			w.Write([]byte("{"))
			m.reflectStructFields(w, v.Index(i))
			w.Write([]byte("}\n"))
		}
	}

	return w.Bytes(), nil
}
