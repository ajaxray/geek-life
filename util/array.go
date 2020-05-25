package util

import "reflect"

// InArray checks is val exists in a Slice
func InArray(val interface{}, array interface{}) bool {
	return AtArrayPosition(val, array) != -1
}

// AtArrayPosition find the int position of val in a Slice
func AtArrayPosition(val interface{}, array interface{}) (index int) {
	index = -1

	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				index = i
				return
			}
		}
	}

	return
}
