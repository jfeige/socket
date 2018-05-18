package socket

import (
	rd "math/rand"
	"reflect"
	"strings"
	"time"
)

func generateAuthKey() string {
	src := RandStr(6)
	return src
}

/**
随机数(all,char,number)
*/
func RandStr(length int, format ...string) string {
	var tp = "all"
	if len(format) > 0 && format[0] != "" {
		tp = strings.ToLower(format[0])
	}
	var bytes []byte
	var r *rd.Rand
	var result []byte
	switch tp {
	case "char":
		bytes = []byte("abcdefghijklmnopqrstuvwxyz")
		if length > len(bytes) {
			length = len(bytes)
		}
		r = rd.New(rd.NewSource(time.Now().UnixNano()))
		for i := 0; i < length; i++ {
			result = append(result, bytes[r.Intn(len(bytes))])
		}
	case "number":
		bytes = []byte("0123456789")
		if length > len(bytes) {
			length = len(bytes)
		}
		r = rd.New(rd.NewSource(time.Now().UnixNano()))
		for i := 0; i < length; i++ {
			result = append(result, bytes[r.Intn(len(bytes))])
		}
	default:
		bytes = []byte("abcdefghijklmnopqrstuvwxyz0123456789")
		if length > len(bytes) {
			length = len(bytes)
		}
		r = rd.New(rd.NewSource(time.Now().UnixNano()))
		for i := 0; i < length; i++ {
			result = append(result, bytes[r.Intn(len(bytes))])
		}
	}

	return string(result)
}

func InArray(obj interface{}, target interface{}) bool {

	target_tp := reflect.TypeOf(target)
	target_vl := reflect.ValueOf(target)

	switch target_tp.Kind() {
	case reflect.Array, reflect.Slice:
		for i := 0; i < target_vl.Len(); i++ {
			if obj == target_vl.Index(i).Interface() {
				return true
			}
		}
	case reflect.Map:
		for _, v := range target_vl.MapKeys() {
			if obj == target_vl.MapIndex(v).Interface() {
				return true
			}
		}
	}

	return false
}
