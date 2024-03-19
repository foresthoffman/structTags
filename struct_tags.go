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

type CustomMarshaler struct {
	TargetTag          string
	IgnoreTagWithValue string
}

func NewCustomMarshaler(targetTag, ignoreTag string) *CustomMarshaler {
	return &CustomMarshaler{
		TargetTag:          targetTag,
		IgnoreTagWithValue: ignoreTag,
	}
}

func (m *CustomMarshaler) reflectStructFields(w io.Writer, v reflect.Value) {
	if v == reflect.Zero(reflect.TypeOf(v)).Interface() {
		return
	}

	fields := map[string]reflect.Value{}
	for j := 0; j < v.NumField(); j++ {
		tagValue := v.Type().Field(j).Tag.Get(m.TargetTag)
		if tagValue == m.IgnoreTagWithValue {
			continue
		}
		fields[tagValue] = v.Field(j)
	}

	i := 0
	for tag, field := range fields {
		switch field.Type().Kind() {
		case reflect.Struct:
			w.Write([]byte(fmt.Sprintf("%q:", tag)))
			w.Write([]byte("{"))
			m.reflectStructFields(w, field)
			w.Write([]byte("}"))
		case reflect.Ptr:
			w.Write([]byte(fmt.Sprintf("%q:", tag)))
			w.Write([]byte("{"))
			m.reflectStructFields(w, field.Elem())
			w.Write([]byte("}"))
		case reflect.Slice:
			w.Write([]byte(fmt.Sprintf("%q:", tag)))
			w.Write([]byte("["))
			numItems := field.Len()
			for x := 0; x < numItems; x++ {
				w.Write([]byte("{"))
				m.reflectStructFields(w, field.Index(x))
				w.Write([]byte("}"))
				if x+1 < numItems {
					w.Write([]byte(","))
				}
			}
			w.Write([]byte("]"))
		default:
			w.Write([]byte(fmt.Sprintf("%q:%q", tag, field.String())))
		}
		if i+1 < len(fields) {
			w.Write([]byte(","))
		}
		i++
	}
}

func (m *CustomMarshaler) Marshal(obj interface{}) ([]byte, error) {
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
