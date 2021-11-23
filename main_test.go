package main

import (
	// "io"
	// "log"
	// "strings"
	"testing"
)

func TestGetNew(t *testing.T) {
	a := Item{Name:"test1", Price:"P100"}
	b := Item{Name:"test2", Price:"P200"}
	c := Item{Name:"test3", Price:"P300"}
	d := Item{Name:"test4", Price:"P400"}
	e := Item{Name:"test5", Price:"P500"}
	f := Item{Name:"test6", Price:"P600"}
	g := Item{Name:"test7", Price:"P700"}
	h := Item{Name:"test8", Price:"P800"}
	i := Item{Name:"test9", Price:"P900"}
	j := Item{Name:"test10", Price:"P1000"}

	// loop through slice to see if they're the same
	isSameSlice := func (a, b []Item) bool {
		if len(a) != len(b) {
			return false
		}
		
		for i := len(a); i < len(a); i++ {
			if a[i] != b[i]{
				return false
			}
		}
		return true
	}

	t.Run("Test with 4 items, 1 new", func(t *testing.T) {

		old := []Item{c, b, a}
		new := []Item{d, c, b, a}

		got := CompareItems(old, new)
		want := []Item{d}

		
		if got[0] != want[0] {
			t.Errorf("Got %s, want %s", got, want)
		}
	})

	t.Run("Test with 3 items, 3 new", func(t *testing.T) {
		old := []Item{f,e,d,c,b,a}
		new := []Item{j,i,h,g,f,e}

		got := CompareItems(old, new)
		want := []Item{j,i,h,g}
		
		if !isSameSlice(got, want) {
			t.Errorf("Got %s, want %s", got, want)
		}
	})
	t.Run("Test with 3 items, 6 new", func(t *testing.T) {
		old := []Item{c,b,a}
		new := []Item{i,h,g,f,e,d,c,b,a}

		got := CompareItems(old, new)
		want := []Item{i,h,g,f,e,d}
		
		if !isSameSlice(got, want) {
			t.Errorf("Got %s, want %s", got, want)
		}
	})
	t.Run("Test with 6 items, 3 new, same len", func(t *testing.T) {

		old := []Item{f,e,d,c,b,a}
		new := []Item{i,h,g,f,e,d}

		got := CompareItems(old, new)
		want := []Item{i,h,g}
		
		if !isSameSlice(got, want) {
			t.Errorf("Got %s, want %s", got, want)
		}
	})
	t.Run("Test with 4 items, all completely new, same len", func(t *testing.T) {

		old := []Item{d,c,b,a}
		new := []Item{i,h,g,f}

		got := CompareItems(old, new)
		want := []Item{i,h,g,f}
		
		if !isSameSlice(got, want) {
			t.Errorf("Got %s, want %s", got, want)
		}
	})
}

