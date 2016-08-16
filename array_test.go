package objconv

import "testing"

func TestArraySlice(t *testing.T) {
	s := []int{1, 2, 3}
	a := NewArraySlice(s)

	if n := a.Len(); n != 3 {
		t.Errorf("invalid array length: %d", n)
	}

	it := a.Iter()

	for i := 0; true; i++ {
		if v, ok := it.Next(); !ok {
			if i != 3 {
				t.Errorf("array iterator returned too early: %d", i)
			}
			break
		} else {
			if v.(int) != s[i] {
				t.Errorf("invalid value produced by array iterator: %d != %d", s[i], v.(int))
			}
		}
	}
}
