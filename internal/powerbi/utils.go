package powerbi

import (
	"net/url"
	"reflect"

	"github.com/codecutout/terraform-provider-powerbi/internal/powerbiapi"
)

func convertStringToPointer(s string) *string {
	return &s
}

func convertBoolToPointer(b bool) *bool {
	return &b
}

func convertStringSliceToPointer(ss []string) *[]string {
	return &ss
}

func convertToStringSlice(interfaceSlice []interface{}) []string {
	stringSlice := make([]string, len(interfaceSlice))
	for i := range interfaceSlice {
		stringSlice[i] = interfaceSlice[i].(string)
	}
	return stringSlice
}

func nilIfFalse(b bool) *bool {
	if !b {
		return nil
	}
	return &b
}

func emptyStringToNil(input string) *string {
	if input == "" {
		return nil
	}
	return &input
}

func isHTTP404Error(err error) bool {
	if httpErr, isHTTPErr := toHTTPUnsuccessfulError(err); isHTTPErr && httpErr.Response.StatusCode == 404 {
		return true
	}
	return false
}

func isHTTP401Error(err error) bool {
	if httpErr, isHTTPErr := toHTTPUnsuccessfulError(err); isHTTPErr && httpErr.Response.StatusCode == 401 {
		return true
	}
	return false
}

func toHTTPUnsuccessfulError(err error) (*powerbiapi.HTTPUnsuccessfulError, bool) {
	if err == nil {
		return nil, false
	}

	if urlErr, isURLErr := err.(*url.Error); isURLErr {
		err = urlErr.Unwrap()
	}

	if httpErr, isHTTPErr := err.(powerbiapi.HTTPUnsuccessfulError); isHTTPErr {
		return &httpErr, true
	}
	return nil, false
}

type wrappedError struct {
	Err          error
	ErrorMessage func(err error) string
}

func (e wrappedError) Error() string {
	return e.ErrorMessage(e.Err)
}

// Courtesy of https://gist.github.com/ParthDesai/5e0f1d4725a644f1e632
func genericMap(arr interface{}, mapFunc interface{}) interface{} {
	funcValue := reflect.ValueOf(mapFunc)
	arrValue := reflect.ValueOf(arr)

	// Retrieve the type, and check if it is one of the array or slice.
	arrType := arrValue.Type()
	arrElemType := arrType.Elem()
	if arrType.Kind() != reflect.Array && arrType.Kind() != reflect.Slice {
		panic("Array parameter's type is neither array nor slice.")
	}

	funcType := funcValue.Type()

	// Checking whether the second argument is function or not.
	// And also checking whether its signature is func ({type A}) {type B}.
	if funcType.Kind() != reflect.Func || funcType.NumIn() != 1 || funcType.NumOut() != 1 {
		panic("Second argument must be map function.")
	}

	// Checking whether element type is convertible to function's first argument's type.
	if !arrElemType.ConvertibleTo(funcType.In(0)) {
		panic("Map function's argument is not compatible with type of array.")
	}

	// Get slice type corresponding to function's return value's type.
	resultSliceType := reflect.SliceOf(funcType.Out(0))

	// MakeSlice takes a slice kind type, and makes a slice.
	resultSlice := reflect.MakeSlice(resultSliceType, 0, arrValue.Len())

	for i := 0; i < arrValue.Len(); i++ {
		resultSlice = reflect.Append(resultSlice, funcValue.Call([]reflect.Value{arrValue.Index(i)})[0])
	}

	// Converting resulting slice back to generic interface.
	return resultSlice.Interface()
}

func intersect(as []interface{}, bs []interface{}, isEqualFn func(a interface{}, b interface{}) bool) []interface{} {
	result := make([]interface{}, 0)
	for _, a := range as {

		for _, b := range bs {
			if isEqualFn(a, b) {
				result = append(result, a)
				break
			}
		}
	}
	return result
}
