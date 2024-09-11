package otelattr

import (
	"testing"

	"go.opentelemetry.io/otel/attribute"
)

type testWithMarshaller struct {
	Value int
}

func (t testWithMarshaller) MarshalOtelAttributes() ([]attribute.KeyValue, error) {
	return []attribute.KeyValue{
		attribute.Int("http.requests", t.Value),
	}, nil
}

func TestMarshalOtelAttributes(t *testing.T) {
	tests := []struct {
		name    string
		args    any
		want    []attribute.KeyValue
		wantErr bool
	}{
		{
			name: "marshal primitive types",
			args: struct {
				BoolValue   bool
				BoolSlice   []bool
				IntValue    int
				IntSlice    []int
				Int64Value  int64
				Int64Slice  []int64
				FloatValue  float64
				FloatSlice  []float64
				StringValue string
				StringSlice []string
			}{
				BoolValue:   true,
				IntValue:    1,
				Int64Value:  2,
				FloatValue:  3.14,
				StringValue: "hello",
				StringSlice: []string{"hello", "world"},
			},
			want: []attribute.KeyValue{
				attribute.Bool("bool_value", true),
				attribute.BoolSlice("bool_slice", []bool{}),
				attribute.Int64("int_value", 1),
				attribute.Int64Slice("int_slice", []int64{}),
				attribute.Int64("int64_value", 2),
				attribute.Int64Slice("int64_slice", []int64{}),
				attribute.Float64("float_value", 3.14),
				attribute.Float64Slice("float_slice", []float64{}),
				attribute.String("string_value", "hello"),
				attribute.StringSlice("string_slice", []string{"hello", "world"}),
			},
			wantErr: false,
		},
		{
			name: "with tag",
			args: struct {
				Bool   bool    `otelattr:"b"`
				Int    int     `otelattr:"i"`
				Int64  int64   `otelattr:"i64"`
				Float  float64 `otelattr:"f"`
				String string  `otelattr:"s"`
			}{
				Bool:   true,
				Int:    1,
				Int64:  2,
				Float:  3.14,
				String: "hello",
			},
			want: []attribute.KeyValue{
				attribute.Bool("b", true),
				attribute.Int64("i", 1),
				attribute.Int64("i64", 2),
				attribute.Float64("f", 3.14),
				attribute.String("s", "hello"),
			},
			wantErr: false,
		},
		{
			name: "with tag and omitempty",
			args: struct {
				Bool   bool    `otelattr:"bool,omitempty"`
				Int    int     `otelattr:"int,omitempty"`
				Int64  int64   `otelattr:"int64,omitempty"`
				Float  float64 `otelattr:"float,omitempty"`
				String string  `otelattr:"string,omitempty"`
			}{},
			want:    []attribute.KeyValue{},
			wantErr: false,
		},
		{
			name: "with struct marshaller",
			args: testWithMarshaller{Value: 10},
			want: []attribute.KeyValue{
				attribute.Int("http.requests", 10),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MarshalOtelAttributes(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalOtelAttributes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("MarshalOtelAttributes() not equal length, got = %v, want %v", len(got), len(tt.want))
				return
			}
			for i, gotkv := range got {
				wantkv := tt.want[i]
				if gotkv.Key != wantkv.Key {
					t.Errorf("MarshalOtelAttributes() key not equal, got = %v, want %v", gotkv.Key, wantkv.Key)
				}
				if gotkv.Value != wantkv.Value {
					t.Errorf("MarshalOtelAttributes() value not equal, got = %v, want %v", gotkv.Value, wantkv.Value)
				}
			}
		})
	}
}
