package main

import "testing"

func TestNewSplitString(t *testing.T) {
	input := "this / is / a /test"
	expectLength := 4
	ss := NewSplitString(input, "/")
	lenTest := ss.len == expectLength
	if !lenTest {
		t.Fatalf("expected length %d got length %d %v", expectLength, ss.len, ss.inner)
	}
}

func TestSplitStringGet(t *testing.T) {
	input := "this / is / a /test"
	expect := []string{"this ", " is ", " a ", "test"}
	ss := NewSplitString(input, "/")
	for i := range expect {
		if ss.Get(i) != expect[i] {
			t.Fatalf("segment mismatch %d expected '%s' got '%s'", i, ss.Get(i), expect[i])
		}
	}
	result := ss.Get(-2)
	expect2 := expect[2]
	if result != expect2 {
		t.Fatalf("negative index fail expected '%s' got '%s'", expect2, result)
	}
}

func TestSplitStringReassemble(t *testing.T) {

}
