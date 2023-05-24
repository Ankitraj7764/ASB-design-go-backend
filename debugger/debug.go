package debugger

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

func Debug(param interface{}) {
	pc, _, _, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc)
	split := strings.Split(fn.Name(), ".")
	_ = split[len(split)-1]

	val := reflect.ValueOf(param)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	fmt.Printf("%s : %v\n", "paramName", val.Interface())
}
