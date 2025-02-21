package vault

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

func sanitised(path string) string {
	return strings.TrimPrefix(strings.TrimSuffix(strings.TrimSpace(path), "/"), "/")
}

func isLeafRoot(data interface{}) bool {
	switch d := data.(type) {

	case map[string]interface{}:
		for _, v := range d {
			switch v.(type) {
			case string:
				return true
			case int:
				return true
			case []interface{}:
				return true
			}
		}
		return false

	default:
		return false
	}
}

func equal(l, r interface{}) bool {
	if reflect.TypeOf(l) != reflect.TypeOf(r) {
		return false
	}

	switch val := l.(type) {

	case int:
		return val == r.(int)

	case string:
		return val == r.(string)

	case []interface{}:
		arr := r.([]interface{})
		if len(val) != len(arr) {
			return false
		}
		for idx, elm := range val {
			if !equal(elm, arr[idx]) {
				return false
			}
		}

	default:
		fmt.Fprintf(os.Stderr, "WARNING: unexpected type: %s\n",
			reflect.TypeOf(l).String(),
		)

	}

	return false
}
