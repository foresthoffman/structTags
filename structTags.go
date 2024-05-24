// Package structTags allows for marshalling non-JSON and third-party struct
// tags. The standard JSON.Marshal() function cannot marshal custom tags, however
// this package can.
//
// e.g. Third-party struct tags could look like the "custom" tag below.
//
//  type MyStruct struct {
//    Field string `json:"api_field" custom:"custom_field"`
//    Ignored string `json:"ignored" custom:"-"`
//  }
package structTags

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"reflect"
	"sort"
	"strconv"
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

func (m *CustomMarshaller) marshal(w io.Writer, obj interface{}, top bool) error {
	if obj == nil {
		return errors.New(ErrNilObject)
	}

	v := reflect.ValueOf(obj)
	t := reflect.TypeOf(obj)
	if val, ok := obj.(reflect.Value); ok {
		v = val
		t = val.Type()
	}
	k := t.Kind()

	if k == reflect.Struct {
		// Remove ignored fields and ensure that the order of the fields won't change.
		fields := map[int]fieldMetadata{}
		count := 0
		for i := 0; i < v.NumField(); i++ {
			tagValue := v.Type().Field(i).Tag.Get(m.TargetTag)
			if tagValue == m.IgnoreTagWithValue {
				continue
			}
			fields[count] = fieldMetadata{
				TagValue: tagValue,
				Value:    v.Field(i),
			}
			count++
		}

		_, err := w.Write([]byte("{"))
		if err != nil {
			return err
		}
		for x := 0; x < len(fields); x++ {
			_, err = w.Write([]byte(fmt.Sprintf("%q:", fields[x].TagValue)))
			if err != nil {
				return err
			}
			err = m.marshal(w, fields[x].Value, false)
			if err != nil {
				return fmt.Errorf("failed to marshal struct field: %s", err.Error())
			}
			if x+1 < len(fields) {
				_, err = w.Write([]byte(","))
				if err != nil {
					return err
				}
			}
		}
		_, err = w.Write([]byte("}"))
		if err != nil {
			return err
		}
	} else if k == reflect.Slice {
		_, err := w.Write([]byte("["))
		if err != nil {
			return err
		}
		for i := 0; i < v.Len(); i++ {
			if err != nil {
				return err
			}
			err = m.marshal(w, v.Index(i), false)
			if err != nil {
				return fmt.Errorf("failed to marshal slice element: %s", err.Error())
			}
			if i+1 < v.Len() {
				_, err = w.Write([]byte(","))
				if err != nil {
					return err
				}
			}
		}
		_, err = w.Write([]byte("]"))
		if err != nil {
			return err
		}
	} else if k == reflect.Map {
		_, err := w.Write([]byte("{"))
		if err != nil {
			return err
		}
		var keys []string
		for _, key := range v.MapKeys() {
			keys = append(keys, key.String())
		}
		sort.Strings(keys)
		for i := 0; i < len(keys); i++ {
			_, err = w.Write([]byte(fmt.Sprintf("%q:", keys[i])))
			if err != nil {
				return err
			}
			err = m.marshal(w, v.MapIndex(reflect.ValueOf(keys[i])), false)
			if err != nil {
				return fmt.Errorf("failed to marshal map field: %s", err.Error())
			}
			if i+1 < v.Len() {
				_, err = w.Write([]byte(","))
				if err != nil {
					return err
				}
			}
		}
		_, err = w.Write([]byte("}"))
		if err != nil {
			return err
		}
	} else if k == reflect.Ptr {
		err := m.marshal(w, v.Elem(), false)
		if err != nil {
			return fmt.Errorf("failed to marshal ptr: %s", err.Error())
		}
	} else if k == reflect.Interface {
		err := m.marshal(w, v.Interface(), false)
		if err != nil {
			return fmt.Errorf("failed to marshal interface: %s", err.Error())
		}
	} else if k == reflect.Int || k == reflect.Int8 || k == reflect.Int16 || k == reflect.Int32 || k == reflect.Int64 {
		_, err := w.Write([]byte(strconv.FormatInt(v.Int(), 10)))
		if err != nil {
			return err
		}
	} else if k == reflect.Uint || k == reflect.Uint8 || k == reflect.Uint16 || k == reflect.Uint32 || k == reflect.Uint64 || k == reflect.Uintptr {
		_, err := w.Write([]byte(strconv.FormatUint(v.Uint(), 10)))
		if err != nil {
			return err
		}
	} else if k == reflect.Float32 {
		_, err := w.Write([]byte(strconv.FormatFloat(v.Float(), 'f', -1, 32)))
		if err != nil {
			return err
		}
	} else if k == reflect.Float64 {
		_, err := w.Write([]byte(strconv.FormatFloat(v.Float(), 'f', -1, 64)))
		if err != nil {
			return err
		}
	} else if k == reflect.Complex64 || k == reflect.Complex128 {
		_, err := w.Write([]byte(fmt.Sprint(v.Complex())))
		if err != nil {
			return err
		}
	} else if k == reflect.String {
		_, err := w.Write([]byte(fmt.Sprintf("%q", v.String())))
		if err != nil {
			return err
		}
	} else if k == reflect.Bool {
		_, err := w.Write([]byte(fmt.Sprintf("%v", v.Bool())))
		if err != nil {
			return err
		}
	} else {
		_, err := w.Write([]byte(v.String()))
		if err != nil {
			return err
		}
	}

	// Perform top-level logic.
	if top {
		_, err := w.Write([]byte("\n"))
		if err != nil {
			return err
		}
		return nil
	}

	return nil
}

// Marshal takes the provided object and JSON-marshals it using the
// pre-configured target tag and ignored tag values.
func (m *CustomMarshaller) Marshal(obj interface{}) ([]byte, error) {
	w := bytes.NewBuffer([]byte{})
	err := m.marshal(w, obj, true)
	if err != nil {
		return nil, err
	}

	return w.Bytes(), nil
}
