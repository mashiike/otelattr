package otelattr_test

import (
	"encoding/json"
	"os"

	"github.com/mashiike/otelattr"
)

func Example() {
	type HTTPContext struct {
		Status int    `otelattr:"http.status_code"`
		Method string `otelattr:"http.method"`
		Path   string `otelattr:"http.path"`
	}
	httpCtx := HTTPContext{
		Status: 200,
		Method: "GET",
		Path:   "/",
	}
	attrs, err := otelattr.MarshalOtelAttributes(httpCtx)
	if err != nil {
		panic(err)
	}
	enc := json.NewEncoder(os.Stdout)
	for _, attr := range attrs {
		if err := enc.Encode(attr); err != nil {
			panic(err)
		}
	}
	// Output:
	//{"Key":"http.status_code","Value":{"Type":"INT64","Value":200}}
	//{"Key":"http.method","Value":{"Type":"STRING","Value":"GET"}}
	//{"Key":"http.path","Value":{"Type":"STRING","Value":"/"}}
}
