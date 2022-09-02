package structx

import (
	"reflect"
	"strings"
	"time"
)

func SyncWithNoZero(src interface{}, dst interface{}) (err error) {
	tSrc := reflect.TypeOf(src)
	vSrc := reflect.ValueOf(src)

	tDst := reflect.TypeOf(dst).Elem()
	vDst := reflect.ValueOf(dst).Elem()
	srcNumField := tSrc.NumField()
	dstNumField := tDst.NumField()
	vis := make([]bool, 100)
	for i := 0; i < srcNumField; i++ {
		if vSrc.Field(i).IsZero() || !tSrc.Field(i).IsExported() {
			continue
		}
		for j := 0; j < dstNumField; j++ {
			if vis[j] {
				continue
			}
			if !tDst.Field(j).IsExported() {
				continue
			}
			if strings.ToLower(tSrc.Field(i).Name) != strings.ToLower(tDst.Field(j).Name) {
				continue
			}
			if tSrc.Field(i).Type != tDst.Field(j).Type {
				if tSrc.Field(i).Type == reflect.TypeOf(time.Time{}) {
					if !vDst.Field(j).CanSet() {
						continue
					}
					vDst.Field(j).Set(reflect.ValueOf(vSrc.Field(i).Interface().(time.Time).String()))
					vis[j] = true
					break
				} else if tDst.Field(j).Type == reflect.TypeOf(time.Time{}) {
					if !vDst.Field(j).CanSet() {
						continue
					}
					t, _ := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", vSrc.Field(i).Interface().(string))
					vDst.Field(j).Set(reflect.ValueOf(t))
					vis[j] = true
					break
				}
			}
			if !vDst.Field(j).CanSet() {
				continue
			}
			vDst.Field(j).Set(vSrc.Field(i))
			vis[j] = true
			break
		}
	}
	return nil
}
