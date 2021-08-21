package lmap

import (
	"fmt"
	"reflect"
)

const (
	flagDebug = 1 << iota
	flagIterating
)

// LinkedMap is a map combined with a linked list. It preserves insertion order and therefore
// iteration order as well.
// Any type of key may be stored, as long as the type is comparable and the map only contains
// one type of key. Any type of value may be stored (including different types in the same map).
// LinkedMap is not safe for concurrent use.
type LinkedMap struct {
	m          map[interface{}]*entryb
	kt         reflect.Type
	head, tail *entryb
	flags      int
}

type entryb struct {
	k, v       interface{}
	prev, next *entryb
}

// Not used, but may be called for debug purposes
// todo: return cycle start and size too
func (e *entryb) cycle() bool {
	tortoise, hare := e, e

	for hare != nil && hare.next != nil {
		tortoise = tortoise.next
		hare = hare.next.next
		if tortoise == hare {
			return true
		}
	}

	return false
}

// New returns a pointer to a new LinkedMap.
func New() *LinkedMap {
	return &LinkedMap{
		m: make(map[interface{}]*entryb),
	}
}

func (l *LinkedMap) assertKeyType(k interface{}) {
	if l.kt != reflect.TypeOf(k) {
		panic(fmt.Sprintf("incompatible types: map key=%s, incoming=%T", l.kt, k))
	}
}

func (l *LinkedMap) assertNotIterating() {
	if l.flags&flagIterating > 0 {
		panic("about to mutate linked list while iterating over it")
	}
}

func (l *LinkedMap) remove(e *entryb) {
	if e == nil {
		panic("nil entry")
	}

	if l.head == nil || l.tail == nil {
		panic("nil head or tail")
	}

	l.assertNotIterating()

	if e.prev != nil {
		e.prev.next = e.next
	} else {
		if l.head != e {
			panic("entry has no previous node but it is not the head")
		}
		l.head = e.next
	}

	if e.next != nil {
		e.next.prev = e.prev
	} else {
		if l.tail != e {
			panic("entry has no next node but it is not the tail")
		}
		l.tail = e.prev
	}
}

func (l *LinkedMap) push(e *entryb) {
	if e == nil {
		panic("nil entry")
	}

	l.assertNotIterating()

	if l.head == nil && l.tail == nil {
		l.head, l.tail = e, e
		return
	}

	e.prev = l.tail
	l.tail.next = e

	e.next = nil
	l.tail = e
}

// KeyType returns the recorded concrete type for this map's keys.
// If a key has not been stored, the type is nil.
func (l *LinkedMap) KeyType() reflect.Type {
	return l.kt
}

// Get behaves like the map access `v, ok := l[k]`. If bump is true and k is in the map,
// k is moved to the tail of the list, as if it were removed and added back to the map.
func (l *LinkedMap) Get(k interface{}, bump bool) (v interface{}, ok bool) {
	if l.kt == nil {
		return nil, false
	}

	l.assertKeyType(k)

	e, ok := l.m[k]
	if !ok {
		return
	}

	if bump {
		l.remove(e)
		l.push(e)
	}

	return e.v, true
}

// Set behaves like the map set `l[k] = v`. If bumpOnExist is true and k is in the map,
// k is moved to the tail of the list, as if it were removed and added back into the map.
// Otherwise, if k is not in the map, it is appended to the tail of the list.
// k must be of a Comparable type (i.e. no maps, slices, functions, or structs embedding those types).
// On the first time any k is added to the map, k's dynamic type is recorded.
// Subsequent additions of other ks must match the recorded type.
// Note that the recorded type persists even when the map is emptied.
//
// It is currently not possible to record an interface type instead of the concrete type.
// One workaround is to define a wrapper struct that embeds the desired interface type.
// While awkward, this solution does not use any extra memory. See example_storeiface_test.go.
func (l *LinkedMap) Set(k, v interface{}, bumpOnExist bool) {
	if l.kt == nil {
		kt := reflect.TypeOf(k)
		if !kt.Comparable() {
			panic(fmt.Sprintf("incomparable type: %s", kt))
		}
		l.kt = kt
	} else {
		l.assertKeyType(k)
	}

	e, exist := l.m[k]
	if exist {
		if e.k != k {
			panic("entry key does not match map key")
		}

		e.v = v
		if bumpOnExist {
			l.remove(e)
			l.push(e)
		}
	} else {
		e = &entryb{
			k: k,
			v: v,
		}

		l.m[k] = e

		l.push(e)
	}
}

// Delete behaves like `delete(l, k)`. If the key was not found, ok will be false.
func (l *LinkedMap) Delete(k interface{}) (ok bool) {
	if l.kt == nil {
		return
	}

	l.assertKeyType(k)

	e, ok := l.m[k]
	if !ok {
		return
	}

	l.remove(e)
	delete(l.m, k)

	return
}

// Iter allows ordered iteration over the map in the same vein as `for k, v := range l {}`.
// The function f is called for every key-value pair in order. If f returns false at any
// iteration, the iteration process is stopped early.
//
// The result of modifying the map while iterating over the map is undefined.
func (l *LinkedMap) Iter(f func(k, v interface{}) bool) {
	if l.head == nil {
		return
	}
	l.flags |= flagIterating

	hare := l.head.next

	for e := l.head; e != nil; e = e.next {
		if e == hare {
			// bug in the map, not in the caller
			panic("cycle detected, iteration will not end")
		}

		if !f(e.k, e.v) {
			break
		}

		if hare != nil && hare.next != nil {
			hare = hare.next.next
		} else {
			// hare has reached the end, iteration will too
			// e will never be nil
			hare = nil
		}
	}

	l.flags &^= flagIterating
}

// Len behaves like `len(l)`. This is a constant-time operation.
func (l *LinkedMap) Len() int {
	return len(l.m)
}

// Head returns the head element of the linked list. If pop is true, the head element
// is also removed from the map and list. If ok is false, no element was found.
func (l *LinkedMap) Head(pop bool) (k, v interface{}, ok bool) {
	if l.head == nil {
		return
	}

	k, v, ok = l.head.k, l.head.v, true

	if pop {
		l.remove(l.head)
		delete(l.m, k)
	}

	return
}

// Tail returns the tail element of the linked list. If pop is true, the tail element
// is also removed from the map and list. If ok is false, no element was found.
func (l *LinkedMap) Tail(pop bool) (k, v interface{}, ok bool) {
	if l.tail == nil {
		return
	}

	k, v, ok = l.tail.k, l.tail.v, true

	if pop {
		l.remove(l.tail)
		delete(l.m, k)
	}

	return
}
