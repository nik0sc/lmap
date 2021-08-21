package lmap

import "reflect"

type Key interface {
	Hash() int
	Equal(Key) bool
}

type LinkedMapCustom interface {
	Get(Key) (interface{}, bool)
	Set(Key, interface{})
	Iter(func(Key, interface{})bool)
}

type lmap struct {
	tab []*entry
	kt reflect.Type
	vt reflect.Type
	head *entry
	tail *entry
}

type entry struct {
	k Key
	v interface{}
	next *entry
}

