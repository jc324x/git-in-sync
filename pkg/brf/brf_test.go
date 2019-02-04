package brf

import (
	"os/user"
	"reflect"
	"strings"
	"testing"
	// "github.com/jychri/git-in-sync/pkg/flags"
)

// func TestPrintv(t *testing.T) {
// 	var tf = flags.Flags{Mode: "oneline"}

// 	if err := Printv(tf, "testing"); err == nil {
// 		t.Errorf("Printv: Print condition flipped")
// 	}

// 	tf.Mode = "verify"

// 	if err := Printv(tf, "testing"); err != nil {
// 		t.Errorf("Printv: Print condition flipped")
// 	}
// }

func TestSingle(t *testing.T) {
	for _, c := range []struct {
		in, want []string
	}{
		{[]string{"a", "b", "b", "c", "c", "c", "c"}, []string{"a", "b", "c"}},
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
	sl3 := []string{"a", "b", "c", "d", "e", "f", "g"}

	for _, c := range []struct {
		sl   []string
		l    int
		want string
	}{
		{sl1, 15, "the, sly, brown, jumped..."},
		{sl2, 25, "Lorem, ipsum, dolor, sit, amet..."},
		{sl3, 20, "a, b, c, d, e, f, g"},
	} {

		got := Summary(c.sl, c.l)

		if len(got) >= (len(got) + 12) {
			t.Errorf("Summary: != len(got) (%v >= %v)", len(got), (len(got) + 12))
		}

		if got != c.want {
			t.Errorf("Summary: (%v != %v)", got, c.want)
		}

	}
}

func TestFirst(t *testing.T) {

	s1 := "First line.\n Second line. \n Third line.\n"
	s2 := "Only one line."
	s3 := ""

	for _, c := range []struct {
		in, want string
	}{
		{s1, "First line."},
		{s2, "Only one line."},
		{s3, ""},
	} {
		got := First(c.in)

		if got != c.want {
			t.Errorf("First: (%v != %v)", got, c.want)
		}
	}
}

func TestAbsUser(t *testing.T) {

	u, err := user.Current()

	if err != nil {
		t.Errorf("Relative: Can't identify current user")
	}

	for _, c := range []struct {
		in, want string
	}{
		{"~/testing", strings.Join([]string{u.HomeDir, "/testing"}, "")},
		{"~/testing/", strings.Join([]string{u.HomeDir, "/testing"}, "")},
		{"/testing/", "/testing"},
	} {
		if got := AbsUser(c.in); err != nil || got != c.want {
			t.Errorf("AbsUser: (%v != %v)", got, c.want)
		}
	}
}
