package game

import "reflect"

func mapsEqual(a, b Map) bool {
	return reflect.DeepEqual(a, b)
}
