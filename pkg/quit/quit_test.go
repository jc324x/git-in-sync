package quit

import (
	"os"
	"testing"
)

// check is an example of how to implement Eval
func check(a int, b int) (bool, string) {
	return Eval(a == b, []string{"%v != %v", "%v == %v"}, a, b)
}

func TestCheck(t *testing.T) {

	os.Setenv("MODE", "TESTING")

	for _, c := range []struct {
		a, b  int
		wantb bool
		wants string
	}{
		{100, 101, false, "100 != 101"},
		{100, 100, true, "100 == 100"},
	} {

		gotb, gots := check(c.a, c.b)

		if gotb != c.wantb {
			t.Errorf("Check: %v != %v", gotb, c.wantb)
		}

		if gots != c.wants {
			t.Errorf("Check: %v != %v", gots, c.wants)
		}
	}
}
