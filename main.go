package main

import (
	"fmt"
	"reflect"
)

func KeyBy(keyBy string, src interface{}, dest interface{}) {
	s := reflect.ValueOf(src)
	d := reflect.ValueOf(dest)
	switch {
	case s.Kind() != reflect.Slice:
		panic("KeyBy() given a non-slice type")
	case d.Kind() != reflect.Map:
		panic("KeyBy() given a non-map type")
	}
	for i := 0; i < s.Len(); i++ {
		value := s.Index(i)
		d.SetMapIndex(value.FieldByName(keyBy), value)
	}
}

type tb struct {
	ID   int
	Name string
}

func main() {
	a := []tb{{1, "11"}, {2, "22"}}
	b := map[int]tb{}
	KeyBy("ID", a, b)
	fmt.Printf("%v\n", b)
}
