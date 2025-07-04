package utils

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"sort"
)

var timestampType = reflect.TypeOf(Timestamp{})

func Stringify(message interface{}) string {
	var buf bytes.Buffer
	v := reflect.ValueOf(message)
	stringifyValue(&buf, v)
	return buf.String()
}

func stringifyValue(w io.Writer, val reflect.Value) {
	if val.Kind() == reflect.Ptr && val.IsNil() {
		_, _ = w.Write([]byte("<nil>"))
		return
	}

	v := reflect.Indirect(val)

	switch v.Kind() {
	case reflect.String:
		fmt.Fprintf(w, `"%s"`, v)
	case reflect.Slice:
		stringifySlice(w, v)
		return
	case reflect.Struct:
		stringifyStruct(w, v)
	case reflect.Map:
		stringifyMap(w, v)
	default:
		if v.CanInterface() {
			fmt.Fprint(w, v.Interface())
		}
	}
}

func stringifySlice(w io.Writer, v reflect.Value) {
	_, _ = w.Write([]byte{'['})
	for i := 0; i < v.Len(); i++ {
		if i > 0 {
			_, _ = w.Write([]byte{' '})
		}

		stringifyValue(w, v.Index(i))
	}

	_, _ = w.Write([]byte{']'})
}

func stringifyMap(w io.Writer, v reflect.Value) {
	_, _ = w.Write([]byte("map["))

	// Sort the keys so that the output is stable
	keys := v.MapKeys()
	sort.Slice(keys, func(i, j int) bool {
		return fmt.Sprintf("%v", keys[i]) < fmt.Sprintf("%v", keys[j])
	})

	for i, key := range keys {
		stringifyValue(w, key)
		_, _ = w.Write([]byte{':'})
		stringifyValue(w, v.MapIndex(key))
		if i < len(keys)-1 {
			_, _ = w.Write([]byte(", "))
		}
	}

	_, _ = w.Write([]byte("]"))
}

func stringifyStruct(w io.Writer, v reflect.Value) {
	if v.Type().Name() != "" {
		_, _ = w.Write([]byte(v.Type().String()))
	}

	// special handling of Timestamp values
	if v.Type() == timestampType {
		fmt.Fprintf(w, "{%s}", v.Interface())
		return
	}

	_, _ = w.Write([]byte{'{'})

	var sep bool
	for i := 0; i < v.NumField(); i++ {
		fv := v.Field(i)
		if fv.Kind() == reflect.Ptr && fv.IsNil() {
			continue
		}
		if fv.Kind() == reflect.Slice && fv.IsNil() {
			continue
		}

		if sep {
			_, _ = w.Write([]byte(", "))
		} else {
			sep = true
		}

		_, _ = w.Write([]byte(v.Type().Field(i).Name))
		_, _ = w.Write([]byte{':'})
		stringifyValue(w, fv)
	}

	_, _ = w.Write([]byte{'}'})
}
