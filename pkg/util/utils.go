package util

import (
	"fmt"
	"reflect"

	"google.golang.org/grpc"
)

func Invoke(any interface{}, name string, args ...interface{}) ([]reflect.Value, error) {
	method := reflect.ValueOf(any).MethodByName(name)
	methodType := method.Type()
	numIn := methodType.NumIn()

	if !methodType.IsVariadic() {
		if numIn != len(args) {
			return []reflect.Value{reflect.ValueOf(nil)}, fmt.Errorf("Method %s must have %d params. Have %d", name, numIn, len(args))
		}
	} else {
		if numIn-1 > len(args) {
			return []reflect.Value{reflect.ValueOf(nil)}, fmt.Errorf("Method %s must have minimum %d params. Have %d", name, numIn-1, len(args))
		}
	}

	in := make([]reflect.Value, len(args))
	for i := 0; i < len(args); i++ {
		var inType reflect.Type
		if methodType.IsVariadic() && i >= numIn-1 {
			inType = methodType.In(numIn - 1).Elem()
		} else {
			inType = methodType.In(i)
		}
		argValue := reflect.ValueOf(args[i])
		if !argValue.IsValid() {
			return []reflect.Value{reflect.ValueOf(nil)}, fmt.Errorf("Method %s. Param[%d] must be %s. Have %s", name, i, inType, argValue.String())
		}
		argType := argValue.Type()
		if argType.ConvertibleTo(inType) {
			in[i] = argValue.Convert(inType)
		} else {
			return []reflect.Value{reflect.ValueOf(nil)}, fmt.Errorf("Method %s. Param[%d] must be %s. Have %s", name, i, inType, argType)
		}
	}
	return method.Call(in), nil
}

// TestConnectivity test the given addresses one by one and return the
// first successful one, if all failed, return empty.
func TestConnectivity(addrs []string, port string) string {
	for _, addr := range addrs {
		var address string
		if port == "" {
			address = addr
		} else {
			address = fmt.Sprintf("%s:%s", addr, port)
		}
		conn, err := grpc.Dial(address, grpc.WithInsecure())
		if err == nil {
			conn.Close()
			return addr
		}
	}
	return ""
}