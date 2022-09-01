package structx

import (
	"fmt"
	"reflect"
)

func SyncWithNoZero(src interface{}, dst interface{}) (err error) {
	tSrc := reflect.TypeOf(src)
	vSrc := reflect.ValueOf(src)

	tDst := reflect.TypeOf(dst).Elem()
	vDst := reflect.ValueOf(dst).Elem()
	for i := 0; i < tSrc.NumField(); i++ {
		if vSrc.Field(i).IsZero() {
			continue
		}
		if tSrc.Field(i).Type != tDst.Field(i).Type {
			return fmt.Errorf("field: %s is not the same", tSrc.Field(i).Name)
		}
		if !vDst.Field(i).CanSet() {
			return fmt.Errorf("field: %s cannot set", tSrc.Field(i).Name)
		}
		vDst.Field(i).Set(vSrc.Field(i))
	}
	return nil
}
