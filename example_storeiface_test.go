package lmap_test

import (
	"fmt"

	"playground-1.16/dsa/lmap"
)

type I interface{
	F() string
}

type A int

func (A) F() string {
	return "A"
}

type B string

func (B) F() string {
	return "B"
}

type IWrapper struct {
	I
}

func ExampleLinkedMap_storeInterfaceType() {
	l := lmap.New()

	a := IWrapper{A('a')}
	b := IWrapper{B("b")}

	l.Set(a, 1, true)
	l.Set(b, 1, true)

	l.Iter(func(k, _ interface{}) bool {
		fmt.Println(k, k.(IWrapper).F())
		return true
	})

	fmt.Println("recorded type:", l.KeyType())

	// Output:
	// {97} A
	// {b} B
	// recorded type: lmap_test.IWrapper
}
