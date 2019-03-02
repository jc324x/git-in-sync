package quit

import (
	"os"
	"testing"
)

func checkB(a int, b int) Out {
	return Bool(a == b, []string{"%v != %v", "%v == %v"}, a, b)
}

func checkE(p string) Out {
	_, te := os.Stat(p)
	return Err(te, []string{"!= %v", "%v"}, p)
}

func TestBool(t *testing.T) {

	os.Setenv("MODE", "TESTING")

	for _, tr := range []struct {
		a, b  int
		wantb bool
		wants string
	}{
		{100, 101, false, "100 != 101"},
		{100, 100, true, "100 == 100"},
	} {

		got := checkB(tr.a, tr.b)

		if got.Status != tr.wantb {
			t.Errorf("Check: %v != %v", got.Status, tr.wantb)
		}

		if got.Summary != tr.wants {
			t.Errorf("Check: %v != %v", got.Summary, tr.wants)
		}
	}
}

func TestErr(t *testing.T) {

	os.Setenv("MODE", "TESTING")

	for _, tr := range []struct {
		in    string
		wantb bool
		wants string
	}{
		{"/a/fake/path", false, "!= /a/fake/path"},
	} {

		got := checkE(tr.in)

		if got.Status != tr.wantb {
			t.Errorf("Check: %v != %v", got.Status, tr.wantb)
		}

		if got.Summary != tr.wants {
			t.Errorf("Check: %v != %v", got.Summary, tr.wants)
		}
	}

}
