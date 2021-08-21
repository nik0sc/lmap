package lmap

import (
	"strings"
	"testing"
)

func TestLinkedMap(t *testing.T) {
	l := New()

	t.Logf("start: %+v", l)

	l.Set("Hello", "World", true)

	t.Logf("hello: %+v", l)

	ok := l.Delete("Hello")

	if !ok {
		t.Error("failed to delete Hello")
	}

	t.Logf("empty: %+v", l)

	l.Set("Hello", "World", true)
	t.Logf("hello: %+v", l)

	k, v, ok := l.Head(true)
	if k != "Hello" || v != "World" || !ok {
		t.Error("failed to pop left")
	}

	t.Logf("empty: %+v", l)

	l.Set("foo", 1, true)
	l.Set("bar", 2, true)
	l.Set("baz", 3, true)

	var sb strings.Builder
	l.Iter(func(k interface{}, v interface{}) bool {
		sb.WriteString(k.(string))
		sb.WriteRune(' ')
		return true
	})
	t.Logf("before bump: %s", sb.String())
	t.Logf("lmap: %+v", l)

	sb.Reset()
	l.Set("foo", 2, true)
	t.Logf("lmap after: %+v", l)

	l.Iter(func(k interface{}, v interface{}) bool {
		sb.WriteString(k.(string))
		sb.WriteRune(' ')
		return true
	})
	t.Logf("after bump: %s", sb.String())
	t.Logf("lmap: %+v", l)

	k, v, ok = l.Head(true)
	if k != "bar" || v != 2 || !ok {
		t.Error("failed to pop left")
	}
	if l.head.cycle() {
		t.Error("cycle detected")
	}

	sb.Reset()
	l.Iter(func(k interface{}, v interface{}) bool {
		sb.WriteString(k.(string))
		sb.WriteRune(' ')
		return true
	})
	t.Logf("after pop left: %s", sb.String())
	t.Logf("lmap: %+v", l)
}

