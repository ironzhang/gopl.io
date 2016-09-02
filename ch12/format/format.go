package format

import (
	"reflect"
	"strconv"
)

func Any(value interface{}) string {
	return formatAtom(reflect.ValueOf(value))
}

func formatAtom(v reflect.Value) string {
	switch v.Kind() {
	case reflect.Invalid:
		return "invalid"
	case reflect.Bool:
		return strconv.FormatBool(v.Bool())
	case reflect.String:
		return v.String()
	case reflect.Float32, reflect.Float64:
		return formatFloat(v.Float())
	case reflect.Complex64, reflect.Complex128:
		return formatComplex(v.Complex())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(v.Uint(), 10)
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.Slice, reflect.UnsafePointer:
		return v.Type().String() + " 0x" + strconv.FormatUint(uint64(v.Pointer()), 16)
	case reflect.Array, reflect.Struct, reflect.Interface:
		return v.Type().String() + " value"
	default:
		return "unknown: " + v.Type().String() + " value"
	}
}

func formatComplex(c complex128) string {
	r := formatFloat(real(c))
	i := formatFloat(imag(c))
	return r + "+" + i + "i"
}

func formatFloat(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}
