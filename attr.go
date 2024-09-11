package otelattr

import (
	"encoding/json"
	"fmt"
	"reflect"

	"go.opentelemetry.io/otel/attribute"
)

type Marshaler interface {
	MarshalOtelAttributes() ([]attribute.KeyValue, error)
}

func MarshalOtelAttributes(v interface{}) ([]attribute.KeyValue, error) {
	if m, ok := v.(Marshaler); ok {
		return m.MarshalOtelAttributes()
	}
	rv := reflect.ValueOf(v)
	return marshalOtelAttributes(rv)
}

func marshalOtelAttributes(rv reflect.Value) ([]attribute.KeyValue, error) {
	switch rv.Kind() {
	case reflect.Struct:
		return marshalStruct(rv)
	case reflect.Ptr:
		return marshalOtelAttributes(rv.Elem())
	case reflect.Interface:
		return marshalOtelAttributes(rv.Elem())
	default:
		return nil, fmt.Errorf("unsupported type %s", rv.Type())
	}
}

func marshalStruct(rv reflect.Value) ([]attribute.KeyValue, error) {
	t := rv.Type()
	fields := getStructFields(t)
	kvs := make([]attribute.KeyValue, 0, len(fields))
	for _, f := range fields {
		fv := rv.Field(f.filedIndex)
		if f.omitEmpty && isEmptyValue(fv) {
			continue
		}
		filedValue, err := marshalField(f, fv)
		if err != nil {
			return nil, err
		}
		kvs = append(kvs, filedValue...)
	}
	return kvs, nil
}

func marshalField(f structFiled, fv reflect.Value) ([]attribute.KeyValue, error) {
	switch fv.Kind() {
	case reflect.Bool:
		return []attribute.KeyValue{attribute.Bool(f.attributeName, fv.Bool())}, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return []attribute.KeyValue{attribute.Int64(f.attributeName, fv.Int())}, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return []attribute.KeyValue{attribute.Int64(f.attributeName, int64(fv.Uint()))}, nil
	case reflect.Float32, reflect.Float64:
		return []attribute.KeyValue{attribute.Float64(f.attributeName, fv.Float())}, nil
	case reflect.String:
		return []attribute.KeyValue{attribute.String(f.attributeName, fv.String())}, nil
	case reflect.Slice:
		return marshalSlice(f, fv)
	default:
		bs, err := json.Marshal(fv.Interface())
		if err != nil {
			return []attribute.KeyValue{}, err
		}
		return []attribute.KeyValue{attribute.String(f.attributeName, string(bs))}, nil
	}
}

func marshalSlice(f structFiled, fv reflect.Value) ([]attribute.KeyValue, error) {
	switch fv.Type().Elem().Kind() {
	case reflect.Bool:
		return []attribute.KeyValue{attribute.BoolSlice(f.attributeName, reflectValueToSlice[bool](fv))}, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return []attribute.KeyValue{attribute.Int64Slice(f.attributeName, reflectValueToSlice[int64](fv))}, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return []attribute.KeyValue{attribute.Int64Slice(f.attributeName, reflectValueToSlice[int64](fv))}, nil
	case reflect.Float32, reflect.Float64:
		return []attribute.KeyValue{attribute.Float64Slice(f.attributeName, reflectValueToSlice[float64](fv))}, nil
	case reflect.String:
		return []attribute.KeyValue{attribute.StringSlice(f.attributeName, reflectValueToSlice[string](fv))}, nil
	default:
		bs, err := json.Marshal(fv.Interface())
		if err != nil {
			return []attribute.KeyValue{}, err
		}
		return []attribute.KeyValue{attribute.String(f.attributeName, string(bs))}, nil
	}
}

func reflectValueToSlice[T any](v reflect.Value) []T {
	slice := make([]T, v.Len())
	for i := 0; i < v.Len(); i++ {
		slice[i] = v.Index(i).Interface().(T)
	}
	return slice
}

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.String:
		return v.String() == ""
	case reflect.Slice:
		return v.Len() == 0
	default:
		return false
	}
}
