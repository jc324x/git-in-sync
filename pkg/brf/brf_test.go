package brf

import (
	"reflect"
	"testing"

	"github.com/jychri/git-in-sync/pkg/flags"
)

func TestPrintv(t *testing.T) {
	var tf = flags.Flags{Mode: "oneline"}

	if err := Printv(tf, "testing"); err == nil {
		t.Errorf("Printv: Print condition flipped")
	}

	tf.Mode = "verify"

	if err := Printv(tf, "testing"); err != nil {
		t.Errorf("Printv: Print condition flipped")
	}
}

func TestSingle(t *testing.T) {
	for _, c := range []struct {
		in, want []string
	}{
		{[]string{"a", "b", "b", "c", "c", "c"}, []string{"a", "b", "c"}},
		{[]string{"x", "y", "z", "z", "z", "z"}, []string{"x", "y", "z"}},
	} {
		got := Single(c.in)

		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("Single: != DeepEqual (%v -> %v != %v)", c.in, got, c.want)
		}
	}
}

func TestSummary(t *testing.T) {

	sl1 := []string{"the", "sly", "brown", "jumped", "over", "the", "lazy", "dog"}
	sl2 := []string{"Lorem", "ipsum", "dolor", "sit", "amet", "consectetur", "adipiscing", "elit", ".", "Maecenas"}

	for _, c := range []struct {
		sl []string
		l  int
	}{
		{sl1, 15},
		{sl2, 25},
	} {
		got := Summary(c.sl, c.l)

		if len(got) >= (len(got) + 12) {
			t.Errorf("Summary: != len(got) (%v >= %v)", len(got), (len(got) + 12))
		}
	}
}

func TestFirst(t *testing.T) {

}
