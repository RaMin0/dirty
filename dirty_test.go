package dirty_test

import (
	"testing"

	"github.com/ramin0/dirty"
)

type User struct {
	ID      int
	IDPtr   *int
	Name    string
	NamePtr *string
	Slice   []bool
	Map     map[int]bool
}

func TestChanged(t *testing.T) {
	u := User{}

	dirty.Track(&u)

	if dirty.Changed(&u) {
		t.Fail()
	}

	u.Name = "John Doe"

	if !dirty.Changed(&u) {
		t.Fail()
	}

	dirty.Track(&u)

	u.Name = "John Doe"

	if dirty.Changed(&u) {
		t.Fail()
	}
}
